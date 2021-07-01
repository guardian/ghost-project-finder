package main

import (
	"github.com/olivere/elastic/config"
	"github.com/olivere/elastic/v6"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	esUrl := os.Getenv("ELASTICSEARCH_URL")
	esInitialIndexName := os.Getenv("ELASTICSEARCH_INDEX")
	esTimeoutString := os.Getenv("ELASTICSEARCH_TIMEOUT")
	if esUrl == "" {
		log.Fatal("You must specify ELASTICSEARCH_URL in the environment")
	}

	esIndexName := "ghost-projects"
	if esInitialIndexName != "" {
		esIndexName = esInitialIndexName
	}

	esTimeout, _ := time.ParseDuration("60s")
	if esTimeoutString != "" {
		var parseErr error
		esTimeout, parseErr = time.ParseDuration(esTimeoutString)
		if parseErr != nil {
			log.Fatalf("Could not parse value %s for duration: %s", esTimeoutString, parseErr)
		}
	}

	esConfig := config.Config{
		URL: esUrl,
	}
	log.Print(esConfig)
	//sniffing is disabled for dev, because it breaks if you run ES in a docker container and the app outside one
	esClient, connectionErr := elastic.NewClient(elastic.SetURL(esUrl), elastic.SetSniff(false))

	if connectionErr != nil {
		log.Fatalf("Could not connect to Elasticsearch at %s: %s", esUrl, connectionErr)
	}

	log.Printf("Established connection to Elasticsearch at %s. Index name is %s.", esUrl, esIndexName)

	healthCheck := HealthcheckHandler{}
	inputHandler := InputHandler{
		elasticSearchClient: esClient,
		indexName:           esIndexName,
		timeout:             esTimeout,
	}
	http.Handle("/healthcheck", healthCheck)
	http.Handle("/foundfile", inputHandler)

	log.Printf("Starting server on port 9000")
	startServerErr := http.ListenAndServe(":9000", nil)

	if startServerErr != nil {
		log.Fatal(startServerErr)
	}
}
