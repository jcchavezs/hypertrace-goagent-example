package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/database/hypersql"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/net/hyperhttp"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = config.String("backend")

	shutdown := hypertrace.Init(cfg)
	defer shutdown()

	db, err := initDB()
	if err != nil {
		log.Fatalf("failed to initialize database connection: %v", err)
	}

	r := mux.NewRouter()
	r.Handle("/", hyperhttp.NewHandler(makeFooHandler(db), "/"))
	log.Fatal(http.ListenAndServe(":9000", r))
}

type person struct {
	Name string `json:"name"`
}

func makeFooHandler(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("failed to read body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		p := &person{}
		err = json.Unmarshal(sBody, p)
		if err != nil {
			log.Printf("failed to unmarshal body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = db.ExecContext(r.Context(), "INSERT INTO `users` (`name`) VALUES (?)", p.Name)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("{\"error\": \"Failed to insert %s\"}", p.Name)))
			fmt.Printf("%s %s - %d\n", r.Method, r.URL.String(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("{\"message\": \"Hello %s\"}", p.Name)))
		fmt.Printf("%s %s - %d\n", r.Method, r.URL.String(), http.StatusOK)
	})
}

const dbPingRetries = 10

func initDB() (*sql.DB, error) {
	var (
		driver driver.Driver
		db     *sql.DB
	)

	dbHost := os.Getenv("MYSQL_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	// Explicitly wrap the MySQLDriver driver with hypersql.
	driver = hypersql.Wrap(&mysql.MySQLDriver{})

	// Register our hypersql wrapper as a database driver.
	sql.Register("ht-mysql", driver)

	// Connect to a mysql database using the hypersql driver wrapper.
	// ?interpolateParams=true will escape the variables for any requests
	// and send ready-for-use queries to the server for github.com/go-sql-driver/mysql.
	// This save us a meaningless span.
	db, err := sql.Open("ht-mysql", "root:root@tcp("+dbHost+")/app?interpolateParams=true")
	if err != nil {
		return nil, fmt.Errorf("failed to connect the DB: %v", err)
	}

	for i := 0; i <= dbPingRetries; i++ {
		if err = db.Ping(); err == nil {
			break
		}

		if i == dbPingRetries {
			return nil, fmt.Errorf("failed to ping the DB after %d retries: %v", dbPingRetries, err)
		}

		fmt.Printf("failed to ping DB: %v. Retrying\n", err)
		time.Sleep(2 * time.Second)
	}

	return db, nil
}
