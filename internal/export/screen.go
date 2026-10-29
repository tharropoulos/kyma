package export

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/cellbuf"
)

// A screen is a frame laid out on cells, the way the terminal lays out what
// kyma writes to it.
type screen struct {
	cols, rows int
	cells      []cell
	scaled     []scaledText
}

// A cell holds a grapheme in a style. A grapheme wider than a cell is held by
// the first of its cells, and the rest hold nothing but its style.
type cell struct {
	text  string
	width int
	style cellbuf.Style
}

// scaledText is a grapheme drawn scale times its size with kitty's text sizing
// protocol, covering width cells of its own width at scale, from x and y down
// and to the right.
type scaledText struct {
	x, y         int
	text         string
	style        cellbuf.Style
	scale, width int
}

// lay lays out a frame on a screen cols by rows cells large. It follows the
// escape codes kyma draws frames with: styles, moving the cursor along a line,
// which [textsize.Render] does around scaled text, and scaled text itself.
// Images in placeholder cells are text like any other here. Whatever falls off
// the screen is dropped.
func lay(frame string, cols, rows int) screen {
	s := screen{cols: cols, rows: rows, cells: make([]cell, cols*rows)}

	p := ansi.GetParser()
	defer ansi.PutParser(p)

	var (
		pen   cellbuf.Style
		state byte
		x, y  int
	)
	for len(frame) > 0 {
		seq, width, n, newState := ansi.DecodeSequence(frame, state, p)
		state, frame = newState, frame[n:]

		switch {
		case width > 0:
			if y < rows && x >= 0 && x+width <= cols {
				row := s.cells[y*cols:]
				row[x] = cell{text: seq, width: width, style: pen}
				for i := 1; i < width; i++ {
					row[x+i] = cell{style: pen}
				}
			}
			x += width
		case seq == "\n":
			x, y = 0, y+1
		case seq == "\r":
			x = 0
		case ansi.HasCsiPrefix(seq):
			cmd := ansi.Cmd(p.Command())
			if cmd.Prefix() != 0 || cmd.Intermediate() != 0 {
				// textsize's markers among them, which only stand in
				// for scaled text until the frame is composed.
				break
			}
			switch cmd.Final() {
			case 'm':
				cellbuf.ReadStyle(p.Params(), &pen)
			case 'C':
				n, _ := p.Param(0, 1)
				x += max(n, 1)
			case 'D':
				n, _ := p.Param(0, 1)
				x = max(x-max(n, 1), 0)
			}
		case ansi.HasOscPrefix(seq):
			t, ok := parseScaled(seq)
			if !ok {
				break
			}
			t.x, t.y, t.style = x, y, pen
			if y+t.scale <= rows && x+t.scale*t.width <= cols {
				s.scaled = append(s.scaled, t)
			}
			x += t.scale * t.width
		}
	}

	return s
}

// parseScaled reads text sized with kitty's text sizing protocol, OSC 66, out
// of an escape code: its scale and the width it is drawn at, which the text
// is measured for when it doesn't say.
func parseScaled(seq string) (scaledText, bool) {
	body, ok := strings.CutPrefix(seq, "\x1b]66;")
	if !ok {
		return scaledText{}, false
	}
	body = strings.TrimSuffix(strings.TrimSuffix(body, "\a"), "\x1b\\")
	meta, text, ok := strings.Cut(body, ";")
	if !ok {
		return scaledText{}, false
	}

	t := scaledText{text: text, scale: 1}
	for kv := range strings.SplitSeq(meta, ":") {
		k, v, _ := strings.Cut(kv, "=")
		n, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		switch k {
		case "s":
			t.scale = min(max(n, 1), 7)
		case "w":
			t.width = n
		}
	}
	if t.width == 0 {
		t.width = ansi.StringWidth(text)
	}
	return t, t.width > 0
}
