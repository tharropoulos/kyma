// Package textsize draws text larger than the body text with kitty's text
// sizing protocol (OSC 66), see
// https://sw.kovidgoyal.net/kitty/text-sizing-protocol/.
//
// Scaled text covers a block of cells several rows tall, which none of the
// layout code in the stack knows how to measure: lipgloss sees the escape code
// as zero cells wide, transitions cut lines apart at arbitrary columns, and
// bubbletea redraws the screen line by line, so the rows under the text get
// overwritten. Text is therefore drawn in two steps. [Block] lays the text out
// as blank placeholder cells, a block of them per grapheme, with every row of
// a block opened by a zero width marker. Every layout step, transitions
// included, measures and moves those like any other text. Once the frame is
// composed, [Render] draws every glyph whose placeholder survived intact and
// leaves the rest blank.
package textsize

import (
	"strconv"
	"strings"
	"sync"

	"github.com/charmbracelet/x/ansi"
)

// MaxScale is the largest scale the protocol supports.
const MaxScale = 7

const (
	// A marker is the private CSI sequence ESC [ < 66 ; id ; row k. Both
	// x/ansi and reflow have to measure it as zero cells wide, and reflow ends
	// an escape sequence at its first letter, so it holds digits only and
	// ends with one.
	markerPrefix = "\x1b[<66;"
	markerFinal  = "k"

	// nbsp fills the placeholder cells. It looks like a space, but code that
	// measures a block's ink by trimming trailing spaces still sees how wide
	// the scaled text will be drawn, and it tells placeholder cells apart from
	// any other blank.
	nbsp = " "

	// maxPayload is the most text the protocol accepts in one escape code.
	maxPayload = 4096
)

// glyph is one grapheme of scaled text, drawn as a block of cells as many rows
// tall as its scale.
type glyph struct {
	// draw holds the escape codes that draw the grapheme in its style.
	draw  string
	cells int
	scale int
}

// glyphs interns every glyph laid out, so that a marker only has to carry an
// id. There are only as many as there are distinct styled graphemes across
// the headings of a deck.
var glyphs = struct {
	sync.Mutex
	ids  map[glyph]int
	list []glyph
}{ids: map[glyph]int{}}

func intern(g glyph) int {
	glyphs.Lock()
	defer glyphs.Unlock()

	if id, ok := glyphs.ids[g]; ok {
		return id
	}
	glyphs.list = append(glyphs.list, g)
	glyphs.ids[g] = len(glyphs.list) - 1
	return len(glyphs.list) - 1
}

func lookup(id int) (glyph, bool) {
	glyphs.Lock()
	defer glyphs.Unlock()

	if id < 0 || id >= len(glyphs.list) {
		return glyph{}, false
	}
	return glyphs.list[id], true
}

// Block lays out a single line of styled text at scale times its size. It
// returns scale rows, each scale times as wide as text, holding a block of
// placeholder cells for every grapheme, see [Render].
func Block(text string, scale int) []string {
	rows := make([]strings.Builder, scale)

	var (
		style string
		state byte
	)
	for len(text) > 0 {
		seq, width, read, newState := ansi.DecodeSequence(text, state, nil)
		state = newState
		text = text[read:]

		if isSGR(seq) {
			// The codes since the last reset make up the style the next
			// grapheme is drawn in.
			if params := seq[2 : len(seq)-1]; params == "" || params == "0" || strings.HasPrefix(params, "0;") {
				style = ""
			}
			style += seq
			continue
		}
		if width == 0 {
			continue
		}

		g := glyph{draw: draw(style, seq, scale, width), cells: scale * width, scale: scale}
		id := intern(g)
		for row := range rows {
			rows[row].WriteString(marker(id, row))
			rows[row].WriteString(strings.Repeat(nbsp, g.cells))
		}
	}

	out := make([]string, scale)
	for i := range rows {
		out[i] = rows[i].String()
	}
	return out
}

// draw returns the escape codes that draw a grapheme at scale times its size.
// ASCII is left to the terminal to measure, but anything else is sent with the
// width kyma measured for it, so the terminal never disagrees with the layout
// about how many cells it covers.
func draw(style, grapheme string, scale, width int) string {
	size := "s=" + strconv.Itoa(scale)
	if len(grapheme) > 1 {
		size += ":w=" + strconv.Itoa(width)
	}
	if len(grapheme) > maxPayload {
		grapheme = grapheme[:maxPayload]
	}
	return ansi.ResetStyle + style + "\x1b]66;" + size + ";" + grapheme + "\a" + ansi.ResetStyle
}

func marker(id, row int) string {
	return markerPrefix + strconv.Itoa(id) + ";" + strconv.Itoa(row) + markerFinal
}

// placement is a marker found in a composed frame.
type placement struct {
	// start is the byte offset of the marker, cells of the placeholder cells
	// that follow it and end of the end of those.
	start, cells, end int
	col               int
	id, row           int
	// intact is whether all of the glyph's placeholder cells follow the
	// marker, uncut and without another marker among them.
	intact bool
}

