package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var requiredIDs = []string{
	"product.docs-puller",
	"product.nicos-catalog",
	"product.openbook",
	"product.agent-ops",
	"product.nicos-hidden-menubar",
	"product.nicos-slot-dock",
	"product.jobkit",
	"product.session-pressure",
}

var moreIDs = []string{
	"product.keepawake",
	"product.wip-commit",
	"product.snapref",
	"product.ngtm",
	"product.nicos-flag-eval",
	"product.nicos-window-switcher",
}

var untaggedIDs = map[string]bool{
	"product.wip-commit":            true,
	"product.ngtm":                  true,
	"product.nicos-flag-eval":       true,
	"product.nicos-window-switcher": true,
}

var forbidden = []string{
	"noise",
	"getnoise.com",
	"Bayer",
	"Enhearten",
	"EduRAIN",
	"EduRain",
	"SmartSpectra",
	"smartspectra",
	"nvault",
	"pw-harness",
	"Garrid",
	"farm-game",
	"runescape-sim",
	"MemeBattle",
	"sol-surfer",
	"idle-time",
	"1,000+",
	"30,000+",
	"3,400+",
	"visitor",
	"github-readme-stats",
	"$50M",
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "catalog.yaml")); err != nil {
		t.Fatalf("catalog.yaml not found from %s: %v", wd, err)
	}
	return root
}

func loadTestCatalog(t *testing.T) Catalog {
	t.Helper()
	root := repoRoot(t)
	cat, err := loadCatalog(filepath.Join(root, "catalog.yaml"), filepath.Join(root, "data", "products.json"))
	if err != nil {
		t.Fatal(err)
	}
	return cat
}

func TestCatalogFeaturedOrderAndProof(t *testing.T) {
	cat := loadTestCatalog(t)
	if len(cat.Featured) != len(requiredIDs) {
		t.Fatalf("featured count %d, want %d", len(cat.Featured), len(requiredIDs))
	}
	for _, p := range cat.Featured {
		if p.PublicState == "" {
			t.Errorf("%s missing public state", p.ID)
		}
	}
	if len(cat.More) != len(moreIDs) {
		t.Fatalf("more count %d, want %d", len(cat.More), len(moreIDs))
	}
	for i, id := range moreIDs {
		if cat.More[i].ID != id {
			t.Errorf("more[%d]=%s, want %s", i, cat.More[i].ID, id)
		}
	}
	for i, id := range requiredIDs {
		p := cat.Featured[i]
		if p.ID != id {
			t.Errorf("featured[%d]=%s, want %s", i, p.ID, id)
		}
		if untaggedIDs[id] {
			if p.Proof != "no public tag" || p.ProofURL != "" {
				t.Errorf("%s must preserve missing tag truth: %#v", id, p)
			}
		} else {
			if !strings.HasPrefix(p.Proof, "v") {
				t.Errorf("%s proof %q should be a version tag", id, p.Proof)
			}
			if !strings.Contains(p.ProofURL, "/tree/"+p.Proof) {
				t.Errorf("%s proof_url %s does not match proof %s", id, p.ProofURL, p.Proof)
			}
		}
		if p.ActionURL == "" || p.ActionLabel == "" {
			t.Errorf("%s is missing its generated primary action", id)
		}
	}
	for _, p := range cat.More {
		if untaggedIDs[p.ID] {
			if p.Proof != "no public tag" || p.ProofURL != "" {
				t.Errorf("%s must preserve missing tag truth: %#v", p.ID, p)
			}
		}
	}
}

func TestCatalogCarriesReferencesInsteadOfProductCopy(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join(repoRoot(t), "catalog.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	if strings.Contains(text, "\nfeatured:\n") || !strings.Contains(text, "\nfeatured_ids:\n") {
		t.Fatalf("catalog must select products by id, not author a featured table:\n%s", text)
	}
	if !strings.Contains(text, "\nextra_products:\n") {
		t.Fatal("catalog must keep extra_products for extracts without passports")
	}
	if !strings.Contains(text, "\nmore_ids:\n") {
		t.Fatal("catalog must keep more_ids for supporting public source")
	}
}

func TestDocsPullerEvalWording(t *testing.T) {
	cat := loadTestCatalog(t)
	var docs Product
	for _, p := range cat.Featured {
		if p.ID == "product.docs-puller" {
			docs = p
		}
	}
	if docs.Metric.Value != "95.8% Hit@1 · 100% Hit@5" {
		t.Fatalf("docs-puller metric value = %q", docs.Metric.Value)
	}
	if !strings.Contains(strings.ToLower(docs.Metric.Label), "sample") {
		t.Fatalf("docs-puller metric must be labeled as a sample: %q", docs.Metric.Label)
	}
}

func TestJobKitSyntheticBoundary(t *testing.T) {
	cat := loadTestCatalog(t)
	for _, p := range cat.Featured {
		if p.ID == "product.jobkit" && !strings.Contains(strings.ToLower(p.Detail), "synthetic") {
			t.Fatalf("jobkit must declare the synthetic fixture boundary")
		}
	}
}

func TestREADMEIsACompactList(t *testing.T) {
	readme := renderREADME(loadTestCatalog(t))
	if strings.Contains(readme, "| Product |") || strings.Contains(readme, "| --- |") {
		t.Fatal("README must stay a compact list, not a markdown table")
	}
	if !strings.Contains(readme, "- **[") {
		t.Fatal("README must render selected work as a bullet list")
	}
	if strings.Contains(readme, "formal release") || strings.Contains(readme, "evidence gates open") {
		t.Fatal("README list must not repeat the evidence spreadsheet columns")
	}
	if strings.Contains(readme, "👋") || strings.Contains(readme, "🚀") {
		t.Fatal("README must not copy emoji catalog voice")
	}
}

