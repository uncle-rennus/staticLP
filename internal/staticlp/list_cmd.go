package staticlp

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// RunList prints campaigns to w.
func RunList(w io.Writer, siteRoot string) error {
	rows, err := ListCampaigns(siteRoot)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		fmt.Fprintln(w, "(no campaigns in data/campaigns/)")
		return nil
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tFORM\tLANDING")
	for _, r := range rows {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", r.ID, r.FormName, r.Landing)
	}
	return tw.Flush()
}
