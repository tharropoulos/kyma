package export

import (
	"bytes"
	"compress/zlib"
	"image"
	"image/color"
	"io"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/cellbuf"

	"github.com/museslabs/kyma/internal/img/placeholder"
	"github.com/museslabs/kyma/internal/textsize"
)

var (
	white = color.RGBA{0xff, 0xff, 0xff, 0xff}
	black = color.RGBA{0x00, 0x00, 0x00, 0xff}
	red   = color.RGBA{0xff, 0x00, 0x00, 0xff}
	blue  = color.RGBA{0x00, 0x00, 0xff, 0xff}
)

// terminal is a small screen with cells 10 by 20 pixels, white on black.
func terminal(cols, rows int) Terminal {
	return Terminal{Cols: cols, Rows: rows, CellWidth: 10, CellHeight: 20, Fg: white, Bg: black}
}

func paint(t *testing.T, term Terminal, frame string) *image.RGBA {
	t.Helper()
	p, err := newPainter(term)
	if err != nil {
		t.Fatalf("newPainter() error = %v", err)
	}
	return p.paint(lay(frame, term.Cols, term.Rows))
}

func at(img *image.RGBA, x, y int) color.RGBA {
	return img.RGBAAt(x, y)
}

// inked reports whether any pixel of r differs from c.
func inked(img *image.RGBA, r image.Rectangle, c color.RGBA) bool {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			if img.RGBAAt(x, y) != c {
				return true
			}
		}
	}
	return false
}

func TestLay(t *testing.T) {
	frame := "\x1b[1;31ma\x1b[0m界b\n" +
		"xy\x1b[2Dz\x1b[3Cw\n" +
		"\x1b[<66;0;0k  \x1b]66;s=2;Q\a"
	s := lay(frame, 8, 4)

	cellAt := func(x, y int) cell { return s.cells[y*s.cols+x] }

	if c := cellAt(0, 0); c.text != "a" || c.style.Attrs&cellbuf.BoldAttr == 0 || c.style.Fg != ansi.Red {
		t.Errorf("cell 0, 0 = %+v, want a bold red a", c)
	}
	if c := cellAt(1, 0); c.text != "界" || c.width != 2 || !c.style.Empty() {
		t.Errorf("cell 1, 0 = %+v, want 界 two cells wide in no style", c)
	}
	if c := cellAt(2, 0); c.text != "" {
		t.Errorf("the second cell of 界 holds %q", c.text)
	}
	if c := cellAt(3, 0); c.text != "b" {
		t.Errorf("cell 3, 0 = %q, want b after the wide 界", c.text)
	}

	// The cursor moves back over y and x and on past them.
	if got := cellAt(0, 1).text + cellAt(1, 1).text + cellAt(4, 1).text; got != "zyw" {
		t.Errorf("row 1 reads %q, want zyw", got)
	}

	if len(s.scaled) != 1 {
		t.Fatalf("found %d scaled graphemes, want 1", len(s.scaled))
	}
	if got := s.scaled[0]; got.x != 2 || got.y != 2 || got.text != "Q" || got.scale != 2 || got.width != 1 {
		t.Errorf("scaled text = %+v, want Q at 2, 2, twice its size", got)
	}
}

func TestParseScaled(t *testing.T) {
	for _, tt := range []struct {
		seq          string
		text         string
		scale, width int
		ok           bool
	}{
		{"\x1b]66;s=3;A\a", "A", 3, 1, true},
		{"\x1b]66;s=2:w=2;界\x1b\\", "界", 2, 2, true},
		{"\x1b]66;w=1;x\a", "x", 1, 1, true},
		{"\x1b]66;s=9;A\a", "A", 7, 1, true},
		{"\x1b]8;;https://example.com\a", "", 0, 0, false},
	} {
		got, ok := parseScaled(tt.seq)
		if ok != tt.ok || ok && (got.text != tt.text || got.scale != tt.scale || got.width != tt.width) {
			t.Errorf("parseScaled(%q) = %+v, %v, want %q at %d, %d wide", tt.seq, got, ok, tt.text, tt.scale, tt.width)
		}
	}
}

