package server

import (
	"fmt"
	"html"
	"net/http"
)

var s *PintaDBServer

// NewHTTPServer creates a new HTTP server for PintaDB
func NewHTTPServer(pinta *PintaDBServer, port uint64) {
	r := http.NewServeMux()
	s = pinta

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "PintaDB v%s is sailing", Version)
	})

	r.HandleFunc("/search", searchHandler)
	r.HandleFunc("/index", indexHandlex)

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Get the query string from the URL
	q := r.URL.Query().Get("query")
	if q == "" {
		q = r.URL.Query().Get("q")
	}

	// Run a search using the query string
	results, err := s.Search(q, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the results to the response
	fmt.Fprintf(w, "%s\n", s.Documents[results[0]].RawText)
}

func indexHandlex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
