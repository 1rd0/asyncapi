package main

import (
	"asyncapi/apiserver"
	"asyncapi/config"
	"asyncapi/store"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	conf, err := config.New()
	if err != nil {
		return err
	}
	Jsonhandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(Jsonhandler)
	JwtManager := apiserver.NewJwtManager(conf)
	db, err := store.NewPostgresDb(conf)
	DataStore := store.New(db)
	fmt.Println("server started")
	server := apiserver.NewApiServer(conf, logger, DataStore, JwtManager)
	if err := server.Start(ctx); err != nil {
		return err
	}

	return nil
}
