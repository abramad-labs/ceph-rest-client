package cephrestclient

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config codegen-config.yaml https://raw.githubusercontent.com/ceph/ceph/reef/src/pybind/mgr/dashboard/openapi.yaml

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
)

type NewAuthenticatedClientOpts struct {
	Username string
	Password string
	Server   string
	SkipTLS  bool
}

func CephAPIVersionHeaderIntercepter(version string) RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Accept", fmt.Sprintf("application/vnd.ceph.api.v%s+json", version))
		return nil
	}
}

func NewAuthenticatedClient(ctx context.Context, opts NewAuthenticatedClientOpts) (*Client, error) {
	client, err := NewClient(
		opts.Server,
		WithHTTPClient(&http.Client{Transport: &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
		}}),
		WithRequestEditorFn(CephAPIVersionHeaderIntercepter("1.0")),
		WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Content-Type", "application/json")
			return nil
		}),
	)
	if err != nil {
		return nil, err
	}

	res, err := client.PostApiAuth(ctx, PostApiAuthJSONRequestBody{
		Username: opts.Username,
		Password: opts.Password,
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code %d. respnse body: %s", res.StatusCode, string(body))
	}

	var AuthResponse struct{ Token string }
	err = json.Unmarshal(body, &AuthResponse)
	if err != nil {
		return nil, err
	}
	if AuthResponse.Token == "" {
		return nil, fmt.Errorf("expected non blank token in body: %s", string(body))
	}

	securityProvicder, err := securityprovider.NewSecurityProviderBearerToken(AuthResponse.Token)
	if err != nil {
		return nil, err
	}

	err = WithRequestEditorFn(securityProvicder.Intercept)(client)
	if err != nil {
		return nil, err
	}

	return client, nil
}
