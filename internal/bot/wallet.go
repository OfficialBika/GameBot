package bot

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"

	"bikagame-go/internal/store"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	dailyMin int64 = 500
	dailyMax int64 = 2000
)

func (a *App) handleBalance(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	msg := update.Message
	if msg == nil || msg.From == nil {
		return
	}

	user, err := store.EnsureUser(
		ctx,
		a.DB,
		msg.From.ID,
		msg.From.Username,
		msg.From.FirstName,
		msg.From.LastName,
		strings.TrimSpace(msg.From.FirstName+" "+msg.From.LastName),
	)
	if err != nil {
		return
	}

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: msg.Chat.ID,
		Text: fmt.Sprintf(
			"💼 <b>Balance</b>\n━━━━━━━━━━━━\nBalance: <b>%s</b> %s\nRank: <b>%s</b>",
			fmtInt(user.Balance),
			esc(a.Cfg.Coin),
			esc(walletRank(user.Balance)),
		),
		ParseMode: models.ParseModeHTML,
	})
}

func (a *App) handleDotBal(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	msg := update.Message
	if msg == nil || msg.From == nil {
		return
	}

	user, err := store.EnsureUser(
		ctx,
		a.DB,
		msg.From.ID,
		msg.From.Username,
		msg.From.FirstName,
		msg.From.LastName,
		strings.TrimSpace(msg.From.FirstName+" "+msg.From.LastName),
	)
	if err != nil {
		return
	}

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: msg.Chat.ID,
		Text: fmt.Sprintf(
			"💼 <b>My Balance</b>\n━━━━━━━━━━━━\nBalance: <b>%s</b> %s\nRank: <b>%s</b>",
			fmtInt(user.Balance),
			esc(a.Cfg.Coin),
			esc(walletRank(user.Balance)),
		),
		ParseMode: models.ParseModeHTML,
	})
}

func (a *App) handleDailyClaim(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	msg := update.Message
	if msg == nil || msg.From == nil {
		return
	}

	user, err := store.EnsureUser(
		ctx,
		a.DB,
		msg.From.ID,
		msg.From.Username,
		msg.From.FirstName,
		msg.From.LastName,
		strings.TrimSpace(msg.From.FirstName+" "+msg.From.LastName),
	)
	if err != nil {
		return
	}

	dateKey := yangonDateKey()
	if user.LastDailyClaimDate != nil && *user.LastDailyClaimDate == dateKey {
		_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID:    msg.Chat.ID,
			Text:      "⏳ ဒီနေ့ Daily Claim ယူပြီးပါပြီ။",
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	amount := dailyMin + rand.Int64N(dailyMax-dailyMin+1)

	err = store.TreasuryPayToUser(ctx, a.DB, msg.From.ID, amount)
	if err != nil {
		_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID:    msg.Chat.ID,
			Text:      "🏦 Treasury မလုံလောက်ပါ။",
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	_ = store.SetLastDailyClaimDate(ctx, a.DB, msg.From.ID, dateKey)

	updatedUser, _ := store.GetUser(ctx, a.DB, msg.From.ID)

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: msg.Chat.ID,
		Text: fmt.Sprintf(
			"🎁 <b>Daily Claim</b>\n━━━━━━━━━━━━\nReceived: <b>%s</b> %s\nBalance: <b>%s</b> %s",
			fmtInt(amount),
			esc(a.Cfg.Coin),
			fmtInt(updatedUser.Balance),
			esc(a.Cfg.Coin),
		),
		ParseMode: models.ParseModeHTML,
	})
}

func (a *App) handleStatus(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	msg := update.Message
	if msg == nil {
		return
	}

	usersCount, _ := a.DB.Users.CountDocuments(ctx, map[string]any{})
	groupsCount, _ := a.DB.Groups.CountDocuments(ctx, map[string]any{})

	t, err := store.GetTreasury(ctx, a.DB)
	if err != nil {
		return
	}

	text := fmt.Sprintf(
		"📊 <b>BIKA Bot Status</b>\n━━━━━━━━━━━━\n👥 Users: <b>%s</b>\n👨‍👩‍👧‍👦 Groups: <b>%s</b>\n🏦 Treasury: <b>%s</b> %s\n📦 Total Supply: <b>%s</b> %s\n🛠 Maintenance: <b>%t</b>\n🎯 VIP WR: <b>%d%%</b>\n🎰 Slot RTP: <b>%.2f%%</b>",
		fmtInt(usersCount),
		fmtInt(groupsCount),
		fmtInt(t.OwnerBalance),
		esc(a.Cfg.Coin),
		fmtInt(t.TotalSupply),
		esc(a.Cfg.Coin),
		t.MaintenanceMode,
		t.VIPWinRate,
		t.SlotRTP*100,
	)

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:    msg.Chat.ID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	})
}
