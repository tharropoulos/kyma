package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

// TestPlaceOverlay_KeepsEscapeCodesWhole covers the line a static slide opens
// with: kitty's command to clear images, which reflow used to measure as text
// and cut in half, leaving the terminal swallowing the rest of the screen.
func TestPlaceOverlay_KeepsEscapeCodesWhole(t *testing.T) {
	const clear = "\x1b_Ga=d\x1b\\"
	bg := "│" + clear + "\x1b[1;35mHello there, slide text\x1b[0m        │\n" +
		"│second line of the slide              │"
	fg := "TIMER\nTIMER"

	out := placeOverlay(3, 0, fg, bg)

	if !strings.Contains(out, clear) {
		t.Errorf("the image command was cut apart:\n%q", out)
	}
	bgLines := strings.Split(bg, "\n")
	for i, line := range strings.Split(out, "\n") {
		if got, want := ansi.StringWidth(line), ansi.StringWidth(bgLines[i]); got != want {
			t.Errorf("line %d is %d cells wide, want %d: %q", i, got, want, line)
		}
		plain := ansi.Strip(line)
		if at := strings.Index(plain, "TIMER"); at < 0 || ansi.StringWidth(plain[:at]) != 3 {
			t.Errorf("line %d does not have the overlay at column 3: %q", i, plain)
		}
	}
}
