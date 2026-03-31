package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	appbot "bikagame-go/internal/bot"
	"bikagame-go/internal/config"
	"bikagame-go/internal/db"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbc, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app, err := appbot.New(ctx, cfg, dbc)
	if err != nil {
		log.Fatal(err)
	}

	app.Start(ctx)
}
