package placeholder

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"strings"
	"testing"
	"unicode"

	"github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"
	reflow "github.com/muesli/reflow/ansi"
)

// enable draws images with placeholders for the rest of a test, uploading them
// to the returned buffer.
func enable(t *testing.T) *bytes.Buffer {
	t.Helper()
	t.Setenv(envOverride, "1")

	var out bytes.Buffer
	SetOutput(&out)
	t.Cleanup(func() { SetOutput(nil) })
	return &out
}

// picture is an image cols by rows cells large, at 2 by 4 pixels a cell, with
// every pixel its own, so that no two parts of it look alike.
func picture(cols, rows int, seed byte) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, 2*cols, 4*rows))
	for i := range img.Pix {
		img.Pix[i] = byte(i) ^ byte(i>>8) ^ seed
	}
	return img
}

// arrival is an image as the terminal receives it.
type arrival struct {
	keys   map[string]string
	pixels []byte
}

// received reads the images out of what was uploaded.
func received(t *testing.T, out string) []arrival {
	t.Helper()

	var (
		images []arrival
		data   strings.Builder
	)
	for code := range strings.SplitSeq(strings.TrimSuffix(out, "\x1b\\"), "\x1b\\") {
		code, ok := strings.CutPrefix(code, "\x1b_G")
		if !ok {
			t.Fatalf("%q is not a graphics escape code", code)
		}
		keys, payload, _ := strings.Cut(code, ";")
		if len(payload) > chunk {
			t.Errorf("an escape code carries %d bytes, more than kitty takes", len(payload))
		}

		kv := map[string]string{}
		for pair := range strings.SplitSeq(keys, ",") {
			k, v, _ := strings.Cut(pair, "=")
			kv[k] = v
		}
		if _, ok := kv["a"]; ok {
			images = append(images, arrival{keys: kv})
		} else if len(kv) != 1 {
			t.Errorf("a chunk after the first has keys %q, want m alone", keys)
		}
		data.WriteString(payload)

		if kv["m"] == "0" {
			pixels, err := base64.StdEncoding.DecodeString(data.String())
			if err != nil {
				t.Fatalf("uploaded pixels are not base64: %v", err)
			}
			images[len(images)-1].pixels = pixels
			data.Reset()
		}
	}
	return images
}

// pixelsOf returns the pixels of part of an image, row after row.
func pixelsOf(img *image.NRGBA, r image.Rectangle) []byte {
	var out []byte
	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := img.PixOffset(r.Min.X, y)
		out = append(out, img.Pix[i:i+4*r.Dx()]...)
	}
	return out
}

func TestDraw_UploadsAnImageOnce(t *testing.T) {
	out := enable(t)
	img := picture(2, 1, 0)

	first, err := Draw(img, 2, 1)
	if err != nil {
		t.Fatalf("Draw() error = %v", err)
	}
	want := "\x1b_Ga=T,U=1,q=2,f=32,s=4,v=4,c=2,r=1,i=1,m=0;" +
		base64.StdEncoding.EncodeToString(img.Pix) + "\x1b\\"
	if got := out.String(); got != want {
		t.Fatalf("uploaded\n%q\nwant\n%q", got, want)
	}

	out.Reset()
	again, err := Draw(img, 2, 1)
	if err != nil {
		t.Fatalf("Draw() error = %v", err)
	}
	if out.Len() > 0 {
		t.Errorf("the same image went up again: %q", out.String())
	}
	if again != first {
		t.Errorf("the same image came back as\n%q\nthen as\n%q", first, again)
	}

	other, err := Draw(picture(2, 1, 1), 2, 1)
	if err != nil {
		t.Fatalf("Draw() error = %v", err)
	}
	if images := received(t, out.String()); len(images) != 1 || images[0].keys["i"] != "2" {
		t.Errorf("another image did not go up under an id of its own: %q", out.String())
	}
	if other == first {
		t.Errorf("another image came back as the same cells: %q", other)
	}
}

func TestDraw_UploadsAgainToANewOutput(t *testing.T) {
	enable(t)
	if _, err := Draw(picture(2, 1, 0), 2, 1); err != nil {
		t.Fatalf("Draw() error = %v", err)
	}

	out := enable(t)
	if _, err := Draw(picture(2, 1, 0), 2, 1); err != nil {
		t.Fatalf("Draw() error = %v", err)
	}
	if out.Len() == 0 {
		t.Error("the image did not go up to the new output")
	}
}

