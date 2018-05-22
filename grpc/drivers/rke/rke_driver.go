package rke

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rancher/rke/cmd"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/k8s"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/tan208123/navigate/grpc/drivers"
	"github.com/tan208123/navigate/grpc/drivers/rke/rkecerts"
	"github.com/tan208123/navigate/grpc/types"
)

const (
	kubeConfigFile = "kube_config_cluster.yml"
	rancherPath    = "./management-state/rke/"
)

type WrapTransportFactory func(config *v3.RancherKubernetesEngineConfig) k8s.WrapTransport

type Driver struct {
	DockerDialer         hosts.DialerFactory
	LocalDialer          hosts.DialerFactory
	WrapTransportFactory WrapTransportFactory
}

func NewDriver() *Driver {
	d := &Driver{}
	return d
}

func (d *Driver) wrapTransport(config *v3.RancherKubernetesEngineConfig) k8s.WrapTransport {
	if d.WrapTransportFactory == nil {
		return nil
	}

	return k8s.WrapTransport(func(rt http.RoundTripper) http.RoundTripper {
		fn := d.WrapTransportFactory(config)
		if fn == nil {
			return rt
		}
		return fn(rt)
	})

}

// SetDriverOptions sets the drivers options to rke driver
func getYAML(driverOptions *types.DriverOptions) (string, error) {
	// first look up the file path then look up raw rkeConfig
	if path, ok := driverOptions.StringOptions["config-file-path"]; ok {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return driverOptions.StringOptions["rkeConfig"], nil
}

// Create creates the rke cluster
func (d *Driver) Create(ctx context.Context, opts *types.DriverOptions) (*types.ClusterInfo, error) {
	yaml, err := getYAML(opts)
	if err != nil {
		return nil, err
	}

	rkeConfig, err := drivers.ConvertToRkeConfig(yaml)
	if err != nil {
		return nil, err
	}

	stateDir, err := d.restore(nil)
	if err != nil {
		return nil, err
	}
	defer d.cleanup(stateDir)

	certsStr := ""
	APIURL, caCrt, clientCert, clientKey, certs, err := cmd.ClusterUp(ctx, &rkeConfig, d.DockerDialer, d.LocalDialer,
		d.wrapTransport(&rkeConfig), false, stateDir, false, false)
	if err == nil {
		certsStr, err = rkecerts.ToString(certs)
	}
	if err != nil {
		return d.save(&types.ClusterInfo{
			Metadata: map[string]string{
				"Config": yaml,
			},
		}, stateDir), err
	}

	return d.save(&types.ClusterInfo{
		Metadata: map[string]string{
			"Endpoint":   APIURL,
			"RootCA":     base64.StdEncoding.EncodeToString([]byte(caCrt)),
			"ClientCert": base64.StdEncoding.EncodeToString([]byte(clientCert)),
			"ClientKey":  base64.StdEncoding.EncodeToString([]byte(clientKey)),
			"Config":     yaml,
			"Certs":      certsStr,
		},
	}, stateDir), nil
}

// Remove removes the cluster
func (d *Driver) Remove(ctx context.Context, clusterInfo *types.ClusterInfo) error {
	rkeConfig, err := drivers.ConvertToRkeConfig(clusterInfo.Metadata["Config"])
	if err != nil {
		return err
	}
	stateDir, _ := d.restore(clusterInfo)
	defer d.save(nil, stateDir)
	return cmd.ClusterRemove(ctx, &rkeConfig, d.DockerDialer, d.wrapTransport(&rkeConfig), false, stateDir)
}

func (d *Driver) restore(info *types.ClusterInfo) (string, error) {
	os.MkdirAll(rancherPath, 0700)
	dir, err := ioutil.TempDir(rancherPath, "rke-")
	if err != nil {
		return "", err
	}

	if info != nil {
		state := info.Metadata["state"]
		if state != "" {
			ioutil.WriteFile(kubeConfig(dir), []byte(state), 0600)
		}
	}

	return filepath.Join(dir, "cluster.yml"), nil
}

func (d *Driver) save(info *types.ClusterInfo, stateDir string) *types.ClusterInfo {
	if info != nil {
		b, err := ioutil.ReadFile(kubeConfig(stateDir))
		if err == nil {
			if info.Metadata == nil {
				info.Metadata = map[string]string{}
			}
			info.Metadata["state"] = string(b)
		}
	}

	d.cleanup(stateDir)

	return info
}

func (d *Driver) cleanup(stateDir string) {
	if strings.HasSuffix(stateDir, "/cluster.yml") && !strings.Contains(stateDir, "..") {
		os.Remove(stateDir)
		os.Remove(kubeConfig(stateDir))
		os.Remove(filepath.Dir(stateDir))
	}
}

func kubeConfig(stateDir string) string {
	if strings.HasSuffix(stateDir, "/cluster.yml") {
		return filepath.Join(filepath.Dir(stateDir), kubeConfigFile)
	}
	return filepath.Join(stateDir, kubeConfigFile)
}
