package export

import (
	"image"
	"image/color"
)

// Terminals draw box drawing characters and block elements themselves rather
// than take them from the font, so that lines meet from one cell to the next
// however tall the font leaves its cells. So does this, for the ones slides
// are drawn with: borders, rules, tables and blocks. The rest, the double
// lines among them, come from the font.

// The weight of a line out of the middle of a cell.
const (
	none = iota
	light
	heavy
)

// lines holds, for each of the box drawing characters from U+2500 made of
// straight lines out of the middle of the cell, the weight of its line up,
// right, down and left. A space stands for a character drawn some other way.
var lines = [...]string{
	"0101", "0202", "1010", "2020", // ─ ━ │ ┃
	" ", " ", " ", " ", " ", " ", " ", " ", // dashes
	"0110", "0210", "0120", "0220", // ┌ ┍ ┎ ┏
	"0011", "0012", "0021", "0022", // ┐ ┑ ┒ ┓
	"1100", "1200", "2100", "2200", // └ ┕ ┖ ┗
	"1001", "1002", "2001", "2002", // ┘ ┙ ┚ ┛
	"1110", "1210", "2110", "1120", "2120", "2210", "1220", "2220", // ├ … ┣
	"1011", "1012", "2011", "1021", "2021", "2012", "1022", "2022", // ┤ … ┫
	"0111", "0112", "0211", "0212", "0121", "0122", "0221", "0222", // ┬ … ┳
	"1101", "1102", "1201", "1202", "2101", "2102", "2201", "2202", // ┴ … ┻
	"1111", "1112", "1211", "1212", "2111", "1121", "2121", "2112", // ┼ … ╃
	"2211", "1122", "1221", "2212", "1222", "2122", "2221", "2222", // ╄ … ╋
}

// halves holds the same for the half lines from U+2574.
var halves = [...]string{
	"0001", "1000", "0100", "0010", // ╴ ╵ ╶ ╷
	"0002", "2000", "0200", "0020", // ╸ ╹ ╺ ╻
	"0201", "1020", "0102", "2010", // ╼ ╽ ╾ ╿
}

// A dash is a dashed line across or down the middle of a cell.
type dash struct {
	weight int
	down   bool
	n      int
}

var dashes = map[rune]dash{
	'┄': {light, false, 3}, '┅': {heavy, false, 3},
	'┆': {light, true, 3}, '┇': {heavy, true, 3},
	'┈': {light, false, 4}, '┉': {heavy, false, 4},
	'┊': {light, true, 4}, '┋': {heavy, true, 4},
	'╌': {light, false, 2}, '╍': {heavy, false, 2},
	'╎': {light, true, 2}, '╏': {heavy, true, 2},
}

// quadrants holds the quarters of the cell the quadrant block elements fill,
// top left, top right, bottom left and bottom right being 1, 2, 4 and 8.
var quadrants = map[rune]int{
	'▖': 4, '▗': 8, '▘': 1, '▙': 13, '▚': 9,
	'▛': 7, '▜': 11, '▝': 2, '▞': 6, '▟': 14,
}

// box draws a box drawing character or block element into r and reports
// whether it was one it draws.
func (p *painter) box(img *image.RGBA, r image.Rectangle, ch rune, fg, bg color.RGBA) bool {
	switch {
	case ch >= 0x2500 && ch < 0x2500+rune(len(lines)) && lines[ch-0x2500] != " ":
		p.arms(img, r, lines[ch-0x2500], fg)
	case ch >= 0x2574 && ch <= 0x257F:
		p.arms(img, r, halves[ch-0x2574], fg)
	case ch >= 0x256D && ch <= 0x2570:
		p.arc(img, r, ch, fg)
	case ch >= 0x2571 && ch <= 0x2573:
		p.diagonals(img, r, ch, fg)
	case ch >= 0x2580 && ch <= 0x259F:
		p.block(img, r, ch, fg, bg)
	default:
		d, ok := dashes[ch]
		if !ok {
			return false
		}
		p.dashed(img, r, d, fg)
	}
	return true
}

