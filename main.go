package main

import (
	"log"

	"github.com/lesnoi-kot/clip-radiot/server"
)

func main() {
	e := server.NewServer()

	if err := e.Start("0.0.0.0:8000"); err != nil {
		log.Printf("Server error: %s. Exiting", err)
	}
}
