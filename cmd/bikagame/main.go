package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"bikagame-go/internal/bot"
	"bikagame-go/internal/config"
	"bikagame-go/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbc, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app, err := bot.New(ctx, cfg, dbc)
	if err != nil {
		log.Fatal(err)
	}

	app.Start(ctx)
}
