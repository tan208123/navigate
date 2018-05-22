package cluster

import (
	"context"

	"github.com/tan208123/navigate/grpc/types"
)

type Cluster struct {
	// The cluster driver to provision cluster
	Driver types.Driver `json:"-"`
	// The name of the cluster driver
	DriverName string `json:"driverName,omitempty" yaml:"driver_name,omitempty"`
	// The name of the cluster
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// The status of the cluster
	Status string `json:"status,omitempty" yaml:"status,omitempty"`

	// specific info about kubernetes cluster
	// Kubernetes cluster version
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	// Service account token to access kubernetes API
	ServiceAccountToken string `json:"serviceAccountToken,omitempty" yaml:"service_account_token,omitempty"`
	// Kubernetes API master endpoint
	Endpoint string `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	// Username for http basic authentication
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	// Password for http basic authentication
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	// Root CaCertificate for API server(base64 encoded)
	RootCACert string `json:"rootCACert,omitempty" yaml:"root_ca_cert,omitempty"`
	// Client Certificate(base64 encoded)
	ClientCertificate string `json:"clientCertificate,omitempty" yaml:"client_certificate,omitempty"`
	// Client private key(base64 encoded)
	ClientKey string `json:"clientKey,omitempty" yaml:"client_key,omitempty"`
	// Node count in the cluster
	NodeCount int64 `json:"nodeCount,omitempty" yaml:"node_count,omitempty"`

	// Metadata store specific driver options per cloud provider
	Metadata map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	ConfigGetter ConfigGetter `json:"-" yaml:"-"`
}

// ConfigGetter defines the interface for getting the driver options.
type ConfigGetter interface {
	GetConfig() (types.DriverOptions, error)
}

func NewCluster(driverName, addr, name string, configGetter ConfigGetter) (*Cluster, error) {
	rpcClient, err := types.NewClient(driverName, addr)
	if err != nil {
		return nil, err
	}
	return &Cluster{
		Driver:       rpcClient,
		DriverName:   driverName,
		Name:         name,
		ConfigGetter: configGetter,
	}, nil
}

func (c *Cluster) Create(ctx context.Context) error {
	driverOpts, err := c.ConfigGetter.GetConfig()
	if err != nil {
		return err
	}
	// create cluster
	info, err := c.Driver.Create(ctx, &driverOpts)
	if err != nil {
		if info != nil {
			transformClusterInfo(c, info)
		}
		return err
	}

	transformClusterInfo(c, info)
	return nil
}

func transformClusterInfo(c *Cluster, clusterInfo *types.ClusterInfo) {
	c.ClientCertificate = clusterInfo.ClientCertificate
	c.ClientKey = clusterInfo.ClientKey
	c.RootCACert = clusterInfo.RootCaCertificate
	c.Username = clusterInfo.Username
	c.Password = clusterInfo.Password
	c.Version = clusterInfo.Version
	c.Endpoint = clusterInfo.Endpoint
	c.NodeCount = clusterInfo.NodeCount
	c.Metadata = clusterInfo.Metadata
	c.ServiceAccountToken = clusterInfo.ServiceAccountToken
	c.Status = clusterInfo.Status
}

// Remove removes a cluster
func (c *Cluster) Remove(ctx context.Context) error {
	return c.Driver.Remove(ctx, toInfo(c))
}

func toInfo(c *Cluster) *types.ClusterInfo {
	return &types.ClusterInfo{
		ClientCertificate:   c.ClientCertificate,
		ClientKey:           c.ClientKey,
		RootCaCertificate:   c.RootCACert,
		Username:            c.Username,
		Password:            c.Password,
		Version:             c.Version,
		Endpoint:            c.Endpoint,
		NodeCount:           c.NodeCount,
		Metadata:            c.Metadata,
		ServiceAccountToken: c.ServiceAccountToken,
		Status:              c.Status,
	}
}