// thickness returns how many pixels thick a line of a weight is.
func (p *painter) thickness(weight int) int {
	switch weight {
	case light:
		return p.line
	case heavy:
		return 2 * p.line
	}
	return 0
}

// across returns the rows of r a line t pixels thick across its middle covers,
// and down the columns one down its middle covers.
func across(r image.Rectangle, t int) (top, bottom int) {
	top = r.Min.Y + (r.Dy()-t)/2
	return top, top + t
}

func down(r image.Rectangle, t int) (left, right int) {
	left = r.Min.X + (r.Dx()-t)/2
	return left, left + t
}

// arms draws lines from the middle of a cell to its edges, up, right, down
// and left at the weights urdl names. Each line runs on through the lines
// across its way, so they join without a notch.
func (p *painter) arms(img *image.RGBA, r image.Rectangle, urdl string, fg color.RGBA) {
	t := [4]int{}
	for i := range t {
		t[i] = p.thickness(int(urdl[i] - '0'))
	}
	up, right, dn, left := t[0], t[1], t[2], t[3]
	vertical, horizontal := max(up, dn), max(left, right)

	if up > 0 {
		x0, x1 := down(r, up)
		_, y1 := across(r, max(horizontal, up))
		p.fill(img, image.Rect(x0, r.Min.Y, x1, y1), fg)
	}
	if dn > 0 {
		x0, x1 := down(r, dn)
		y0, _ := across(r, max(horizontal, dn))
		p.fill(img, image.Rect(x0, y0, x1, r.Max.Y), fg)
	}
	if left > 0 {
		y0, y1 := across(r, left)
		_, x1 := down(r, max(vertical, left))
		p.fill(img, image.Rect(r.Min.X, y0, x1, y1), fg)
	}
	if right > 0 {
		y0, y1 := across(r, right)
		x0, _ := down(r, max(vertical, right))
		p.fill(img, image.Rect(x0, y0, r.Max.X, y1), fg)
	}
}

func (p *painter) dashed(img *image.RGBA, r image.Rectangle, d dash, fg color.RGBA) {
	t := p.thickness(d.weight)
	if d.down {
		x0, x1 := down(r, t)
		for i := range d.n {
			y0, y1 := r.Min.Y+i*r.Dy()/d.n, r.Min.Y+(i+1)*r.Dy()/d.n
			gap := (y1 - y0) / 4
			p.fill(img, image.Rect(x0, y0+gap/2, x1, y1-gap+gap/2), fg)
		}
		return
	}
	y0, y1 := across(r, t)
	for i := range d.n {
		x0, x1 := r.Min.X+i*r.Dx()/d.n, r.Min.X+(i+1)*r.Dx()/d.n
		gap := (x1 - x0) / 4
		p.fill(img, image.Rect(x0+gap/2, y0, x1-gap+gap/2, y1), fg)
	}
}

// arc draws a rounded corner, ╭ ╮ ╯ or ╰: a light line in from one edge and
// out to the next, turning through a quarter circle on the way.
func (p *painter) arc(img *image.RGBA, r image.Rectangle, ch rune, fg color.RGBA) {
	// Which way the line goes across and down from the corner, +1 for right
	// and down, -1 for left and up.
	sx, sy := float32(1), float32(1)
	switch ch {
	case '╮':
		sx = -1
	case '╯':
		sx, sy = -1, -1
	case '╰':
		sy = -1
	}

	w, h := r.Dx(), r.Dy()
	t := p.line
	ht := float32(t) / 2
	// The middle of the lines the corner joins, matching those of │ and ─.
	left, _ := down(image.Rect(0, 0, w, h), t)
	top, _ := across(image.Rect(0, 0, w, h), t)
	cx, cy := float32(left)+ht, float32(top)+ht

	// As round as the cell allows, out to the nearest edge, so the curve
	// meets the next cell's line on either side alike.
	radius := min(cx, float32(w)-cx, cy, float32(h)-cy)
	// The center of the circle the corner turns along, and the edges the
	// line goes out through.
	ox, oy := cx+sx*radius, cy+sy*radius
	ex, ey := float32(0), float32(0)
	if sx > 0 {
		ex = float32(w)
	}
	if sy > 0 {
		ey = float32(h)
	}

	// k places the control points of a cubic that bends like a quarter
	// circle.
	const k = 0.5523
	outer, inner := radius+ht, radius-ht

	z := &p.raster
	z.Reset(w, h)
	z.MoveTo(ex, cy-sy*ht)
	z.LineTo(ox, cy-sy*ht)
	z.CubeTo(ox-sx*k*outer, oy-sy*outer, ox-sx*outer, oy-sy*k*outer, ox-sx*outer, oy)
	z.LineTo(cx-sx*ht, ey)
	z.LineTo(cx+sx*ht, ey)
	z.LineTo(ox-sx*inner, oy)
	z.CubeTo(ox-sx*inner, oy-sy*k*inner, ox-sx*k*inner, oy-sy*inner, ox, oy-sy*inner)
	z.LineTo(ex, cy+sy*ht)
	z.ClosePath()
	z.Draw(img, r, image.NewUniform(fg), image.Point{})
}

