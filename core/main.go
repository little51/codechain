package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	abciserver "github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectDb() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	db, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open mongodb: %v", err)
		os.Exit(1)
	}
	err = db.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to ping mongodb: %v", err)
		os.Exit(1)
	}
	return db
}

func main() {
	// define abciapp & mongodb
	var app types.Application
	var db *mongo.Client
	// connect to mongodb,disconnect when app shutdown
	db = connectDb()
	defer db.Disconnect(context.TODO())
	app = NewAssetsApplication(db)
	flag.Parse()
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	server := abciserver.NewSocketServer("tcp://0.0.0.0:26658", app)
	// start abciserver stop when app shutdown
	server.SetLogger(logger)
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting socket server: %v", err)
		os.Exit(1)
	}
	defer server.Stop()
	// loop
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)
}