func TestPaint_Colors(t *testing.T) {
	term := terminal(4, 1)
	term.Palette[1] = color.RGBA{0x12, 0x34, 0x56, 0xff}

	img := paint(t, term, "\x1b[48;2;0;0;255m \x1b[0m \x1b[31m█\x1b[0m\x1b[7m \x1b[0m")

	if got := at(img, 5, 10); got != blue {
		t.Errorf("a cell with its own background is %v, want blue", got)
	}
	if got := at(img, 15, 10); got != black {
		t.Errorf("a blank cell is %v, want the terminal's background", got)
	}
	if got := at(img, 25, 10); got != term.Palette[1] {
		t.Errorf("a block in palette color 1 is %v, want the terminal's %v", got, term.Palette[1])
	}
	if got := at(img, 35, 10); got != white {
		t.Errorf("a reversed blank is %v, want the terminal's foreground", got)
	}
}

func TestPaint_DrawsText(t *testing.T) {
	img := paint(t, terminal(3, 1), "A \x1b[4m \x1b[0m")

	if !inked(img, image.Rect(0, 0, 10, 20), black) {
		t.Errorf("A left no ink")
	}
	if inked(img, image.Rect(10, 0, 20, 20), black) {
		t.Errorf("a space left ink")
	}
	if !inked(img, image.Rect(20, 0, 30, 20), black) {
		t.Errorf("an underlined space left no underline")
	}
}

func TestPaint_BoxLinesMeetAcrossCells(t *testing.T) {
	img := paint(t, terminal(3, 3), "╭─╮\n│ │\n╰─╯")
	p, _ := newPainter(terminal(3, 3))

	// Down the middle of the left column, the line runs unbroken from where
	// the top corner turns down to where the bottom one turns away, the
	// corners' pixels only shaded off where the curve starts.
	x0, x1 := down(image.Rect(0, 0, 10, 20), p.line)
	for y := 13; y < 46; y++ {
		for x := x0; x < x1; x++ {
			if got := at(img, x, y); got.R < 0xc0 {
				t.Fatalf("the left side is broken at %d, %d: %v", x, y, got)
			}
		}
	}
	// And across the top, from the left corner to the right one.
	y0, y1 := across(image.Rect(0, 0, 10, 20), p.line)
	for x := 8; x < 21; x++ {
		for y := y0; y < y1; y++ {
			if got := at(img, x, y); got.R < 0xc0 {
				t.Fatalf("the top is broken at %d, %d: %v", x, y, got)
			}
		}
	}
	// The corners round off: the outer corner of each cell stays clear.
	if got := at(img, 0, 0); got != black {
		t.Errorf("the corner of ╭ is %v, want it left blank", got)
	}
}

func TestPaint_Blocks(t *testing.T) {
	img := paint(t, terminal(3, 1), "▀▐▚")

	for _, tt := range []struct {
		x, y int
		want color.RGBA
	}{
		{5, 5, white}, {5, 15, black}, // ▀
		{12, 10, black}, {17, 10, white}, // ▐
		{22, 5, white}, {27, 5, black}, {22, 15, black}, {27, 15, white}, // ▚
	} {
		if got := at(img, tt.x, tt.y); got != tt.want {
			t.Errorf("pixel %d, %d = %v, want %v", tt.x, tt.y, got, tt.want)
		}
	}
}

func TestPaint_Images(t *testing.T) {
	placeholder.Keep()
	t.Cleanup(func() { placeholder.SetOutput(nil) })

	// An image 2 by 1 cells large, at 3 by 5 pixels a cell, its left half
	// red and its right half blue.
	src := image.NewNRGBA(image.Rect(0, 0, 6, 5))
	for y := range 5 {
		for x := range 6 {
			c := color.NRGBA{0xff, 0, 0, 0xff}
			if x >= 3 {
				c = color.NRGBA{0, 0, 0xff, 0xff}
			}
			src.SetNRGBA(x, y, c)
		}
	}
	cells, err := placeholder.Draw(src, 2, 1)
	if err != nil {
		t.Fatalf("Draw() error = %v", err)
	}

	img := paint(t, terminal(3, 1), " "+cells)
	if got := at(img, 5, 10); got != black {
		t.Errorf("the cell before the image is %v", got)
	}
	if got := at(img, 15, 10); got != red {
		t.Errorf("the image's first cell is %v, want red", got)
	}
	if got := at(img, 25, 10); got != blue {
		t.Errorf("the image's second cell is %v, want blue", got)
	}
}

