// Package placeholder draws images as text, with the Unicode placeholders of
// kitty's graphics protocol, see
// https://sw.kovidgoyal.net/kitty/graphics-protocol/#unicode-placeholders.
//
// An image placed straight onto the screen sits on top of the text, out of
// reach of everything that lays text out: a transition can neither move it nor
// cut it, so slides used to swap their images for symbols while one ran.
// Uploaded to the terminal once instead, an image is shown by a block of
// placeholder cells, each of which draws the part of the image that its row
// and column marks name. To everything that lays out text, transitions
// included, the cells are text like any other, and the image goes wherever
// they do, at its full resolution.
package placeholder

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// cell is the character placeholder cells are drawn with.
const cell = "\U0010EEEE"

// tile is the most rows or columns of cells one uploaded image spans, one for
// each of the marks that number them. Larger images go up in tiles this size.
const tile = len(diacritics)

// chunk is the most base64 kitty takes in one escape code.
const chunk = 4096

// envOverride names the environment variable that turns placeholders on or off
// whatever the terminal.
const envOverride = "KYMA_IMAGE_PLACEHOLDERS"

// unplace deletes every image placed on the screen. It leaves the images that
// placeholders show alone, since those have no place on the screen.
const unplace = "\x1b_Ga=d\x1b\\"

var enabled atomic.Bool

// uploads holds where images are uploaded to, and the id each image went up
// under, by what went up.
var uploads = struct {
	sync.Mutex
	out io.Writer
	ids map[[sha256.Size]byte]uint32
}{ids: map[[sha256.Size]byte]uint32{}}

// SetOutput sets where images are uploaded to, which draws them with
// placeholders from then on, if the terminal supports them. A nil out places
// them on the screen again. Images are uploaded from any goroutine, so out has
// to keep each write whole, the way the program's synchronized output does, or
// an upload could land in the middle of a frame.
func SetOutput(out io.Writer) {
	uploads.Lock()
	defer uploads.Unlock()

	if out == nil || !supported() {
		out = nil
	}
	uploads.out = out
	// What went up to the last output has to go up again to this one.
	clear(uploads.ids)
	enabled.Store(out != nil)
}

// Enabled reports whether images are drawn with placeholders.
func Enabled() bool {
	return enabled.Load()
}

// supported reports whether the terminal draws placeholders: kitty, which
// introduced them, and Ghostty. Anything else gets the images placed on the
// screen, and so do these inside tmux or screen, which chafa wraps images for.
//
// Setting KYMA_IMAGE_PLACEHOLDERS to a boolean forces the answer.
func supported() bool {
	if force, err := strconv.ParseBool(os.Getenv(envOverride)); err == nil {
		return force
	}
	switch os.Getenv("TERM") {
	case "xterm-kitty", "xterm-ghostty":
		return os.Getenv("TMUX") == ""
	}
	return false
}

// Clear returns kitty's command to delete the images placed on the screen,
// which stay up whatever text is drawn over them until they are deleted.
// Images drawn with placeholders come and go with their text, so while
// placeholders are enabled there is nothing to delete and it returns "".
func Clear() string {
	if Enabled() {
		return ""
	}
	return unplace
}

// Draw uploads an image to the terminal and returns the placeholder cells that
// show it, cols wide and rows tall. The image comes scaled to the size of those
// cells already, and the terminal shows it pixel for pixel. Drawing the same
// image again uploads nothing and comes back as the same cells.
func Draw(img *image.NRGBA, cols, rows int) (string, error) {
	b := img.Bounds()
	if cols <= 0 || rows <= 0 || b.Dx() < cols || b.Dy() < rows {
		return "", fmt.Errorf("%dx%d pixels do not fill %dx%d cells", b.Dx(), b.Dy(), cols, rows)
	}

	// Where the cells at column x and row y begin, in pixels.
	left := func(x int) int { return b.Min.X + x*b.Dx()/cols }
	top := func(y int) int { return b.Min.Y + y*b.Dy()/rows }

	across, down := (cols+tile-1)/tile, (rows+tile-1)/tile
	ids := make([]uint32, across*down)
	for ty := range down {
		for tx := range across {
			x0, y0 := tx*tile, ty*tile
			x1, y1 := min(x0+tile, cols), min(y0+tile, rows)
			part := img.SubImage(image.Rect(left(x0), top(y0), left(x1), top(y1))).(*image.NRGBA)

			id, err := upload(part, x1-x0, y1-y0)
			if err != nil {
				return "", err
			}
			ids[ty*across+tx] = id
		}
	}
	return text(ids, across, cols, rows), nil
}

