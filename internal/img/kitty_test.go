package img

import (
	"bytes"
	"encoding/base64"
	"strconv"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"github.com/museslabs/kyma/internal/img/placeholder"
)

func TestKittyPixels_ReadsChunkedPixels(t *testing.T) {
	pixels := make([]byte, 4*3*2)
	for i := range pixels {
		pixels[i] = byte(i)
	}
	data := base64.StdEncoding.EncodeToString(pixels)
	codes := "\x1b_Ga=T,f=32,s=3,v=2,c=1,r=1,m=1\x1b\\" +
		"\x1b_Gm=1;" + data[:16] + "\x1b\\" +
		"\x1b_Gm=1;" + data[16:] + "\x1b\\" +
		"\x1b_Gm=0\x1b\\"

	img, err := kittyPixels(codes)
	if err != nil {
		t.Fatalf("kittyPixels() error = %v", err)
	}
	if b := img.Bounds(); b.Dx() != 3 || b.Dy() != 2 {
		t.Errorf("image is %dx%d, want 3x2", b.Dx(), b.Dy())
	}
	if !bytes.Equal(img.Pix, pixels) {
		t.Errorf("pixels are %v, want %v", img.Pix, pixels)
	}
}

func TestKittyPixels_Refuses(t *testing.T) {
	tests := map[string]string{
		"sixels":         "\x1bPq#0;2;0;0;0~-\x1b\\",
		"png":            "\x1b_Ga=T,f=100,s=1,v=1,m=0;AAAAAA==\x1b\\",
		"missing pixels": "\x1b_Ga=T,f=32,s=2,v=2,m=0;AAAAAA==\x1b\\",
		"cut short":      "\x1b_Ga=T,f=32,s=1,v=1,m=0;AAAAAA==",
	}
	for name, codes := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := kittyPixels(codes); err == nil {
				t.Error("kittyPixels() returned no error")
			}
		})
	}
}

// With placeholders enabled, an image comes out as placeholder cells in the box
// its symbols take, and goes up to the terminal at the size of those cells.
func TestChafaBackendDrawsPlaceholders(t *testing.T) {
	const path = "../../docs/kyma_logo.png"
	cell := cellSize{18, 36}

	t.Setenv("KYMA_IMAGE_PLACEHOLDERS", "1")
	var uploads bytes.Buffer
	placeholder.SetOutput(&uploads)
	t.Cleanup(func() { placeholder.SetOutput(nil) })

	b := NewChafaBackend()
	symbols, err := b.render(path, 43, 10, cell, true)
	if err != nil {
		t.Fatalf("render() symbols error = %v", err)
	}
	cells, err := b.render(path, 43, 10, cell, false)
	if err != nil {
		t.Fatalf("render() error = %v", err)
	}

	want := strings.Split(strings.TrimRight(symbols, "\n"), "\n")
	got := strings.Split(cells, "\n")
	if len(got) != len(want) {
		t.Fatalf("placeholders are %d rows tall, symbols %d", len(got), len(want))
	}
	for i := range got {
		if g, w := ansi.StringWidth(got[i]), ansi.StringWidth(want[i]); g != w {
			t.Errorf("row %d is %d cells wide, symbols %d", i, g, w)
		}
		if !placeholder.Is(ansi.Strip(got[i])) {
			t.Errorf("row %d is not placeholder cells: %q", i, got[i])
		}
	}

	cols := ansi.StringWidth(got[0])
	upload := uploads.String()
	if !strings.HasPrefix(upload, "\x1b_Ga=T,U=1,q=2,f=32,") {
		t.Fatalf("uploaded %.60q, want a quiet virtual placement", upload)
	}
	for _, key := range []string{
		",s=" + strconv.Itoa(cols*int(cell.width)) + ",",
		",v=" + strconv.Itoa(len(got)*int(cell.height)) + ",",
	} {
		if !strings.Contains(upload[:strings.Index(upload, ";")], key) {
			t.Errorf("uploaded %.80q, want %s", upload, key)
		}
	}
}
