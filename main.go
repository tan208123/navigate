package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/rancher/norman/pkg/dump"
	"github.com/rancher/norman/signal"
	"github.com/tan208123/navigate/app"
	"github.com/tan208123/navigate/k8s"
        "github.com/ehazlett/simplelog"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	VERSION = "dev"
)

func main() {

	var config app.Config

	app := cli.NewApp()
	app.Version = VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "kubeconfig",
			Usage:       "Kube config for accessing k8s cluster",
			EnvVar:      "KUBECONFIG",
			Destination: &config.KubeConfig,
		},
		cli.StringFlag{
			Name:        "k8s-mode",
			Usage:       "Mode to run or access k8s API server for management API (embedded, external, auto)",
			Value:       "auto",
			Destination: &config.K8sMode,
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Enable debug logs",
			Destination: &config.Debug,
		},
		cli.StringFlag{
			Name:        "log-format",
			Usage:       "Log formatter used (json, text, simple)",
			Value:       "simple",
		},
	}

	app.Action = func(c *cli.Context) error {

		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
        initLogs(c, config)
		return run(config)
	}

	app.Run(os.Args)
}

func initLogs(c *cli.Context, cfg app.Config) {
	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	switch c.String("log-format") {
	case "simple":
		logrus.SetFormatter(&simplelog.StandardFormatter{})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	logrus.SetOutput(os.Stdout)
}

func run(cfg app.Config) error {
	dump.GoroutineDumpOn(syscall.SIGUSR1, syscall.SIGILL)
	ctx := signal.SigTermCancelContext(context.Background())
	logrus.Infof("main.go ctx is %v ", ctx)

	embedded, ctx, kubeConfig, err := k8s.GetConfig(ctx, cfg.K8sMode, cfg.KubeConfig)
        logrus.Infof("main.go kubeConfig is %v ", kubeConfig)
	if err != nil {
                logrus.Errorf("main.go  ", err)
		//return err
	}
	cfg.Embedded = embedded

	return app.Run(&cfg)
}
