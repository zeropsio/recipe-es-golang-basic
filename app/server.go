package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", ElasticSdk)
	const port = "8080"
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("... failed listening on port %s: %e, application not started.", port, err)
	}
	fmt.Printf("... listening on port %s, application started.\n", port)
}

func ElasticSdk(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "... Request: %s %s", r.Host, r.URL.Path)
	fmt.Println(w, "... Hello! A new document was inserted into Elasticsearch!")
}
