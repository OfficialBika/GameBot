package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	appmodels "bikagame-go/internal/models"
	"bikagame-go/internal/store"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a *App) handleTop10(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	msg:=update.Message; if msg==nil { return }
	users,err:=store.TopUsersByBalance(ctx,a.DB,10); if err!=nil{return}
	if len(users)==0 { _,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:"📭 No users yet.",ParseMode:models.ParseModeHTML}); return }
	lines:=[]string{"🏆 <b>Top 10 Richest</b>","━━━━━━━━━━━━"}
	for i,u := range users { name:=strings.TrimSpace(u.FullName); if name==""{name="User"}; lines=append(lines, fmt.Sprintf("%d. <b>%s</b> — <b>%s</b> %s", i+1, esc(name), fmtInt(u.Balance), esc(a.Cfg.Coin))) }
	_,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:strings.Join(lines,"\n"),ParseMode:models.ParseModeHTML})
}

func parseGiftAmount(text string) (int64, bool) { parts:=strings.Fields(text); if len(parts)<2{return 0,false}; n,err:=strconv.ParseInt(parts[len(parts)-1],10,64); if err!=nil||n<=0{return 0,false}; return n,true }

func (a *App) handleGift(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	msg:=update.Message; if msg==nil || msg.From==nil { return }
	amount,ok:=parseGiftAmount(msg.Text); if !ok { _,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:"🎁 <b>Gift Usage</b>\n━━━━━━━━━━━━\n• Reply + <code>/gift 500</code>\n• Reply + <code>.gift 500</code>\n• <code>/gift @username 500</code>",ParseMode:models.ParseModeHTML}); return }
	fromUser,err:=store.EnsureUser(ctx,a.DB,msg.From.ID,msg.From.Username,msg.From.FirstName,msg.From.LastName,strings.TrimSpace(msg.From.FirstName+" "+msg.From.LastName)); if err!=nil{return}
	var toUser *appmodels.User
	if msg.ReplyToMessage!=nil && msg.ReplyToMessage.From!=nil && !msg.ReplyToMessage.From.IsBot {
		r:=msg.ReplyToMessage.From; if r.ID==msg.From.ID { _,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:"😅 ကိုယ့်ကိုကိုယ် gift မပို့နိုင်ပါ။",ParseMode:models.ParseModeHTML}); return }
		toUser,err=store.EnsureUser(ctx,a.DB,r.ID,r.Username,r.FirstName,r.LastName,strings.TrimSpace(r.FirstName+" "+r.LastName)); if err!=nil{return}
	} else {
		parts:=strings.Fields(msg.Text); if len(parts)<3 { _,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:"👤 Reply သို့ <code>/gift @username 500</code> သုံးပါ။",ParseMode:models.ParseModeHTML}); return }
		target:=parts[1]; if !strings.HasPrefix(target,"@") { _,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:"👤 Target username ကို <code>@username</code> ပုံစံနဲ့ပေးပါ။",ParseMode:models.ParseModeHTML}); return }
		toUser,err=store.GetUserByUsername(ctx,a.DB,target); if err!=nil || toUser==nil { _,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:"⚠️ ဒီ username ကို မတွေ့ပါ။",ParseMode:models.ParseModeHTML}); return }
		if toUser.UserID==msg.From.ID { _,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:"😅 ကိုယ့်ကိုကိုယ် gift မပို့နိုင်ပါ။",ParseMode:models.ParseModeHTML}); return }
	}
	if err:=store.TransferBalance(ctx,a.DB,fromUser.UserID,toUser.UserID,amount); err!=nil { _,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:"❌ လက်ကျန်ငွေ မလုံလောက်ပါ။",ParseMode:models.ParseModeHTML}); return }
	updated,_:=store.GetUser(ctx,a.DB,fromUser.UserID); toName:=strings.TrimSpace(toUser.FullName); if toName==""{toName="User"}
	_,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:fmt.Sprintf("🎁 <b>Gift Success</b>\n━━━━━━━━━━━━\nTo: <b>%s</b>\nAmount: <b>%s</b> %s\nYour Balance: <b>%s</b> %s",esc(toName),fmtInt(amount),esc(a.Cfg.Coin),fmtInt(updated.Balance),esc(a.Cfg.Coin)),ParseMode:models.ParseModeHTML})
}
