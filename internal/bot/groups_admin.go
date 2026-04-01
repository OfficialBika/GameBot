package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"bikagame-go/internal/store"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func isOwner(a *App, userID int64) bool { return userID == a.Cfg.OwnerID }
func parseGroupIDArg(text string) (int64, bool) { parts:=strings.Fields(text); if len(parts)<2{return 0,false}; n,err:=strconv.ParseInt(parts[1],10,64); if err!=nil{return 0,false}; return n,true }

func (a *App) handlePendingGroups(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	msg:=update.Message; if msg==nil || msg.From==nil || !isOwner(a,msg.From.ID) { return }
	groups,err:=store.ListPendingGroups(ctx,a.DB,50); if err!=nil{return}
	if len(groups)==0 { _,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:"📭 <b>Pending Groups</b>\n━━━━━━━━━━━━\nNo pending groups.",ParseMode:models.ParseModeHTML}); return }
	lines:=[]string{"🕒 <b>Pending Groups</b>","━━━━━━━━━━━━"}
	for i,g := range groups { lines=append(lines, fmt.Sprintf("%d. %s",i+1,groupLabel(&g))) }
	lines=append(lines,"","Approve: <code>/approve -1001234567890</code>","Reject: <code>/reject -1001234567890</code>")
	_,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:strings.Join(lines,"\n"),ParseMode:models.ParseModeHTML})
}

func (a *App) handleApproveCmd(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	msg:=update.Message; if msg==nil || msg.From==nil || !isOwner(a,msg.From.ID) { return }
	groupID,ok:=parseGroupIDArg(msg.Text); if !ok { _,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:"✅ Usage: <code>/approve -1001234567890</code>",ParseMode:models.ParseModeHTML}); return }
	if err:=store.ApproveGroup(ctx,a.DB,groupID,msg.From.ID); err!=nil{return}
	g,_:=store.GetGroup(ctx,a.DB,groupID)
	_,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:"✅ <b>Group Approved</b>\n━━━━━━━━━━━━\n"+groupLabel(g),ParseMode:models.ParseModeHTML})
	_,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:groupID,Text:"✅ <b>This group has been approved by owner.</b>",ParseMode:models.ParseModeHTML})
}

func (a *App) handleRejectCmd(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	msg:=update.Message; if msg==nil || msg.From==nil || !isOwner(a,msg.From.ID) { return }
	groupID,ok:=parseGroupIDArg(msg.Text); if !ok { _,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:"❌ Usage: <code>/reject -1001234567890</code>",ParseMode:models.ParseModeHTML}); return }
	if err:=store.RejectGroup(ctx,a.DB,groupID,msg.From.ID); err!=nil{return}
	g,_:=store.GetGroup(ctx,a.DB,groupID)
	_,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:"❌ <b>Group Rejected</b>\n━━━━━━━━━━━━\n"+groupLabel(g),ParseMode:models.ParseModeHTML})
	_,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:groupID,Text:"❌ <b>This group was rejected by owner.</b>",ParseMode:models.ParseModeHTML})
}
