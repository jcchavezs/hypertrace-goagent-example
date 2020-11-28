package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/github.com/gorilla/hypermux"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/net/hyperhttp"
)

func main() {
	cfg := config.LoadFromFile("./config.yml")
	cfg.ServiceName = config.String("frontend")

	shutdown := hypertrace.Init(cfg)
	defer shutdown()

	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:9000"
	}

	r := mux.NewRouter()
	r.Use(hypermux.NewMiddleware())
	r.Handle("/", makeHandler(backendURL))
	log.Fatal(http.ListenAndServe(":8081", r))
}

func makeHandler(backendURL string) http.HandlerFunc {
	client := http.Client{
		Transport: hyperhttp.NewTransport(http.DefaultTransport),
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest(
			"GET",
			backendURL,
			bytes.NewBufferString(`{"name":"Dave"}`),
		)
		req = req.WithContext(r.Context())
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

		fmt.Printf("%s %s - %s\n", r.Method, r.URL.String(), res.Status)
	})
}
