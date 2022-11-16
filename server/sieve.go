package sieve

import (
	"sync"
	"time"
)

// OperationsStoreListener represents a function signature that any listeners must adhere to.
type OperationsStoreListener = func(o *Operation)

// OperationsStore defines all of the different ways to persist and retrieve operations.
type OperationsStore struct {
	mutex     sync.RWMutex
	listeners map[uint]OperationsStoreListener
}

// NewOperationsStore creates a new instance of an operations store.
func NewOperationsStore() *OperationsStore {
	return &OperationsStore{
		listeners: make(map[uint]func(o *Operation)),
	}
}

// AddListener adds a listener to the store, returning a function that can be used to remove that newly added listener.
func (s *OperationsStore) AddListener(listener OperationsStoreListener) func() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	lid := uint(len(s.listeners) + 1)

	s.listeners[lid] = listener

	return func() {
		s.removeListener(lid)
	}
}

// Save persists an operation and updates all listeners.
func (s *OperationsStore) Save(o *Operation) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, l := range s.listeners {
		l(o)
	}
}

// removeListener removes a listener from the store.
func (s *OperationsStore) removeListener(lid uint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.listeners, lid)
}

// Operation represents a wrapper around a request and response coupling and contains additional fields that are used
// to build the application.
type Operation struct {
	Id        string
	Host      string
	Request   *Request
	Response  *Response
	CreatedAt time.Time
}

// Request represents the original request that was made.
type Request struct {
	Host    string
	Path    string
	FullUrl string
}

// Response represents the response retrieved from the request made to the origin server.
type Response struct {
	Latency   *time.Duration
	TotalTime *time.Duration
}
