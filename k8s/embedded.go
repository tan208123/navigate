package k8s

import (
	"context"
	"fmt"

	"github.com/docker/docker/pkg/reexec"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

func init() {
	logrus.Infof("reexec init... ")
	if reexec.Init() {
		logrus.Infof("reexec error ")
	}
}

func getEmbedded(ctx context.Context) (bool, context.Context, *rest.Config, error) {
	//return err
	return false, ctx, nil, fmt.Errorf("embedded support is not compiled in, rebuild with -tags k8s")
}
