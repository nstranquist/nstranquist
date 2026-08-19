package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Catalog struct {
	SchemaVersion int                   `yaml:"schema_version"`
	Identity      Identity              `yaml:"identity"`
	FeaturedIDs   []string              `yaml:"featured_ids"`
	Featured      []Product             `yaml:"-"`
	PublicCopy    map[string]PublicCopy `yaml:"public_copy"`
	Toolbox       []string              `yaml:"toolbox"`
	Principles    []Principle           `yaml:"principles"`
	Footnote      string                `yaml:"footnote"`
	Glossary      []GlossaryEntry       `yaml:"glossary"`
}

type PublicCopy struct {
	Name    string `yaml:"name"`
	Summary string `yaml:"summary"`
}

type Identity struct {
	Name     string `yaml:"name"`
	Role     string `yaml:"role"`
	Location string `yaml:"location"`
	Thesis   string `yaml:"thesis"`
	Intro    string `yaml:"intro"`
	GitHub   string `yaml:"github"`
	LinkedIn string `yaml:"linkedin"`
	Site     string `yaml:"site"`
}

type GlossaryEntry struct {
	Term    string `yaml:"term"`
	Meaning string `yaml:"meaning"`
}

type Product struct {
	ID          string
	Name        string
	Repo        string
	URL         string
	ProofURL    string
	License     string
	Language    string
	Lane        string
	Proof       string
	Release     string
	ReleaseURL  string
	ActionLabel string
	ActionURL   string
	DemoURL     string
	PublicState string
	Ready       bool
	Summary     string
	Detail      string
	Metric      Metric
}

type Metric struct {
	Label string `yaml:"label"`
	Value string `yaml:"value"`
}

type Principle struct {
	Title string `yaml:"title"`
	Body  string `yaml:"body"`
}

type productSnapshot struct {
	SchemaVersion int                     `json:"schema_version"`
	Products      []productSnapshotRecord `json:"products"`
}

type productSnapshotRecord struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Version      string `json:"version"`
	VersionURL   string `json:"version_url"`
	License      string `json:"license"`
	Ready        bool   `json:"ready"`
	Presentation struct {
		Summary  string `json:"summary"`
		Detail   string `json:"detail"`
		Language string `json:"language"`
		Lane     string `json:"lane"`
		Metric   *struct {
			Label string `json:"label"`
			Value string `json:"value"`
		} `json:"metric"`
	} `json:"presentation"`
	Access struct {
		RepositoryURL string `json:"repository_url"`
		DemoURL       string `json:"demo_url"`
		PrimaryAction struct {
			Label string `json:"label"`
			URL   string `json:"url"`
		} `json:"primary_action"`
	} `json:"access"`
	Lifecycle struct {
		AccessVerified struct {
			State string `json:"state"`
		} `json:"access_verified"`
	} `json:"lifecycle"`
	GitHub *struct {
		RepositoryPublic bool   `json:"repository_public"`
		LatestTag        string `json:"latest_tag"`
		LatestTagURL     string `json:"latest_tag_url"`
		LatestReleaseTag string `json:"latest_release_tag"`
		LatestReleaseURL string `json:"latest_release_url"`
		CIState          string `json:"ci_state"`
	} `json:"github"`
}

func loadCatalog(catalogPath, productsPath string) (Catalog, error) {
	raw, err := os.ReadFile(catalogPath)
	if err != nil {
		return Catalog{}, err
	}
	var cat Catalog
	decoder := yaml.NewDecoder(bytes.NewReader(raw))
	decoder.KnownFields(true)
	if err := decoder.Decode(&cat); err != nil {
		return Catalog{}, err
	}
	productsRaw, err := os.ReadFile(productsPath)
	if err != nil {
		return Catalog{}, err
	}
	var snapshot productSnapshot
	if err := json.Unmarshal(productsRaw, &snapshot); err != nil {
		return Catalog{}, fmt.Errorf("parse Product Passport snapshot: %w", err)
	}
	if snapshot.SchemaVersion != 1 {
		return Catalog{}, fmt.Errorf("unsupported Product Passport schema_version %d", snapshot.SchemaVersion)
	}
	byID := make(map[string]productSnapshotRecord, len(snapshot.Products))
	for _, product := range snapshot.Products {
		if _, exists := byID[product.ID]; exists {
			return Catalog{}, fmt.Errorf("duplicate Product Passport id %s", product.ID)
		}
		byID[product.ID] = product
	}
	for _, id := range cat.FeaturedIDs {
		product, ok := byID[id]
		if !ok {
			return Catalog{}, fmt.Errorf("featured id %s is missing from data/products.json", id)
		}
		mapped, err := mapProductPassport(product)
		if err != nil {
			return Catalog{}, err
		}
		if copy, ok := cat.PublicCopy[id]; ok {
			if copy.Name != "" {
				mapped.Name = copy.Name
			}
			if copy.Summary != "" {
				mapped.Summary = copy.Summary
			}
		}
		cat.Featured = append(cat.Featured, mapped)
	}
	if err := cat.validate(); err != nil {
		return Catalog{}, err
	}
	return cat, nil
}

