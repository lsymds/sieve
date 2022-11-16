package main

import (
	"github.com/lsymds/sieve"
	"github.com/lsymds/sieve/http"
)

func main() {
	server, err := http.NewHttpServer(sieve.NewOperationsStore())
	if err != nil {
		panic(err)
	}

	panic(server.ListenAndServe(":8080"))
}
