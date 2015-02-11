package server_test

import "github.com/toastdriven/takeanumber/server"

func ExampleServer() {
    port := 13331

    // Create a Server.
    s := server.New(port)

    // Run the server.
    s.Run()
}
