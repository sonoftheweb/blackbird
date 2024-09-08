package main

import (
	"auth/internal/server"
	"fmt"
)

func main() {

	s := server.NewServer()

	err := s.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
