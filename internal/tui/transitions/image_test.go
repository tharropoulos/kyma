package transitions

import (
	"bytes"
	"fmt"
	"image"
	"maps"
	"strings"
	"testing"
	"unicode"

	"github.com/charmbracelet/x/cellbuf"

	"github.com/museslabs/kyma/internal/img/placeholder"
)

// drawImage draws an image cols by rows cells large with placeholders, the
// way a slide gets it, with every cell of it and every seed looking different.
func drawImage(t *testing.T, cols, rows int, seed byte) []string {
	t.Helper()

	if !placeholder.Enabled() {
		t.Setenv("KYMA_IMAGE_PLACEHOLDERS", "1")
		placeholder.SetOutput(&bytes.Buffer{})
		t.Cleanup(func() { placeholder.SetOutput(nil) })
	}

	img := image.NewNRGBA(image.Rect(0, 0, 2*cols, 4*rows))
	for i := range img.Pix {
		img.Pix[i] = byte(i) ^ byte(i>>8) ^ seed
	}
	out, err := placeholder.Draw(img, cols, rows)
	if err != nil {
		t.Fatalf("placeholder.Draw() error = %v", err)
	}
	return strings.Split(out, "\n")
}

// imageSlide lays a slide out with a title on its first line and an image with
// its top left corner at x, y.
func imageSlide(width, height int, title string, img []string, x, y int) string {
	lines := make([]string, height)
	lines[0] = title
	for i, row := range img {
		lines[y+i] = strings.Repeat(" ", x) + row
	}
	return slide(width, lines...)
}

// An imageCell is a placeholder cell as the terminal reads it: the cell with
// its marks, and the color that says which image it shows.
type imageCell struct {
	text, fg string
}

// imageCells returns the placeholder cells of a frame by where they are. It
// fails the test on any cell that the terminal would draw wrong: one that lost
// a mark or got a stray one, or a mark that came apart from its cell.
func imageCells(t *testing.T, frame string, width, height int) map[image.Point]imageCell {
	t.Helper()

	buf := cellbuf.NewBuffer(width, height)
	cellbuf.SetContent(buf, frame)

	out := map[image.Point]imageCell{}
	for y := range height {
		for x := range width {
			c := buf.Cell(x, y)
			if c == nil || c.Width == 0 {
				continue
			}
			text := c.String()
			switch {
			case placeholder.Is(text):
				if len(c.Comb) != 2 {
					t.Fatalf("the placeholder cell at %d,%d has %d marks, want a row and a column: %q", x, y, len(c.Comb), text)
				}
				out[image.Pt(x, y)] = imageCell{text, fmt.Sprint(c.Style.Fg)}
			case len(c.Comb) > 0 || unicode.IsMark(c.Rune):
				t.Fatalf("the cell at %d,%d carries marks of a placeholder: %q", x, y, text)
			}
		}
	}
	return out
}

// Transitions cut lines apart and put them back together, and every cell an
// image is drawn with has to come through with its marks.
func TestTransitions_KeepImageCellsWhole(t *testing.T) {
	const width, height = 60, 10
	// Wider than a tile, so its rows change color part of the way along.
	prev := imageSlide(width, height, "Before", drawImage(t, 40, 3, 0), 4, 2)
	next := imageSlide(width, height, "After", drawImage(t, 40, 3, 1), 12, 5)

	names := []string{"slideUp", "slideDown", "swipeLeft", "swipeRight", "flip", "collapse", "expand", "fade"}
	for _, name := range names {
		for _, dir := range []direction{Forwards, Backwards} {
			t.Run(name, func(t *testing.T) {
				tr := Get(name, Fps).Start(width, height, dir)
				for frame := 0; tr.Animating() && frame < 600; frame++ {
					imageCells(t, tr.View(prev, next), width, height)
					tr, _ = tr.Update()
				}
			})
		}
	}

	for kind, s := range map[string]scene{
		"morph":     newScene(prev, next),
		"crossfade": newStillScene(prev, next),
	} {
		t.Run(kind, func(t *testing.T) {
			for step := range 21 {
				imageCells(t, s.frame(float64(step)/20), width, height)
			}
		})
	}
}

// An image both slides have glides from one place to the other in one piece.
func TestMorph_MovesImagesWhole(t *testing.T) {
	const width, height = 50, 10
	img := drawImage(t, 12, 3, 0)
	prev := imageSlide(width, height, "Title", img, 2, 1)
	next := imageSlide(width, height, "Title", img, 30, 6)

	from := imageCells(t, prev, width, height)
	corner := from[image.Pt(2, 1)]

	s := newScene(prev, next)
	var last image.Point
	for step := range 21 {
		p := float64(step) / 20
		got := imageCells(t, s.frame(p), width, height)

		var at image.Point
		found := false
		for pt, c := range got {
			if c == corner {
				at, found = pt, true
			}
		}
		if !found {
			t.Fatalf("at %.2f the image's corner is not drawn", p)
		}
		if len(got) != len(from) {
			t.Errorf("at %.2f the frame draws %d image cells, want %d", p, len(got), len(from))
		}
		shift := at.Sub(image.Pt(2, 1))
		for pt, c := range from {
			if got[pt.Add(shift)] != c {
				t.Fatalf("at %.2f the image does not hold together around %v", p, at)
			}
		}
		if step > 0 && (at.X < last.X || at.Y < last.Y) {
			t.Errorf("at %.2f the image went back from %v to %v", p, last, at)
		}
		last = at
	}
	if last != image.Pt(30, 6) {
		t.Errorf("the image ends up at %v, want 30,6", last)
	}
}

// One image giving way to another of the same size in the same place is not
// the same image recolored: its color is which image it is, and a color in
// between would show a third. Neither can fade either, so one swaps for the
// other whole.
func TestMorph_SwapsImagesWhole(t *testing.T) {
	const width, height = 30, 6
	prev := imageSlide(width, height, "Title", drawImage(t, 8, 3, 0), 4, 2)
	// Images uploaded in between leave colors between the two for a blend to
	// land on.
	for seed := range byte(4) {
		drawImage(t, 8, 3, 10+seed)
	}
	next := imageSlide(width, height, "Title", drawImage(t, 8, 3, 1), 4, 2)
	before, after := imageCells(t, prev, width, height), imageCells(t, next, width, height)

	for kind, s := range map[string]scene{
		"morph":     newScene(prev, next),
		"crossfade": newStillScene(prev, next),
	} {
		t.Run(kind, func(t *testing.T) {
			for step := range 21 {
				p := float64(step) / 20
				got := imageCells(t, s.frame(p), width, height)

				switch {
				case p < 0.5 && !maps.Equal(got, before):
					t.Errorf("at %.2f the frame does not draw the previous image whole", p)
				case p > 0.5 && !maps.Equal(got, after):
					t.Errorf("at %.2f the frame does not draw the next image whole", p)
				case !maps.Equal(got, before) && !maps.Equal(got, after):
					t.Errorf("at %.2f the frame draws neither image whole", p)
				}
			}
		})
	}
}
