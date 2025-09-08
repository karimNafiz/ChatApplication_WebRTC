package main

import (
	"flag"
	"fmt"
	"net/http"
)

const version = "1.0.0"

type configs struct {
	port int
	env  string

	/*
		required for the database connection pool
	*/
	db struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	cfg configs
	// need a logger
}

func main() {
	var cfg configs

	// take in the command line arguments
	flag.IntVar(&cfg.port, "port", 4000, "web application port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.Parse()

	app := &application{

		cfg: cfg,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	fmt.Printf("%s", err)

}
