package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type ElasticDoc struct {
	Service string
	Version string
	Message string
}

// The chosen hostname of the Elasticsearch service.
const hostname = "recipees"

// The requested environment variable name.
const connectionString = "connectionString"

// For example, the result of the <host> would be: ["http://recipees:9200"]
var host, _ = os.LookupEnv(hostname + "_" + connectionString)

var cfg = elasticsearch.Config{
	Addresses: []string{host},
	// Sniffing should be disabled.
	DiscoverNodesOnStart: false,
}
var esClient, _ = elasticsearch.NewClient(cfg)

func main() {
	http.HandleFunc("/", ElasticSdk)
	const port = "8080"
	fmt.Printf("... binding on port %s, the application is being started.\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("... failed listening on port %s: %e, the application had been stopped.\n", port, err)
	}
}

func Insert(esClient *elasticsearch.Client) (*esapi.Response, error) {
	doc := ElasticDoc{}
	doc.Service = "Golang"
	doc.Version = "1.16.3"
	doc.Message = "es-golang-basic"
	// jsonDoc, _ := json.Marshal(doc)
	return esClient.Index(
		"zerops-recipes",
		// strings.NewReader(string(jsonDoc)),
		strings.NewReader(`{
			"service": "Golang",
			"version": "1.16.3",
			"message": "es-golang-basic"
		}`),
	)
}

func ElasticSdk(w http.ResponseWriter, r *http.Request) {
	type Result struct {
		_id string `json:"Id"`
	}
	var result Result
	if r.URL.Path == "/" {
		insertResult, err := Insert(esClient)
		if err != nil {
			log.Fatalf("... Error! Elasticsearch insert operation failed: %e", err)
		}
		defer insertResult.Body.Close()
		body, _ := io.ReadAll(insertResult.Body)
		json.Unmarshal(body, &result)
		if insertResult.StatusCode == 201 {
			fmt.Fprintf(w, "... Hello! A new document was inserted into Elasticsearch!\n")
			fmt.Printf("... created document id: %s\n", result)
		} else {
			fmt.Fprintf(w, "... Error! Elasticsearch insert operation failed: %d\n", insertResult.StatusCode)
			fmt.Printf("... document creation failed: %d\n", insertResult.StatusCode)
		}
	}
}
