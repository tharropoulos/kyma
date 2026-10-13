package markdown

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/x/ansi"

	"github.com/museslabs/kyma/internal/textsize"
)

// renderHeading draws a heading at the font size set for its level. Glamour
// still styles it, so it keeps the theme's colours and prefixes, but wrapped to
// the width left once the text is scaled up. Every line it wraps to becomes a
// block of rows as tall as the scale, see [textsize.Block].
func (r *Renderer) renderHeading(n *HeadingNode, animating bool, width int, b *strings.Builder) error {
	md := strings.Repeat("#", n.Level) + " " + n.Text
	scale := r.options.fontSize[n.Level]

	// A glamour renderer built for width indents its text by the margin and
	// wraps it a margin short of the wrap width it was given, which is itself a
	// margin short of width. The heading gets the same room as the body text.
	room := width - 3*glamourMargin

	lines, err := r.wrapHeading(md, scale, room)
	if err != nil {
		return err
	}

	// Without room for the scaled text the heading keeps its normal size.
	if lines == nil {
		tr, err := r.termRenderer(width)
		if err != nil {
			return err
		}
		out, err := tr.Render(md)
		if err != nil {
			return err
		}
		b.WriteString(out)
		return nil
	}

	// Glamour opens every document with a blank line. Inside the surrounding
	// markdown the heading would only get it when it opens the slide or column.
	if isFirstChild(n) {
		b.WriteString("\n")
	}

	for _, l := range lines {
		pad := strings.Repeat(" ", room-scale*l.width)
		for _, row := range textsize.Block(l.text, scale, !animating) {
			b.WriteString(l.margin + row + pad + "\n")
		}
	}

	return nil
}

// headingLine is one line of a heading as glamour styled it, split into the
// margin glamour indents it by and the text itself, padding left out.
type headingLine struct {
	margin string
	text   string
	width  int
}

// wrapHeading renders md with glamour, wrapped so that every line still fits in
// room cells once it is scaled up. It returns nil when not even a single cell
// of scaled text fits.
func (r *Renderer) wrapHeading(md string, scale, room int) ([]headingLine, error) {
	for wrap := room / scale; wrap > 0; {
		tr, err := r.termRenderer(wrap + 3*glamourMargin)
		if err != nil {
			return nil, err
		}
		out, err := tr.Render(md)
		if err != nil {
			return nil, err
		}

		var lines []headingLine
		widest := 0
		for _, line := range strings.Split(out, "\n") {
			margin, rest := splitCells(line, glamourMargin)
			width := inkWidth(rest)
			if width == 0 {
				continue
			}
			text, _ := splitCells(rest, width)
			lines = append(lines, headingLine{margin: margin, text: text, width: width})
			widest = max(widest, width)
		}

		if widest*scale <= room {
			return lines, nil
		}

		// Some themes pad a heading past the width glamour wraps it at, like
		// the coloured spaces either side of the dark theme's h1, so wrap it
		// tighter by however much it overflowed.
		wrap -= widest - room/scale
	}

	return nil, nil
}

func isFirstChild(n Node) bool {
	p := n.Parent()
	return p == nil || len(p.Children()) == 0 || p.Children()[0] == n
}

// splitCells splits s after its first n cells. Escape codes that follow the
// last of those cells stay with the second half, which they style.
func splitCells(s string, n int) (string, string) {
	var state byte
	i := 0
	for n > 0 && i < len(s) {
		_, width, read, newState := ansi.DecodeSequence(s[i:], state, nil)
		state = newState
		i += read
		n -= width
	}
	return s[:i], s[i:]
}

// inkWidth is the width of s up to its last cell that shows anything: a
// character other than a space, or a space with a background colour, such as
// the padding some themes put around a heading. The plain spaces glamour pads
// a line out with do not count.
func inkWidth(s string) int {
	var (
		state    byte
		cells    int
		ink      int
		coloured bool
	)

	for len(s) > 0 {
		seq, width, read, newState := ansi.DecodeSequence(s, state, nil)
		state = newState
		s = s[read:]

		if width == 0 {
			if strings.HasPrefix(seq, "\x1b[") && strings.HasSuffix(seq, "m") {
				coloured = setsBackground(seq, coloured)
			}
			continue
		}

		cells += width
		if seq != " " || coloured {
			ink = cells
		}
	}

	return ink
}

// setsBackground reports whether text drawn after the SGR sequence seq has a
// background colour, given whether it had one before.
func setsBackground(seq string, bg bool) bool {
	params := strings.Split(seq[len("\x1b["):len(seq)-1], ";")
	for i := 0; i < len(params); i++ {
		p, _ := strconv.Atoi(params[i])
		switch {
		case p == 0 || p == 49:
			bg = false
		case p >= 40 && p <= 47, p >= 100 && p <= 107:
			bg = true
		case p == 38 || p == 48 || p == 58:
			// Extended colours carry their value in the parameters that follow.
			if p == 48 {
				bg = true
			}
			if i+1 < len(params) && params[i+1] == "5" {
				i += 2
			} else if i+1 < len(params) && params[i+1] == "2" {
				i += 4
			}
		}
	}
	return bg
}
