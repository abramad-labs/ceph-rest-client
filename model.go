package cephrestclient

type Image struct {
	Size            int64    `json:"size"`
	ObjSize         int64    `json:"obj_size"`
	NumObjs         int      `json:"num_objs"`
	Order           int      `json:"order"`
	BlockNamePrefix string   `json:"block_name_prefix"`
	MirrorMode      string   `json:"mirror_mode"`
	Name            string   `json:"name"`
	UniqueID        string   `json:"unique_id"`
	ID              string   `json:"id"`
	ImageFormat     int      `json:"image_format"`
	PoolName        string   `json:"pool_name"`
	Namespace       *string  `json:"namespace"`
	Features        int      `json:"features"`
	FeaturesName    []string `json:"features_name"`
	Timestamp       string   `json:"timestamp"`
	StripeCount     int      `json:"stripe_count"`
	StripeUnit      int64    `json:"stripe_unit"`
	DataPool        *string  `json:"data_pool"`
	Parent          struct {
		PoolName      string `json:"pool_name"`
		PoolNamespace string `json:"pool_namespace"`
		ImageName     string `json:"image_name"`
		SnapName      string `json:"snap_name"`
	} `json:"parent"`
	Snapshots      []any `json:"snapshots"`
	TotalDiskUsage int64 `json:"total_disk_usage"`
	DiskUsage      int64 `json:"disk_usage"`
	Configuration  []struct {
		Name   string `json:"name"`
		Value  string `json:"value"`
		Source int    `json:"source"`
	} `json:"configuration"`
}
