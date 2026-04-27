package staticlp

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

// RunEditInteractive loads a campaign and edits common fields with huh.
func RunEditInteractive(siteRoot, id string) error {
	c, path, err := LoadCampaignFile(siteRoot, id)
	if err != nil {
		return err
	}
	if c.Presenter == nil {
		c.Presenter = &Presenter{}
	}

	fileID := c.Delivery.FileID
	submit := c.SubmitLabel
	consent := c.ConsentText
	pName := c.Presenter.Name
	pHead := c.Presenter.Headline
	pBio := c.Presenter.Bio
	pPhoto := c.Presenter.Photo
	pLI := c.Presenter.LinkedinURL
	aBody := c.Analytics.BodyClass
	aCamp := c.Analytics.Campaign

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("delivery.file_id (Google Drive)").Value(&fileID),
			huh.NewInput().Title("submit_label").Value(&submit),
			huh.NewText().Title("consent_text").Value(&consent).Lines(5),
		),
		huh.NewGroup(
			huh.NewInput().Title("presenter.name").Value(&pName),
			huh.NewInput().Title("presenter.headline").Value(&pHead),
			huh.NewText().Title("presenter.bio").Value(&pBio).Lines(4),
			huh.NewInput().Title("presenter.photo (path or URL under static/)").Value(&pPhoto),
			huh.NewInput().Title("presenter.linkedin_url").Value(&pLI),
		),
		huh.NewGroup(
			huh.NewInput().Title("analytics.body_class").Value(&aBody),
			huh.NewInput().Title("analytics.campaign").Value(&aCamp),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		return err
	}

	c.Delivery.FileID = fileID
	c.SubmitLabel = submit
	c.ConsentText = consent
	c.Presenter.Name = pName
	c.Presenter.Headline = pHead
	c.Presenter.Bio = pBio
	c.Presenter.Photo = pPhoto
	c.Presenter.LinkedinURL = pLI
	c.Analytics.BodyClass = aBody
	c.Analytics.Campaign = aCamp

	if err := SaveCampaign(path, c); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	fmt.Println("Saved:", path)
	return nil
}