// upload sends an image to the terminal, cols by rows cells large, unless it
// already went, and returns the id it went up under.
func upload(img *image.NRGBA, cols, rows int) (uint32, error) {
	b := img.Bounds()

	// The rows of a tile of a larger image lie apart from one another.
	pixels := make([]byte, 0, 4*b.Dx()*b.Dy())
	for y := b.Min.Y; y < b.Max.Y; y++ {
		i := img.PixOffset(b.Min.X, y)
		pixels = append(pixels, img.Pix[i:i+4*b.Dx()]...)
	}

	// Transmitted and displayed, but quietly, as kitty would otherwise answer
	// on stdin, and as a virtual placement, which shows nothing until
	// placeholder cells ask for it.
	keys := fmt.Sprintf("a=T,U=1,q=2,f=32,s=%d,v=%d,c=%d,r=%d", b.Dx(), b.Dy(), cols, rows)

	h := sha256.New()
	h.Write([]byte(keys))
	h.Write(pixels)
	var sum [sha256.Size]byte
	h.Sum(sum[:0])

	uploads.Lock()
	defer uploads.Unlock()

	if uploads.out == nil {
		return 0, errors.New("placeholders are not enabled")
	}
	if id, ok := uploads.ids[sum]; ok {
		return id, nil
	}

	// Ids count from 1, as 0 stands for no image.
	id := uint32(len(uploads.ids) + 1)
	if _, err := io.WriteString(uploads.out, transmission(id, keys, pixels)); err != nil {
		return 0, err
	}
	uploads.ids[sum] = id
	return id, nil
}

// transmission returns the escape codes that send pixels to the terminal under
// id: keys and the first chunk of base64 in the first code, and the rest of it
// in codes of their own, every one but the last saying more follows.
func transmission(id uint32, keys string, pixels []byte) string {
	data := base64.StdEncoding.EncodeToString(pixels)

	var b strings.Builder
	b.WriteString("\x1b_G" + keys + ",i=" + strconv.FormatUint(uint64(id), 10) + ",")
	for {
		n := min(len(data), chunk)
		if n < len(data) {
			b.WriteString("m=1;")
		} else {
			b.WriteString("m=0;")
		}
		b.WriteString(data[:n])
		b.WriteString("\x1b\\")

		data = data[n:]
		if data == "" {
			return b.String()
		}
		b.WriteString("\x1b_G")
	}
}

// text returns the placeholder cells that show an image cols wide and rows
// tall, uploaded as tiles with the given ids, across of them to a row. Every
// cell names its own row and column in its tile, so each one shows its part of
// the image wherever it is moved to and however the block is cut. The tile's
// id goes in the foreground color, which kitty reads as a 24-bit number.
func text(ids []uint32, across, cols, rows int) string {
	var b strings.Builder
	for y := range rows {
		if y > 0 {
			b.WriteByte('\n')
		}
		for x := range cols {
			if x%tile == 0 {
				id := ids[y/tile*across+x/tile]
				fmt.Fprintf(&b, "\x1b[38;2;%d;%d;%dm", byte(id>>16), byte(id>>8), byte(id))
			}
			b.WriteString(cell)
			b.WriteRune(diacritics[y%tile])
			b.WriteRune(diacritics[x%tile])
		}
		b.WriteString("\x1b[39m")
	}
	return b.String()
}

// Is reports whether a grapheme is a placeholder cell.
func Is(grapheme string) bool {
	return strings.HasPrefix(grapheme, cell)
}
