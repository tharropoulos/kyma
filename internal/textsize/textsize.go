// Package textsize draws text larger than the body text with kitty's text
// sizing protocol (OSC 66), see
// https://sw.kovidgoyal.net/kitty/text-sizing-protocol/.
//
// Scaled text covers a block of cells several rows tall, which none of the
// layout code in the stack knows how to measure: lipgloss sees the escape code
// as zero cells wide, and bubbletea redraws the screen line by line, so the
// rows under the text get overwritten. Text is therefore drawn in two steps.
// [Block] reserves the cells with ordinary placeholder text, tagged with zero
// width markers, that every layout step measures and moves like any other
// text. Once the frame is composed, [Render] swaps each placeholder for the
// escape codes that draw the text, while [Strip] drops the markers instead and
// leaves the placeholder showing the text at its normal size.
package textsize

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// MaxScale is the largest scale the protocol supports.
const MaxScale = 7

const (
	markerPrefix = "\x1b_kyma;"
	markerSuffix = "\x1b\\"

	// kindText marks the first row of a block, where the text is drawn.
	kindText = "t"
	// kindSkip marks a row beneath it, which the text hangs down over.
	kindSkip = "s"

	// nbsp fills the blank cells of a block. It looks like a space, but code
	// that measures a block's ink by trimming trailing spaces still sees how
	// wide the scaled text will be.
	nbsp = "\u00a0"

	// maxPayload is the most text the protocol accepts in one escape code.
	maxPayload = 4096
)

// markerPattern matches a whole marker, payload included.
var markerPattern = regexp.MustCompile(regexp.QuoteMeta(markerPrefix) + `[^\x1b]*` + regexp.QuoteMeta(markerSuffix))

// Block lays out a single line of styled text at scale times its size. It
// returns scale rows, each scale times as wide as text. The first row holds the
// text at its normal size followed by blank cells and the others are blank, so
// a block whose markers are stripped still reads as a plain heading. When
// marked is false the markers are left out altogether, for frames that are
// going to be cut up by a transition.
func Block(text string, scale int, marked bool) []string {
	width := ansi.StringWidth(text)
	cells := scale * width

	rows := make([]string, scale)
	rows[0] = text + strings.Repeat(nbsp, cells-width)
	for i := 1; i < scale; i++ {
		rows[i] = strings.Repeat(nbsp, cells)
	}

	if !marked || width == 0 {
		return rows
	}

	rows[0] = marker(kindText, cells, sequence(text, scale)) + rows[0]
	for i := 1; i < scale; i++ {
		rows[i] = marker(kindSkip, cells, "") + rows[i]
	}

	return rows
}

// Render replaces the placeholder of every block in a composed frame with the
// escape codes that draw its text. The first row of a block becomes blank
// cells, so that redrawing the line wipes the text drawn there before, then a
// cursor move back over them and the text itself. The rows beneath become a
// cursor move over the cells the text hangs down into, since writing anything
// there, even a space, would push the rest of the line to the right.
func Render(view string) string {
	return rewrite(view, func(kind string, cells int, payload string) string {
		if kind == kindSkip {
			return ansi.CursorForward(cells)
		}
		return strings.Repeat(" ", cells) + ansi.CursorBackward(cells) + payload
	})
}

// Strip removes the markers from a frame, leaving each block's text drawn at
// its normal size. It is for frames that get composed further after rendering,
// such as the slide under the timer overlay.
func Strip(view string) string {
	if !strings.Contains(view, markerPrefix) {
		return view
	}
	return markerPattern.ReplaceAllString(view, "")
}

// rewrite swaps every marked placeholder in view for whatever replace returns.
// A placeholder that has lost cells, because the line was cut short, is left
// as it is with only its marker removed.
func rewrite(view string, replace func(kind string, cells int, payload string) string) string {
	if !strings.Contains(view, markerPrefix) {
		return view
	}

	var b strings.Builder
	for {
		start := strings.Index(view, markerPrefix)
		if start < 0 {
			b.WriteString(view)
			return b.String()
		}
		end := strings.Index(view[start:], markerSuffix)
		if end < 0 {
			b.WriteString(view)
			return b.String()
		}

		b.WriteString(view[:start])
		body := view[start+len(markerPrefix) : start+end]
		view = view[start+end+len(markerSuffix):]

		kind, cells, payload, ok := parseMarker(body)
		if !ok {
			continue
		}

		rest, ok := skipCells(view, cells)
		if !ok {
			continue
		}

		b.WriteString(replace(kind, cells, payload))
		view = rest
	}
}

func marker(kind string, cells int, payload string) string {
	m := markerPrefix + kind + ";" + strconv.Itoa(cells)
	if payload != "" {
		m += ";" + base64.StdEncoding.EncodeToString([]byte(payload))
	}
	return m + markerSuffix
}

func parseMarker(body string) (kind string, cells int, payload string, ok bool) {
	parts := strings.SplitN(body, ";", 3)
	if len(parts) < 2 {
		return "", 0, "", false
	}

	cells, err := strconv.Atoi(parts[1])
	if err != nil || cells <= 0 {
		return "", 0, "", false
	}

	if len(parts) == 3 {
		decoded, err := base64.StdEncoding.DecodeString(parts[2])
		if err != nil {
			return "", 0, "", false
		}
		payload = string(decoded)
	}

	return parts[0], cells, payload, true
}

// skipCells returns what is left of s after its first n cells, dropping the
// escape codes that styled them. It reports false when s is narrower than n.
func skipCells(s string, n int) (string, bool) {
	var state byte
	for n > 0 && len(s) > 0 {
		_, width, read, newState := ansi.DecodeSequence(s, state, nil)
		state = newState
		s = s[read:]
		n -= width
	}

	return s, n <= 0
}

// sequence returns the escape codes that draw styled text at scale times its
// size. Styles carry over from the SGR codes in between, which set the cursor
// attributes the terminal draws the text with. ASCII runs are sent whole and
// left to the terminal to split into cells, but every other grapheme gets its
// own escape code with the width kyma measured for it, so the terminal never
// disagrees with the layout about how many cells the text covers.
func sequence(text string, scale int) string {
	var (
		b     strings.Builder
		run   strings.Builder
		state byte
	)

	flush := func() {
		if run.Len() > 0 {
			fmt.Fprintf(&b, "\x1b]66;s=%d;%s\a", scale, run.String())
			run.Reset()
		}
	}

	for len(text) > 0 {
		seq, width, read, newState := ansi.DecodeSequence(text, state, nil)
		state = newState
		text = text[read:]

		switch {
		case width == 1 && len(seq) == 1:
			run.WriteString(seq)
			if run.Len() == maxPayload {
				flush()
			}
		case width > 0:
			flush()
			fmt.Fprintf(&b, "\x1b]66;s=%d:w=%d;%s\a", scale, width, seq)
		case strings.HasPrefix(seq, "\x1b[") && strings.HasSuffix(seq, "m"):
			flush()
			b.WriteString(seq)
		}
	}
	flush()

	b.WriteString(ansi.ResetStyle)
	return b.String()
}
