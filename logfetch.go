package main

import (
	"log"
	"net/http"

	"github.com/jdowling/logfetch/api"
)

func handleRequests() {
	s := api.NewServer()
	http.HandleFunc("/events", s.GetEvents)
	// TODO: dynanmic port
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
