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

// The chosen hostname of the Elasticsearch service.
const hostname = "recipees"

// Declaration of the Elasticsearch SDK API client.
var esClient *elasticsearch.Client

func main() {
	http.HandleFunc("/", ElasticSdk)
	const port = "8080"
	fmt.Printf("... binding on port %s, the application is being started.\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("... failed listening on port %s: %e, the application had been stopped.\n", port, err)
	}
}

func getConnectionString(hostname string) (string, bool) {
	// The requested environment variable name.
	const connectionString = "connectionString"
	return os.LookupEnv(hostname + "_" + connectionString)
}

func getEsClient(host string) *elasticsearch.Client {
	var cfg = elasticsearch.Config{
		Addresses: []string{host},
		// Sniffing should be disabled.
		DiscoverNodesOnStart: false,
	}
	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil
	}
	return esClient
}

func initialization(hostname string) {
	// For example, the result of the <host> would be: ["http://recipees:9200"]
	host, found := getConnectionString(hostname)
	if found {
		esClient = getEsClient(host)
	}
}

func Insert(esClient *elasticsearch.Client) (*esapi.Response, error) {
	return esClient.Index(
		"zerops-recipes",
		strings.NewReader(`{
			"service": "Golang",
			"version": "1.16.3",
			"message": "es-golang-basic"
		}`),
	)
}

func ElasticSdk(w http.ResponseWriter, r *http.Request) {
	type Result struct {
		Id string `json:"_id"`
	}
	var result Result
	if esClient == nil {
		initialization(hostname)
		if esClient == nil {
			fmt.Fprintf(w, "... Error! Elasticsearch SDK API client not initialized.")
			log.Fatal("... Error! Elasticsearch SDK API client not initialized.")
		}
	}
	if r.URL.Path == "/" {
		insertResult, err := Insert(esClient)
		if err != nil {
			fmt.Fprintf(w, "... Error! Elasticsearch insert operation failed.")
			log.Fatalf("... Error! Elasticsearch insert operation failed: %e", err)
		}
		defer insertResult.Body.Close()
		body, _ := io.ReadAll(insertResult.Body)
		json.Unmarshal(body, &result)
		if insertResult.StatusCode == 201 {
			fmt.Fprintf(w, "... Hello! A new document was inserted into Elasticsearch!\n")
			fmt.Printf("... created document id: %s\n", result.Id)
		} else {
			fmt.Fprintf(w, "... Error! Elasticsearch insert operation failed: %d\n", insertResult.StatusCode)
			fmt.Printf("... document creation failed: %d\n", insertResult.StatusCode)
		}
	}
}
