package http

import (
	"net/http"
	"strings"

	"github.com/lsymds/sieve"
)

// HttpServer is a wrapper struct that contains all of the required dependencies to run a fully
// functioning Sieve server.
type HttpServer struct {
	store *sieve.OperationsStore
}

// NewHttpServer creates a new HttpServer instance or returns an error if there are issues creating
// it.
func NewHttpServer(store *sieve.OperationsStore) (*HttpServer, error) {
	server := &HttpServer{
		store,
	}

	return server, nil
}

// ListenAndServe listens on the provided port, serving any relevant endpoints for its lifetime.
func (h *HttpServer) ListenAndServe(addr string) error {
	http.DefaultServeMux.HandleFunc("/_/ws", h.handleWebsocketEndpoint)

	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// If internal API (ala /_/) let the multiplexer handle it - else, proxy it.
			if strings.HasPrefix(r.RequestURI, "/_/") {
				h, _ := http.DefaultServeMux.Handler(r)
				h.ServeHTTP(rw, r)
			} else if strings.HasPrefix(r.RequestURI, "/http://") || strings.HasPrefix(r.RequestURI, "/https://") {
				h.handleProxyEndpoint(rw, r)
			}
		}),
	}

	return server.ListenAndServe()
}

// respondBadGateway writes a 502 BAD GATEWAY response to the response writer.
func respondBadGateway(w http.ResponseWriter) {
	w.WriteHeader(502)
}