func TestREADMEContainsProofAndOmitsForbidden(t *testing.T) {
	cat := loadTestCatalog(t)
	readme := renderREADME(cat)
	if !strings.Contains(readme, "## More on GitHub") {
		t.Fatal("README must keep a second list for supporting public source")
	}
	if strings.Contains(strings.ToLower(readme), "missing tags stay visible") || strings.Contains(strings.ToLower(readme), "stand behind") {
		t.Fatal("README still uses catalog-manifesto phrasing")
	}
	for _, p := range append(append([]Product{}, cat.Featured...), cat.More...) {
		needles := []string{p.Name, p.URL, p.License}
		if p.ProofURL != "" {
			needles = append(needles, p.ProofURL)
		} else if !strings.Contains(readme, "no public tag") {
			t.Errorf("%s is missing its missing-tag marker", p.ID)
		}
		for _, needle := range needles {
			if !strings.Contains(readme, needle) {
				t.Errorf("README missing %q", needle)
			}
		}
	}
	lower := strings.ToLower(readme)
	for _, bad := range forbidden {
		if strings.Contains(readme, bad) || strings.Contains(lower, strings.ToLower(bad)) {
			t.Errorf("README contains forbidden token %q", bad)
		}
	}
	if !strings.Contains(readme, "public live-page BM25 sample") {
		t.Error("README must keep the docs-puller sample wording")
	}
	if !strings.Contains(readme, "https://docs-puller-demo.nstranquist.workers.dev") {
		t.Error("README must link the docs-puller live demo")
	}
	if !strings.Contains(readme, "docs-puller/tree/v0.7.6") {
		t.Error("README must cite the live docs-puller tag")
	}
	if !strings.Contains(readme, "example data only") {
		t.Error("README must keep the JobKit synthetic-data boundary")
	}
	for _, needle := range []string{
		"SessionPressure", "keepawake", "SnapRef", "Nicos GTM",
		"nicos-flag-eval", "Nicos Window Switcher",
	} {
		if !strings.Contains(readme, needle) {
			t.Errorf("README missing %s", needle)
		}
	}
	if strings.Contains(readme, "agent-native") || strings.Contains(readme, "evidence-backed") {
		t.Error("README must use public one-liners, not passport hiring copy")
	}
	if !strings.Contains(readme, "https://nstranquist.github.io") {
		t.Error("README must point at the public catalog site")
	}
	if strings.Contains(readme, "Hi, I'm") || strings.Contains(readme, "👋") {
		t.Error("README should not use the generic profile greeting")
	}
	if strings.Contains(readme, "assets/cards/") || strings.Contains(readme, "assets/banner") || strings.Contains(readme, "assets/catalog") {
		t.Error("README should stay condensed prose, not embed banner, board, or card SVGs")
	}
	if strings.Contains(readme, "<img ") {
		t.Error("README should not embed images on the profile page")
	}
	if strings.Contains(readme, "## How I work") || strings.Contains(readme, "## Toolbox") {
		t.Error("README should stay condensed: no How I work or Toolbox section")
	}
	for _, phrase := range []string{
		"inspectable extract",
		"claim boundary",
		"this surface",
		"public extracts",
		"AI infrastructure",
	} {
		if strings.Contains(strings.ToLower(readme), phrase) {
			t.Errorf("README still uses factory phrasing %q", phrase)
		}
	}
}

func TestSVGsAreWellFormedAndThemed(t *testing.T) {
	cat := loadTestCatalog(t)
	dark := renderBanner(cat, themeDark())
	light := renderBanner(cat, themeLight())
	board := renderCatalogBoard(cat, themeDark())
	for name, svg := range map[string]string{"banner-dark": dark, "banner-light": light, "catalog": board} {
		if !strings.HasPrefix(svg, "<svg ") {
			t.Errorf("%s: missing svg open", name)
		}
		if !strings.HasSuffix(strings.TrimSpace(svg), "</svg>") {
			t.Errorf("%s: missing svg close", name)
		}
		if strings.Contains(svg, "<script") {
			t.Errorf("%s: script tags are not allowed", name)
		}
		if strings.Contains(svg, `font-family=""`) {
			t.Errorf("%s: broken font-family quotes", name)
		}
	}
	if dark == light {
		t.Fatal("dark and light banners should differ")
	}
	for _, p := range cat.Featured {
		if !strings.Contains(board, p.Name) || !strings.Contains(board, p.Repo) || !strings.Contains(board, p.Proof) {
			t.Errorf("catalog board missing %s", p.ID)
		}
	}
}

func TestRenderCheckRoundTrip(t *testing.T) {
	root := t.TempDir()
	src := repoRoot(t)
	raw, err := os.ReadFile(filepath.Join(src, "catalog.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "catalog.yaml"), raw, 0o644); err != nil {
		t.Fatal(err)
	}
	products, err := os.ReadFile(filepath.Join(src, "data", "products.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "data"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "data", "products.json"), products, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := run(root, false); err != nil {
		t.Fatalf("render: %v", err)
	}
	if err := run(root, true); err != nil {
		t.Fatalf("render --check after render: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "README.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "assets", "cards", "docs-puller-dark.svg")); err != nil {
		t.Fatal(err)
	}
}

func TestMarkIsAnN(t *testing.T) {
	svg := markN(0, 0, themeDark())
	// 4x4 N uses 10 filled cells. Count mark-on fill occurrences.
	if got := strings.Count(svg, themeDark().MarkOn); got != 10 {
		t.Fatalf("expected 10 filled N cells, got %d", got)
	}
}