func mapProductPassport(source productSnapshotRecord) (Product, error) {
	if source.GitHub == nil || !source.GitHub.RepositoryPublic {
		return Product{}, fmt.Errorf("%s is not verified as a public GitHub repository", source.ID)
	}
	if source.GitHub.CIState == "" {
		return Product{}, fmt.Errorf("%s public CI is unmeasured", source.ID)
	}
	if source.Lifecycle.AccessVerified.State != "verified" {
		return Product{}, fmt.Errorf("%s access is %q, expected verified", source.ID, source.Lifecycle.AccessVerified.State)
	}
	if source.Version != source.GitHub.LatestTag || source.VersionURL != source.GitHub.LatestTagURL {
		return Product{}, fmt.Errorf("%s generated version does not match its GitHub observation", source.ID)
	}
	repo := strings.TrimPrefix(source.Access.RepositoryURL, "https://github.com/")
	if !strings.HasPrefix(repo, "nstranquist/") || strings.Contains(strings.TrimPrefix(repo, "nstranquist/"), "/") {
		return Product{}, fmt.Errorf("%s repository is outside nstranquist", source.ID)
	}
	proof := source.Version
	if proof == "" {
		proof = "no public tag"
	}
	state := "source public"
	if source.Version != "" {
		state += " · tag " + source.Version
	} else {
		state += " · no tag"
	}
	if source.GitHub.LatestReleaseTag != "" {
		state += " · release " + source.GitHub.LatestReleaseTag
	} else {
		state += " · no formal release"
	}
	if source.GitHub.CIState != "green" {
		state += " · CI " + source.GitHub.CIState
	}
	if !source.Ready {
		state += " · evidence gates open"
	}
	product := Product{
		ID: source.ID, Name: source.Name, Repo: repo, URL: source.Access.RepositoryURL,
		ProofURL: source.VersionURL, License: source.License,
		Language: source.Presentation.Language, Lane: source.Presentation.Lane,
		Proof: proof, Release: source.GitHub.LatestReleaseTag, ReleaseURL: source.GitHub.LatestReleaseURL,
		ActionLabel: source.Access.PrimaryAction.Label, ActionURL: source.Access.PrimaryAction.URL,
		DemoURL: source.Access.DemoURL, PublicState: state, Ready: source.Ready,
		Summary: source.Presentation.Summary, Detail: source.Presentation.Detail,
	}
	if source.Presentation.Metric != nil {
		product.Metric = Metric{Label: source.Presentation.Metric.Label, Value: source.Presentation.Metric.Value}
	}
	return product, nil
}

func (c Catalog) validate() error {
	if c.SchemaVersion != 2 {
		return fmt.Errorf("unsupported schema_version %d", c.SchemaVersion)
	}
	if strings.TrimSpace(c.Identity.Name) == "" {
		return fmt.Errorf("identity.name is required")
	}
	if strings.TrimSpace(c.Identity.Intro) == "" {
		return fmt.Errorf("identity.intro is required")
	}
	if len(c.FeaturedIDs) == 0 || len(c.Featured) != len(c.FeaturedIDs) {
		return fmt.Errorf("featured catalog is empty")
	}
	if len(c.Glossary) == 0 {
		return fmt.Errorf("glossary is required")
	}
	seen := map[string]struct{}{}
	for i, p := range c.Featured {
		if p.ID == "" || p.Name == "" || p.Repo == "" || p.URL == "" || p.ActionURL == "" || p.PublicState == "" {
			return fmt.Errorf("featured[%d] is missing required fields", i)
		}
		if !strings.HasPrefix(p.ID, "product.") {
			return fmt.Errorf("featured[%d].id %q must start with product.", i, p.ID)
		}
		if !strings.HasPrefix(p.Repo, "nstranquist/") {
			return fmt.Errorf("featured[%d].repo %q must be under nstranquist/", i, p.Repo)
		}
		if !strings.HasPrefix(p.URL, "https://github.com/nstranquist/") {
			return fmt.Errorf("featured[%d].url %q is not a public nstranquist GitHub URL", i, p.URL)
		}
		if p.ProofURL != "" && !strings.HasPrefix(p.ProofURL, p.URL+"/tree/") {
			return fmt.Errorf("featured[%d].proof_url must be a source tag tree on %s", i, p.URL)
		}
		if p.ProofURL == "" && p.Proof != "no public tag" {
			return fmt.Errorf("featured[%d] has inconsistent tag evidence", i)
		}
		if !strings.HasPrefix(p.ActionURL, "https://") {
			return fmt.Errorf("featured[%d].action_url must be an https URL", i)
		}
		if strings.Contains(p.ActionURL, "github.com/") && !strings.HasPrefix(p.ActionURL, p.URL) {
			return fmt.Errorf("featured[%d].action_url github target must stay on %s", i, p.URL)
		}
		if p.DemoURL != "" && !strings.HasPrefix(p.DemoURL, "https://") {
			return fmt.Errorf("featured[%d].demo_url must be an https URL", i)
		}
		publicText := strings.ToLower(strings.Join([]string{p.ID, p.Name, p.Summary, p.Detail}, " "))
		if strings.Contains(publicText, "edurain") {
			return fmt.Errorf("featured[%d] contains a permanently excluded product marker", i)
		}
		if _, ok := seen[p.ID]; ok {
			return fmt.Errorf("duplicate featured id %s", p.ID)
		}
		seen[p.ID] = struct{}{}
	}
	return nil
}
