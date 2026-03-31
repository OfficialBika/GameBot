package slot

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
)

type SymbolWeight struct {
	S string
	W int
}

type Engine struct {
	MinBet     int64
	MaxBet     int64
	CooldownMS int64
	CapPercent float64
	Reels      [][]SymbolWeight
	BasePayout map[string]float64
}

func New() *Engine {
	return &Engine{
		MinBet:     50,
		MaxBet:     5000,
		CooldownMS: 700,
		CapPercent: 0.30,
		Reels: [][]SymbolWeight{
			{
				{S: "🍒", W: 3200},
				{S: "🍋", W: 2200},
				{S: "🍉", W: 1500},
				{S: "🔔", W: 900},
				{S: "⭐", W: 450},
				{S: "BAR", W: 200},
				{S: "7", W: 100},
			},
			{
				{S: "🍒", W: 3200},
				{S: "🍋", W: 2200},
				{S: "🍉", W: 1500},
				{S: "🔔", W: 900},
				{S: "⭐", W: 450},
				{S: "BAR", W: 200},
				{S: "7", W: 100},
			},
			{
				{S: "🍒", W: 3200},
				{S: "🍋", W: 2200},
				{S: "🍉", W: 1500},
				{S: "🔔", W: 900},
				{S: "⭐", W: 450},
				{S: "BAR", W: 200},
				{S: "7", W: 100},
			},
		},
		BasePayout: map[string]float64{
			"7,7,7":       20.0,
			"BAR,BAR,BAR": 15.0,
			"⭐,⭐,⭐":       12.0,
			"🔔,🔔,🔔":       9.0,
			"🍉,🍉,🍉":       7.0,
			"🍋,🍋,🍋":       5.0,
			"🍒,🍒,🍒":       3.0,
			"ANY2":        1.5,
		},
	}
}

func (e *Engine) WeightedPick(items []SymbolWeight) string {
	total := 0
	for _, it := range items {
		total += it.W
	}
	if total <= 0 {
		return items[len(items)-1].S
	}

	r := rand.Float64() * float64(total)
	for _, it := range items {
		r -= float64(it.W)
		if r <= 0 {
			return it.S
		}
	}
	return items[len(items)-1].S
}

func (e *Engine) RandomSymbol(reel []SymbolWeight) string {
	if len(reel) == 0 {
		return "?"
	}
	return reel[rand.IntN(len(reel))].S
}

func IsAnyTwo(a, b, c string) bool {
	return (a == b && a != c) || (a == c && a != b) || (b == c && b != a)
}

func (e *Engine) CalcMultiplier(a, b, c string, payout map[string]float64) float64 {
	key := fmt.Sprintf("%s,%s,%s", a, b, c)
	if v, ok := payout[key]; ok {
		return v
	}
	if IsAnyTwo(a, b, c) {
		return payout["ANY2"]
	}
	return 0
}

func (e *Engine) CurrentPayouts(rtp float64) map[string]float64 {
	if rtp <= 0 {
		rtp = 0.90
	}
	if rtp > 1 {
		rtp = rtp / 100.0
	}
	rtp = math.Max(0.50, math.Min(0.98, rtp))

	base := 0.90
	factor := rtp / base

	out := make(map[string]float64, len(e.BasePayout))
	for k, v := range e.BasePayout {
		out[k] = math.Round((v*factor)*10000) / 10000
	}
	return out
}

func (e *Engine) SpinNormal(rtp float64, payout map[string]float64) (string, string, string) {
	for i := 0; i < 25; i++ {
		a := e.WeightedPick(e.Reels[0])
		b := e.WeightedPick(e.Reels[1])
		c := e.WeightedPick(e.Reels[2])
		m := e.CalcMultiplier(a, b, c, payout)
		if m <= 0 {
			return a, b, c
		}
		if rand.Float64() < rtp {
			return a, b, c
		}
	}

	for {
		a := e.RandomSymbol(e.Reels[0])
		b := e.RandomSymbol(e.Reels[1])
		c := e.RandomSymbol(e.Reels[2])
		if e.CalcMultiplier(a, b, c, payout) <= 0 {
			return a, b, c
		}
	}
}

func (e *Engine) SpinVIP(vipRate int, payout map[string]float64) (string, string, string) {
	chance := float64(vipRate) / 100.0
	if chance < 0 {
		chance = 0
	}
	if chance > 1 {
		chance = 1
	}

	if rand.Float64() < chance {
		winning := make([][3]string, 0)
		for k, v := range payout {
			if k == "ANY2" || v <= 0 {
				continue
			}
			parts := strings.Split(k, ",")
			if len(parts) != 3 {
				continue
			}
			winning = append(winning, [3]string{parts[0], parts[1], parts[2]})
		}
		if len(winning) > 0 {
			w := winning[rand.IntN(len(winning))]
			return w[0], w[1], w[2]
		}
	}
	return e.SpinNormal(0.90, payout)
}

func (e *Engine) Spin(isVIP bool, vipRate int, rtp float64) (string, string, string, float64) {
	payouts := e.CurrentPayouts(rtp)
	if isVIP {
		a, b, c := e.SpinVIP(vipRate, payouts)
		return a, b, c, e.CalcMultiplier(a, b, c, payouts)
	}
	a, b, c := e.SpinNormal(rtp, payouts)
	return a, b, c, e.CalcMultiplier(a, b, c, payouts)
}

func Box(x string) string {
	switch x {
	case "BAR":
		return "BAR"
	case "7":
		return "7️⃣"
	default:
		return x
	}
}

func Art(a, b, c string) string {
	return fmt.Sprintf("┏━━━━━━━━━━━━━━━━━━┓\n┃  %s  |  %s  |  %s  ┃\n┗━━━━━━━━━━━━━━━━━━┛", Box(a), Box(b), Box(c))
}

func Headline(a, b, c string, payout int64) string {
	if a == "7" && b == "7" && c == "7" {
		return "🏆 JACKPOT 777!"
	}
	if payout > 0 {
		return "✅ WIN"
	}
	return "❌ LOSE"
}
