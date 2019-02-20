package dialer

import "net"

type Dialer func(network, address string) (net.Conn, error)

type Factory interface {
	DockerDialer(clusterName, machineName string) (Dialer, error)
	NodeDialer(clusterName, machineName string) (Dialer, error)
}
