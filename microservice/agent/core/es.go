package core

import (
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

const (
	agentESAddrEnv = "AGENT_ES_ADDR"
	agentESUserEnv = "AGENT_ES_USER"
	agentESPassEnv = "AGENT_ES_PASS"
)

var ES *elasticsearch.Client

func ESClientInit() {
	addr := os.Getenv(agentESAddrEnv)
	if addr == "" {
		addr = "http://localhost:9200"
	}

	cfg := elasticsearch.Config{Addresses: []string{addr}}
	if user := os.Getenv(agentESUserEnv); user != "" {
		cfg.Username = user
		cfg.Password = os.Getenv(agentESPassEnv)
	}

	var err error
	ES, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal("Failed to create Elasticsearch client: ", err)
	}
}
