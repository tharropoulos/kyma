package transitions

import (
	"strings"
	"testing"
)

// TestTransitions_KeepTheSlideHeight checks every frame of every transition
// against the height of the slides it moves between. A taller frame does not
// fit the terminal, which the app reports instead of drawing the slide.
func TestTransitions_KeepTheSlideHeight(t *testing.T) {
	const width, height = 40, 12

	slide := func(fill string) string {
		lines := make([]string, height)
		for i := range lines {
			lines[i] = strings.Repeat(fill, width)
		}
		return strings.Join(lines, "\n")
	}
	prev, next := slide("a"), slide("b")

	names := []string{"slideUp", "slideDown", "swipeLeft", "swipeRight", "flip", "collapse", "expand", "fade"}
	for _, name := range names {
		for _, dir := range []direction{Forwards, Backwards} {
			t.Run(name, func(t *testing.T) {
				tr := Get(name, Fps).Start(width, height, dir)
				for frame := 0; tr.Animating() && frame < 600; frame++ {
					if got := len(strings.Split(tr.View(prev, next), "\n")); got != height {
						t.Fatalf("frame %d is %d lines tall, want %d", frame, got, height)
					}
					tr, _ = tr.Update()
				}
			})
		}
	}
}
