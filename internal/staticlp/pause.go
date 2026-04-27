package staticlp

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

// WaitEnterForMenu blocks until the user presses Enter, so the main TUI can show
// results (create/edit/delete) before the menu redraws. Skipped when stdin is not a TTY.
func WaitEnterForMenu() {
	if !isatty.IsTerminal(os.Stdin.Fd()) {
		return
	}
	_, _ = fmt.Fprintln(os.Stdout, "\nPress Enter to return to the menu…")
	_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
}