// Render draws the scaled text in a composed frame that is width cells wide.
// Like the terminal, it only shows the last height lines of a taller frame.
//
// A glyph is drawn only where its placeholder survived every layout step
// whole: each of its rows has to sit at the same column, one below the other,
// fully on screen. Whatever a transition or an overlay cut through is left
// blank. The first row of a glyph becomes blank cells, which wipe whatever
// was drawn there before when the line is redrawn, then a cursor move back
// over them and the escape codes that draw the glyph. The rows beneath become
// a cursor move over the cells the glyph hangs down into, since writing
// anything there, even a space, would push the rest of the line right.
//
// Lines holding scaled text are cut to width here rather than by bubbletea,
// which cuts lines that overflow the terminal, as a swipe's do by a cell.
// Once a row has been turned into cursor moves bubbletea measures it shorter
// than it is, would leave it whole and the terminal would wrap it.
func Render(view string, width, height int) string {
	if !strings.Contains(view, markerPrefix) {
		return view
	}

	lines := strings.Split(view, "\n")
	found := make([][]placement, len(lines))
	at := map[[2]int]placement{}
	for i, line := range lines {
		if width > 0 && strings.Contains(line, markerPrefix) {
			line = ansi.Truncate(line, width, "")
			lines[i] = line
		}
		found[i] = scan(line)
		// Markers cut off from their cells pile up at the column the cut
		// left, so only intact ones are told apart by where they sit.
		for _, p := range found[i] {
			if p.intact {
				at[[2]int{i, p.col}] = p
			}
		}
	}

	first := max(len(lines)-height, 0)
	drawn := map[[2]int]bool{}
	for i := first; i < len(lines); i++ {
		for _, p := range found[i] {
			g, ok := lookup(p.id)
			if !ok || p.row != 0 || !p.intact || p.col+g.cells > width || i+g.scale > len(lines) {
				continue
			}

			whole := true
			for row := 1; row < g.scale && whole; row++ {
				q, ok := at[[2]int{i + row, p.col}]
				whole = ok && q.intact && q.id == p.id && q.row == row
			}
			if !whole {
				continue
			}

			for row := range g.scale {
				drawn[[2]int{i + row, p.col}] = true
			}
		}
	}

	for i, line := range lines {
		if len(found[i]) == 0 {
			continue
		}

		var b strings.Builder
		last := 0
		for _, p := range found[i] {
			b.WriteString(line[last:p.start])
			last = p.end

			switch {
			case !p.intact || !drawn[[2]int{i, p.col}]:
				// Drop the marker and leave the placeholder cells blank.
				b.WriteString(line[p.start+len(marker(p.id, p.row)) : p.end])
			case p.row == 0:
				g, _ := lookup(p.id)
				b.WriteString(strings.Repeat(" ", p.cells))
				b.WriteString(ansi.CursorBackward(p.cells))
				b.WriteString(g.draw)
			default:
				b.WriteString(ansi.CursorForward(p.cells))
			}
		}
		b.WriteString(line[last:])
		lines[i] = b.String()
	}

	return strings.Join(lines, "\n")
}

// scan finds the markers in a line, the column each sits at and whether the
// placeholder cells after it are intact.
func scan(line string) []placement {
	var (
		found []placement
		state byte
		col   int
	)

	for i := 0; i < len(line); {
		seq, width, read, newState := ansi.DecodeSequence(line[i:], state, nil)
		state = newState

		id, row, ok := parseMarker(seq)
		if !ok {
			i += read
			col += width
			continue
		}

		p := placement{start: i, col: col, id: id, row: row}
		i += read

		g, known := lookup(id)
		for known && p.cells < g.cells && i < len(line) {
			seq, width, read, newState := ansi.DecodeSequence(line[i:], state, nil)
			if _, _, ok := parseMarker(seq); ok || (width > 0 && seq != nbsp) {
				break
			}
			state = newState
			i += read
			p.cells += width
		}
		p.end = i
		p.intact = known && p.cells == g.cells
		col += p.cells

		found = append(found, p)
	}

	return found
}

// IsMarker reports whether seq is one of the zero width markers [Block] opens
// the rows of its placeholders with. Code that takes a frame apart cell by
// cell has to put them back in front of the cell they came before, or [Render]
// won't find the text they stand for.
func IsMarker(seq string) bool {
	_, _, ok := parseMarker(seq)
	return ok
}

func parseMarker(seq string) (id, row int, ok bool) {
	body, found := strings.CutPrefix(seq, markerPrefix)
	if !found {
		return 0, 0, false
	}
	body, found = strings.CutSuffix(body, markerFinal)
	if !found {
		return 0, 0, false
	}

	ids, rows, found := strings.Cut(body, ";")
	if !found {
		return 0, 0, false
	}
	id, err := strconv.Atoi(ids)
	if err != nil {
		return 0, 0, false
	}
	row, err = strconv.Atoi(rows)
	if err != nil {
		return 0, 0, false
	}

	return id, row, true
}

func isSGR(seq string) bool {
	return strings.HasPrefix(seq, "\x1b[") && strings.HasSuffix(seq, "m")
}
