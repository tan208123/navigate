package app

import (
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	KubeConfig string
}

func Run(config *Config) error {
	fmt.Printf("app run ... \n")
	fmt.Printf("config is %v \n", config)
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		fmt.Printf("ListenAndServe: %v \n", err)
	}
	return nil
}

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}
