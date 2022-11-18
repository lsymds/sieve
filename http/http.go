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
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.RequestURI, "/http://") || strings.HasPrefix(r.RequestURI, "/https://") {
				h.handleProxy(rw, r)
			} else {
				respondNotFound(rw)
			}
		}),
	}

	return server.ListenAndServe()
}

// respondBadGateway writes a 502 BAD GATEWAY response to the response writer.
func respondBadGateway(w http.ResponseWriter) {
	w.WriteHeader(502)
}

// respondNotFound writes a 404 NOT FOUND response to the response writer.
func respondNotFound(w http.ResponseWriter) {
	w.WriteHeader(404)
}
