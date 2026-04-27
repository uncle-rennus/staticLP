package staticlp

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
)

// selectCancel is a value that never collides with a real campaign id.
const selectCancel = "__staticlp_cancel__"

// RunTUI shows the main menu and dispatches to list/new/edit/rm.
func RunTUI(siteRoot string) error {
	sniff, err := LoadHugoSniff(siteRoot)
	if err != nil {
		return err
	}

	for {
		// Default "list" must be set before Options() runs (see pickCampaignID / huh Select).
		action := "list"
		form := BuildMainMenuForm(siteRoot, &action)

		if err := form.Run(); err != nil {
			return err
		}

		switch action {
		case "list":
			if err := RunList(os.Stdout, siteRoot); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			WaitEnterForMenu()
		case "new":
			if err := runTUINew(siteRoot, sniff); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			WaitEnterForMenu()
		case "edit":
			if err := runTUIEdit(siteRoot); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			WaitEnterForMenu()
		case "rm":
			if err := runTUIRemove(siteRoot, sniff); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			WaitEnterForMenu()
		case "quit":
			return nil
		}
	}
}

func runTUINew(siteRoot string, sniff *HugoSniff) error {
	fields := &NewCampaignFormFields{Lang: sniff.DefaultContentLanguage}
	form := BuildNewCampaignWizardForm(fields)

	if err := form.Run(); err != nil {
		return err
	}
	title := fields.Title
	desc := fields.Description
	if strings.TrimSpace(title) == "" {
		title = "Event landing"
	}
	if strings.TrimSpace(desc) == "" {
		desc = "Enter your details to access the material."
	}
	opts := NewOpts{
		SiteRoot:             siteRoot,
		Lang:                 fields.Lang,
		Slug:                 fields.Slug,
		CampaignID:           "",
		FormName:             "",
		Title:                title,
		Description:          desc,
		Draft:                fields.Draft,
		FileID:               fields.FileID,
		SubmitLabel:          fields.SubmitLabel,
		ConsentText:          fields.ConsentText,
		PresenterName:        fields.PresenterName,
		PresenterHeadline:    fields.PresenterHeadline,
		PresenterBio:         fields.PresenterBio,
		PresenterPhoto:       fields.PresenterPhoto,
		PresenterLinkedinURL: fields.PresenterLinkedinURL,
		AnalyticsBodyClass:   fields.AnalyticsBodyClass,
		AnalyticsCampaign:    fields.AnalyticsCampaign,
	}
	if err := NewCampaign(sniff, opts); err != nil {
		return err
	}
	fmt.Println("Created campaign. Default id = slug. Adjust data/campaigns/<id>.toml if needed before deploy.")
	return nil
}

func runTUIEdit(siteRoot string) error {
	id, err := pickCampaignID(siteRoot, "Campaign to edit")
	if err != nil {
		return err
	}
	if id == "" {
		return nil
	}
	return RunEditInteractive(siteRoot, id)
}

func runTUIRemove(siteRoot string, sniff *HugoSniff) error {
	id, err := pickCampaignID(siteRoot, "Campaign to remove")
	if err != nil {
		return err
	}
	if id == "" {
		return nil
	}
	var ok bool
	confirm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Delete " + id + " from disk?").
				Description("Removes data/campaigns/" + id + ".toml and the content section folder.").
				Value(&ok),
		),
	).WithTheme(huh.ThemeCharm())
	if err := confirm.Run(); err != nil {
		return err
	}
	if !ok {
		return nil
	}
	if err := RemoveCampaign(siteRoot, sniff, id, "", "", true); err != nil {
		return err
	}
	fmt.Println("Removed:", id)
	return nil
}

// PickCampaignIDForCLI prompts for a campaign id (huh select).
func PickCampaignIDForCLI(siteRoot, title string) (string, error) {
	return pickCampaignID(siteRoot, title)
}

func pickCampaignID(siteRoot, title string) (string, error) {
	rows, err := ListCampaigns(siteRoot)
	if err != nil {
		return "", err
	}
	if len(rows) == 0 {
		fmt.Println("No campaigns in data/campaigns/")
		return "", nil
	}
	opts := make([]huh.Option[string], 0, len(rows)+1)
	for _, r := range rows {
		label := r.ID
		if r.Landing != "" {
			label += " — " + r.Landing
		}
		opts = append(opts, huh.NewOption(label, r.ID))
	}
	opts = append(opts, huh.NewOption("— Cancel —", selectCancel))

	// Value before Options: Options() initial sync uses the accessor; "" matched Cancel and left the viewport on the last row.
	id := rows[0].ID
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(title).
				Height(12).
				Value(&id).
				Options(opts...),
		),
	).WithTheme(huh.ThemeCharm())
	if err := form.Run(); err != nil {
		return "", err
	}
	if id == selectCancel {
		return "", nil
	}
	return id, nil
}
