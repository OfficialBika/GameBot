package bot

import (
	"context"
	"fmt"
	"strings"

	"bikagame-go/internal/store"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a *App) handleMyChatMember(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	m:=update.MyChatMember; if m==nil { return }
	if m.Chat.Type != models.ChatTypeGroup && m.Chat.Type != models.ChatTypeSupergroup { return }
	g,err:=store.EnsureGroup(ctx,a.DB,m.Chat.ID,m.Chat.Title,m.Chat.Username,false); if err!=nil{return}
	if g.ApprovalStatus=="approved" { return }
	text := "🔔 <b>New Group Approval Request</b>\n━━━━━━━━━━━━\n" + fmt.Sprintf("👥 Group: %s\n\n", groupLabel(g)) + "ဒီ group ကို approve ပေးမလား?"
	kb := &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{
		{Text:"✅ Approve", CallbackData: fmt.Sprintf("groupapprove:%d", g.GroupID)},
		{Text:"❌ Reject", CallbackData: fmt.Sprintf("groupreject:%d", g.GroupID)},
	}}}
	_,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:a.Cfg.OwnerID,Text:text,ParseMode:models.ParseModeHTML,ReplyMarkup:kb})
}

func (a *App) handleApproveGroup(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	cb:=update.CallbackQuery; if cb==nil { return }
	if cb.From.ID != a.Cfg.OwnerID { _,_ = b.AnswerCallbackQuery(ctx,&tgbot.AnswerCallbackQueryParams{CallbackQueryID:cb.ID,Text:"Owner only"}); return }
	var groupID int64; _,_ = fmt.Sscanf(cb.Data, "groupapprove:%d", &groupID)
	_ = store.ApproveGroup(ctx,a.DB,groupID,a.Cfg.OwnerID)
	g,_ := store.GetGroup(ctx,a.DB,groupID)
	_,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:a.Cfg.OwnerID,Text:"✅ <b>Group Approved</b>\n━━━━━━━━━━━━\n"+groupLabel(g),ParseMode:models.ParseModeHTML})
	_,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:groupID,Text:"✅ <b>This group has been approved by owner.</b>",ParseMode:models.ParseModeHTML})
	_,_ = b.AnswerCallbackQuery(ctx,&tgbot.AnswerCallbackQueryParams{CallbackQueryID:cb.ID,Text:"Approved"})
}

func (a *App) handleRejectGroup(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	cb:=update.CallbackQuery; if cb==nil { return }
	if cb.From.ID != a.Cfg.OwnerID { _,_ = b.AnswerCallbackQuery(ctx,&tgbot.AnswerCallbackQueryParams{CallbackQueryID:cb.ID,Text:"Owner only"}); return }
	var groupID int64; _,_ = fmt.Sscanf(cb.Data, "groupreject:%d", &groupID)
	_ = store.RejectGroup(ctx,a.DB,groupID,a.Cfg.OwnerID)
	g,_ := store.GetGroup(ctx,a.DB,groupID)
	_,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:a.Cfg.OwnerID,Text:"❌ <b>Group Rejected</b>\n━━━━━━━━━━━━\n"+groupLabel(g),ParseMode:models.ParseModeHTML})
	_,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:groupID,Text:"❌ <b>This group was rejected by owner.</b>",ParseMode:models.ParseModeHTML})
	_,_ = b.AnswerCallbackQuery(ctx,&tgbot.AnswerCallbackQueryParams{CallbackQueryID:cb.ID,Text:"Rejected"})
}

func fullName(first,last string) string { return strings.TrimSpace(first+" "+last) }
