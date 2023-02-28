package main

import (
	"flag"
	"github.com/distrD/boltX/db"
	"github.com/distrD/boltX/web"
	"log"
	"net/http"
)

var (
	dbLocation = flag.String("db-location", "", "the path to the bolt database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP post and host")
)

func parseFlag() {
	flag.Parse()
	if *dbLocation == "" {
		log.Fatalf("Must be Provide  db-location")
	}
}

func main() {
	// Open the XXXX.db data file in your current directory.
	// It will be created if it doesn't exist.
	parseFlag()
	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("NewDatabase(%q) : %v", *dbLocation, err)
	}

	defer close()

	srv := web.NewServer(db)
	http.HandleFunc("/set", srv.SetHandler)
	http.HandleFunc("/get", srv.GetHandler)

	srv.ListenAndServe(*httpAddr)
}
