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
	"product.jobkit",
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
	cat, err := loadCatalog(filepath.Join(repoRoot(t), "catalog.yaml"))
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
	for i, id := range requiredIDs {
		p := cat.Featured[i]
		if p.ID != id {
			t.Errorf("featured[%d]=%s, want %s", i, p.ID, id)
		}
		if !strings.HasPrefix(p.Proof, "v") {
			t.Errorf("%s proof %q should be a version tag", id, p.Proof)
		}
		if !strings.Contains(p.ProofURL, "/releases/tag/"+p.Proof) {
			t.Errorf("%s proof_url %s does not match proof %s", id, p.ProofURL, p.Proof)
		}
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
	if !strings.Contains(docs.Detail, "public BM25 sample") {
		t.Fatalf("docs-puller detail must name the public BM25 sample")
	}
}

func TestJobKitSyntheticBoundary(t *testing.T) {
	cat := loadTestCatalog(t)
	for _, p := range cat.Featured {
		if p.ID == "product.jobkit" && !strings.Contains(p.Detail, "synthetic-fixture") {
			t.Fatalf("jobkit must declare the synthetic-fixture boundary")
		}
	}
}

func TestREADMEContainsProofAndOmitsForbidden(t *testing.T) {
	cat := loadTestCatalog(t)
	readme := renderREADME(cat)
	for _, p := range cat.Featured {
		for _, needle := range []string{p.Name, p.URL, p.ProofURL, p.License} {
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
	if !strings.Contains(readme, "public BM25 sample") {
		t.Error("README must keep the docs-puller sample wording")
	}
	if strings.Contains(readme, "Hi, I'm") || strings.Contains(readme, "👋") {
		t.Error("README should not use the generic profile greeting")
	}
	if strings.Contains(readme, "assets/cards/") {
		t.Error("README should stay a catalog + prose surface, not a card dump")
	}
	if !strings.Contains(readme, "## Terms") || !strings.Contains(readme, "Fail-closed") {
		t.Error("README must include a glossary for remaining terms")
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
