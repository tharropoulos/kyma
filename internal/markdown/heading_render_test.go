package markdown

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"github.com/museslabs/kyma/internal/textsize"
)

func TestRenderer_HeadingWithFontSize(t *testing.T) {
	const width = 40

	for _, theme := range []string{"dark", "dracula", "notty"} {
		t.Run(theme, func(t *testing.T) {
			r, err := NewRenderer(theme, WithFontSize(map[int]int{1: 3}))
			if err != nil {
				t.Fatalf("NewRenderer() failed: %v", err)
			}

			static, err := r.Render("# A title long enough to wrap\n\nbody", false, width, 20)
			if err != nil {
				t.Fatalf("Render() failed: %v", err)
			}
			animating, err := r.Render("# A title long enough to wrap\n\nbody", true, width, 20)
			if err != nil {
				t.Fatalf("Render() failed: %v", err)
			}

			lines := strings.Split(static, "\n")
			for i, line := range lines {
				if got := ansi.StringWidth(line); got > width {
					t.Errorf("line %d is %d cells wide, more than %d: %q", i, got, width, ansi.Strip(line))
				}
			}

			if n := strings.Count(static, "\x1b_kyma;t;"); n < 2 {
				t.Errorf("heading was drawn in %d scaled lines, want it wrapped over at least 2", n)
			}
			if n := strings.Count(static, "\x1b_kyma;s;"); n != 2*strings.Count(static, "\x1b_kyma;t;") {
				t.Errorf("found %d reserved rows, want 2 under every scaled line", n)
			}

			// Transitions cut frames apart, so they get no markers, but the
			// heading has to take up the same rows or the slide would jump
			// when the transition ends.
			if strings.Contains(animating, "\x1b_kyma") {
				t.Errorf("animating render has markers:\n%q", animating)
			}
			if got, want := ansi.Strip(animating), ansi.Strip(textsize.Strip(static)); got != want {
				t.Errorf("animating render differs from the static one:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}

func TestRenderer_HeadingKeepsUnscaledLevels(t *testing.T) {
	scaled, err := NewRenderer("dark", WithFontSize(map[int]int{1: 2, 2: 1}))
	if err != nil {
		t.Fatalf("NewRenderer() failed: %v", err)
	}
	plain, err := NewRenderer("dark")
	if err != nil {
		t.Fatalf("NewRenderer() failed: %v", err)
	}

	in := "## Section\n\nbody"
	got, err := scaled.Render(in, false, 40, 20)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}
	want, err := plain.Render(in, false, 40, 20)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	if got != want {
		t.Errorf("Render() =\n%q\nwant\n%q", got, want)
	}
}

func TestRenderer_HeadingTooNarrowToScale(t *testing.T) {
	r, err := NewRenderer("dark", WithFontSize(map[int]int{1: 7}))
	if err != nil {
		t.Fatalf("NewRenderer() failed: %v", err)
	}

	out, err := r.Render("# Title", false, 10, 20)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	if strings.Contains(out, "\x1b_kyma") {
		t.Errorf("heading was scaled with no room for it:\n%q", out)
	}
	if !strings.Contains(ansi.Strip(out), "Title") {
		t.Errorf("heading is missing:\n%q", out)
	}
}

func TestInkWidth(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{name: "trailing padding", in: "Title     ", want: 5},
		{name: "styled padding", in: "\x1b[38;5;39;1mTitle\x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m", want: 5},
		{name: "padding with a background", in: "\x1b[38;5;228;48;5;63;1m Title \x1b[0m  ", want: 7},
		{name: "true colour background", in: "\x1b[48;2;1;2;3m Title \x1b[49m  ", want: 7},
		{name: "background reset", in: "\x1b[44m Title\x1b[0m   ", want: 6},
		{name: "blank", in: "     ", want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inkWidth(tt.in); got != tt.want {
				t.Errorf("inkWidth(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}
