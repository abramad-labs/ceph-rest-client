package cephrestclient

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config codegen-config.yaml https://raw.githubusercontent.com/ceph/ceph/reef/src/pybind/mgr/dashboard/openapi.yaml

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
)

func NewAuthenticatedClient(server, username, password string) (*Client, error) {
	res, err := NewPostApiAuthRequest(server, PostApiAuthJSONRequestBody{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var AuthResponse struct{ Token string }
	err = json.Unmarshal(body, &AuthResponse)
	if err != nil {
		return nil, err
	}

	provider, err := securityprovider.NewSecurityProviderBearerToken(AuthResponse.Token)
	if err != nil {
		return nil, err
	}

	return NewClient(
		server,
		WithHTTPClient(&http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}),
		WithRequestEditorFn(provider.Intercept),
		WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Accept", "application/vnd.ceph.api.v1.0+json")
			req.Header.Set("Content-Type", "application/json")
			return nil
		}),
	)
}
