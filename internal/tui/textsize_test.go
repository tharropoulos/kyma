package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"github.com/museslabs/kyma/internal/config"
	"github.com/museslabs/kyma/internal/textsize"
	"github.com/museslabs/kyma/internal/tui/transitions"
)

// forceTextSizing makes slides created during the test scale their headings,
// as they would in a terminal that supports it.
func forceTextSizing(t *testing.T) {
	t.Helper()

	// Cleanups run last in, first out, so this one runs once the variable is
	// back to how it was, and detection settles back to its usual answer.
	t.Cleanup(textsize.Detect)
	t.Setenv("KYMA_TEXT_SIZING", "1")
	textsize.Detect()
}

func scaledProps() config.Properties {
	return config.Properties{
		Style:      config.StyleConfig{FontSize: config.FontSize{1: 3, 2: 2}},
		Transition: transitions.Get("swipeLeft", transitions.Fps),
	}
}

// TestSlide_ScaledHeadingsFitTheFrame holds scaled headings to the same frame
// invariant as everything else. Before the final pass swaps the placeholders
// for escape codes, every layout step has to see them as plain text.
func TestSlide_ScaledHeadingsFitTheFrame(t *testing.T) {
	forceTextSizing(t)

	bodies := map[string]string{
		"plain":   "# Title\n\nsome text\n\n## Section\n\nmore",
		"long":    "# A heading long enough to wrap onto several lines\n\nbody",
		"columns": "[row]\n[col]\n## Left\n\nleft\n[/col]\n[col align=center]\n## Right\n\nright\n[/col]\n[/row]",
	}

	for _, size := range [][2]int{{80, 24}, {120, 40}, {60, 20}} {
		width, height := size[0], size[1]

		for name, body := range bodies {
			t.Run(name, func(t *testing.T) {
				slide, err := NewSlide(body, scaledProps())
				if err != nil {
					t.Fatalf("NewSlide() failed: %v", err)
				}
				slide.Style = config.StyleConfig{}.Apply(width, height)

				view := slide.View(false, width, height)
				if !strings.Contains(view, "\x1b_kyma;") {
					t.Fatalf("no heading was scaled:\n%s", ansi.Strip(view))
				}
				for i, line := range strings.Split(view, "\n") {
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

// TestModel_ViewDrawsScaledText checks the final pass: scaled text reaches the
// terminal only on frames nothing else will cut apart, and the markers never
// reach it at all.
func TestModel_ViewDrawsScaledText(t *testing.T) {
	forceTextSizing(t)

	first, err := NewSlide("# Title\n\nbody", scaledProps())
	if err != nil {
		t.Fatalf("NewSlide() failed: %v", err)
	}
	second, err := NewSlide("## Next\n\nbody", scaledProps())
	if err != nil {
		t.Fatalf("NewSlide() failed: %v", err)
	}
	first.Next, second.Prev = second, first

	m := model{
		slide:        first,
		rootSlide:    first,
		keys:         keys,
		width:        80,
		height:       24,
		timerDisplay: NewTimerDisplay(),
		globalTimer:  NewTimer(),
	}

	check := func(name string, wantScaled bool) {
		t.Helper()

		view := m.View()
		if strings.Contains(view, "\x1b_kyma") {
			t.Errorf("%s: a marker reached the terminal", name)
		}
		if got := strings.Contains(view, "\x1b]66;"); got != wantScaled {
			t.Errorf("%s: drew scaled text: %t, want %t", name, got, wantScaled)
		}
	}

	check("static slide", true)

	m.timerDisplay = m.timerDisplay.ToggleVisible()
	check("timer shown", false)
	m.timerDisplay = m.timerDisplay.ToggleVisible()

	m.navigateToSlide(second)
	m.slide.ActiveTransition = second.Properties.Transition.Start(m.width, m.height, transitions.Forwards)
	m.slide.ActiveTransition, _ = m.slide.ActiveTransition.Update()
	check("transition", false)

	m.slide.ActiveTransition = nil
	check("next slide", true)
}
