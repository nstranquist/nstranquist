package main

import (
	"fmt"
	"html"
	"strings"
)

type theme struct {
	Name    string
	BG      string
	Panel   string
	Ink     string
	Muted   string
	Faint   string
	Line    string
	Signal  string
	Ember   string
	Chip    string
	AltRow  string
	Dot     string
	MarkOn  string
	MarkOff string
}

func themeDark() theme {
	return theme{
		Name:    "dark",
		BG:      "#10100f",
		Panel:   "#171716",
		Ink:     "#f0eee7",
		Muted:   "#918f88",
		Faint:   "#6b6962",
		Line:    "rgba(240,238,231,0.16)",
		Signal:  "#5c7cff",
		Ember:   "#f09462",
		Chip:    "#1d1d1b",
		AltRow:  "rgba(240,238,231,0.03)",
		Dot:     "rgba(240,238,231,0.07)",
		MarkOn:  "#5c7cff",
		MarkOff: "rgba(240,238,231,0.08)",
	}
}

func themeLight() theme {
	return theme{
		Name:    "light",
		BG:      "#f4f1ea",
		Panel:   "#ebe6dc",
		Ink:     "#161614",
		Muted:   "#6f6c64",
		Faint:   "#8a877e",
		Line:    "rgba(22,22,20,0.12)",
		Signal:  "#3d5fe0",
		Ember:   "#c05d2e",
		Chip:    "#e3ddd2",
		AltRow:  "rgba(22,22,20,0.035)",
		Dot:     "rgba(22,22,20,0.07)",
		MarkOn:  "#3d5fe0",
		MarkOff: "rgba(22,22,20,0.08)",
	}
}

const (
	serif = `'Iowan Old Style', Palatino, 'Times New Roman', serif`
	sans  = `'Avenir Next', 'Segoe UI', system-ui, sans-serif`
	mono  = `'SFMono-Regular', ui-monospace, Menlo, Consolas, monospace`
)

func renderBanner(cat Catalog, t theme) string {
	var b strings.Builder
	const w, h = 880, 228
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d" role="img">`+"\n", w, h, w, h)
	fmt.Fprintf(&b, `  <title>%s — %s</title>`+"\n", xml(cat.Identity.Name), xml(cat.Identity.Role))
	fmt.Fprintf(&b, `  <rect width="%d" height="%d" rx="22" fill="%s"/>`+"\n", w, h, t.BG)
	fmt.Fprintf(&b, `  <rect x="18" y="18" width="%d" height="%d" rx="16" fill="none" stroke="%s" stroke-width="1"/>`+"\n", w-36, h-36, t.Line)

	// Catalog-grid N
	b.WriteString(markN(40, 70, t))

	fmt.Fprintf(&b, `  <text x="148" y="94" fill="%s" font-family="%s" font-size="30" letter-spacing="1.4">%s</text>`+"\n", t.Ink, serif, xml(strings.ToUpper(cat.Identity.Name)))
	fmt.Fprintf(&b, `  <text x="148" y="126" fill="%s" font-family="%s" font-size="16">%s</text>`+"\n", t.Muted, sans, xml(cat.Identity.Role))
	fmt.Fprintf(&b, `  <text x="148" y="154" fill="%s" font-family="%s" font-size="12" letter-spacing="0.6">%s  ·  public catalog  ·  tagged releases</text>`+"\n", t.Signal, mono, xml(cat.Identity.Location))

	chips := []string{"Platforms", "Local tools", "Full-stack"}
	x := 148.0
	for _, chip := range chips {
		cw := 7.2*float64(len(chip)) + 26
		fmt.Fprintf(&b, `  <rect x="%.1f" y="182" width="%.1f" height="22" rx="11" fill="%s" stroke="%s"/>`+"\n", x, cw, t.Chip, t.Line)
		fmt.Fprintf(&b, `  <text x="%.1f" y="197" fill="%s" font-family="%s" font-size="11">%s</text>`+"\n", x+13, t.Muted, sans, xml(chip))
		x += cw + 8
	}
	b.WriteString("</svg>\n")
	return b.String()
}

