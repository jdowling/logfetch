package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	reversescan "github.com/jdowling/logfetch/pkg"
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

	path := filepath.Join(s.logFilePrefix, r.URL.Query().Get("file"))
	file, err := os.Open(path)
	if err != nil {
		// TODO: differentiate based on err probably
		// could be a 500 instead in some cases.
		log.Println("Error opening file:", path, " err:", err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error opening file:%v", path)
		return
	}
	defer file.Close()

	regex_input := r.URL.Query().Get("filter")
	matcher, err := regexp.Compile(regex_input)
	if err != nil {
		log.Println("Error invalid regex:", regex_input, " err:", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid regex:%v", regex_input)
		return
	}

	match_limit := r.URL.Query().Get("n")
	matches_size, err := strconv.ParseInt(match_limit, 10, 64)
	if err != nil {
		log.Println("Error converting:", match_limit, " to int.", " err:", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error converting:%v to int.", match_limit)
		return
	}
	matches := make([]string, 0, matches_size)

	scanner := reversescan.New(file)
	for scanner.Scan() {
		line := scanner.Text()
		if matcher.MatchString(line) {
			matches = append(matches, line)
			if int64(len(matches)) == matches_size {
				break
			}
		}
	}

	result := EventsResponse{matches}
	json, err := json.Marshal(result)
	if err != nil {
		// TODO: unit test this path.
		log.Println("Error marshalling result:", result, " to JSON: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshalling result:%v to JSON", result)
		return
	}
	fmt.Fprint(w, string(json))
}