func TestPaint_ScaledText(t *testing.T) {
	block := textsize.Block("\x1b[1mH\x1b[0m", 2)
	frame := textsize.Render(" "+block[0]+" \n "+block[1]+" ", 4, 2)

	img := paint(t, terminal(4, 2), frame)

	// The glyph reaches down into the second row of its block, which text
	// at its own size never does.
	if !inked(img, image.Rect(10, 20, 30, 40), black) {
		t.Errorf("the scaled H left no ink in the lower half of its block")
	}
	if inked(img, image.Rect(0, 0, 10, 40), black) || inked(img, image.Rect(30, 0, 40, 40), black) {
		t.Errorf("the scaled H left ink outside its block")
	}
}

func TestPaint_Wallpaper(t *testing.T) {
	wallpaper := image.NewRGBA(image.Rect(0, 0, 400, 100))
	for i := 0; i < len(wallpaper.Pix); i += 4 {
		copy(wallpaper.Pix[i:], []byte{0, 0, 0xff, 0xff})
	}

	term := terminal(2, 1)
	term.Wallpaper, term.Opacity = wallpaper, 0.5
	p, err := newPainter(term)
	if err != nil {
		t.Fatalf("newPainter() error = %v", err)
	}
	img := p.paint(lay(" \x1b[41m \x1b[0m", 2, 1))

	if got := at(img, 5, 10); got.A < 0x7f || got.A > 0x80 {
		t.Errorf("the terminal's background is %v, want it half opaque", got)
	}
	if got := at(img, 15, 10); got.A != 0xff || got.R == 0 {
		t.Errorf("a cell with its own background is %v, want it solid red", got)
	}
	if got := p.wallpaper(); got == nil || got.RGBAAt(5, 10) != blue {
		t.Errorf("the wallpaper doesn't show behind the frame")
	}

	term.Opacity = 1
	p, _ = newPainter(term)
	if p.wallpaper() != nil {
		t.Errorf("an opaque terminal shows its wallpaper")
	}
}

func TestPDF(t *testing.T) {
	var out bytes.Buffer
	doc := newPDF(&out)

	// An opaque page, and one where half of it lets a wallpaper through.
	opaque, translucent := image.NewRGBA(image.Rect(0, 0, 7, 3)), image.NewRGBA(image.Rect(0, 0, 4, 4))
	for i := range opaque.Pix {
		opaque.Pix[i] = byte(i * 31)
		if i%4 == 3 {
			opaque.Pix[i] = 0xff
		}
	}
	for i := 0; i < len(translucent.Pix); i += 4 {
		a := byte(0xff)
		if i%8 == 0 {
			a = 0x80
		}
		copy(translucent.Pix[i:], []byte{a / 2, a / 4, 0, a})
	}

	wallpaper, err := doc.photo(image.NewRGBA(image.Rect(0, 0, 4, 4)))
	if err != nil {
		t.Fatalf("photo() error = %v", err)
	}
	if err := doc.page(opaque, 0, 960, 540); err != nil {
		t.Fatalf("page() error = %v", err)
	}
	if err := doc.page(translucent, wallpaper, 960, 960); err != nil {
		t.Fatalf("page() error = %v", err)
	}
	if err := doc.close(); err != nil {
		t.Fatalf("close() error = %v", err)
	}
	file := out.String()

	if !strings.HasPrefix(file, "%PDF-1.4\n") || !strings.HasSuffix(file, "%%EOF\n") {
		t.Fatalf("not a PDF:\n%q", file)
	}
	if !strings.Contains(file, "/Count 2") {
		t.Errorf("the page tree doesn't count 2 pages")
	}
	if n := strings.Count(file, "/SMask"); n != 1 {
		t.Errorf("%d images have a soft mask, want only the translucent page's", n)
	}
	if !strings.Contains(file, "/Bg Do Q q 960 0 0 960 0 0 cm /Im Do Q") {
		t.Errorf("the translucent page isn't drawn over the wallpaper")
	}

	// Every offset in the table of objects lands on the object it is for.
	start, err := strconv.Atoi(regexp.MustCompile(`startxref\n(\d+)\n`).FindStringSubmatch(file)[1])
	if err != nil || !strings.HasPrefix(file[start:], "xref\n") {
		t.Fatalf("startxref doesn't point at the table of objects")
	}
	entries := regexp.MustCompile(`(\d{10}) 00000 n \n`).FindAllStringSubmatch(file[start:], -1)
	for i, e := range entries {
		offset, _ := strconv.Atoi(e[1])
		if want := strconv.Itoa(i+1) + " 0 obj\n"; !strings.HasPrefix(file[offset:], want) {
			t.Errorf("object %d is listed at %d, which holds %q", i+1, offset, file[offset:min(offset+10, len(file))])
		}
	}

	// The images come back out pixel for pixel: the opaque page's colors,
	// then the translucent page's alpha and its colors, no longer multiplied
	// by it.
	want := [][]byte{nil, nil, nil}
	for i := 0; i < len(opaque.Pix); i += 4 {
		want[0] = append(want[0], opaque.Pix[i:i+3]...)
	}
	for i := 0; i < len(translucent.Pix); i += 4 {
		a := translucent.Pix[i+3]
		want[1] = append(want[1], a)
		// 0x40 and 0x20 at half alpha, and 0x7f and 0x3f at full.
		if a == 0x80 {
			want[2] = append(want[2], 0x80, 0x40, 0)
		} else {
			want[2] = append(want[2], 0x7f, 0x3f, 0)
		}
	}
	streams := regexp.MustCompile(`(?s)/Colors (\d) /BitsPerComponent 8 /Columns (\d+) >>(?: /SMask \d+ 0 R)? /Length (\d+) >>\nstream\n`).FindAllStringSubmatchIndex(file, -1)
	if len(streams) != len(want) {
		t.Fatalf("found %d deflated images, want %d", len(streams), len(want))
	}
	for i, m := range streams {
		channels, _ := strconv.Atoi(file[m[2]:m[3]])
		width, _ := strconv.Atoi(file[m[4]:m[5]])
		length, _ := strconv.Atoi(file[m[6]:m[7]])
		if got := unfilter(t, []byte(file[m[1]:m[1]+length]), width, channels); !bytes.Equal(got, want[i]) {
			t.Errorf("image %d came back as\n%v\nwant\n%v", i+1, got, want[i])
		}
	}
}

