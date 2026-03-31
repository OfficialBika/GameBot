package bot

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	slotgame "bikagame-go/internal/games/slot"
	"bikagame-go/internal/store"
	tgbot "github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"
)

var (
	slotEngine    = slotgame.New()
	slotMu        sync.Mutex
	activeSlots   = map[int64]bool{}
	lastSlotAt    = map[int64]int64{}
	maxActiveSlot = 5
)

func canSpin(userID int64) (bool, string) {
	slotMu.Lock()
	defer slotMu.Unlock()

	if len(activeSlots) >= maxActiveSlot && !activeSlots[userID] {
		return false, "⛔ <b>Slot Busy</b>\n━━━━━━━━━━━━\nအခုတလော တစ်ပြိုင်နက် ဆော့နေသူများလို့ ခဏနားပြီး ပြန်ကြိုးစားပါ။"
	}

	now := time.Now().UnixMilli()
	last := lastSlotAt[userID]
	if now-last < slotEngine.CooldownMS {
		waitSec := int(((slotEngine.CooldownMS - (now - last)) + 999) / 1000)
		return false, fmt.Sprintf("⏳ ခဏစောင့်ပါ… (%ds) နောက်တစ်ခါ spin လုပ်နိုင်ပါမယ်။", waitSec)
	}

	activeSlots[userID] = true
	lastSlotAt[userID] = now
	return true, ""
}

func releaseSpin(userID int64) {
	slotMu.Lock()
	defer slotMu.Unlock()
	delete(activeSlots, userID)
}

func parseBet(text string) (int64, bool) {
	parts := strings.Fields(text)
	if len(parts) != 2 {
		return 0, false
	}
	n, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || n <= 0 {
		return 0, false
	}
	return n, true
}

func slotFrame(title, note, a, b, c string) string {
	return fmt.Sprintf(
		"🎰 <b>%s</b>\n━━━━━━━━━━━━\n<pre>%s</pre>\n━━━━━━━━━━━━\n%s",
		esc(title),
		esc(slotgame.Art(a, b, c)),
		esc(note),
	)
}

