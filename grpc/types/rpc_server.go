package types

import (
	"net"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	listenAddr = "127.0.0.1"
)

type GrpcServer struct {
	driver  Driver
	address chan string
}

// NewServer creates a grpc server for a specific plugin
func NewServer(driver Driver, addr chan string) *GrpcServer {
	return &GrpcServer{
		driver:  driver,
		address: addr,
	}
}

// Create implements grpc method
func (s *GrpcServer) Create(ctx context.Context, opts *DriverOptions) (*ClusterInfo, error) {
	return s.driver.Create(ctx, opts)
}

// Remove implements grpc method
func (s *GrpcServer) Remove(ctx context.Context, clusterInfo *ClusterInfo) (*Empty, error) {
	return &Empty{}, s.driver.Remove(ctx, clusterInfo)
}

// Serve serves a grpc server
func (s *GrpcServer) Serve() {
	listen, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logrus.Fatal(err)
	}
	addr := listen.Addr().String()
	s.address <- addr
	grpcServer := grpc.NewServer()
	RegisterDriverServer(grpcServer, s)
	reflection.Register(grpcServer)
	logrus.Debugf("RPC GrpcServer listening on address %s", addr)
	if err := grpcServer.Serve(listen); err != nil {
		logrus.Fatal(err)
	}
	return
}
