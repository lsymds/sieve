package http

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/lsymds/sieve"
)

// upgrader is a configured websocket upgrader instance.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// websocketMessage represents a websocket message sent over the wire to any websocket connections.
type websocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// handleWebsocket handles the websocket endpoint, passing events over the wire when they are received.
func (s *HttpServer) handleWebsocket(rw http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		respondBadGateway(rw)
		return
	}
	defer ws.Close()

	// The websocket client can only write one message at a time. Thus, we need a mtx to prevent the concurrent
	// writes from happening.
	mtx := sync.Mutex{}

	// Subscribe to the operations store, publishing a message whenever an operation is saved.
	removeOperationsStoreSubscription := s.store.AddListener(func(o *sieve.Operation) {
		mtx.Lock()
		defer mtx.Unlock()

		ws.WriteJSON(websocketMessage{
			Type: "operation",
			Data: struct {
				OperationId string `json:"operationId"`
			}{
				OperationId: o.Id,
			},
		})
	})
	defer removeOperationsStoreSubscription()

	// Continuously loop for the lifetime of the websocket connection, sending a 'ping' message every 10 seconds.
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			func() {
				mtx.Lock()
				defer mtx.Unlock()

				if err := ws.WriteJSON(websocketMessage{Type: "ping"}); err != nil {
					return
				}
			}()
		}
	}
}

// handleGetOperation returns an operation by its requested identifier.
func (s *HttpServer) handleGetOperation(rw http.ResponseWriter, r *http.Request) {
	operationId := mux.Vars(r)["operationId"]

	operation := s.store.GetOperationById(operationId)
	if operation == nil {
		respondNotFound(rw)
		return
	}

	rw.Header().Add("Content-Type", "application/json")

	json.NewEncoder(rw).Encode(operation)
}
