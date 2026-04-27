package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/uncle-rennus/staticlp/internal/staticlp"
)

var flagRoot string

func siteRootFromFlag() string {
	if flagRoot != "" {
		return flagRoot
	}
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func main() {
	rootCmd := &cobra.Command{
		Use:           "staticlp",
		Short:         "staticLP — create and manage Hugo Netlify campaign landings",
		Long:          "Companion CLI for the staticLP Hugo module.\nRun with no arguments to open the TUI; use subcommands for scripts.",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return cmd.Help()
			}
			siteRoot, err := staticlp.FindSiteRoot(siteRootFromFlag())
			if err != nil {
				return err
			}
			return staticlp.RunTUI(siteRoot)
		},
	}
	rootCmd.PersistentFlags().StringVar(&flagRoot, "root", "", "Hugo site root (default: current working directory)")

	var (
		newLang, newSlug, newID, newForm, newTitle, newDesc string
		newDraft                                            bool
		newFileID, newSubmit, newConsent                    string
		newPName, newPHead, newPBio, newPPhoto, newPLI      string
		newABody, newACamp                                  string
	)
	newCmd := &cobra.Command{
		Use:   "new",
		Short: "Scaffold data/campaigns/<id>.toml and content section",
		RunE: func(cmd *cobra.Command, args []string) error {
			siteRoot, err := staticlp.FindSiteRoot(siteRootFromFlag())
			if err != nil {
				return err
			}
			sniff, err := staticlp.LoadHugoSniff(siteRoot)
			if err != nil {
				return err
			}
			if newSlug == "" {
				return fmt.Errorf("new requires --slug")
			}
			if newTitle == "" {
				newTitle = "Event landing"
			}
			if newDesc == "" {
				newDesc = "Enter your details to access the material."
			}
			return staticlp.NewCampaign(sniff, staticlp.NewOpts{
				SiteRoot:             siteRoot,
				Lang:                 newLang,
				Slug:                 newSlug,
				CampaignID:           newID,
				FormName:             newForm,
				Title:                newTitle,
				Description:          newDesc,
				Draft:                newDraft,
				FileID:               newFileID,
				SubmitLabel:          newSubmit,
				ConsentText:          newConsent,
				PresenterName:        newPName,
				PresenterHeadline:    newPHead,
				PresenterBio:         newPBio,
				PresenterPhoto:       newPPhoto,
				PresenterLinkedinURL: newPLI,
				AnalyticsBodyClass:   newABody,
				AnalyticsCampaign:    newACamp,
			})
		},
	}
	newCmd.Flags().StringVar(&newLang, "lang", "", "language code (default: site defaultContentLanguage)")
	newCmd.Flags().StringVar(&newSlug, "slug", "", "URL slug / folder name under contentDir")
	newCmd.Flags().StringVar(&newID, "id", "", "campaign id (default: same as --slug)")
	newCmd.Flags().StringVar(&newForm, "form-name", "", "Netlify form name (default: <id>-leads)")
	newCmd.Flags().StringVar(&newTitle, "title", "", "landing page title (default: Event landing)")
	newCmd.Flags().StringVar(&newDesc, "description", "", "meta description")
	newCmd.Flags().BoolVar(&newDraft, "draft", false, "mark generated pages as draft")
	newCmd.Flags().StringVar(&newFileID, "file-id", "", "Google Drive PDF file id (default: TOML placeholder)")
	newCmd.Flags().StringVar(&newSubmit, "submit-label", "", "form submit label")
	newCmd.Flags().StringVar(&newConsent, "consent", "", "consent text (long)")
	newCmd.Flags().StringVar(&newPName, "presenter-name", "", "presenter name")
	newCmd.Flags().StringVar(&newPHead, "presenter-headline", "", "presenter headline")
	newCmd.Flags().StringVar(&newPBio, "presenter-bio", "", "presenter bio")
	newCmd.Flags().StringVar(&newPPhoto, "presenter-photo", "", "profile image path/URL")
	newCmd.Flags().StringVar(&newPLI, "presenter-linkedin", "", "LinkedIn profile URL")
	newCmd.Flags().StringVar(&newABody, "analytics-body-class", "", "body CSS class (default from slug in template)")
	newCmd.Flags().StringVar(&newACamp, "analytics-campaign", "", "analytics campaign label (default: <slug>-funnel)")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List campaigns from data/campaigns/",
		RunE: func(cmd *cobra.Command, args []string) error {
			siteRoot, err := staticlp.FindSiteRoot(siteRootFromFlag())
			if err != nil {
				return err
			}
			return staticlp.RunList(os.Stdout, siteRoot)
		},
	}

	editCmd := &cobra.Command{
		Use:   "edit [id]",
		Short: "Edit a campaign TOML (interactive)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			siteRoot, err := staticlp.FindSiteRoot(siteRootFromFlag())
			if err != nil {
				return err
			}
			var id string
			if len(args) == 1 {
				id = args[0]
			} else {
				id, err = staticlp.PickCampaignIDForCLI(siteRoot, "Campaign to edit")
				if err != nil {
					return err
				}
				if id == "" {
					return nil
				}
			}
			return staticlp.RunEditInteractive(siteRoot, id)
		},
	}

	var rmSlug, rmLang string
	var rmForce bool
	rmCmd := &cobra.Command{
		Use:   "rm [id]",
		Short: "Remove campaign TOML and content section",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			siteRoot, err := staticlp.FindSiteRoot(siteRootFromFlag())
			if err != nil {
				return err
			}
			sniff, err := staticlp.LoadHugoSniff(siteRoot)
			if err != nil {
				return err
			}
			var id string
			if len(args) == 1 {
				id = args[0]
			} else {
				id, err = staticlp.PickCampaignIDForCLI(siteRoot, "Campaign to remove")
				if err != nil {
					return err
				}
				if id == "" {
					return nil
				}
			}
			return staticlp.RemoveCampaign(siteRoot, sniff, id, rmSlug, rmLang, rmForce)
		},
	}
	rmCmd.Flags().StringVar(&rmSlug, "slug", "", "content slug path (default: infer from paths.landing)")
	rmCmd.Flags().StringVar(&rmLang, "lang", "", "language code for contentDir (default: infer or site default)")
	rmCmd.Flags().BoolVar(&rmForce, "force", false, "skip confirmation prompt")

	rootCmd.AddCommand(newCmd, listCmd, editCmd, rmCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
