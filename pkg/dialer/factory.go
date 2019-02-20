package dialer

import (
	"fmt"
	"net"
	"time"

	"github.com/tan208123/navigate/pkg/config"
	"github.com/tan208123/navigate/pkg/config/dialer"
	"github.com/tan208123/navigate/pkg/encryptedstore"
	"github.com/tan208123/navigate/pkg/remotedialer"
	"github.com/tan208123/navigate/pkg/tunnelserver"
)

func NewFactory(management *config.ManagementContext) (dialer.Factory, error) {
	tunneler := tunnelserver.NewTunnelServer()
	secretStore := encryptedstore.NewGenericEncrypedStore("mc-", "", management.K8sClient.CoreV1().Namespaces(), management.K8sClient.CoreV1())

	return &Factory{
		TunnelServer: tunneler,
		store:        secretStore,
	}, nil
}

type Factory struct {
	TunnelServer *remotedialer.Server
	store        *encryptedstore.GenericEncryptedStore
}

func (f *Factory) NodeDialer(clusterName, machineName string) (dialer.Dialer, error) {
	return func(network, address string) (net.Conn, error) {
		d, err := f.nodeDialer(clusterName, machineName)
		if err != nil {
			return nil, err
		}
		return d(network, address)
	}, nil
}

func (f *Factory) nodeDialer(clusterName, machineName string) (dialer.Dialer, error) {

	if f.TunnelServer.HasSession(machineName) {
		d := f.TunnelServer.Dialer(machineName, 15*time.Second)
		return dialer.Dialer(d), nil
	}

	return nil, fmt.Errorf("can not build dialer to %s:%s", clusterName, machineName)
}

func (f *Factory) DockerDialer(clusterName, machineName string) (dialer.Dialer, error) {

	if f.TunnelServer.HasSession(machineName) {
		d := f.TunnelServer.Dialer(machineName, 15*time.Second)
		return func(string, string) (net.Conn, error) {
			return d("unix", "/var/run/docker.sock")
		}, nil
	}

	return nil, fmt.Errorf("can not build dialer to %s:%s", clusterName, machineName)
}
