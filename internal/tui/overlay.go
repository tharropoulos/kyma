// This code is from https://github.com/charmbracelet/lipgloss/pull/102 & https://github.com/yorukot/superfile
package tui

import (
	"strings"

	charmansi "github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
)

type whitespace struct {
	style termenv.Style
	chars string
}

type WhitespaceOption func(*whitespace)

func (w whitespace) render(width int) string {
	if w.chars == "" {
		w.chars = " "
	}

	r := []rune(w.chars)
	j := 0
	b := strings.Builder{}

	// Cycle through runes and print them into the whitespace.
	for i := 0; i < width; {
		b.WriteRune(r[j])
		j++
		if j >= len(r) {
			j = 0
		}
		i += charmansi.StringWidth(string(r[j]))
	}

	// Fill any extra gaps white spaces. This might be necessary if any runes
	// are more than one cell wide, which could leave a one-rune gap.
	short := width - charmansi.StringWidth(b.String())
	if short > 0 {
		b.WriteString(strings.Repeat(" ", short))
	}

	return w.style.Styled(b.String())
}

// placeOverlay places fg on top of bg.
func placeOverlay(x, y int, fg, bg string, opts ...WhitespaceOption) string {
	fgLines, fgWidth := getLines(fg)
	bgLines, bgWidth := getLines(bg)
	bgHeight := len(bgLines)
	fgHeight := len(fgLines)

	if fgWidth >= bgWidth && fgHeight >= bgHeight {
		// FIXME: return fg or bg?
		return fg
	}
	// TODO: allow placement outside of the bg box?
	x = clamp(x, 0, bgWidth-fgWidth)
	y = clamp(y, 0, bgHeight-fgHeight)

	ws := &whitespace{}
	for _, opt := range opts {
		opt(ws)
	}

	var b strings.Builder

	for i, bgLine := range bgLines {
		if i > 0 {
			b.WriteByte('\n')
		}
		if i < y || i >= y+fgHeight {
			b.WriteString(bgLine)
			continue
		}

		pos := 0
		if x > 0 {
			left := cutLeft(bgLine, x)
			pos = charmansi.StringWidth(left)
			b.WriteString(left)
			if pos < x {
				b.WriteString(ws.render(x - pos))
				pos = x
			}
		}

		fgLine := fgLines[i-y]
		b.WriteString(fgLine)
		pos += charmansi.StringWidth(fgLine)

		right := charmansi.TruncateLeft(bgLine, pos, "")
		bgWidth = charmansi.StringWidth(bgLine)
		rightWidth := charmansi.StringWidth(right)
		if rightWidth <= bgWidth-pos {
			b.WriteString(ws.render(bgWidth - rightWidth - pos))
		}

		b.WriteString(right)
	}

	return b.String()
}

// cutLeft returns the first n cells of s, with the style they were drawn in
// reset after them. Escape codes past the cut are dropped along with the text,
// so the part of the line the overlay covers cannot leak styles or images
// into it. Every escape code is measured as x/ansi parses it: reflow ends one
// at its first letter, which cut the kitty graphics codes slides carry in
// half and left the terminal swallowing the rest of the screen.
func cutLeft(s string, n int) string {
	var (
		b     strings.Builder
		state byte
		width int
	)
	for len(s) > 0 {
		seq, w, read, newState := charmansi.DecodeSequence(s, state, nil)
		if width+w > n {
			break
		}
		state = newState
		width += w
		b.WriteString(seq)
		s = s[read:]
	}
	b.WriteString(charmansi.ResetStyle)
	return b.String()
}

func clamp(v, lower, upper int) int {
	return min(max(v, lower), upper)
}

// Split a string into lines, additionally returning the size of the widest
// line.
func getLines(s string) ([]string, int) {
	lines := strings.Split(s, "\n")
	widest := 0
	for _, l := range lines {
		w := charmansi.StringWidth(l)
		if widest < w {
			widest = w
		}
	}

	return lines, widest
}
