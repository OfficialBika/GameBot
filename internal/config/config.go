package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	BotToken   string
	MongoURI   string
	DBName     string
	OwnerID    int64
	Coin       string
	StartBonus int64
}

func mustEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		v = fallback
	}
	if v == "" {
		log.Fatalf("missing env: %s", key)
	}
	return v
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func Load() Config {
	ownerID, err := strconv.ParseInt(mustEnv("OWNER_ID", ""), 10, 64)
	if err != nil {
		log.Fatalf("invalid OWNER_ID: %v", err)
	}
	startBonus, err := strconv.ParseInt(getEnv("START_BONUS", "30000"), 10, 64)
	if err != nil {
		startBonus = 30000
	}
	return Config{
		BotToken:   mustEnv("BOT_TOKEN", ""),
		MongoURI:   mustEnv("MONGODB_URI", ""),
		DBName:     getEnv("DB_NAME", "bika_slot"),
		OwnerID:    ownerID,
		Coin:       getEnv("STORE_CURRENCY", "MMK"),
		StartBonus: startBonus,
	}
}
