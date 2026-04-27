package staticlp

import (
	"io"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/x/ansi"
)

// bootFormForTest applies a terminal size and one Init message so View/layout work without a real TTY.
func bootFormForTest(f *huh.Form) tea.Model {
	m := tea.Model(f)
	m, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 44})
	if c := m.Init(); c != nil {
		if msg := c(); msg != nil {
			m, _ = m.Update(msg)
		}
	}
	return m
}

// applyKey applies one key and drains a bounded number of follow-up msgs (avoid tick loops).
func applyKey(m tea.Model, msg tea.Msg) tea.Model {
	m2, cmd := m.Update(msg)
	for i := 0; i < 4 && cmd != nil; i++ {
		next := cmd()
		if next == nil {
			break
		}
		m2, cmd = m2.Update(next)
	}
	return m2
}

func stripView(m tea.Model) string {
	f, ok := m.(*huh.Form)
	if !ok {
		return ""
	}
	return ansi.Strip(f.View())
}

func keyRunes(r ...rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: r}
}

func requireContains(tb testing.TB, view, needle string) {
	tb.Helper()
	if !strings.Contains(view, needle) {
		tb.Fatalf("expected view to contain %q; got:\n%s", needle, view)
	}
}

// TestMainMenu_SelectNewCampaign drives the root menu with Bubble Tea messages (no real TTY).
func TestMainMenu_SelectNewCampaign(t *testing.T) {
	action := "list"
	f := BuildMainMenuForm("/tmp/hugo-site", &action).
		WithOutput(io.Discard).
		WithInput(nil).
		WithAccessible(false)

	m := bootFormForTest(f)
	requireContains(t, stripView(m), "New campaign")

	m = applyKey(m, tea.KeyMsg{Type: tea.KeyDown})
	m = applyKey(m, tea.KeyMsg{Type: tea.KeyEnter})

	if action != "new" {
		t.Fatalf("expected action=new, got %q; view:\n%s", action, stripView(m))
	}
}

// TestNewCampaignWizard_FirstScreenShowsOnlyStep1 ensures the 4-group wizard does not collapse to a single step in the UI model.
func TestNewCampaignWizard_FirstScreenShowsOnlyStep1(t *testing.T) {
	fields := &NewCampaignFormFields{Lang: "en"}
	f := BuildNewCampaignWizardForm(fields).
		WithOutput(io.Discard).
		WithInput(nil).
		WithAccessible(false)

	m := bootFormForTest(f)
	v := stripView(m)
	requireContains(t, v, "1/4")
	if strings.Contains(v, "2/4") {
		t.Fatalf("step 2 header must not appear on the first screen:\n%s", v)
	}
}
