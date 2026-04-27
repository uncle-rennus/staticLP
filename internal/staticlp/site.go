package staticlp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// HugoSniff holds minimal Hugo config needed for paths and content dirs.
type HugoSniff struct {
	ConfigPath                     string
	DefaultContentLanguage         string
	DefaultContentLanguageInSubdir bool
	Languages                      map[string]languageEntry
}

type languageEntry struct {
	ContentDir string `toml:"contentDir"`
}

type hugoConfigFile struct {
	DefaultContentLanguage         string                   `toml:"defaultContentLanguage"`
	DefaultContentLanguageInSubdir bool                     `toml:"defaultContentLanguageInSubdir"`
	Languages                      map[string]languageEntry `toml:"languages"`
}

// FindSiteRoot returns dir containing hugo.toml or config.toml and a content/ directory.
func FindSiteRoot(root string) (string, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	if err := validateSiteRoot(abs); err == nil {
		return abs, nil
	}
	return "", fmt.Errorf("not a Hugo site root (need hugo.toml or config.toml and content/): %s", abs)
}

func validateSiteRoot(dir string) error {
	st, err := os.Stat(filepath.Join(dir, "content"))
	if err != nil || !st.IsDir() {
		return fmt.Errorf("missing content/")
	}
	_, hugo := os.Stat(filepath.Join(dir, "hugo.toml"))
	_, cfg := os.Stat(filepath.Join(dir, "config.toml"))
	if hugo == nil || cfg == nil {
		return nil
	}
	return fmt.Errorf("missing hugo.toml and config.toml")
}

// LoadHugoSniff reads hugo.toml or config.toml from site root.
func LoadHugoSniff(siteRoot string) (*HugoSniff, error) {
	var path string
	for _, name := range []string{"hugo.toml", "config.toml"} {
		p := filepath.Join(siteRoot, name)
		if _, err := os.Stat(p); err == nil {
			path = p
			break
		}
	}
	if path == "" {
		return nil, fmt.Errorf("no hugo.toml or config.toml in %s", siteRoot)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg hugoConfigFile
	if err := toml.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if cfg.DefaultContentLanguage == "" {
		cfg.DefaultContentLanguage = "en"
	}
	return &HugoSniff{
		ConfigPath:                     path,
		DefaultContentLanguage:         cfg.DefaultContentLanguage,
		DefaultContentLanguageInSubdir: cfg.DefaultContentLanguageInSubdir,
		Languages:                      cfg.Languages,
	}, nil
}

// ContentDirForLang returns the Hugo contentDir for a language code (e.g. "pt" -> "content/pt").
// When there is no [languages] section, returns "content".
func (s *HugoSniff) ContentDirForLang(lang string) (string, error) {
	lang = strings.TrimSpace(lang)
	if lang == "" {
		lang = s.DefaultContentLanguage
	}
	if len(s.Languages) == 0 {
		return "content", nil
	}
	entry, ok := s.Languages[lang]
	if !ok {
		return "", fmt.Errorf("unknown language %q (not in [languages])", lang)
	}
	if entry.ContentDir == "" {
		return "", fmt.Errorf("language %q has no contentDir", lang)
	}
	return filepath.ToSlash(strings.TrimSpace(entry.ContentDir)), nil
}

// LangPlaceholderForTemplate returns the segment used in campaign URLs for __LANG__ in the template.
func (s *HugoSniff) LangPlaceholderForTemplate(lang string) string {
	if s.DefaultContentLanguageInSubdir {
		return lang
	}
	return ""
}

// ApplyCampaignTemplate substitutes placeholders in the scaffold TOML template.
func ApplyCampaignTemplate(tpl, slug, campaignID, formName, langSeg string) string {
	var out string
	if langSeg == "" {
		out = strings.ReplaceAll(tpl, "/__LANG__/", "/")
	} else {
		out = strings.ReplaceAll(tpl, "__LANG__", langSeg)
	}
	out = strings.ReplaceAll(out, "__SLUG__", slug)
	out = strings.ReplaceAll(out, "__CAMPAIGN_ID__", campaignID)
	out = strings.ReplaceAll(out, "__FORM_NAME__", formName)
	return out
}