// diagonals draws ╱, ╲ or both, ╳, from corner to corner of the cell.
func (p *painter) diagonals(img *image.RGBA, r image.Rectangle, ch rune, fg color.RGBA) {
	w, h := float32(r.Dx()), float32(r.Dy())
	// Half the line's thickness, measured across the cell.
	d := float32(p.line) / 2 * (1 + h/w) / 1.4

	z := &p.raster
	z.Reset(r.Dx(), r.Dy())
	if ch != '╲' {
		z.MoveTo(w-d, 0)
		z.LineTo(w+d, 0)
		z.LineTo(d, h)
		z.LineTo(-d, h)
		z.ClosePath()
	}
	if ch != '╱' {
		z.MoveTo(-d, 0)
		z.LineTo(d, 0)
		z.LineTo(w+d, h)
		z.LineTo(w-d, h)
		z.ClosePath()
	}
	z.Draw(img, r, image.NewUniform(fg), image.Point{})
}

// block draws one of the block elements, U+2580 to U+259F.
func (p *painter) block(img *image.RGBA, r image.Rectangle, ch rune, fg, bg color.RGBA) {
	w, h := r.Dx(), r.Dy()
	eighths := func(n, of int) int { return (n*of + 4) / 8 }

	switch {
	case ch == '▀':
		p.fill(img, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+h/2), fg)
	case ch >= '▁' && ch <= '█':
		n := int(ch-'▁') + 1
		p.fill(img, image.Rect(r.Min.X, r.Max.Y-eighths(n, h), r.Max.X, r.Max.Y), fg)
	case ch >= '▉' && ch <= '▏':
		n := 7 - int(ch-'▉')
		p.fill(img, image.Rect(r.Min.X, r.Min.Y, r.Min.X+eighths(n, w), r.Max.Y), fg)
	case ch == '▐':
		p.fill(img, image.Rect(r.Min.X+w/2, r.Min.Y, r.Max.X, r.Max.Y), fg)
	case ch >= '░' && ch <= '▓':
		p.fill(img, r, bg)
		shade := color.NRGBA{fg.R, fg.G, fg.B, uint8(0x40 * (int(ch-'░') + 1))}
		p.fill(img, r, shade)
	case ch == '▔':
		p.fill(img, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+eighths(1, h)), fg)
	case ch == '▕':
		p.fill(img, image.Rect(r.Max.X-eighths(1, w), r.Min.Y, r.Max.X, r.Max.Y), fg)
	default:
		q := quadrants[ch]
		mx, my := r.Min.X+w/2, r.Min.Y+h/2
		for i, part := range []image.Rectangle{
			image.Rect(r.Min.X, r.Min.Y, mx, my),
			image.Rect(mx, r.Min.Y, r.Max.X, my),
			image.Rect(r.Min.X, my, mx, r.Max.Y),
			image.Rect(mx, my, r.Max.X, r.Max.Y),
		} {
			if q&(1<<i) != 0 {
				p.fill(img, part, fg)
			}
		}
	}
}
