package bot

import (
	"fmt"
	"html"
	"strings"
	"time"

	"bikagame-go/internal/models"
)

func esc(s string) string { return html.EscapeString(s) }
func fmtInt(n int64) string { in:=fmt.Sprintf("%d",n); if len(in)<=3{return in}; var out []byte; c:=0; for i:=len(in)-1;i>=0;i--{ out=append(out,in[i]); c++; if c%3==0 && i!=0 { out=append(out,',') } }; for i,j:=0,len(out)-1;i<j;i,j=i+1,j-1{ out[i],out[j]=out[j],out[i] }; return string(out) }
func normalizeVIPRate(v int) int { if v<0{return 0}; if v>100{return 100}; return v }
func normalizeRTP(v float64) float64 { if v>1{v=v/100}; if v<0.50{return 0.50}; if v>0.98{return 0.98}; return v }
func yangonNow() time.Time { loc,err:=time.LoadLocation("Asia/Yangon"); if err!=nil{return time.Now()}; return time.Now().In(loc) }
func yangonDateKey() string { return yangonNow().Format("2006-01-02") }
func walletRank(balance int64) string { switch { case balance<=0: return "ဖင်ပြောင်ငမွဲ"; case balance<=500: return "ဆင်းရဲသား အိမ်ခြေမဲ့"; case balance<=1000: return "အိမ်ပိုင်ဝန်းပိုင် ဆင်းရဲသား"; case balance<=5000: return "လူလတ်တန်းစား"; case balance<=10000: return "သူဌေးပေါက်စ"; case balance<=100000: return "သိန်းကြွယ်သူဌေး"; case balance<=1000000: return "သန်းကြွယ်သူဌေးအကြီးစား"; case balance<=50000000: return "ကုဋေရှစ်ဆယ် သူဌေးကြီး"; default: return "အာကာသသူဌေး" } }
func groupLabel(g *models.Group) string { if g==nil{return "Unknown Group"}; if g.Username!="" { return fmt.Sprintf("<b>%s</b> (@%s) — <code>%d</code>", esc(g.Title), esc(strings.TrimPrefix(g.Username,"@")), g.GroupID)}; return fmt.Sprintf("<b>%s</b> — <code>%d</code>", esc(g.Title), g.GroupID)}
