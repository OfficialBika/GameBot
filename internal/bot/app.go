package bot

import (
	"context"
	"log"

	"bikagame-go/internal/config"
	"bikagame-go/internal/db"
	"bikagame-go/internal/store"
	tgbot "github.com/go-telegram/bot"
)

type App struct {
	Bot *tgbot.Bot
	Cfg config.Config
	DB  *db.DB
}

func New(ctx context.Context, cfg config.Config, dbc *db.DB) (*App, error) {
	b, err := tgbot.New(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	app := &App{
		Bot: b,
		Cfg: cfg,
		DB:  dbc,
	}

	app.registerHandlers()

	_, err = store.EnsureTreasury(ctx, dbc, cfg)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Start(ctx context.Context) {
	log.Println("bot polling started")
	a.Bot.Start(ctx)
}

func (a *App) registerHandlers() {
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/start", tgbot.MatchTypePrefix, a.handleStart)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/ping", tgbot.MatchTypePrefix, a.handlePing)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/pendinggroups", tgbot.MatchTypePrefix, a.handlePendingGroups)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/groupstatus", tgbot.MatchTypePrefix, a.handleGroupStatus)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMessageText, ".slot", tgbot.MatchTypePrefix, a.handleSlot)
	a.Bot.RegisterHandler(tgbot.HandlerTypeCallbackQueryData, "groupapprove:", tgbot.MatchTypePrefix, a.handleApproveGroup)
	a.Bot.RegisterHandler(tgbot.HandlerTypeCallbackQueryData, "groupreject:", tgbot.MatchTypePrefix, a.handleRejectGroup)
	a.Bot.RegisterHandler(tgbot.HandlerTypeMyChatMember, "", tgbot.MatchTypeContains, a.handleMyChatMember)
}
