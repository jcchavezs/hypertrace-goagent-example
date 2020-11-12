package main

import (
	"bytes"
	"context"
	"log"
	"net/http"

	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/net/hyperhttp"
)

func main() {
	cfg := config.LoadFromFile("./config.yml")
	cfg.ServiceName = config.String("client")

	shutdown := hypertrace.Init(cfg)
	defer shutdown()

	client := http.Client{
		Transport: hyperhttp.NewTransport(http.DefaultTransport),
	}

	req, err := http.NewRequest("GET", "http://localhost:8081/foo", bytes.NewBufferString(`{"name":"Dave"}`))
	req = req.WithContext(context.Background())
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatalf("failed to create the request: %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to perform the request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("failed to signup")
	}
}
