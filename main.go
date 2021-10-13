package main

import (
	"fuze/cmd"
	"log"
)

func main() {
	log.Println("Welcome to fuze, homie!")

	// there are 2 parts of Fuze: server-side and client-side.
	// server-side listen and retrieve the data, checking authorization, existed IP and port of incomming client.
	// the cmd triggers to start server.
	cmd.Execute()
}
