package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	logFilePrefix string
}

type EventsResponse struct {
	Events []string
}

func NewServer() *Server {
	return &Server{"/var/log"}
}

func (s *Server) GetEvents(w http.ResponseWriter, r *http.Request) {
	log.Println("/events GET params:", r.URL.Query())

	result := EventsResponse{[]string{fmt.Sprintf("prefix:%v", s.logFilePrefix)}}
	json, err := json.Marshal(result)
	if err != nil {
		// TODO: unit test this path.
		log.Println("Error marshalling result:", result, " to JSON: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshalling result:%v to JSON", result)
	}
	fmt.Fprint(w, string(json))
}
