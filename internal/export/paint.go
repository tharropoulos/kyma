package export

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/cellbuf"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/vector"

	"github.com/museslabs/kyma/internal/img/placeholder"
)

// Terminal is what exported frames are drawn as: a screen of cells, in the
// colors and fonts of the terminal the deck is exported from.
type Terminal struct {
	// Cols and Rows are the size of the screen in cells, and CellWidth and
	// CellHeight the size of a cell in pixels.
	Cols, Rows            int
	CellWidth, CellHeight int

	// Fg and Bg are the colors text and the background are drawn in where a
	// cell doesn't set its own, and Palette the 16 colors that escape codes
	// name by number. A nil color is xterm's.
	Fg, Bg  color.Color
	Palette [16]color.Color

	// Fonts draw text, each one what those before it have no glyph for. Go
	// Mono and the system's fallback fonts come after them. Text is sized to
	// fit the first font into a cell.
	Fonts []Family

	// Wallpaper shows through a terminal drawn with an Opacity below 1, the
	// way it does behind a transparent terminal window. Only the background
	// cells don't set lets it through, as in kitty and Ghostty.
	Wallpaper image.Image
	Opacity   float64
}

// The colors a terminal draws in when it doesn't say, the same ones the morph
// transition assumes.
var (
	defaultFg = color.RGBA{0xbf, 0xbf, 0xbf, 0xff}
	defaultBg = color.RGBA{0x00, 0x00, 0x00, 0xff}
)

// painter draws screens into images the way the terminal draws them.
type painter struct {
	t        Terminal
	families []Family
	fg, bg   color.RGBA

	// ppem is the size text is drawn at, in pixels to the em, baseline how
	// far down a cell its baseline sits, and line how thick a light line of
	// a box is, all for text at its own size.
	ppem     float64
	baseline int
	line     int

	// base is the background of every frame.
	base *image.RGBA

	buf    sfnt.Buffer
	faces  map[faceKey]font.Face
	family map[rune]int
	images map[uint32]*image.RGBA
	raster vector.Rasterizer
}

type faceKey struct {
	font  *sfnt.Font
	scale int
}

func newPainter(t Terminal) (*painter, error) {
	if t.Cols <= 0 || t.Rows <= 0 || t.CellWidth <= 0 || t.CellHeight <= 0 {
		return nil, errors.New("the screen has no cells to draw")
	}

	p := &painter{
		t:        t,
		families: append(append([]Family(nil), t.Fonts...), fallbacks()...),
		fg:       defaultFg,
		bg:       defaultBg,
		faces:    map[faceKey]font.Face{},
		family:   map[rune]int{},
		images:   map[uint32]*image.RGBA{},
	}
	if t.Fg != nil {
		p.fg = rgba(t.Fg)
	}
	if t.Bg != nil {
		p.bg = rgba(t.Bg)
	}
	p.measure()
	p.base = p.background()
	return p, nil
}

// measure sizes text so that the first font's glyphs fill a cell across, the
// way a terminal sizes its cells to its font, and no taller than a cell, then
// sits its baseline where the font's ascent and descent are centered in it.
func (p *painter) measure() {
	f := p.families[0].Regular
	cw, ch := float64(p.t.CellWidth), float64(p.t.CellHeight)
	p.line = max(1, (p.t.CellWidth+4)/9)

	p.ppem = ch * 0.75
	const em = 1000
	if i, err := f.GlyphIndex(&p.buf, 'M'); err == nil && i != 0 {
		adv, err := f.GlyphAdvance(&p.buf, i, fixed.I(em), font.HintingNone)
		if err == nil && adv > 0 {
			p.ppem = cw * em / (float64(adv) / 64)
		}
	}

	m, err := f.Metrics(&p.buf, fixed.Int26_6(p.ppem*64), font.HintingNone)
	if err != nil {
		p.baseline = int(ch * 0.8)
		return
	}
	ascent, descent := float64(m.Ascent)/64, float64(m.Descent)/64
	if height := ascent + descent; height > ch {
		p.ppem *= ch / height
		ascent, descent = ascent*ch/height, descent*ch/height
	}
	p.baseline = int(math.Round((ch-ascent-descent)/2 + ascent))
}

// background draws what every frame starts out as: the terminal's background,
// only as opaque as the terminal draws it when a wallpaper shows through.
func (p *painter) background() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, p.t.Cols*p.t.CellWidth, p.t.Rows*p.t.CellHeight))
	bg := color.NRGBA{p.bg.R, p.bg.G, p.bg.B, 0xff}
	if p.translucent() {
		bg.A = uint8(math.Round(max(p.t.Opacity, 0) * 0xff))
	}
	draw.Draw(img, img.Bounds(), image.NewUniform(bg), image.Point{}, draw.Src)
	return img
}

