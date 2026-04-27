package staticlp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/uncle-rennus/staticlp/internal/staticlp/scaffold"
)

// NewOpts configures scaffold of a campaign. Non-empty TOML fields override the embed template.
type NewOpts struct {
	SiteRoot    string
	Lang        string
	Slug        string
	CampaignID  string
	FormName    string
	Title       string
	Description string
	Draft       bool

	FileID      string
	SubmitLabel string
	ConsentText string

	PresenterName        string
	PresenterHeadline    string
	PresenterBio         string
	PresenterPhoto       string
	PresenterLinkedinURL string

	AnalyticsBodyClass string
	AnalyticsCampaign  string
}

// NewCampaign writes data/campaigns/<id>.toml and content files.
func NewCampaign(sniff *HugoSniff, opts NewOpts) error {
	slug := strings.Trim(opts.Slug, "/")
	if slug == "" {
		return fmt.Errorf("slug is required")
	}
	lang := strings.TrimSpace(opts.Lang)
	if lang == "" {
		lang = sniff.DefaultContentLanguage
	}
	campaignID := strings.TrimSpace(opts.CampaignID)
	if campaignID == "" {
		campaignID = slug
	}
	formName := strings.TrimSpace(opts.FormName)
	if formName == "" {
		formName = campaignID + "-leads"
	}

	dataPath := filepath.Join(opts.SiteRoot, "data", "campaigns", campaignID+".toml")
	if _, err := os.Stat(dataPath); err == nil {
		return fmt.Errorf("campaign file already exists: %s", dataPath)
	}

	contentDirRel, err := sniff.ContentDirForLang(lang)
	if err != nil {
		return err
	}
	sectionDir := filepath.Join(opts.SiteRoot, filepath.FromSlash(contentDirRel), slug)
	if err := os.MkdirAll(sectionDir, 0o755); err != nil {
		return err
	}

	langSeg := sniff.LangPlaceholderForTemplate(lang)
	tomlOut := ApplyCampaignTemplate(scaffold.CampaignTomlTpl, slug, campaignID, formName, langSeg)
	var c Campaign
	if err := toml.Unmarshal([]byte(tomlOut), &c); err != nil {
		return fmt.Errorf("parse scaffold toml: %w", err)
	}
	mergeNewCampaignOpts(&c, &opts)
	if err := SaveCampaign(dataPath, &c); err != nil {
		return err
	}

	commonFM := `ShowBreadCrumbs: false
hideMeta: true
ShowPostNavLinks: false
campaign: ` + campaignID
	if opts.Draft {
		commonFM += "\ndraft: true"
	}

	indexMd := "---\ntitle: \"" + escapeYAMLDouble(opts.Title) + "\"\ndescription: \"" + escapeYAMLDouble(opts.Description) + "\"\n" + commonFM + "\n---\n"
	if err := WriteFileUTF8(filepath.Join(sectionDir, "_index.md"), indexMd, 0o644); err != nil {
		return err
	}
	thanksMd := "---\ntitle: \"Thank you\"\nslug: thank-you\n" + commonFM + "\n---\n"
	if err := WriteFileUTF8(filepath.Join(sectionDir, "thank-you.md"), thanksMd, 0o644); err != nil {
		return err
	}
	materialMd := "---\ntitle: \"Material\"\nslug: material\n" + commonFM + "\n---\n"
	if err := WriteFileUTF8(filepath.Join(sectionDir, "material.md"), materialMd, 0o644); err != nil {
		return err
	}
	qrMd := "---\ntitle: \"Redirecting\"\nslug: qr\n" + commonFM + "\n_build:\n  list: never\n---\n"
	if err := WriteFileUTF8(filepath.Join(sectionDir, "qr.md"), qrMd, 0o644); err != nil {
		return err
	}

	return nil
}

// mergeNewCampaignOpts applies user-provided NewOpts over the unmarshaled template. Empty string = keep template value.
func mergeNewCampaignOpts(c *Campaign, opts *NewOpts) {
	if s := strings.TrimSpace(opts.FileID); s != "" {
		c.Delivery.FileID = s
	}
	if s := strings.TrimSpace(opts.SubmitLabel); s != "" {
		c.SubmitLabel = s
	}
	if s := strings.TrimSpace(opts.ConsentText); s != "" {
		c.ConsentText = s
	}
	if s := strings.TrimSpace(opts.AnalyticsBodyClass); s != "" {
		c.Analytics.BodyClass = s
	}
	if s := strings.TrimSpace(opts.AnalyticsCampaign); s != "" {
		c.Analytics.Campaign = s
	}
	if c.Presenter == nil {
		c.Presenter = &Presenter{}
	}
	if s := strings.TrimSpace(opts.PresenterName); s != "" {
		c.Presenter.Name = s
	}
	if s := strings.TrimSpace(opts.PresenterHeadline); s != "" {
		c.Presenter.Headline = s
	}
	if s := strings.TrimSpace(opts.PresenterBio); s != "" {
		c.Presenter.Bio = s
	}
	if s := strings.TrimSpace(opts.PresenterPhoto); s != "" {
		c.Presenter.Photo = s
	}
	if s := strings.TrimSpace(opts.PresenterLinkedinURL); s != "" {
		c.Presenter.LinkedinURL = s
	}
}

func escapeYAMLDouble(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}
