package bot

import (
	"context"
	"strings"

	"bikagame-go/internal/store"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a *App) handleStart(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	msg:=update.Message; if msg==nil || msg.From==nil { return }
	_,_ = store.EnsureUser(ctx,a.DB,msg.From.ID,msg.From.Username,msg.From.FirstName,msg.From.LastName,strings.TrimSpace(msg.From.FirstName+" "+msg.From.LastName))
	_,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:"👋 <b>Welcome to BIKA Game Bot</b>\n━━━━━━━━━━━━\nUse <code>/balance</code>, <code>/dailyclaim</code>, <code>.slot 100</code>",ParseMode:models.ParseModeHTML})
}

func (a *App) handlePing(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	msg:=update.Message; if msg==nil { return }
	_,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:"🏓 <b>PONG</b>",ParseMode:models.ParseModeHTML})
}
