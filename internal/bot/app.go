package bot

import (
	"context"
	"log"

	"bikagame-go/internal/config"
	"bikagame-go/internal/db"
	"bikagame-go/internal/store"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type App struct { Bot *tgbot.Bot; Cfg config.Config; DB *db.DB }

func New(ctx context.Context, cfg config.Config, dbc *db.DB) (*App, error) {
	app := &App{Cfg: cfg, DB: dbc}
	b, err := tgbot.New(cfg.BotToken, tgbot.WithDefaultHandler(func(ctx context.Context, b *tgbot.Bot, update *models.Update) {
		if update.MyChatMember != nil { app.handleMyChatMember(ctx, b, update) }
	}))
	if err != nil { return nil, err }
	app.Bot = b
	if _, err = store.EnsureTreasury(ctx, dbc, cfg); err != nil { return nil, err }
	app.registerHandlers()
	return app, nil
}

func (a *App) Start(ctx context.Context) { log.Println("bot polling started"); a.Bot.Start(ctx) }

func (a *App) registerHandlers() {
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/start", tgbot.MatchTypePrefix, a.handleStart)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/ping", tgbot.MatchTypePrefix, a.handlePing)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/status", tgbot.MatchTypePrefix, a.handleStatus)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/balance", tgbot.MatchTypePrefix, a.handleBalance)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/bal", tgbot.MatchTypePrefix, a.handleBalance)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, ".bal", tgbot.MatchTypePrefix, a.handleBalance)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/dailyclaim", tgbot.MatchTypePrefix, a.handleDailyClaim)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, ".dailyclaim", tgbot.MatchTypePrefix, a.handleDailyClaim)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/top10", tgbot.MatchTypePrefix, a.handleTop10)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, ".top10", tgbot.MatchTypePrefix, a.handleTop10)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/gift", tgbot.MatchTypePrefix, a.handleGift)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, ".gift", tgbot.MatchTypePrefix, a.handleGift)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/pendinggroups", tgbot.MatchTypePrefix, a.handlePendingGroups)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/groupstatus", tgbot.MatchTypePrefix, a.handleGroupStatus)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/approve", tgbot.MatchTypePrefix, a.handleApproveCmd)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/reject", tgbot.MatchTypePrefix, a.handleRejectCmd)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, ".slot", tgbot.MatchTypePrefix, a.handleSlot)
	a.Bot.RegisterHandler(tgbot.HandlerTypeCallbackQueryData, "groupapprove:", tgbot.MatchTypePrefix, a.handleApproveGroup)
	a.Bot.RegisterHandler(tgbot.HandlerTypeCallbackQueryData, "groupreject:", tgbot.MatchTypePrefix, a.handleRejectGroup)
}
