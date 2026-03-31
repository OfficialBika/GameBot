package bot

import (
	"context"
	"fmt"

	"bikagame-go/internal/store"
	tgbot "github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"
)

func (a *App) handleMyChatMember(ctx context.Context, b *tgbot.Bot, update *tgmodels.Update) {
	m := update.MyChatMember
	if m == nil {
		return
	}

	if m.Chat.Type != tgmodels.ChatTypeGroup && m.Chat.Type != tgmodels.ChatTypeSupergroup {
		return
	}

	oldStatus := ""
	newStatus := ""
	if m.OldChatMember != nil {
		oldStatus = string(m.OldChatMember.Status)
	}
	if m.NewChatMember != nil {
		newStatus = string(m.NewChatMember.Status)
	}

	becameMember := (oldStatus == "left" || oldStatus == "kicked") && (newStatus == "member" || newStatus == "administrator")
	if !becameMember {
		return
	}

	g, err := store.EnsureGroup(ctx, a.DB, m.Chat.ID, m.Chat.Title, m.Chat.Username, newStatus == "administrator")
	if err != nil {
		return
	}
	if g.ApprovalStatus == "approved" {
		return
	}

	text := "🔔 <b>New Group Approval Request</b>\n━━━━━━━━━━━━\n" +
		fmt.Sprintf("👥 Group: %s\n", groupLabel(g)) +
		fmt.Sprintf("🤖 Bot Admin: <b>%t</b>\n\n", g.BotIsAdmin) +
		"ဒီ group ကို approve ပေးမလား?"

	kb := &tgmodels.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgmodels.InlineKeyboardButton{
			{
				{Text: "✅ Approve", CallbackData: fmt.Sprintf("groupapprove:%d", g.GroupID)},
				{Text: "❌ Reject", CallbackData: fmt.Sprintf("groupreject:%d", g.GroupID)},
			},
		},
	}

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:      a.Cfg.OwnerID,
		Text:        text,
		ParseMode:   tgmodels.ParseModeHTML,
		ReplyMarkup: kb,
	})
}

func (a *App) handleApproveGroup(ctx context.Context, b *tgbot.Bot, update *tgmodels.Update) {
	cb := update.CallbackQuery
	if cb == nil || cb.From == nil || cb.From.ID != a.Cfg.OwnerID {
		return
	}

	var groupID int64
	_, _ = fmt.Sscanf(cb.Data, "groupapprove:%d", &groupID)

	_ = store.ApproveGroup(ctx, a.DB, groupID, a.Cfg.OwnerID)
	g, _ := store.GetGroup(ctx, a.DB, groupID)

	if cb.Message != nil {
		_, _ = b.EditMessageText(ctx, &tgbot.EditMessageTextParams{
			ChatID:    cb.Message.Message.Chat.ID,
			MessageID: cb.Message.Message.ID,
			Text:      "✅ <b>Group Approved</b>\n━━━━━━━━━━━━\n" + fmt.Sprintf("👥 Group: %s", groupLabel(g)),
			ParseMode: tgmodels.ParseModeHTML,
		})
	}

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:    groupID,
		Text:      "✅ <b>This group has been approved by owner.</b>",
		ParseMode: tgmodels.ParseModeHTML,
	})

	_, _ = b.AnswerCallbackQuery(ctx, &tgbot.AnswerCallbackQueryParams{
		CallbackQueryID: cb.ID,
		Text:            "Approved",
	})
}

func (a *App) handleRejectGroup(ctx context.Context, b *tgbot.Bot, update *tgmodels.Update) {
	cb := update.CallbackQuery
	if cb == nil || cb.From == nil || cb.From.ID != a.Cfg.OwnerID {
		return
	}

	var groupID int64
	_, _ = fmt.Sscanf(cb.Data, "groupreject:%d", &groupID)

	_ = store.RejectGroup(ctx, a.DB, groupID, a.Cfg.OwnerID)
	g, _ := store.GetGroup(ctx, a.DB, groupID)

	if cb.Message != nil {
		_, _ = b.EditMessageText(ctx, &tgbot.EditMessageTextParams{
			ChatID:    cb.Message.Message.Chat.ID,
			MessageID: cb.Message.Message.ID,
			Text:      "❌ <b>Group Rejected</b>\n━━━━━━━━━━━━\n" + fmt.Sprintf("👥 Group: %s", groupLabel(g)),
			ParseMode: tgmodels.ParseModeHTML,
		})
	}

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:    groupID,
		Text:      "❌ <b>This group was rejected by owner.</b>",
		ParseMode: tgmodels.ParseModeHTML,
	})

	_, _ = b.AnswerCallbackQuery(ctx, &tgbot.AnswerCallbackQueryParams{
		CallbackQueryID: cb.ID,
		Text:            "Rejected",
	})
}
