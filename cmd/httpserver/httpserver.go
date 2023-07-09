package httpserver

import (
	"kayn-form/cmd/adapters"
	"log"
	"net/http"
)

const (
	ServerPort       = ":7000"
	UserDataEndpoint = "/send-user-data"
	ResultEndpoint   = "/result"
)

// KaynFormServer is an interface to an HTTP server which handles requests
type KaynFormServer struct {
	Mux                *http.ServeMux
	Logger             *log.Logger
	HelpPixHTTPAdapter *adapters.KaynFormAdapter
}

// SetupRoutes configures the routes of the API
func (srv *KaynFormServer) SetupRoutes() {
	srv.Mux.Handle(UserDataEndpoint, http.HandlerFunc(srv.HelpPixHTTPAdapter.GetSummonerInfo))
}

// Start sets up the HTTP webserver to listen and handle traffic. It
// takes the port number to listen on as a parameter in the form ":PORT_NUMBER"
func (srv *KaynFormServer) Start(port string) error {
	return http.ListenAndServe(port, srv.Mux)
}

// NewParrotServer returns an instance of a configured ParrotServer
func NewParrotServer(logger *log.Logger, adapter *adapters.KaynFormAdapter) *KaynFormServer {
	httpServer := &KaynFormServer{
		Mux:                http.NewServeMux(),
		Logger:             logger,
		HelpPixHTTPAdapter: adapter,
	}
	return httpServer
}
