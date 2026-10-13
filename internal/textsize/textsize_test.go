package textsize

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	reflow "github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/truncate"

	"github.com/museslabs/kyma/internal/skip"
)

func TestBlock(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		scale int
	}{
		{name: "plain text", text: "Title", scale: 2},
		{name: "styled text", text: "\x1b[1;35mTitle\x1b[0m", scale: 3},
		{name: "wide characters", text: "你好", scale: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.scale * ansi.StringWidth(tt.text)

			rows := Block(tt.text, tt.scale)
			if len(rows) != tt.scale {
				t.Fatalf("Block() returned %d rows, want %d", len(rows), tt.scale)
			}
			for i, row := range rows {
				// Transitions measure with reflow and everything else with
				// x/ansi, so the markers have to be invisible to both.
				if got := ansi.StringWidth(row); got != want {
					t.Errorf("row %d is %d cells wide to x/ansi, want %d", i, got, want)
				}
				if got := reflow.PrintableRuneWidth(row); got != want {
					t.Errorf("row %d is %d cells wide to reflow, want %d", i, got, want)
				}
				if got := strings.Count(row, markerPrefix); got != len([]rune(ansi.Strip(tt.text))) {
					t.Errorf("row %d has %d markers, want one per grapheme", i, got)
				}
			}
		})
	}
}

// frame puts rows inside a border, the way a slide sits in the terminal.
func frame(rows []string) string {
	lines := make([]string, len(rows))
	for i, row := range rows {
		lines[i] = "│" + row + "│"
	}
	return strings.Join(lines, "\n")
}

func TestRender(t *testing.T) {
	view := frame(Block("Hi", 2))

	want := "│" +
		"  \x1b[2D\x1b[m\x1b]66;s=2;H\a\x1b[m" +
		"  \x1b[2D\x1b[m\x1b]66;s=2;i\a\x1b[m" +
		"│\n│\x1b[2C\x1b[2C│"
	if got := Render(view, 6, 2); got != want {
		t.Errorf("Render() =\n%q\nwant\n%q", got, want)
	}
}

func TestRender_Styles(t *testing.T) {
	out := Render(frame(Block("\x1b[1;35mH\x1b[0m!你", 2)), 10, 2)

	for _, want := range []string{
		"\x1b[m\x1b[1;35m\x1b]66;s=2;H\a",
		"\x1b[m\x1b[0m\x1b]66;s=2;!\a",
		"\x1b]66;s=2:w=2;你\a",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Render() =\n%q\nwant it to contain %q", out, want)
		}
	}
}

func TestRender_LeavesCutGlyphsBlank(t *testing.T) {
	rows := Block("Hello", 2)

	cut := func(f func(string) string) []string {
		out := make([]string, len(rows))
		for i, row := range rows {
			out[i] = f(row)
		}
		return out
	}

	tests := []struct {
		name   string
		view   string
		width  int
		height int
		want   int
	}{
		{
			name:  "whole",
			view:  frame(rows),
			width: 12, height: 2,
			want: 5,
		},
		{
			name: "cut on the right, as a swipe does",
			view: frame(cut(func(s string) string { return truncate.String(s, 5) })),
			// H and e fit in the first four cells, l is cut in half.
			width: 12, height: 2,
			want: 2,
		},
		{
			name: "cut on the left, as a swipe does",
			view: frame(cut(func(s string) string { return skip.String(s, 3) })),
			// H and half of e are gone, leaving l, l and o.
			width: 12, height: 2,
			want: 3,
		},
		{
			name:  "cut by x/ansi",
			view:  frame(cut(func(s string) string { return ansi.TruncateLeft(s, 4, "") })),
			width: 12, height: 2,
			want: 3,
		},
		{
			name:  "the row beneath is missing",
			view:  frame(rows[:1]),
			width: 12, height: 2,
			want: 0,
		},
		{
			name:  "the row beneath has moved",
			view:  frame([]string{rows[0], " " + rows[1]}),
			width: 13, height: 2,
			want: 0,
		},
		{
			name:  "past the right edge",
			view:  frame(rows),
			width: 9, height: 2,
			// The border takes the first column, so o is the one that ends
			// past the ninth.
			want: 4,
		},
		{
			name:  "above the top of the screen",
			view:  frame(rows),
			width: 12, height: 1,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := Render(tt.view, tt.width, tt.height)

			if strings.Contains(out, markerPrefix) {
				t.Errorf("a marker was left in the frame: %q", out)
			}
			if got := strings.Count(out, "\x1b]66;"); got != tt.want {
				t.Errorf("drew %d glyphs, want %d:\n%q", got, tt.want, out)
			}
			// What is not drawn stays blank, so every line keeps its width.
			for i, line := range strings.Split(out, "\n") {
				if got, want := ansi.StringWidth(line), ansi.StringWidth(strings.Split(tt.view, "\n")[i]); tt.want == 0 && got != want {
					t.Errorf("line %d is %d cells wide, want %d", i, got, want)
				}
			}
		})
	}
}

func TestRender_WithoutMarkers(t *testing.T) {
	view := "│ plain │\n│ text  │"
	if got := Render(view, 9, 2); got != view {
		t.Errorf("Render() = %q, want the frame unchanged", got)
	}
}

func TestRender_CutsOverflowingLines(t *testing.T) {
	// A swipe joins the outgoing and incoming slides with a cell to spare.
	view := frame(Block("Hi", 2))
	lines := strings.Split(view, "\n")
	for i := range lines {
		lines[i] += "x"
	}

	out := Render(strings.Join(lines, "\n"), 6, 2)

	cursorForward := regexp.MustCompile(`\x1b\[(\d*)C`)
	for i, line := range strings.Split(out, "\n") {
		// Cursor moves over the rows beneath a glyph take up cells that
		// x/ansi does not count.
		cells := ansi.StringWidth(line)
		for _, m := range cursorForward.FindAllStringSubmatch(line, -1) {
			n, _ := strconv.Atoi(m[1])
			cells += max(n, 1)
		}
		if cells != 6 {
			t.Errorf("line %d takes up %d cells, want 6: %q", i, cells, line)
		}
	}
}
