package types

import "context"

type Driver interface {
	// Create creates the cluster
	Create(ctx context.Context, opts *DriverOptions) (*ClusterInfo, error)

	// Update updates the cluster
	Update(ctx context.Context, clusterInfo *ClusterInfo, opts *DriverOptions) (*ClusterInfo, error)

	// PostCheck does post action after provisioning
	PostCheck(ctx context.Context, clusterInfo *ClusterInfo) (*ClusterInfo, error)

	// Remove removes the cluster
	Remove(ctx context.Context, clusterInfo *ClusterInfo) error
}
