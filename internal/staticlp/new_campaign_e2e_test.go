package staticlp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

// TestNewCampaign_WritesAllSections verifies scaffold + merge writes TOML and four markdown files (no TUI).
func TestNewCampaign_WritesAllSections(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "content"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "data", "campaigns"), 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := filepath.Join(root, "hugo.toml")
	if err := os.WriteFile(cfg, []byte("defaultContentLanguage = \"en\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	sniff, err := LoadHugoSniff(root)
	if err != nil {
		t.Fatal(err)
	}

	opts := NewOpts{
		SiteRoot:             root,
		Lang:                 "en",
		Slug:                 "e2e-slug",
		Title:                "E2E Title",
		Description:          "E2E Desc",
		FileID:               "drive-file-123",
		SubmitLabel:          "Get PDF",
		ConsentText:          "Consent line.",
		PresenterName:        "Pat",
		PresenterHeadline:    "Hi",
		PresenterBio:         "Bio here.",
		PresenterPhoto:       "p.jpg",
		PresenterLinkedinURL: "https://linkedin.com/in/x",
		AnalyticsBodyClass:   "analytics-campaign-e2e-slug",
		AnalyticsCampaign:    "e2e-slug-funnel",
	}
	if err := NewCampaign(sniff, opts); err != nil {
		t.Fatal(err)
	}

	tomlPath := filepath.Join(root, "data", "campaigns", "e2e-slug.toml")
	raw, err := os.ReadFile(tomlPath)
	if err != nil {
		t.Fatal(err)
	}
	var c Campaign
	if err := toml.Unmarshal(raw, &c); err != nil {
		t.Fatal(err)
	}
	if c.Delivery.FileID != "drive-file-123" {
		t.Fatalf("file_id: %q", c.Delivery.FileID)
	}
	if c.Presenter == nil || c.Presenter.Name != "Pat" {
		t.Fatalf("presenter: %+v", c.Presenter)
	}
	if c.Presenter.LinkedinURL != "https://linkedin.com/in/x" {
		t.Fatalf("linkedin: %q", c.Presenter.LinkedinURL)
	}
	if c.Presenter.Photo != "p.jpg" {
		t.Fatalf("photo: %q", c.Presenter.Photo)
	}

	for _, name := range []string{"_index.md", "obrigado.md", "material.md", "qr.md"} {
		p := filepath.Join(root, "content", "e2e-slug", name)
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), "campaign: e2e-slug") {
			t.Fatalf("%s missing campaign front matter", name)
		}
	}
}
