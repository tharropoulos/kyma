package textsize

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
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

			for _, marked := range []bool{true, false} {
				rows := Block(tt.text, tt.scale, marked)
				if len(rows) != tt.scale {
					t.Fatalf("Block() returned %d rows, want %d", len(rows), tt.scale)
				}
				for i, row := range rows {
					if got := ansi.StringWidth(row); got != want {
						t.Errorf("row %d is %d cells wide, want %d: %q", i, got, want, row)
					}
					if got := strings.Contains(row, markerPrefix); got != marked {
						t.Errorf("row %d has a marker: %t, want %t", i, got, marked)
					}
				}
				if !strings.Contains(ansi.Strip(rows[0]), ansi.Strip(tt.text)) {
					t.Errorf("first row %q does not show the text %q", rows[0], tt.text)
				}
			}
		})
	}
}

func TestRender(t *testing.T) {
	rows := Block("Hi", 2, true)
	view := "│" + rows[0] + "  │\n│" + rows[1] + "  │"

	want := "│    \x1b[4D\x1b]66;s=2;Hi\a\x1b[m  │\n│\x1b[4C  │"
	if got := Render(view); got != want {
		t.Errorf("Render() =\n%q\nwant\n%q", got, want)
	}
}

func TestRender_LeavesCutPlaceholders(t *testing.T) {
	rows := Block("Hi", 2, true)
	// A line cut short partway through the placeholder, as a narrow column
	// would, must not have scaled text drawn over whatever follows it.
	cut := ansi.Truncate(rows[0], 3, "")

	if got, want := Render(cut), "Hi"+nbsp; got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

func TestStrip(t *testing.T) {
	rows := Block("\x1b[1mHi\x1b[0m", 2, true)
	view := strings.Join(rows, "\n")

	want := strings.Join(Block("\x1b[1mHi\x1b[0m", 2, false), "\n")
	if got := Strip(view); got != want {
		t.Errorf("Strip() =\n%q\nwant\n%q", got, want)
	}
}

func TestSequence(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		scale int
		want  string
	}{
		{
			name:  "ascii is sent in one piece",
			text:  "Hello world",
			scale: 2,
			want:  "\x1b]66;s=2;Hello world\a\x1b[m",
		},
		{
			name:  "styles sit between the escape codes",
			text:  "\x1b[1mBold\x1b[0m plain",
			scale: 3,
			want:  "\x1b[1m\x1b]66;s=3;Bold\a\x1b[0m\x1b]66;s=3; plain\a\x1b[m",
		},
		{
			name:  "wide graphemes carry their width",
			text:  "a你b",
			scale: 2,
			want:  "\x1b]66;s=2;a\a\x1b]66;s=2:w=2;你\a\x1b]66;s=2;b\a\x1b[m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sequence(tt.text, tt.scale); got != tt.want {
				t.Errorf("sequence() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSequence_SplitsLongText(t *testing.T) {
	got := sequence(strings.Repeat("a", maxPayload+1), 2)
	if n := strings.Count(got, "\x1b]66;"); n != 2 {
		t.Errorf("sequence() used %d escape codes, want 2", n)
	}
}
