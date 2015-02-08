package main

import (
	"fmt"
	"github.com/toastdriven/takeanumber/server"
)

const Version = "1.0.0"

func main() {
	port := 13331

	fmt.Printf("takeanumber v%v\n", Version)

	// FIXME: This will need to accept command-line arguments.
	s := server.New(port)

	fmt.Printf("Listening on port %v\n", port)
	s.Run()
}
