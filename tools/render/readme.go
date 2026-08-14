package main

import (
	"fmt"
	"strings"
)

func renderREADME(cat Catalog) string {
	var b strings.Builder
	id := cat.Identity

	b.WriteString("# Public catalog\n\n")
	fmt.Fprintf(&b, "%s — %s, %s.\n\n", id.Name, strings.ToLower(id.Role), id.Location)
	for _, para := range splitParas(id.Intro) {
		b.WriteString(para + "\n\n")
	}

	b.WriteString(figure("assets/banner-dark.svg", "assets/banner-light.svg", id.Name+" — "+id.Role, 880))
	b.WriteString("\n")
	b.WriteString(figure("assets/catalog-dark.svg", "assets/catalog-light.svg", "Public catalog of products", 880))
	b.WriteString("\n")

	b.WriteString("## Selected work\n\n")
	for _, p := range cat.Featured {
		fmt.Fprintf(&b, "### [%s](%s)\n\n", p.Name, p.URL)
		b.WriteString(collapseWS(p.Detail) + "\n\n")
		fmt.Fprintf(&b, "`%s` · %s · [%s](%s)\n\n", p.Language, p.License, p.Proof, p.ProofURL)
		if p.Metric.Value != "" {
			label := p.Metric.Label
			if label != "" {
				label = strings.ToUpper(label[:1]) + label[1:]
			}
			fmt.Fprintf(&b, "%s: **%s**\n\n", label, p.Metric.Value)
		}
	}

	b.WriteString("## How I work\n\n")
	for _, p := range cat.Principles {
		fmt.Fprintf(&b, "- **%s.** %s\n", p.Title, p.Body)
	}
	b.WriteString("\n")

	if len(cat.Glossary) > 0 {
		b.WriteString("## Terms\n\n")
		b.WriteString("Short definitions for words used above.\n\n")
		b.WriteString("| Term | Meaning |\n| --- | --- |\n")
		for _, g := range cat.Glossary {
			fmt.Fprintf(&b, "| %s | %s |\n", g.Term, g.Meaning)
		}
		b.WriteString("\n")
	}

	b.WriteString("## Toolbox\n\n")
	b.WriteString(strings.Join(cat.Toolbox, " · ") + "\n\n")

	b.WriteString("## Contact\n\n")
	fmt.Fprintf(&b, "[LinkedIn](%s)\n\n", id.LinkedIn)

	b.WriteString("<!-- Generated from catalog.yaml by tools/render. Edit the catalog, then `make render`. -->\n")
	return b.String()
}

func renderPreview(cat Catalog) string {
	var b strings.Builder
	b.WriteString("<!doctype html>\n<meta charset=\"utf-8\">\n<title>nstranquist profile preview</title>\n")
	b.WriteString("<style>body{margin:0;background:#10100f;color:#f0eee7;font:16px/1.5 'Avenir Next',sans-serif}main{max-width:920px;margin:0 auto;padding:48px 24px}h1,h2{font-weight:500}img{display:block;width:100%;height:auto;margin:20px 0}section{margin:48px 0}.pair{display:grid;gap:16px}@media(min-width:900px){.pair{grid-template-columns:1fr 1fr}}</style>\n")
	b.WriteString("<main>\n<h1>Local preview</h1>\n<p>Dark and light assets for the GitHub profile README. This file is not the published profile.</p>\n")
	b.WriteString("<section><h2>Banner</h2><div class=\"pair\"><img src=\"assets/banner-dark.svg\" alt=\"banner dark\"><img src=\"assets/banner-light.svg\" alt=\"banner light\"></div></section>\n")
	b.WriteString("<section><h2>Catalog</h2><div class=\"pair\"><img src=\"assets/catalog-dark.svg\" alt=\"catalog dark\"><img src=\"assets/catalog-light.svg\" alt=\"catalog light\"></div></section>\n")
	b.WriteString("<section><h2>Cards</h2><div class=\"pair\">\n")
	for _, p := range cat.Featured {
		slug := slugFor(p)
		fmt.Fprintf(&b, "<img src=\"assets/cards/%s-dark.svg\" alt=\"%s dark\">\n", slug, xml(p.Name))
		fmt.Fprintf(&b, "<img src=\"assets/cards/%s-light.svg\" alt=\"%s light\">\n", slug, xml(p.Name))
	}
	b.WriteString("</div></section>\n</main>\n")
	return b.String()
}

func figure(dark, light, alt string, width int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(
		"<p align=\"center\">\n  <img src=\"%s#gh-dark-mode-only\" alt=\"%s\" width=\"%d\" />\n  <img src=\"%s#gh-light-mode-only\" alt=\"%s\" width=\"%d\" />\n</p>\n",
		dark, xml(alt), width, light, xml(alt), width,
	))
	return b.String()
}

func collapseWS(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func splitParas(s string) []string {
	var out []string
	for _, para := range strings.Split(s, "\n\n") {
		para = collapseWS(para)
		if para != "" {
			out = append(out, para)
		}
	}
	if len(out) == 0 && collapseWS(s) != "" {
		return []string{collapseWS(s)}
	}
	return out
}