func markN(x, y float64, t theme) string {
	// 4x4 catalog cells forming an N
	on := [][2]int{
		{0, 0}, {0, 1}, {0, 2}, {0, 3},
		{1, 1},
		{2, 2},
		{3, 0}, {3, 1}, {3, 2}, {3, 3},
	}
	filled := map[[2]int]bool{}
	for _, p := range on {
		filled[p] = true
	}
	const cell, gap, n = 16.0, 5.0, 4
	var b strings.Builder
	fmt.Fprintf(&b, `  <g aria-hidden="true">`+"\n")
	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			fill := t.MarkOff
			if filled[[2]int{col, row}] {
				fill = t.MarkOn
			}
			cx := x + float64(col)*(cell+gap)
			cy := y + float64(row)*(cell+gap)
			fmt.Fprintf(&b, `    <rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="3.5" fill="%s"/>`+"\n", cx, cy, cell, cell, fill)
		}
	}
	fmt.Fprintf(&b, `  </g>`+"\n")
	return b.String()
}

func renderCatalogBoard(cat Catalog, t theme) string {
	const (
		w      = 880
		pad    = 22
		header = 58
		colH   = 30
		rowH   = 52
	)
	h := pad + header + colH + rowH*len(cat.Featured) + 44 + pad
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d" role="img">`+"\n", w, h, w, h)
	fmt.Fprintf(&b, `  <title>Public catalog — %d products</title>`+"\n", len(cat.Featured))
	fmt.Fprintf(&b, `  <rect width="%d" height="%d" rx="22" fill="%s"/>`+"\n", w, h, t.BG)
	fmt.Fprintf(&b, `  <rect x="14" y="14" width="%d" height="%d" rx="16" fill="%s" stroke="%s"/>`+"\n", w-28, h-28, t.Panel, t.Line)

	fmt.Fprintf(&b, `  <circle cx="42" cy="48" r="5" fill="%s"/>`+"\n", t.Signal)
	fmt.Fprintf(&b, `  <text x="58" y="53" fill="%s" font-family="%s" font-size="16" letter-spacing="1.8">PUBLIC CATALOG</text>`+"\n", t.Ink, sans)
	fmt.Fprintf(&b, `  <text x="836" y="53" text-anchor="end" fill="%s" font-family="%s" font-size="12">%d products  ·  tagged releases</text>`+"\n", t.Muted, mono, len(cat.Featured))
	fmt.Fprintf(&b, `  <line x1="32" y1="72" x2="848" y2="72" stroke="%s" stroke-width="1"/>`+"\n", t.Line)

	headers := []struct {
		x    int
		text string
	}{
		{36, "PRODUCT"},
		{300, "KIND"},
		{478, "STACK"},
		{612, "LICENSE"},
		{748, "RELEASE"},
	}
	hy := pad + header + 18
	for _, col := range headers {
		fmt.Fprintf(&b, `  <text x="%d" y="%d" fill="%s" font-family="%s" font-size="10" letter-spacing="1.4">%s</text>`+"\n", col.x, hy, t.Faint, mono, col.text)
	}

	for i, p := range cat.Featured {
		y := pad + header + colH + i*rowH
		if i%2 == 0 {
			fmt.Fprintf(&b, `  <rect x="28" y="%d" width="824" height="%d" rx="8" fill="%s"/>`+"\n", y, rowH-4, t.AltRow)
		}
		ty := y + 32
		fmt.Fprintf(&b, `  <text x="36" y="%d" fill="%s" font-family="%s" font-size="15">%s</text>`+"\n", ty, t.Ink, sans, xml(p.Name))
		fmt.Fprintf(&b, `  <text x="36" y="%d" fill="%s" font-family="%s" font-size="10">%s</text>`+"\n", ty+14, t.Faint, mono, xml(p.Repo))
		fmt.Fprintf(&b, `  <text x="300" y="%d" fill="%s" font-family="%s" font-size="13">%s</text>`+"\n", ty, t.Muted, sans, xml(p.Lane))
		fmt.Fprintf(&b, `  <text x="478" y="%d" fill="%s" font-family="%s" font-size="13">%s</text>`+"\n", ty, t.Muted, sans, xml(p.Language))
		fmt.Fprintf(&b, `  <rect x="612" y="%d" width="112" height="22" rx="11" fill="%s" stroke="%s"/>`+"\n", ty-15, t.Chip, t.Line)
		fmt.Fprintf(&b, `  <text x="668" y="%d" text-anchor="middle" fill="%s" font-family="%s" font-size="11">%s</text>`+"\n", ty, t.Ink, mono, xml(p.License))
		fmt.Fprintf(&b, `  <text x="748" y="%d" fill="%s" font-family="%s" font-size="13">%s</text>`+"\n", ty, t.Signal, mono, xml(p.Proof))
	}

	fy := h - pad - 16
	fmt.Fprintf(&b, `  <line x1="32" y1="%d" x2="848" y2="%d" stroke="%s" stroke-width="1"/>`+"\n", fy-18, fy-18, t.Line)
	note := cat.Footnote
	if note == "" {
		note = "Public products only. Unproven claims are rejected. No invented numbers."
	}
	fmt.Fprintf(&b, `  <text x="36" y="%d" fill="%s" font-family="%s" font-size="11">%s</text>`+"\n", fy, t.Faint, mono, xml(note))
	b.WriteString("</svg>\n")
	return b.String()
}

