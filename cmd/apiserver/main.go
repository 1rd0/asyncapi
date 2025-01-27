package main

import (
	"asyncapi/apiserver"
	"asyncapi/config"
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
	Jsonhandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(Jsonhandler)
	defer cancel()
	conf, err := config.New()
	if err != nil {
		return err
	}
	fmt.Println("server started")
	server := apiserver.NewApiServer(conf, logger)
	if err := server.Start(ctx); err != nil {
		return err
	}

	return nil
}
