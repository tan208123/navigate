package k8s

import (
	"context"

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
	return true, ctx, nil, err
}
