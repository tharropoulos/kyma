package markdown

import (
	"bytes"
	"io"
	"strings"
)

// maxHeadingIndent is how far a heading may be indented before markdown reads
// the line as an indented code block instead.
const maxHeadingIndent = 3

// HeadingParser lifts ATX headings such as "## Title" out of the markdown so
// they can be drawn at a larger font size. It only claims the levels it was
// given: every other heading is left for glamour to render.
type HeadingParser struct {
	levels map[int]bool
}

func NewHeadingParser(levels ...int) *HeadingParser {
	p := &HeadingParser{levels: map[int]bool{}}
	for _, l := range levels {
		p.levels[l] = true
	}
	return p
}

func (p HeadingParser) Trigger() []byte {
	return []byte{'#'}
}

// Parse attempts to extract a [HeadingNode] from the input. The leading '#'
// has to open its line, and the run of them has to be followed by a space, so
// that a #hashtag stays plain text. A closing run of '#' is dropped, as
// markdown does.
func (p HeadingParser) Parse(r *bytes.Reader, _ *MarkdownParser) Node {
	if indent, ok := lineIndent(r); !ok || indent > maxHeadingIndent {
		return nil
	}

	level := 1
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil
		}
		if b == '#' {
			level++
			continue
		}
		if b != ' ' && b != '\t' {
			return nil
		}
		break
	}

	if !p.levels[level] {
		return nil
	}

	var line bytes.Buffer
	for {
		b, err := r.ReadByte()
		if err != nil || b == '\n' {
			break
		}
		line.WriteByte(b)
	}

	text := strings.TrimSpace(line.String())
	if closed := strings.TrimRight(text, "#"); closed == "" || strings.HasSuffix(closed, " ") || strings.HasSuffix(closed, "\t") {
		text = strings.TrimSpace(closed)
	}
	if text == "" {
		return nil
	}

	return &HeadingNode{Level: level, Text: text}
}

// lineIndent reports how far the byte just consumed is indented, with tabs
// stopping every four columns, and whether only whitespace comes before it on
// its line.
func lineIndent(r *bytes.Reader) (int, bool) {
	pos, _ := r.Seek(0, io.SeekCurrent)

	// pos-1 holds the byte that was just consumed, so its line starts right
	// after the last newline before that.
	start := pos - 1
	for ; start > 0; start-- {
		if byteAt(r, start-1) == '\n' {
			break
		}
	}

	indent := 0
	for i := start; i < pos-1; i++ {
		switch byteAt(r, i) {
		case ' ':
			indent++
		case '\t':
			indent += 4 - indent%4
		default:
			return 0, false
		}
	}

	return indent, true
}

// byteAt returns the byte at off without moving the reader.
func byteAt(r *bytes.Reader, off int64) byte {
	var b [1]byte
	_, _ = r.ReadAt(b[:], off)
	return b[0]
}
