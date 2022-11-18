package http

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lsymds/sieve"
)

// handleProxy is the HTTP handler for proxying requests to their intended destinations. It captures interesting
// information, forwards or proxies the request, reads the response, then returns it to the user like this
// server did not exist 👻.
func (s *HttpServer) handleProxy(rw http.ResponseWriter, r *http.Request) {
	newRequest, err := buildIntendedRequest(r)
	if err != nil {
		respondBadGateway(rw)
		return
	}

	operation := buildOperation(newRequest)

	s.store.Save(operation)

	timeBeforeRequest := time.Now()

	originResponse, err := http.DefaultClient.Do(newRequest)
	if err != nil {
		respondBadGateway(rw)
		return
	}

	response, err := mapProxiedResponse(originResponse, rw, timeBeforeRequest)
	if err != nil {
		respondBadGateway(rw)
		return
	}

	operation.Response = response

	s.store.Save(operation)
}

// buildIntendedRequest takes a request made to an endpoint and generates a new *http.Request of what the caller
// originally intended to be sent to the origin.
func buildIntendedRequest(r *http.Request) (*http.Request, error) {
	replacementUrl, err := url.Parse(strings.TrimPrefix(r.URL.RequestURI(), "/"))
	if err != nil {
		return nil, err
	}

	originRequest := r.Clone(context.Background())
	originRequest.RequestURI = ""
	originRequest.URL = replacementUrl
	originRequest.Host = replacementUrl.Host

	return originRequest, nil
}

// buildOperation constructs a the wrapper representation of a given request and response.
func buildOperation(r *http.Request) sieve.Operation {
	operation := sieve.Operation{
		Id:   uuid.NewString(),
		Host: r.Host,
		Request: sieve.Request{
			Host:    r.Host,
			Path:    r.URL.Path,
			FullUrl: r.URL.RequestURI(),
		},
		CreatedAt: time.Now().UTC(),
	}

	return operation
}

// mapProxiedResponse takes the response from the origin server and maps it to the response this proxy API will send
// back to the requestor. It also returns a summary object to be used to inspect the response via the front-end
// application.
func mapProxiedResponse(
	originResponse *http.Response,
	rw http.ResponseWriter,
	timeBeforeRequest time.Time,
) (sieve.Response, error) {
	latency := time.Since(timeBeforeRequest)

	defer originResponse.Body.Close()

	// Headers.
	for h, hvs := range originResponse.Header {
		for _, hv := range hvs {
			rw.Header().Add(h, hv)
		}
	}

	// Cookies.
	for _, cv := range originResponse.Cookies() {
		http.SetCookie(rw, cv)
	}

	// Write status code.
	rw.WriteHeader(originResponse.StatusCode)

	// Body content.
	if _, err := io.Copy(rw, originResponse.Body); err != nil {
		return sieve.Response{}, err
	}

	totalTime := time.Since(timeBeforeRequest)

	return sieve.Response{Latency: latency, TotalTime: totalTime}, nil
}
