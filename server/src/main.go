package main

import (
	"log"
	"net/http"

	"./smile"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// Websocket server
	server := smile.NewServer()
	go server.Listen()

	// Static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
