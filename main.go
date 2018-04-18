package main

import (	
	"context"
	"fmt"
	"os"
	"log"
	"net/http"
    "syscall"

    "github.com/rancher/norman/pkg/dump"
    "github.com/rancher/norman/signal"
	"github.com/tan208123/navigate/app"
	"github.com/urfave/cli"
)

var (
	VERSION = "dev"
)

func main() {

	var config  app.Config

	app := cli.NewApp()
	app.Version = VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "kubeconfig",
			Usage:       "Kube config for accessing k8s cluster",
			EnvVar:      "KUBECONFIG",
			Destination: &config.KubeConfig,
		},
	}

	app.Action = func(c *cli.Context) error {
		
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
		
		return run(config)
	}

	app.Run(os.Args)
}

func run(config app.Config) error {
	dump.GoroutineDumpOn(syscall.SIGUSR1, syscall.SIGILL)
	ctx := signal.SigTermCancelContext(context.Background())
	fmt.Printf("main.go ctx is %v \n", ctx)
	return app.Run(&config) 
}