// An image wider or taller than the marks can number goes up in tiles, every
// one of which its cells show pixel for pixel.
func TestDraw_TilesLargeImages(t *testing.T) {
	out := enable(t)
	const cols, rows = tile + 5, tile + 2
	img := picture(cols, rows, 0)

	cells, err := Draw(img, cols, rows)
	if err != nil {
		t.Fatalf("Draw() error = %v", err)
	}

	images := received(t, out.String())
	tiles := []image.Rectangle{
		image.Rect(0, 0, tile, tile), image.Rect(tile, 0, cols, tile),
		image.Rect(0, tile, tile, rows), image.Rect(tile, tile, cols, rows),
	}
	if len(images) != len(tiles) {
		t.Fatalf("uploaded %d images, want %d tiles", len(images), len(tiles))
	}
	ids := map[string]image.Rectangle{}
	for i, r := range tiles {
		got := images[i]
		want := map[string]string{
			"c": fmt.Sprint(r.Dx()), "r": fmt.Sprint(r.Dy()),
			"s": fmt.Sprint(2 * r.Dx()), "v": fmt.Sprint(4 * r.Dy()),
		}
		for k, v := range want {
			if got.keys[k] != v {
				t.Errorf("tile %v has %s=%s, want %s", r, k, got.keys[k], v)
			}
		}
		px := image.Rect(2*r.Min.X, 4*r.Min.Y, 2*r.Max.X, 4*r.Max.Y)
		if !bytes.Equal(got.pixels, pixelsOf(img, px)) {
			t.Errorf("tile %v does not hold its part of the image", r)
		}
		ids[got.keys["i"]] = r
	}

	lines := strings.Split(cells, "\n")
	if len(lines) != rows {
		t.Fatalf("Draw() returned %d rows, want %d", len(lines), rows)
	}
	for y, line := range lines {
		x := 0
		for _, part := range strings.Split(line, "\x1b[38;2;")[1:] {
			color, rest, _ := strings.Cut(part, "m")
			var r, g, b int
			fmt.Sscanf(color, "%d;%d;%d", &r, &g, &b)
			at, ok := ids[fmt.Sprint(r<<16|g<<8|b)]
			if !ok {
				t.Fatalf("row %d shows image %s, which never went up", y, color)
			}
			for c := range strings.SplitSeq(strings.TrimSuffix(rest, "\x1b[39m"), cell) {
				if c == "" {
					continue
				}
				want := string(diacritics[y-at.Min.Y]) + string(diacritics[x-at.Min.X])
				if !image.Pt(x, y).In(at) || c != want {
					t.Fatalf("cell %d,%d shows %q of tile %v", x, y, c, at)
				}
				x++
			}
		}
		if x != cols {
			t.Errorf("row %d has %d cells, want %d", y, x, cols)
		}
	}
}

func TestDraw_Refuses(t *testing.T) {
	tests := []struct {
		name       string
		enabled    bool
		img        *image.NRGBA
		cols, rows int
	}{
		{name: "while disabled", img: picture(2, 1, 0), cols: 2, rows: 1},
		{name: "fewer pixels than cells", enabled: true, img: image.NewNRGBA(image.Rect(0, 0, 1, 1)), cols: 2, rows: 1},
		{name: "no cells", enabled: true, img: picture(2, 1, 0), cols: 0, rows: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out *bytes.Buffer
			if tt.enabled {
				out = enable(t)
			}
			if got, err := Draw(tt.img, tt.cols, tt.rows); err == nil {
				t.Errorf("Draw() = %q, want an error", got)
			}
			if out != nil && out.Len() > 0 {
				t.Errorf("uploaded %q", out.String())
			}
		})
	}
}

// Kitty takes no more than 4096 bytes of base64 in one escape code.
func TestTransmission_SendsLargeImagesInChunks(t *testing.T) {
	pixels := picture(30, 30, 0).Pix
	images := received(t, transmission(7, "a=T,f=32", pixels))
	if len(images) != 1 {
		t.Fatalf("received %d images, want 1", len(images))
	}
	if images[0].keys["i"] != "7" {
		t.Errorf("the image went up as %q, want 7", images[0].keys["i"])
	}
	if !bytes.Equal(images[0].pixels, pixels) {
		t.Error("the pixels did not arrive whole")
	}
}

// Layout measures with x/ansi and transitions with reflow, so both have to see
// one cell for each placeholder, whatever marks number it.
func TestText_IsOneCellAPlaceholder(t *testing.T) {
	for _, mark := range diacritics {
		if !unicode.Is(unicode.Mn, mark) || runewidth.RuneWidth(mark) != 0 {
			t.Errorf("%U does not combine with the cell before it", mark)
		}
	}

	ids := []uint32{1, 2, 3, 4}
	const cols, rows = 2*tile - 1, 2 * tile
	for y, line := range strings.Split(text(ids, 2, cols, rows), "\n") {
		if got := ansi.StringWidth(line); got != cols {
			t.Fatalf("row %d is %d cells wide to x/ansi, want %d", y, got, cols)
		}
		if got := reflow.PrintableRuneWidth(line); got != cols {
			t.Fatalf("row %d is %d cells wide to reflow, want %d", y, got, cols)
		}
	}
}

