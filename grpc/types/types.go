package types

import "context"

type Driver interface {
	Create(ctx context.Context, opts *DriverOptions) (*ClusterInfo, error)
	Remove(ctx context.Context, clusterInfo *ClusterInfo) error
}