// translucent reports whether the wallpaper shows through the background.
func (p *painter) translucent() bool {
	return p.t.Wallpaper != nil && p.t.Opacity < 1
}

// wallpaper returns the wallpaper as it shows behind the frames, or nil when
// it doesn't. It covers the screen with its middle, the way a desktop shows a
// wallpaper of other proportions than the screen.
func (p *painter) wallpaper() *image.RGBA {
	if !p.translucent() {
		return nil
	}
	img := image.NewRGBA(p.base.Bounds())
	src, w, h := p.t.Wallpaper.Bounds(), img.Bounds().Dx(), img.Bounds().Dy()
	crop := image.Rect(0, 0, src.Dx(), src.Dx()*h/w)
	if crop.Dy() > src.Dy() {
		crop = image.Rect(0, 0, src.Dy()*w/h, src.Dy())
	}
	crop = crop.Add(src.Min).Add(image.Pt((src.Dx()-crop.Dx())/2, (src.Dy()-crop.Dy())/2))
	xdraw.CatmullRom.Scale(img, img.Bounds(), p.t.Wallpaper, crop, xdraw.Src, nil)
	return img
}

// paint draws a screen. Backgrounds go first, so that the glyphs that reach
// past their cells, as some do in every font, aren't painted over by the next
// cell's background, then the cells' contents, then scaled text on top. Where
// the wallpaper shows through, the frame is only as opaque as the terminal's
// background, to be laid over the wallpaper.
func (p *painter) paint(s screen) *image.RGBA {
	img := image.NewRGBA(p.base.Bounds())
	copy(img.Pix, p.base.Pix)

	for i, c := range s.cells {
		if _, bg, solid := p.colors(c.style); solid {
			p.fill(img, p.rect(i%s.cols, i/s.cols, 1, 1), bg)
		}
	}

	for i, c := range s.cells {
		x, y := i%s.cols, i/s.cols
		fg, bg, _ := p.colors(c.style)
		if c.text != "" && c.style.Attrs&cellbuf.ConcealAttr == 0 {
			p.glyph(img, p.rect(x, y, max(c.width, 1), 1), c, fg, bg)
		}
		p.decorate(img, p.rect(x, y, 1, 1), c.style, fg, 1)
	}

	for _, t := range s.scaled {
		if t.style.Attrs&cellbuf.ConcealAttr != 0 {
			continue
		}
		fg, _, _ := p.colors(t.style)
		r := p.rect(t.x, t.y, t.scale*t.width, t.scale)
		p.text(img, r.Min, t.text, t.style, fg, t.scale)
		p.decorate(img, r, t.style, fg, t.scale)
	}

	return img
}

// rect returns the pixels covered by w by h cells from column x and row y.
func (p *painter) rect(x, y, w, h int) image.Rectangle {
	cw, ch := p.t.CellWidth, p.t.CellHeight
	return image.Rect(x*cw, y*ch, (x+w)*cw, (y+h)*ch)
}

func (p *painter) fill(img *image.RGBA, r image.Rectangle, c color.Color) {
	draw.Draw(img, r, image.NewUniform(c), image.Point{}, draw.Over)
}

// glyph draws what a cell shows into r, the cells it covers.
func (p *painter) glyph(img *image.RGBA, r image.Rectangle, c cell, fg, bg color.RGBA) {
	switch {
	case c.text == " " || c.text == " ":
		return
	case placeholder.Is(c.text):
		p.image(img, r, c)
		return
	}
	if ch, size := utf8.DecodeRuneInString(c.text); size == len(c.text) && p.box(img, r, ch, fg, bg) {
		return
	}
	p.text(img, r.Min, c.text, c.style, fg, 1)
}

// text draws a grapheme at scale times its size, from the top left corner of
// its cells at.
func (p *painter) text(img *image.RGBA, at image.Point, text string, style cellbuf.Style, fg color.RGBA, scale int) {
	base, _ := utf8.DecodeRuneInString(text)
	f, fakeBold := p.families[p.familyOf(base)].style(
		style.Attrs&cellbuf.BoldAttr != 0,
		style.Attrs&cellbuf.ItalicAttr != 0,
	)
	face := p.face(f, scale)
	if face == nil {
		return
	}

	// Leave out what the font has no glyph for, such as the selector asking
	// for an emoji's colored presentation, rather than draw it as a box.
	var b strings.Builder
	for i, r := range text {
		if g, err := f.GlyphIndex(&p.buf, r); i == 0 || err == nil && g != 0 {
			b.WriteRune(r)
		}
	}

	d := font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(fg),
		Face: face,
		Dot:  fixed.P(at.X, at.Y+p.baseline*scale),
	}
	d.DrawString(b.String())
	if fakeBold {
		d.Dot = fixed.P(at.X+max(1, int(p.ppem*float64(scale)/24)), at.Y+p.baseline*scale)
		d.DrawString(b.String())
	}
}