func TestIs(t *testing.T) {
	if !Is(cell + string(diacritics[0]) + string(diacritics[1])) {
		t.Error("a placeholder cell is not one")
	}
	if Is("a") || Is("") {
		t.Error("text is a placeholder cell")
	}
}

func TestClear(t *testing.T) {
	if got := Clear(); got != unplace {
		t.Errorf("Clear() = %q with placeholders disabled, want %q", got, unplace)
	}
	enable(t)
	if got := Clear(); got != "" {
		t.Errorf("Clear() = %q with placeholders enabled, want nothing", got)
	}
}

func TestSetOutput_OnlyEnablesWhereSupported(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		supports bool
	}{
		{name: "kitty", env: map[string]string{"TERM": "xterm-kitty"}, supports: true},
		{name: "ghostty", env: map[string]string{"TERM": "xterm-ghostty"}, supports: true},
		{name: "kitty in tmux", env: map[string]string{"TERM": "xterm-kitty", "TMUX": "/tmp/tmux"}},
		{name: "other terminals", env: map[string]string{"TERM": "xterm-256color"}},
		{name: "forced on", env: map[string]string{"TERM": "xterm-256color", envOverride: "1"}, supports: true},
		{name: "forced off", env: map[string]string{"TERM": "xterm-kitty", envOverride: "0"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, k := range []string{"TERM", "TMUX", envOverride} {
				t.Setenv(k, tt.env[k])
			}
			SetOutput(&bytes.Buffer{})
			t.Cleanup(func() { SetOutput(nil) })

			if got := Enabled(); got != tt.supports {
				t.Errorf("Enabled() = %v, want %v", got, tt.supports)
			}
		})
	}
}

func TestKeep_KeepsImagesForTheirCells(t *testing.T) {
	t.Setenv(envOverride, "0")
	Keep()
	t.Cleanup(func() { SetOutput(nil) })

	if !Enabled() || !Keeping() {
		t.Fatalf("Enabled() = %v, Keeping() = %v while keeping images, want both", Enabled(), Keeping())
	}

	img := picture(3, 2, 0)
	cells, err := Draw(img, 3, 2)
	if err != nil {
		t.Fatalf("Draw() error = %v", err)
	}

	kept, cols, rows, ok := Image(1)
	if !ok || cols != 3 || rows != 2 {
		t.Fatalf("Image(1) = %v, %d by %d cells, want the image at 3 by 2", ok, cols, rows)
	}
	if !bytes.Equal(pixelsOf(kept, kept.Bounds()), img.Pix) {
		t.Errorf("the image kept is not the one drawn")
	}
	if _, _, _, ok := Image(2); ok {
		t.Errorf("Image(2) found an image that was never drawn")
	}

	// Every cell reads back as the image and the cell of it it shows.
	for y, line := range strings.Split(cells, "\n") {
		x := 0
		for _, grapheme := range graphemes(line) {
			id, row, col, ok := At(grapheme, color.RGBA{0, 0, 1, 0xff})
			if !ok || id != 1 || row != y || col != x {
				t.Errorf("cell %d, %d reads as image %d at %d, %d (%v)", x, y, id, col, row, ok)
			}
			x++
		}
		if x != 3 {
			t.Errorf("row %d has %d cells, want 3", y, x)
		}
	}

	SetOutput(nil)
	if Keeping() {
		t.Errorf("still keeping images after SetOutput")
	}
}

func TestAt_RejectsOtherText(t *testing.T) {
	for _, grapheme := range []string{"a", cell, cell + "̅"} {
		if _, _, _, ok := At(grapheme, color.White); ok {
			t.Errorf("At(%q) read it as a placeholder cell", grapheme)
		}
	}
	if _, _, _, ok := At(cell+"̅̅", nil); ok {
		t.Errorf("At() read a cell with no image id as a placeholder cell")
	}
}

// graphemes returns the graphemes of a line, leaving out escape codes.
func graphemes(line string) []string {
	var out []string
	var state byte
	for line != "" {
		seq, width, n, newState := ansi.DecodeSequence(line, state, nil)
		state, line = newState, line[n:]
		if width > 0 {
			out = append(out, seq)
		}
	}
	return out
}