func (a *App) handleSlot(ctx context.Context, b *tgbot.Bot, update *tgmodels.Update) {
	msg := update.Message
	if msg == nil || msg.From == nil || msg.Text == "" {
		return
	}

	if msg.Chat.Type != tgmodels.ChatTypeGroup && msg.Chat.Type != tgmodels.ChatTypeSupergroup {
		return
	}

	bet, ok := parseBet(msg.Text)
	if !ok {
		return
	}

	if bet < slotEngine.MinBet || bet > slotEngine.MaxBet {
		_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: msg.Chat.ID,
			Text: fmt.Sprintf(
				"🎰 <b>BIKA Pro Slot</b>\n━━━━━━━━━━━━\nUsage: <code>.slot 1000</code>\nMin: <b>%s</b>\nMax: <b>%s</b>",
				fmtInt(slotEngine.MinBet),
				fmtInt(slotEngine.MaxBet),
			),
			ParseMode: tgmodels.ParseModeHTML,
		})
		return
	}

	okSpin, reason := canSpin(msg.From.ID)
	if !okSpin {
		_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID:    msg.Chat.ID,
			Text:      reason,
			ParseMode: tgmodels.ParseModeHTML,
		})
		return
	}
	defer releaseSpin(msg.From.ID)

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

	if user.Balance < bet {
		_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: msg.Chat.ID,
			Text: fmt.Sprintf(
				"❌ <b>Balance မလုံလောက်ပါ</b>\n━━━━━━━━━━━━\nBet: <b>%s</b>\nYour Balance: <b>%s</b>",
				fmtInt(bet),
				fmtInt(user.Balance),
			),
			ParseMode: tgmodels.ParseModeHTML,
		})
		return
	}

	treasury, err := store.GetTreasury(ctx, a.DB)
	if err != nil {
		return
	}

	initA := slotEngine.RandomSymbol(slotEngine.Reels[0])
	initB := slotEngine.RandomSymbol(slotEngine.Reels[1])
	initC := slotEngine.RandomSymbol(slotEngine.Reels[2])

	sent, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:    msg.Chat.ID,
		Text:      slotFrame("BIKA Pro Slot", "reels spinning…", initA, initB, initC),
		ParseMode: tgmodels.ParseModeHTML,
	})
	if err != nil {
		return
	}

	err = store.UserPayToTreasury(ctx, a.DB, msg.From.ID, bet)
	if err != nil {
		if errors.Is(err, store.ErrUserInsufficient) {
			_, _ = b.EditMessageText(ctx, &tgbot.EditMessageTextParams{
				ChatID:    msg.Chat.ID,
				MessageID: sent.ID,
				Text:      "❌ <b>Balance မလုံလောက်ပါ</b>\n━━━━━━━━━━━━\nSpin မလုပ်နိုင်ပါ။",
				ParseMode: tgmodels.ParseModeHTML,
			})
			return
		}
		return
	}

	time.Sleep(220 * time.Millisecond)

	lockA := slotEngine.RandomSymbol(slotEngine.Reels[0])
	midB := slotEngine.RandomSymbol(slotEngine.Reels[1])
	midC := slotEngine.RandomSymbol(slotEngine.Reels[2])

	_, _ = b.EditMessageText(ctx, &tgbot.EditMessageTextParams{
		ChatID:    msg.Chat.ID,
		MessageID: sent.ID,
		Text:      slotFrame("BIKA Pro Slot", "locking first reel…", lockA, midB, midC),
		ParseMode: tgmodels.ParseModeHTML,
	})

	time.Sleep(260 * time.Millisecond)

	vipRate := normalizeVIPRate(treasury.VIPWinRate)
	rtp := normalizeRTP(treasury.SlotRTP)

	finalA, finalB, finalC, mult := slotEngine.Spin(user.IsVIP, vipRate, rtp)

	rawPayout := int64(float64(bet) * mult)
	payout := rawPayout

	treasuryNow, err := store.GetTreasury(ctx, a.DB)
	if err == nil {
		maxPay := int64(float64(treasuryNow.OwnerBalance) * slotEngine.CapPercent)
		if payout > maxPay {
			payout = maxPay
		}
		if payout > treasuryNow.OwnerBalance {
			payout = treasuryNow.OwnerBalance
		}
	}

	if payout > 0 {
		err = store.TreasuryPayToUser(ctx, a.DB, msg.From.ID, payout)
		if err != nil {
			_ = store.TreasuryPayToUser(ctx, a.DB, msg.From.ID, bet)
			_, _ = b.EditMessageText(ctx, &tgbot.EditMessageTextParams{
				ChatID:    msg.Chat.ID,
				MessageID: sent.ID,
				Text:      "⚠️ <b>Payout error</b>\n━━━━━━━━━━━━\nRefund ပြန်ပေးလိုက်ပါတယ်။",
				ParseMode: tgmodels.ParseModeHTML,
			})
			return
		}
	}

	updatedUser, _ := store.GetUser(ctx, a.DB, msg.From.ID)
	headline := slotgame.Headline(finalA, finalB, finalC, payout)

	finalText := fmt.Sprintf(
		"🎰 <b>BIKA Pro Slot</b>\n━━━━━━━━━━━━\n<pre>%s</pre>\n━━━━━━━━━━━━\n<b>%s</b>\nBet: <b>%s</b>\nPayout: <b>%s</b>\nNet: <b>%s</b>\nBalance: <b>%s</b>",
		esc(slotgame.Art(finalA, finalB, finalC)),
		esc(headline),
		fmtInt(bet),
		fmtInt(payout),
		fmtInt(payout-bet),
		fmtInt(updatedUser.Balance),
	)

	_, _ = b.EditMessageText(ctx, &tgbot.EditMessageTextParams{
		ChatID:    msg.Chat.ID,
		MessageID: sent.ID,
		Text:      finalText,
		ParseMode: tgmodels.ParseModeHTML,
	})
}
