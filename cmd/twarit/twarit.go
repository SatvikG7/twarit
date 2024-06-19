package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/satvikg7/twarit/internal/server"
)

func main() {
	// get port
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	// start server
	server.SetupServer()
	go server.Listen(*port)
	println()

	// wait for user to quit
	var input string
	for {
		fmt.Scanln(&input)
		if input == "q" {
			os.Exit(0)
		}
		println("Invalid input")
	}
}
