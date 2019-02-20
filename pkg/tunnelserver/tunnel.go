package tunnelserver

import (
	"net/http"

	"github.com/tan208123/navigate/pkg/remotedialer"
)

func NewTunnelServer() *remotedialer.Server {
	return remotedialer.New(func(rw http.ResponseWriter, req *http.Request, code int, err error) {
		rw.WriteHeader(code)
		rw.Write([]byte(err.Error()))
	})
}