func renderCard(p Product, t theme) string {
	const w, h = 420, 168
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d" role="img">`+"\n", w, h, w, h)
	fmt.Fprintf(&b, `  <title>%s — %s</title>`+"\n", xml(p.Name), xml(p.Summary))
	fmt.Fprintf(&b, `  <rect width="%d" height="%d" rx="18" fill="%s" stroke="%s"/>`+"\n", w, h, t.BG, t.Line)
	fmt.Fprintf(&b, `  <text x="22" y="36" fill="%s" font-family="%s" font-size="18">%s</text>`+"\n", t.Ink, serif, xml(p.Name))
	fmt.Fprintf(&b, `  <rect x="286" y="18" width="112" height="22" rx="11" fill="%s" stroke="%s"/>`+"\n", t.Chip, t.Line)
	fmt.Fprintf(&b, `  <text x="342" y="33" text-anchor="middle" fill="%s" font-family="%s" font-size="11">%s</text>`+"\n", t.Ink, mono, xml(p.License))
	fmt.Fprintf(&b, `  <text x="22" y="58" fill="%s" font-family="%s" font-size="12">%s  ·  %s</text>`+"\n", t.Muted, mono, xml(p.Lane), xml(p.Language))

	summary := wrapSVGText(p.Summary, 46)
	y := 88
	for _, line := range summary {
		fmt.Fprintf(&b, `  <text x="22" y="%d" fill="%s" font-family="%s" font-size="14">%s</text>`+"\n", y, t.Ink, sans, xml(line))
		y += 20
	}
	if p.Metric.Value != "" {
		fmt.Fprintf(&b, `  <text x="22" y="148" fill="%s" font-family="%s" font-size="11">%s</text>`+"\n", t.Ember, mono, xml(p.Metric.Value))
	}
	fmt.Fprintf(&b, `  <text x="398" y="148" text-anchor="end" fill="%s" font-family="%s" font-size="12">%s</text>`+"\n", t.Signal, mono, xml(p.Proof))
	b.WriteString("</svg>\n")
	return b.String()
}

func wrapSVGText(s string, width int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return nil
	}
	var lines []string
	cur := words[0]
	for _, w := range words[1:] {
		if len(cur)+1+len(w) > width {
			lines = append(lines, cur)
			cur = w
			continue
		}
		cur += " " + w
	}
	return append(lines, cur)
}

func xml(s string) string {
	return html.EscapeString(s)
}
