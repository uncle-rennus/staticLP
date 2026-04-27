package staticlp

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// InferLangSlugFromLanding parses paths.landing into language code and section slug path.
func InferLangSlugFromLanding(sniff *HugoSniff, landing string) (lang, slug string, err error) {
	p := strings.Trim(landing, "/")
	if p == "" {
		return "", "", fmt.Errorf("empty landing path")
	}
	var segs []string
	for _, x := range strings.Split(p, "/") {
		if x != "" {
			segs = append(segs, x)
		}
	}
	if len(segs) == 0 {
		return "", "", fmt.Errorf("empty landing path")
	}
	if sniff.DefaultContentLanguageInSubdir {
		if len(segs) < 2 {
			return "", "", fmt.Errorf("landing %q: expected /lang/slug/ for this site (defaultContentLanguageInSubdir)", landing)
		}
		return segs[0], strings.Join(segs[1:], "/"), nil
	}
	return sniff.DefaultContentLanguage, strings.Join(segs, "/"), nil
}

// RemoveCampaign deletes data/campaigns/<id>.toml and the content section directory.
func RemoveCampaign(siteRoot string, sniff *HugoSniff, id string, slugOverride, langOverride string, force bool) error {
	c, dataPath, err := LoadCampaignFile(siteRoot, id)
	if err != nil {
		return err
	}
	lang := strings.TrimSpace(langOverride)
	slug := strings.TrimSpace(slugOverride)
	if slug == "" {
		var inferErr error
		lang, slug, inferErr = InferLangSlugFromLanding(sniff, c.Paths.Landing)
		if inferErr != nil {
			return fmt.Errorf("infer slug from landing: %w (pass --slug and --lang)", inferErr)
		}
	} else if lang == "" {
		lang = sniff.DefaultContentLanguage
	}

	contentDirRel, err := sniff.ContentDirForLang(lang)
	if err != nil {
		return err
	}
	sectionDir := filepath.Join(siteRoot, filepath.FromSlash(contentDirRel), filepath.FromSlash(slug))

	if !force {
		fmt.Fprintf(os.Stderr, "Delete campaign %q?\n", id)
		fmt.Fprintf(os.Stderr, "  %s\n", dataPath)
		fmt.Fprintf(os.Stderr, "  %s\n", sectionDir)
		fmt.Fprint(os.Stderr, "Type YES to confirm: ")
		sc := bufio.NewScanner(os.Stdin)
		if !sc.Scan() {
			return fmt.Errorf("aborted")
		}
		if strings.TrimSpace(sc.Text()) != "YES" {
			return fmt.Errorf("aborted")
		}
	}

	if err := os.Remove(dataPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.RemoveAll(sectionDir); err != nil {
		return err
	}
	return nil
}
