package bot

import (
	"context"
	"fmt"
	"strings"

	"bikagame-go/internal/models"
	"bikagame-go/internal/store"
	tgbot "github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (a *App) handleStart(ctx context.Context, b *tgbot.Bot, update *tgmodels.Update) {
	msg := update.Message
	if msg == nil || msg.From == nil {
		return
	}

	_, _ = store.EnsureUser(ctx, a.DB, msg.From.ID, msg.From.Username, msg.From.FirstName, msg.From.LastName, strings.TrimSpace(msg.From.FirstName+" "+msg.From.LastName))

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:    msg.Chat.ID,
		Text:      "👋 Welcome to BIKA Game Bot (Go version, Mongo-compatible)",
		ParseMode: tgmodels.ParseModeHTML,
	})
}

func (a *App) handlePing(ctx context.Context, b *tgbot.Bot, update *tgmodels.Update) {
	msg := update.Message
	if msg == nil {
		return
	}

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:    msg.Chat.ID,
		Text:      "🏓 PONG",
		ParseMode: tgmodels.ParseModeHTML,
	})
}

func (a *App) handlePendingGroups(ctx context.Context, b *tgbot.Bot, update *tgmodels.Update) {
	msg := update.Message
	if msg == nil || msg.From == nil {
		return
	}
	if msg.From.ID != a.Cfg.OwnerID {
		return
	}

	cur, err := a.DB.Groups.Find(ctx, bson.M{"approvalStatus": "pending"})
	if err != nil {
		return
	}
	defer cur.Close(ctx)

	lines := []string{"🕒 <b>Pending Groups</b>", "━━━━━━━━━━━━"}
	i := 0
	for cur.Next(ctx) {
		i++
		var g models.Group
		if err := cur.Decode(&g); err != nil {
			continue
		}
		lines = append(lines, fmt.Sprintf("%d. %s", i, groupLabel(&g)))
	}

	if i == 0 {
		lines = append(lines, "No pending groups.")
	}

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:    msg.Chat.ID,
		Text:      strings.Join(lines, "\n"),
		ParseMode: tgmodels.ParseModeHTML,
	})
}

func (a *App) handleGroupStatus(ctx context.Context, b *tgbot.Bot, update *tgmodels.Update) {
	msg := update.Message
	if msg == nil {
		return
	}
	if msg.Chat.Type != tgmodels.ChatTypeGroup && msg.Chat.Type != tgmodels.ChatTypeSupergroup {
		return
	}

	g, err := store.EnsureGroup(ctx, a.DB, msg.Chat.ID, msg.Chat.Title, msg.Chat.Username, false)
	if err != nil {
		return
	}

	text := "👥 <b>Group Status</b>\n━━━━━━━━━━━━\n" +
		fmt.Sprintf("Group: %s\n", groupLabel(g)) +
		fmt.Sprintf("Approved: <b>%s</b>\n", esc(g.ApprovalStatus)) +
		fmt.Sprintf("Bot Admin: <b>%t</b>", g.BotIsAdmin)

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:    msg.Chat.ID,
		Text:      text,
		ParseMode: tgmodels.ParseModeHTML,
	})
}
