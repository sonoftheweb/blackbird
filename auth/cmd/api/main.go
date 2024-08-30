package main

import (
	"auth/internal/server"
	"fmt"
)

func main() {

	newServer := server.NewServer()

	err := newServer.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start newServer: %s", err))
	}
}
