package main

import (
	"fmt"
	"strings"
)

func renderREADME(cat Catalog) string {
	var b strings.Builder
	id := cat.Identity

	fmt.Fprintf(&b, "# %s\n\n", id.Name)
	for _, para := range splitParas(id.Intro) {
		b.WriteString(para + "\n\n")
	}

	writeList := func(heading string, products []Product) {
		if len(products) == 0 {
			return
		}
		fmt.Fprintf(&b, "## %s\n\n", heading)
		for _, p := range products {
			what := collapseWS(p.Summary)
			if p.Metric.Value != "" {
				what = strings.TrimSuffix(what, ".") + " (" + p.Metric.Label + ": " + p.Metric.Value + ")."
			}
			meta := []string{"`" + p.Language + "`", p.License}
			if p.ProofURL != "" {
				meta = append(meta, fmt.Sprintf("[%s](%s)", p.Proof, p.ProofURL))
			} else {
				meta = append(meta, "no public tag")
			}
			if p.DemoURL != "" {
				meta = append(meta, fmt.Sprintf("[live demo](%s)", p.DemoURL))
			}
			fmt.Fprintf(&b, "- **[%s](%s)** — %s %s\n", p.Name, p.URL, what, strings.Join(meta, " · "))
		}
		b.WriteString("\n")
	}
	writeList("Selected work", cat.Featured)
	writeList("More on GitHub", cat.More)
	if id.Site != "" {
		fmt.Fprintf(&b, "Longer write-ups: [%s](%s).\n\n", strings.TrimPrefix(id.Site, "https://"), id.Site)
	}

	b.WriteString("<!-- Generated from catalog.yaml + data/products.json by tools/render. Edit the source contracts, then `make render`. -->\n")
	return b.String()
}

func renderPreview(cat Catalog) string {
	var b strings.Builder
	b.WriteString("<!doctype html>\n<meta charset=\"utf-8\">\n<title>nstranquist profile preview</title>\n")
	b.WriteString("<style>body{margin:0;background:#10100f;color:#f0eee7;font:16px/1.5 'Avenir Next',sans-serif}main{max-width:920px;margin:0 auto;padding:48px 24px}h1,h2{font-weight:500}img{display:block;width:100%;height:auto;margin:20px 0}section{margin:48px 0}.pair{display:grid;gap:16px}@media(min-width:900px){.pair{grid-template-columns:1fr 1fr}}</style>\n")
	b.WriteString("<main>\n<h1>Local preview</h1>\n<p>Banner, catalog board, and product cards stay in this repo. They are not embedded on the published GitHub profile README.</p>\n")
	b.WriteString("<section><h2>Banner</h2><div class=\"pair\"><img src=\"assets/banner-dark.svg\" alt=\"banner dark\"><img src=\"assets/banner-light.svg\" alt=\"banner light\"></div></section>\n")
	b.WriteString("<section><h2>Catalog</h2><div class=\"pair\"><img src=\"assets/catalog-dark.svg\" alt=\"catalog dark\"><img src=\"assets/catalog-light.svg\" alt=\"catalog light\"></div></section>\n")
	b.WriteString("<section><h2>Cards</h2><div class=\"pair\">\n")
	for _, p := range append(append([]Product{}, cat.Featured...), cat.More...) {
		slug := slugFor(p)
		fmt.Fprintf(&b, "<img src=\"assets/cards/%s-dark.svg\" alt=\"%s dark\">\n", slug, xml(p.Name))
		fmt.Fprintf(&b, "<img src=\"assets/cards/%s-light.svg\" alt=\"%s light\">\n", slug, xml(p.Name))
	}
	b.WriteString("</div></section>\n</main>\n")
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
