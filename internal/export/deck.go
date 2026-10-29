// Package export draws a presentation's frames the way the terminal shows
// them, and writes them out as the pages of a PDF.
//
// A frame is the text kyma writes to the terminal. It is laid out on cells
// the way the terminal lays it out, then every cell is painted: its
// background, its glyph in the terminal's font or the first fallback font
// that has one, box drawing drawn by hand the way terminals draw it, images
// from the pixels kept behind their placeholder cells, see
// [placeholder.Keep], and scaled text at its scale.
package export

import "io"

// PageWidth is how wide a page is, in points: 13.33 inches, as wide as a
// widescreen slide. Pages are as tall as the terminal's proportions make them.
const PageWidth = 960

// A Deck writes frames out as the pages of a PDF.
type Deck struct {
	painter *painter
	pdf     *pdf

	// wallpaper is the object every page shows the wallpaper with, which
	// goes in once for them all, or 0 when it doesn't show.
	wallpaper int
}

// NewDeck starts a PDF on w, its pages drawn as t draws them.
func NewDeck(w io.Writer, t Terminal) (*Deck, error) {
	p, err := newPainter(t)
	if err != nil {
		return nil, err
	}
	d := &Deck{painter: p, pdf: newPDF(w)}
	if wallpaper := p.wallpaper(); wallpaper != nil {
		if d.wallpaper, err = d.pdf.photo(wallpaper); err != nil {
			return nil, err
		}
	}
	return d, nil
}

// Add draws a frame the way the terminal would show it, as the deck's next
// page.
func (d *Deck) Add(frame string) error {
	t := d.painter.t
	img := d.painter.paint(lay(frame, t.Cols, t.Rows))
	b := img.Bounds()
	return d.pdf.page(img, d.wallpaper, PageWidth, PageWidth*float64(b.Dy())/float64(b.Dx()))
}

// Close finishes the PDF. It doesn't close the writer the deck went to.
func (d *Deck) Close() error {
	return d.pdf.close()
}
