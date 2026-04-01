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
	"github.com/go-telegram/bot/models"
)

var (
	slotEngine = slotgame.New()
	slotMu sync.Mutex
	activeSlots = map[int64]bool{}
	lastSlotAt = map[int64]int64{}
	maxActiveSlot = 5
)

func canSpin(userID int64)(bool,string){ slotMu.Lock(); defer slotMu.Unlock(); if len(activeSlots)>=maxActiveSlot && !activeSlots[userID] { return false,"⛔ <b>Slot Busy</b>\n━━━━━━━━━━━━\nခဏနားပြီး ပြန်ကြိုးစားပါ။" }; now:=time.Now().UnixMilli(); if now-lastSlotAt[userID] < slotEngine.CooldownMS { waitSec:=int(((slotEngine.CooldownMS-(now-lastSlotAt[userID]))+999)/1000); return false, fmt.Sprintf("⏳ ခဏစောင့်ပါ… (%ds)", waitSec) }; activeSlots[userID]=true; lastSlotAt[userID]=now; return true, "" }
func releaseSpin(userID int64){ slotMu.Lock(); defer slotMu.Unlock(); delete(activeSlots,userID) }
func parseBet(text string)(int64,bool){ parts:=strings.Fields(text); if len(parts)!=2{return 0,false}; n,err:=strconv.ParseInt(parts[1],10,64); if err!=nil||n<=0{return 0,false}; return n,true }
func slotFrame(title,note,a,b,c string) string { return fmt.Sprintf("🎰 <b>%s</b>\n━━━━━━━━━━━━\n<pre>%s</pre>\n━━━━━━━━━━━━\n%s", esc(title), esc(slotgame.Art(a,b,c)), esc(note)) }

func (a *App) handleSlot(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	msg:=update.Message; if msg==nil || msg.From==nil || msg.Text=="" { return }
	if msg.Chat.Type != models.ChatTypeGroup && msg.Chat.Type != models.ChatTypeSupergroup { return }
	bet,ok:=parseBet(msg.Text); if !ok { return }
	if bet<slotEngine.MinBet || bet>slotEngine.MaxBet { _,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:fmt.Sprintf("🎰 <b>BIKA Pro Slot</b>\n━━━━━━━━━━━━\nUsage: <code>.slot 1000</code>\nMin: <b>%s</b>\nMax: <b>%s</b>",fmtInt(slotEngine.MinBet),fmtInt(slotEngine.MaxBet)),ParseMode:models.ParseModeHTML}); return }
	okSpin,reason:=canSpin(msg.From.ID); if !okSpin { _,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:reason,ParseMode:models.ParseModeHTML}); return }
	defer releaseSpin(msg.From.ID)
	user,err:=store.EnsureUser(ctx,a.DB,msg.From.ID,msg.From.Username,msg.From.FirstName,msg.From.LastName,strings.TrimSpace(msg.From.FirstName+" "+msg.From.LastName)); if err!=nil{return}
	if user.Balance<bet { _,_ = b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:fmt.Sprintf("❌ <b>Balance မလုံလောက်ပါ</b>\n━━━━━━━━━━━━\nBet: <b>%s</b>\nYour Balance: <b>%s</b>",fmtInt(bet),fmtInt(user.Balance)),ParseMode:models.ParseModeHTML}); return }
	treasury,err:=store.GetTreasury(ctx,a.DB); if err!=nil{return}
	initA,initB,initC:=slotEngine.RandomSymbol(slotEngine.Reels[0]),slotEngine.RandomSymbol(slotEngine.Reels[1]),slotEngine.RandomSymbol(slotEngine.Reels[2])
	sent,err:=b.SendMessage(ctx,&tgbot.SendMessageParams{ChatID:msg.Chat.ID,Text:slotFrame("BIKA Pro Slot","reels spinning…",initA,initB,initC),ParseMode:models.ParseModeHTML}); if err!=nil{return}
	if err:=store.UserPayToTreasury(ctx,a.DB,msg.From.ID,bet); err!=nil { if errors.Is(err,store.ErrUserInsufficient){ _,_ = b.EditMessageText(ctx,&tgbot.EditMessageTextParams{ChatID:msg.Chat.ID,MessageID:sent.ID,Text:"❌ <b>Balance မလုံလောက်ပါ</b>",ParseMode:models.ParseModeHTML}) }; return }
	time.Sleep(220*time.Millisecond)
	_,_ = b.EditMessageText(ctx,&tgbot.EditMessageTextParams{ChatID:msg.Chat.ID,MessageID:sent.ID,Text:slotFrame("BIKA Pro Slot","locking first reel…",slotEngine.RandomSymbol(slotEngine.Reels[0]),slotEngine.RandomSymbol(slotEngine.Reels[1]),slotEngine.RandomSymbol(slotEngine.Reels[2])),ParseMode:models.ParseModeHTML})
	time.Sleep(260*time.Millisecond)
	finalA,finalB,finalC,mult:=slotEngine.Spin(user.IsVIP, normalizeVIPRate(treasury.VIPWinRate), normalizeRTP(treasury.SlotRTP))
	payout:=int64(float64(bet)*mult)
	if treasuryNow,err:=store.GetTreasury(ctx,a.DB); err==nil { maxPay:=int64(float64(treasuryNow.OwnerBalance)*slotEngine.CapPercent); if payout>maxPay{payout=maxPay}; if payout>treasuryNow.OwnerBalance{payout=treasuryNow.OwnerBalance} }
	if payout>0 { if err:=store.TreasuryPayToUser(ctx,a.DB,msg.From.ID,payout); err!=nil { _ = store.TreasuryPayToUser(ctx,a.DB,msg.From.ID,bet); _,_ = b.EditMessageText(ctx,&tgbot.EditMessageTextParams{ChatID:msg.Chat.ID,MessageID:sent.ID,Text:"⚠️ <b>Payout error</b>\n━━━━━━━━━━━━\nRefund ပြန်ပေးလိုက်ပါတယ်။",ParseMode:models.ParseModeHTML}); return } }
	updated,_:=store.GetUser(ctx,a.DB,msg.From.ID)
	finalText:=fmt.Sprintf("🎰 <b>BIKA Pro Slot</b>\n━━━━━━━━━━━━\n<pre>%s</pre>\n━━━━━━━━━━━━\n<b>%s</b>\nBet: <b>%s</b>\nPayout: <b>%s</b>\nNet: <b>%s</b>\nBalance: <b>%s</b>", esc(slotgame.Art(finalA,finalB,finalC)), esc(slotgame.Headline(finalA,finalB,finalC,payout)), fmtInt(bet), fmtInt(payout), fmtInt(payout-bet), fmtInt(updated.Balance))
	_,_ = b.EditMessageText(ctx,&tgbot.EditMessageTextParams{ChatID:msg.Chat.ID,MessageID:sent.ID,Text:finalText,ParseMode:models.ParseModeHTML})
}
