package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/karimNafiz/ChatApplication_WebRTC/internal/jsonlog"
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
	cfg    configs
	logger *jsonlog.Logger
	cancel chan struct{}

	// need a logger
}

func main() {
	var cfg configs
	var logger *jsonlog.Logger = jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	// take in the command line arguments
	flag.IntVar(&cfg.port, "port", 4000, "web application port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.Parse()

	app := &application{

		cfg:    cfg,
		logger: logger,

		// make the cancel channel
		cancel: make(chan struct{}),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go monitorMainGoroutine(ctx, app.cancel)

	// before we listen for requests, we need to start the monitorMainGoroutine
	err := srv.ListenAndServe()
	fmt.Printf("%s", err)

}

func monitorMainGoroutine(mainContext context.Context, applicationCancelChannel chan<- struct{}) {
	/*
		until the main context is cancelled this sentence will block this go-routine
		but the moment this is unblocked
		we will send a signal to applicationCancelChannel
	*/
	<-mainContext.Done()
	/*
		need to signal application that the main channel is closed
	*/
	applicationCancelChannel <- struct{}{}
}

// func main() {
// 	migration_tools.ReorderMigrationFiles(2, "migration_test", ".sql", "./migrationTest")
// }
