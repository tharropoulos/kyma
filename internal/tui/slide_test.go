package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"github.com/museslabs/kyma/internal/config"
)

// TestSlide_ViewFitsTheFrame is the invariant a presenter actually notices: the
// slide, border included, is exactly as wide as the terminal, so a grid never
// pushes content past the frame.
func TestSlide_ViewFitsTheFrame(t *testing.T) {
	bodies := map[string]string{
		"plain":   "# Title\n\nsome text",
		"columns": "[row]\n[col]left[/col]\n[col]right[/col]\n[/row]",
		"master": "# Title\n\n[grid]\n[col span=2]\n## Main\n\nbody\n[/col]\n" +
			"[col]\n[row]one[/row]\n[row]two[/row]\n[/col]\n[/grid]",
		"bordered columns": "[row gap=2]\n[col border=rounded pad=1]a[/col]\n" +
			"[col border=rounded pad=1]b[/col]\n[/row]",
	}

	for _, size := range [][2]int{{80, 24}, {120, 40}, {60, 20}} {
		width, height := size[0], size[1]

		for name, body := range bodies {
			t.Run(name, func(t *testing.T) {
				slide, err := NewSlide(body, config.Properties{})
				if err != nil {
					t.Fatalf("NewSlide() failed: %v", err)
				}
				slide.Style = config.StyleConfig{}.Apply(width, height)

				for i, line := range strings.Split(slide.View(false, width, height), "\n") {
					if got := ansi.StringWidth(line); got != width {
						t.Errorf(
							"%dx%d: line %d is %d cells wide, want %d:\n%q",
							width, height, i, got, width, ansi.Strip(line),
						)
					}
				}
			})
		}
	}
}

// TestSlide_LayoutAlignsText checks that a slide's layout moves its text on a
// terminal wider than the text: center puts it midway between where left and
// right do, and those two are far apart. Text padded to the slide's full width
// only shifts by the few cells of glamour's margin.
func TestSlide_LayoutAlignsText(t *testing.T) {
	const width, height = 120, 30

	indent := map[string]int{}
	for _, layout := range []string{"left", "center", "right"} {
		props, err := config.NewProperties("style:\n  layout: " + layout)
		if err != nil {
			t.Fatalf("NewProperties() failed: %v", err)
		}
		slide, err := NewSlide("# Title\n\nA short line.", props)
		if err != nil {
			t.Fatalf("NewSlide() failed: %v", err)
		}
		slide.Style = props.Style.Apply(width, height)

		for _, line := range strings.Split(ansi.Strip(slide.View(false, width, height)), "\n") {
			if i := strings.Index(line, "A short line."); i >= 0 {
				indent[layout] = len([]rune(line[:i]))
			}
		}
	}

	left, center, right := indent["left"], indent["center"], indent["right"]
	if right-left < width/4 || center-left != right-center {
		t.Errorf(
			"text starts at column %d for left, %d for center and %d for right, "+
				"want left and right at least %d apart with center midway",
			left, center, right, width/4,
		)
	}
}

// TestSlide_LayoutAlignsGrid checks where a centered slide puts a grid: in line
// with the text around it, or across the whole slide when the slide is nothing
// but the grid, as a master layout is.
func TestSlide_LayoutAlignsGrid(t *testing.T) {
	const width, height = 160, 30
	grid := "[row]\n[col]\nleft column\n[/col]\n[col]\nright column\n[/col]\n[/row]\n"

	indents := func(t *testing.T, body string) map[string]int {
		t.Helper()

		props, err := config.NewProperties("style:\n  layout: center")
		if err != nil {
			t.Fatalf("NewProperties() failed: %v", err)
		}
		slide, err := NewSlide(body, props)
		if err != nil {
			t.Fatalf("NewSlide() failed: %v", err)
		}
		slide.Style = props.Style.Apply(width, height)

		found := map[string]int{}
		for _, line := range strings.Split(ansi.Strip(slide.View(false, width, height)), "\n") {
			for _, text := range []string{"Some text.", "left column"} {
				if i := strings.Index(line, text); i >= 0 {
					found[text] = len([]rune(line[:i]))
				}
			}
		}
		return found
	}

	t.Run("among text", func(t *testing.T) {
		got := indents(t, "Some text.\n\n"+grid)
		text, column := got["Some text."], got["left column"]
		if text < width/4 || column != text {
			t.Errorf(
				"text starts at column %d and the grid's text at %d, want both "+
					"the same and centered past column %d",
				text, column, width/4,
			)
		}
	})

	t.Run("on its own", func(t *testing.T) {
		if column := indents(t, grid)["left column"]; column >= width/4 {
			t.Errorf("grid's text starts at column %d, want the grid to fill the slide", column)
		}
	})
}

