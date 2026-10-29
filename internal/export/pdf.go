package export

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"math"
	"strconv"
)

// pdf writes a PDF one page at a time, each page an image that fills it. Pages
// go out as they come, so a long deck never has to be held whole.
type pdf struct {
	w   *bufio.Writer
	n   int64
	err error

	// offsets holds where each object starts, object i at i-1, and pages
	// the object numbers of the pages.
	offsets []int64
	pages   []int
}

// The page tree and the catalog are written last, once the pages are known,
// but every page refers to the tree, so their numbers are kept aside first.
const (
	pagesObject = 1 + iota
	catalogObject
)

func newPDF(w io.Writer) *pdf {
	p := &pdf{w: bufio.NewWriter(w), offsets: make([]int64, catalogObject)}
	// A comment of bytes above 127 tells tools the file is binary.
	p.printf("%%PDF-1.4\n%%\xe2\xe3\xcf\xd3\n")
	return p
}

func (p *pdf) printf(format string, args ...any) {
	if p.err != nil {
		return
	}
	n, err := fmt.Fprintf(p.w, format, args...)
	p.n += int64(n)
	p.err = err
}

func (p *pdf) write(b []byte) {
	if p.err != nil {
		return
	}
	n, err := p.w.Write(b)
	p.n += int64(n)
	p.err = err
}

// begin starts the next object and returns its number.
func (p *pdf) begin() int {
	p.offsets = append(p.offsets, p.n)
	id := len(p.offsets)
	p.printf("%d 0 obj\n", id)
	return id
}

// stream writes an object holding data, after the entries of its dictionary.
func (p *pdf) stream(entries string, data []byte) int {
	id := p.begin()
	p.printf("<< %s /Length %d >>\nstream\n", entries, len(data))
	p.write(data)
	p.printf("\nendstream\nendobj\n")
	return id
}

// page adds a page width by height points large, img filling it. A
// background, added with [pdf.photo] already, shows through wherever img
// isn't opaque, and 0 stands for none.
func (p *pdf) page(img *image.RGBA, background int, width, height float64) error {
	picture, err := p.image(img)
	if err != nil {
		return err
	}

	w, h := number(width), number(height)
	place := "q " + w + " 0 0 " + h + " 0 0 cm "
	contents, resources := place+"/Im Do Q", fmt.Sprintf("/Im %d 0 R", picture)
	if background != 0 {
		contents = place + "/Bg Do Q " + contents
		resources += fmt.Sprintf(" /Bg %d 0 R", background)
	}
	content := p.stream("", []byte(contents))

	id := p.begin()
	p.printf(
		"<< /Type /Page /Parent %d 0 R /MediaBox [0 0 %s %s] /Resources << /XObject << %s >> >> /Contents %d 0 R >>\nendobj\n",
		pagesObject, w, h, resources, content,
	)
	p.pages = append(p.pages, id)
	return p.err
}

// image adds an image and returns its object. Its colors and its alpha go in
// apart, the alpha as a soft mask, and only when the image isn't opaque all
// over.
func (p *pdf) image(img *image.RGBA) (int, error) {
	b := img.Bounds()
	rgb, alpha := make([]byte, 0, 3*b.Dx()*b.Dy()), make([]byte, 0, b.Dx()*b.Dy())
	opaque := true
	for y := b.Min.Y; y < b.Max.Y; y++ {
		row := img.Pix[img.PixOffset(b.Min.X, y):]
		for x := range b.Dx() {
			px := row[4*x : 4*x+4]
			a := px[3]
			opaque = opaque && a == 0xff
			alpha = append(alpha, a)
			// The image holds its colors multiplied by its alpha, which
			// PDF wants apart.
			for _, c := range px[:3] {
				switch a {
				case 0xff:
				case 0:
					c = 0
				default:
					c = byte((int(c)*0xff + int(a)/2) / int(a))
				}
				rgb = append(rgb, c)
			}
		}
	}

	entries := func(colorSpace string, colors int) string {
		return fmt.Sprintf(
			"/Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace %s /BitsPerComponent 8 "+
				"/Filter /FlateDecode /DecodeParms << /Predictor 15 /Colors %d /BitsPerComponent 8 /Columns %d >>",
			b.Dx(), b.Dy(), colorSpace, colors, b.Dx(),
		)
	}

	data, err := deflate(rgb, b.Dx(), 3)
	if err != nil {
		return 0, err
	}
	if opaque {
		return p.stream(entries("/DeviceRGB", 3), data), p.err
	}

	mask, err := deflate(alpha, b.Dx(), 1)
	if err != nil {
		return 0, err
	}
	soft := p.stream(entries("/DeviceGray", 1), mask)
	return p.stream(entries("/DeviceRGB", 3)+fmt.Sprintf(" /SMask %d 0 R", soft), data), p.err
}

