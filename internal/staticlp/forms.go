package staticlp

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
)

// NewCampaignFormFields binds all wizard inputs (used by TUI and tests).
type NewCampaignFormFields struct {
	Slug, Lang, Title, Description string
	Draft                          bool
	FileID, SubmitLabel, ConsentText string
	PresenterName, PresenterHeadline, PresenterBio, PresenterPhoto, PresenterLinkedinURL string
	AnalyticsBodyClass, AnalyticsCampaign string
}

// BuildMainMenuForm returns the root staticLP menu (caller runs .Run() or drives .Update in tests).
func BuildMainMenuForm(siteRoot string, action *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("staticLP").
				Description("Hugo Netlify campaign landings at "+siteRoot).
				Height(8).
				Value(action).
				Options(
					huh.NewOption("List campaigns", "list"),
					huh.NewOption("New campaign", "new"),
					huh.NewOption("Edit campaign", "edit"),
					huh.NewOption("Remove campaign", "rm"),
					huh.NewOption("Quit", "quit"),
				),
		),
	).WithTheme(huh.ThemeCharm())
}

// BuildNewCampaignWizardForm returns the 4-step new-campaign wizard.
func BuildNewCampaignWizardForm(fields *NewCampaignFormFields) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("URL slug (folder name)").Value(&fields.Slug).Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("required")
				}
				return nil
			}),
			huh.NewInput().Title("Language code").Description("Must match [languages] key, e.g. pt").Value(&fields.Lang),
			huh.NewInput().Title("Landing page title (H1)").Value(&fields.Title),
			huh.NewInput().Title("Meta description (SEO)").Value(&fields.Description),
			huh.NewConfirm().Title("Mark pages as draft?").Value(&fields.Draft),
		).Title("New campaign 1/4 — Basics").
			Description("Tab / Shift+Tab: move fields. Enter: next field. After the last field here, the next step opens."),

		huh.NewGroup(
			huh.NewInput().Title("Google Drive file_id (PDF)").Description("Empty = keep TOML placeholder (SUBSTITUA…).").Value(&fields.FileID),
			huh.NewText().Title("Consent text (GDPR)").Value(&fields.ConsentText).Lines(4),
			huh.NewInput().Title("Submit button label").Value(&fields.SubmitLabel),
		).Title("New campaign 2/4 — Form & PDF").
			Description("In the consent box, Enter adds a line; use Tab to move to the submit label, then Enter to go to step 3."),

		huh.NewGroup(
			huh.NewInput().Title("Presenter name").Value(&fields.PresenterName),
			huh.NewInput().Title("Headline (short)").Value(&fields.PresenterHeadline),
			huh.NewText().Title("Presenter bio").Value(&fields.PresenterBio).Lines(4),
			huh.NewInput().Title("Profile image").Description("e.g. profile.jpg under static/ or a full URL").Value(&fields.PresenterPhoto),
			huh.NewInput().Title("LinkedIn profile URL").Value(&fields.PresenterLinkedinURL),
		).Title("New campaign 3/4 — Presenter").
			Description("In bio, Enter = new line; Tab to the next field. Last field: Enter → step 4."),

		huh.NewGroup(
			huh.NewInput().Title("analytics.body_class (CSS)").Description("Empty = slug default from scaffold.").Value(&fields.AnalyticsBodyClass),
			huh.NewInput().Title("analytics.campaign (label)").Value(&fields.AnalyticsCampaign),
		).Title("New campaign 4/4 — Analytics"),
	).WithTheme(huh.ThemeCharm())
}