// TestSlide_GridFillsTheSlide checks that a grid claiming the slide's height
// with span ends on the slide's last row rather than past it, blank lines
// around its tags included.
func TestSlide_GridFillsTheSlide(t *testing.T) {
	const width, height = 120, 40
	grid := "[grid]\n[row span=1 border=rounded]\nonly row\n[/row]\n[/grid]"

	for name, body := range map[string]string{
		"bare":               grid,
		"blank lines around": "\n" + grid + "\n\n",
	} {
		t.Run(name, func(t *testing.T) {
			slide, err := NewSlide(body, config.Properties{})
			if err != nil {
				t.Fatalf("NewSlide() failed: %v", err)
			}
			slide.Style = config.StyleConfig{}.Apply(width, height)

			lines := strings.Split(ansi.Strip(slide.View(false, width, height)), "\n")
			if len(lines) != height {
				t.Fatalf("slide is %d lines, want %d:\n%s", len(lines), height, strings.Join(lines, "\n"))
			}
			if !strings.Contains(lines[height-2], "╰") {
				t.Errorf("grid ends above the slide's last row:\n%s", strings.Join(lines, "\n"))
			}
		})
	}
}

// A slide keeps what it draws for animations between frames, but has to draw
// itself again when the terminal changes size under it.
func TestSlide_AnimatedViewFollowsTheSize(t *testing.T) {
	slide, err := NewSlide("# Title\n\nsome text", config.Properties{})
	if err != nil {
		t.Fatalf("NewSlide() failed: %v", err)
	}

	for _, size := range [][2]int{{80, 24}, {80, 24}, {100, 30}, {80, 24}} {
		width, height := size[0], size[1]
		slide.Style = config.StyleConfig{}.Apply(width, height)

		lines := strings.Split(slide.View(true, width, height), "\n")
		if len(lines) != height || ansi.StringWidth(lines[0]) != width {
			t.Errorf("at %dx%d the slide is %dx%d",
				width, height, ansi.StringWidth(lines[0]), len(lines))
		}
	}
}

// Slides are rendered ahead of time in the background while the presentation
// keeps showing them, and have to come out the same either way.
func TestSlide_PrerenderAlongsideView(t *testing.T) {
	const width, height = 80, 24

	var slides []*Slide
	for _, body := range []string{"# One\n\nfirst", "# Two\n\n- a\n- b", "# Three\n\n```go\nx := 1\n```"} {
		slide, err := NewSlide(body, config.Properties{})
		if err != nil {
			t.Fatalf("NewSlide() failed: %v", err)
		}
		slide.Style = config.StyleConfig{}.Apply(width, height)
		slides = append(slides, slide)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for _, s := range slides {
			s.Prerender(width, height)
		}
	}()
	for range 20 {
		for _, s := range slides {
			s.View(true, width, height)
			s.View(false, width, height)
		}
	}
	<-done

	for i, s := range slides {
		fresh, err := NewSlide(s.Data, config.Properties{})
		if err != nil {
			t.Fatalf("NewSlide() failed: %v", err)
		}
		fresh.Style = s.Style
		for _, animating := range []bool{true, false} {
			if got, want := s.View(animating, width, height), fresh.View(animating, width, height); got != want {
				t.Errorf("slide %d, animating %v: the prerendered view differs from a fresh one", i, animating)
			}
		}
	}
}

// TestSlide_MasterLayout runs a slide through the templating layer the way the
// presentation loader does.
func TestSlide_MasterLayout(t *testing.T) {
	config.GlobalConfig.Masters = map[string]string{
		"two-col": "[row]\n[col span=2]\n[slot main]\n[/col]\n[col]\n[slot side]\n[/col]\n[/row]",
	}
	t.Cleanup(func() { config.GlobalConfig.Masters = nil })

	slide, err := NewSlide(
		"[slot main]\n# Headline\n[/slot]\n[slot side]\n- a note\n[/slot]",
		config.Properties{Master: "two-col"},
	)
	if err != nil {
		t.Fatalf("NewSlide() failed: %v", err)
	}
	slide.Style = config.StyleConfig{}.Apply(80, 24)

	out := ansi.Strip(slide.View(false, 80, 24))
	for _, want := range []string{"Headline", "a note"} {
		if !strings.Contains(out, want) {
			t.Errorf("rendered slide is missing %q:\n%s", want, out)
		}
	}

	// The layout put them side by side, so both land on the same line.
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "Headline") && strings.Contains(line, "a note") {
			return
		}
	}
	t.Errorf("the slots were not laid out side by side:\n%s", out)
}

func TestSlide_UnknownMasterLayout(t *testing.T) {
	if _, err := NewSlide("# Title", config.Properties{Master: "nope"}); err == nil {
		t.Error("NewSlide() accepted a layout that does not exist")
	}
}
