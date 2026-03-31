package bot

import (
	"fmt"
	"html"
	"strings"

	"bikagame-go/internal/models"
)

func esc(s string) string {
	return html.EscapeString(s)
}

func fmtInt(n int64) string {
	in := fmt.Sprintf("%d", n)
	if len(in) <= 3 {
		return in
	}
	var out []byte
	c := 0
	for i := len(in) - 1; i >= 0; i-- {
		out = append(out, in[i])
		c++
		if c%3 == 0 && i != 0 {
			out = append(out, ',')
		}
	}
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return string(out)
}

func groupLabel(g *models.Group) string {
	if g == nil {
		return "Unknown Group"
	}
	if g.Username != "" {
		return fmt.Sprintf("<b>%s</b> (@%s) — <code>%d</code>", esc(g.Title), esc(strings.TrimPrefix(g.Username, "@")), g.GroupID)
	}
	return fmt.Sprintf("<b>%s</b> — <code>%d</code>", esc(g.Title), g.GroupID)
}

func normalizeVIPRate(v int) int {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

func normalizeRTP(v float64) float64 {
	if v > 1 {
		v = v / 100.0
	}
	if v < 0.50 {
		return 0.50
	}
	if v > 0.98 {
		return 0.98
	}
	return v
}
