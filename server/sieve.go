package sieve

import (
	"sync"
	"time"
)

// OperationsStoreListener represents a function signature that any listeners must adhere to.
type OperationsStoreListener = func(o *Operation)

// OperationsStore defines all of the different ways to persist and retrieve operations.
type OperationsStore struct {
	mtx        sync.RWMutex
	listeners  map[uint]OperationsStoreListener
	operations map[string]*Operation
}

// NewOperationsStore creates a new instance of an operations store.
func NewOperationsStore() *OperationsStore {
	return &OperationsStore{
		listeners:  make(map[uint]func(o *Operation)),
		operations: make(map[string]*Operation),
	}
}

// AddListener adds a listener to the store, returning a function that can be used to remove that newly added listener.
func (s *OperationsStore) AddListener(listener OperationsStoreListener) func() {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	lid := uint(len(s.listeners) + 1)

	s.listeners[lid] = listener

	return func() {
		s.removeListener(lid)
	}
}

// Save persists an operation and updates all listeners.
func (s *OperationsStore) Save(o *Operation) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	for _, l := range s.listeners {
		go l(o)
	}

	s.operations[o.Id] = o
}

// GetOperationById retrieves an operation by its identifier.
func (s *OperationsStore) GetOperationById(id string) *Operation {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return s.operations[id]
}

// removeListener removes a listener from the store.
func (s *OperationsStore) removeListener(lid uint) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.listeners, lid)
}

// Operation represents a wrapper around a request and response coupling and contains additional fields that are used
// to build the application.
type Operation struct {
	Id        string    `json:"id"`
	Host      string    `json:"host"`
	Request   *Request  `json:"request"`
	Response  *Response `json:"response"`
	CreatedAt time.Time `json:"createdAt"`
}

// Request represents the original request that was made.
type Request struct {
	Host    string `json:"host"`
	Path    string `json:"path"`
	FullUrl string `json:"fullUrl"`
}

// Response represents the response retrieved from the request made to the origin server.
type Response struct {
	Latency   *time.Duration `json:"latency"`
	TotalTime *time.Duration `json:"totalTime"`
}