// photo adds an opaque image as a JPEG, which keeps a photograph far smaller
// than deflating it does, and returns its object.
func (p *pdf) photo(img image.Image) (int, error) {
	var data bytes.Buffer
	if err := jpeg.Encode(&data, img, &jpeg.Options{Quality: 90}); err != nil {
		return 0, err
	}
	b := img.Bounds()
	return p.stream(fmt.Sprintf(
		"/Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace /DeviceRGB /BitsPerComponent 8 /Filter /DCTDecode",
		b.Dx(), b.Dy(),
	), data.Bytes()), p.err
}

// close writes the page tree, the catalog and the table of where every object
// starts, which a reader opens the file with.
func (p *pdf) close() error {
	p.offsets[pagesObject-1] = p.n
	p.printf("%d 0 obj\n<< /Type /Pages /Kids [", pagesObject)
	for _, id := range p.pages {
		p.printf(" %d 0 R", id)
	}
	p.printf(" ] /Count %d >>\nendobj\n", len(p.pages))

	p.offsets[catalogObject-1] = p.n
	p.printf("%d 0 obj\n<< /Type /Catalog /Pages %d 0 R >>\nendobj\n", catalogObject, pagesObject)

	xref := p.n
	p.printf("xref\n0 %d\n0000000000 65535 f \n", len(p.offsets)+1)
	for _, offset := range p.offsets {
		p.printf("%010d 00000 n \n", offset)
	}
	p.printf("trailer\n<< /Size %d /Root %d 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(p.offsets)+1, catalogObject, xref)

	if p.err != nil {
		return p.err
	}
	return p.w.Flush()
}

// deflate returns rows of pixels, width wide at channels bytes a pixel, every
// row filtered the way PNG filters them, by how each pixel differs from its
// neighbours, and deflated. Slides are mostly runs of a few flat colors,
// which the filter turns into runs of zeros.
func deflate(pix []byte, width, channels int) ([]byte, error) {
	stride := width * channels
	prev := make([]byte, stride)

	var out bytes.Buffer
	z := zlib.NewWriter(&out)
	filtered := make([]byte, 1+stride)
	filtered[0] = 4 // Paeth
	for len(pix) >= stride {
		row := pix[:stride]
		for i := range row {
			var left, upLeft byte
			if i >= channels {
				left, upLeft = row[i-channels], prev[i-channels]
			}
			filtered[1+i] = row[i] - paeth(left, prev[i], upLeft)
		}
		if _, err := z.Write(filtered); err != nil {
			return nil, err
		}
		prev, pix = row, pix[stride:]
	}
	if err := z.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// paeth predicts a byte from the ones left of it, above it and above and left
// of it, the way PNG's Paeth filter does.
func paeth(a, b, c byte) byte {
	p := int(a) + int(b) - int(c)
	pa, pb, pc := abs(p-int(a)), abs(p-int(b)), abs(p-int(c))
	switch {
	case pa <= pb && pa <= pc:
		return a
	case pb <= pc:
		return b
	}
	return c
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// number writes a length in points the way PDF reads numbers.
func number(f float64) string {
	return strconv.FormatFloat(math.Round(f*100)/100, 'f', -1, 64)
}