// unfilter inflates an image's pixels and undoes PNG's filters on each row.
func unfilter(t *testing.T, data []byte, width, channels int) []byte {
	t.Helper()
	z, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("not deflated: %v", err)
	}
	raw, err := io.ReadAll(z)
	if err != nil {
		t.Fatalf("not deflated: %v", err)
	}

	stride := channels * width
	var out, prev []byte
	prev = make([]byte, stride)
	for len(raw) > 0 {
		if raw[0] != 4 {
			t.Fatalf("a row is filtered with %d, want Paeth", raw[0])
		}
		row := make([]byte, stride)
		for i := range row {
			var left, upLeft byte
			if i >= channels {
				left, upLeft = row[i-channels], prev[i-channels]
			}
			row[i] = raw[1+i] + paeth(left, prev[i], upLeft)
		}
		out = append(out, row...)
		prev, raw = row, raw[1+stride:]
	}
	return out
}

func TestParseColors(t *testing.T) {
	answers := "\x1b]10;rgb:ffff/ffff/ffff\x1b\\" +
		"\x1b]11;rgb:08/00/1c\a" +
		"\x1b]4;1;rgb:cdcd/0000/0000\x1b\\" +
		"\x1b]4;15;rgb:f/8/0\x1b\\" +
		"\x1b[?62;22c"

	c := parseColors(answers)
	if c.Fg != white {
		t.Errorf("Fg = %v, want white", c.Fg)
	}
	if want := (color.RGBA{0x08, 0x00, 0x1c, 0xff}); c.Bg != want {
		t.Errorf("Bg = %v, want %v", c.Bg, want)
	}
	if want := (color.RGBA{0xcd, 0x00, 0x00, 0xff}); c.Palette[1] != want {
		t.Errorf("Palette[1] = %v, want %v", c.Palette[1], want)
	}
	if want := (color.RGBA{0xff, 0x88, 0x00, 0xff}); c.Palette[15] != want {
		t.Errorf("Palette[15] = %v, want %v", c.Palette[15], want)
	}
	if c.Palette[0] != nil {
		t.Errorf("Palette[0] = %v, want nil as the terminal didn't say", c.Palette[0])
	}
}
