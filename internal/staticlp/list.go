package staticlp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// CampaignListRow is one line for list output.
type CampaignListRow struct {
	ID       string
	FormName string
	Landing  string
	Path     string
}

// ListCampaigns reads data/campaigns/*.toml.
func ListCampaigns(siteRoot string) ([]CampaignListRow, error) {
	dir := filepath.Join(siteRoot, "data", "campaigns")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var rows []CampaignListRow
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".toml") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var c Campaign
		if err := toml.Unmarshal(raw, &c); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		rows = append(rows, CampaignListRow{
			ID:       c.ID,
			FormName: c.FormName,
			Landing:  c.Paths.Landing,
			Path:     path,
		})
	}
	return rows, nil
}

// FindCampaignPath returns the filesystem path to data/campaigns/<id>.toml or fs.ErrNotExist.
func FindCampaignPath(siteRoot, id string) (string, error) {
	path := filepath.Join(siteRoot, "data", "campaigns", id+".toml")
	if _, err := os.Stat(path); err != nil {
		return "", err
	}
	return path, nil
}

// LoadCampaignFile loads a campaign TOML by id.
func LoadCampaignFile(siteRoot, id string) (*Campaign, string, error) {
	path, err := FindCampaignPath(siteRoot, id)
	if err != nil {
		return nil, "", err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	var c Campaign
	if err := toml.Unmarshal(raw, &c); err != nil {
		return nil, path, err
	}
	return &c, path, nil
}

// SaveCampaign writes campaign TOML (UTF-8, no BOM).
func SaveCampaign(path string, c *Campaign) error {
	b, err := toml.Marshal(c)
	if err != nil {
		return err
	}
	return WriteFileUTF8(path, string(b), 0o644)
}