// familyOf returns which of the fonts to draw a character in: the first one
// with a glyph for it, or the first font of all when none has one.
func (p *painter) familyOf(r rune) int {
	if i, ok := p.family[r]; ok {
		return i
	}
	i := 0
	for j, f := range p.families {
		if g, err := f.Regular.GlyphIndex(&p.buf, r); err == nil && g != 0 {
			i = j
			break
		}
	}
	p.family[r] = i
	return i
}

func (p *painter) face(f *sfnt.Font, scale int) font.Face {
	k := faceKey{f, scale}
	if face, ok := p.faces[k]; ok {
		return face
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    p.ppem * float64(scale),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		face = nil
	}
	p.faces[k] = face
	return face
}

// image draws the part of an image that a placeholder cell shows.
func (p *painter) image(img *image.RGBA, r image.Rectangle, c cell) {
	id, row, col, ok := placeholder.At(c.text, c.style.Fg)
	if !ok {
		return
	}
	pic := p.picture(id)
	if pic == nil {
		return
	}
	draw.Draw(img, r, pic, image.Pt(col*p.t.CellWidth, row*p.t.CellHeight), draw.Over)
}

// picture returns a kept image scaled to the cells its placeholders cover,
// scaled whole rather than cell by cell so no seams show between cells.
func (p *painter) picture(id uint32) *image.RGBA {
	if pic, ok := p.images[id]; ok {
		return pic
	}
	var pic *image.RGBA
	if src, cols, rows, ok := placeholder.Image(id); ok {
		pic = image.NewRGBA(image.Rect(0, 0, cols*p.t.CellWidth, rows*p.t.CellHeight))
		xdraw.CatmullRom.Scale(pic, pic.Bounds(), src, src.Bounds(), xdraw.Src, nil)
	}
	p.images[id] = pic
	return pic
}

// decorate underlines and strikes through r, the cells text in style covers,
// at scale times their size.
func (p *painter) decorate(img *image.RGBA, r image.Rectangle, style cellbuf.Style, fg color.RGBA, scale int) {
	t := p.line * scale
	if style.UlStyle != cellbuf.NoUnderline {
		ul := fg
		if style.Ul != nil {
			ul = p.color(style.Ul, fg)
		}
		y := min(r.Min.Y+p.baseline*scale+t, r.Max.Y-t)
		p.fill(img, image.Rect(r.Min.X, y, r.Max.X, y+t), ul)
		if style.UlStyle == cellbuf.DoubleUnderline && y+3*t <= r.Max.Y {
			p.fill(img, image.Rect(r.Min.X, y+2*t, r.Max.X, y+3*t), ul)
		}
	}
	if style.Attrs&cellbuf.StrikethroughAttr != 0 {
		y := r.Min.Y + r.Dy()/2 - t/2
		p.fill(img, image.Rect(r.Min.X, y, r.Max.X, y+t), fg)
	}
}

// colors returns the colors a cell in style is drawn in, and whether its
// background is its own rather than the terminal's.
func (p *painter) colors(style cellbuf.Style) (fg, bg color.RGBA, solid bool) {
	fg, bg = p.color(style.Fg, p.fg), p.color(style.Bg, p.bg)
	solid = style.Bg != nil
	if style.Attrs&cellbuf.ReverseAttr != 0 {
		fg, bg, solid = bg, fg, true
	}
	if style.Attrs&cellbuf.FaintAttr != 0 {
		fg = mix(fg, bg, 0.5)
	}
	return fg, bg, solid
}

// color returns the color c stands for, in the terminal's palette when it
// names one of its first 16 colors, or def when there is no c.
func (p *painter) color(c color.Color, def color.RGBA) color.RGBA {
	switch c := c.(type) {
	case nil:
		return def
	case ansi.BasicColor:
		if c < 16 && p.t.Palette[c] != nil {
			return rgba(p.t.Palette[c])
		}
	case ansi.IndexedColor:
		if c < 16 && p.t.Palette[c] != nil {
			return rgba(p.t.Palette[c])
		}
	}
	return rgba(c)
}

func rgba(c color.Color) color.RGBA {
	r, g, b, _ := c.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 0xff}
}

// mix returns the color t of the way from a to b.
func mix(a, b color.RGBA, t float64) color.RGBA {
	m := func(x, y uint8) uint8 { return uint8(math.Round(float64(x) + (float64(y)-float64(x))*t)) }
	return color.RGBA{m(a.R, b.R), m(a.G, b.G), m(a.B, b.B), 0xff}
}
