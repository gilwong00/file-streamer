package main

import (
	"context"
	"log"

	"github.com/gilwong00/file-streamer/internal/pkg/config"
	"github.com/gilwong00/file-streamer/internal/server"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Errpr unmarshalling config")
	}
	// TODO: will need to add logging and maybe a better implementation
	if err := server.StartServer(ctx, cfg); err != nil {
		panic(err)
	}
}
