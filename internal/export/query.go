package export

import (
	"fmt"
	"image/color"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/muesli/cancelreader"
	"golang.org/x/term"
)

// queryTimeout bounds how long to wait on a terminal that doesn't answer.
const queryTimeout = 500 * time.Millisecond

// colorQuery asks for the colors text and the background are drawn in, then
// each of the 16 palette colors, then for the terminal's attributes. Every
// terminal answers the last, after the others, so its answer says that all
// the colors the terminal is going to report are in.
var colorQuery = func() string {
	var b strings.Builder
	b.WriteString("\x1b]10;?\x1b\\\x1b]11;?\x1b\\")
	for i := range 16 {
		fmt.Fprintf(&b, "\x1b]4;%d;?\x1b\\", i)
	}
	b.WriteString(ansi.RequestPrimaryDeviceAttributes)
	return b.String()
}()

var (
	colorReport      = regexp.MustCompile(`\x1b\](10|11|4;(\d+));rgba?:([0-9a-fA-F]+)/([0-9a-fA-F]+)/([0-9a-fA-F]+)`)
	attributesReport = regexp.MustCompile(`\x1b\[\?[\d;]*c`)
)

// Colors are the colors a terminal draws in, nil where it didn't say.
type Colors struct {
	Fg, Bg  color.Color
	Palette [16]color.Color
}

// QueryColors asks the terminal on stdin and stdout which colors it draws
// text, its background and its 16 palette colors in. It has to run before
// anything else reads stdin, which is where the answers arrive.
func QueryColors() Colors {
	in, out := os.Stdin, os.Stdout
	if !term.IsTerminal(int(in.Fd())) || !term.IsTerminal(int(out.Fd())) {
		return Colors{}
	}

	state, err := term.MakeRaw(int(in.Fd()))
	if err != nil {
		return Colors{}
	}
	defer func() { _ = term.Restore(int(in.Fd()), state) }()

	r, err := cancelreader.NewReader(in)
	if err != nil {
		return Colors{}
	}
	defer func() { _ = r.Close() }()

	if _, err := io.WriteString(out, colorQuery); err != nil {
		return Colors{}
	}

	answers := make(chan string, 1)
	go func() {
		var buf []byte
		chunk := make([]byte, 256)
		for {
			n, err := r.Read(chunk)
			buf = append(buf, chunk[:n]...)
			if err != nil || attributesReport.Match(buf) {
				answers <- string(buf)
				return
			}
		}
	}()

	select {
	case answer := <-answers:
		return parseColors(answer)
	case <-time.After(queryTimeout):
		// A reader that cannot be cancelled stays blocked until the next
		// keypress, so only wait for it to give up when it says it will.
		if r.Cancel() {
			return parseColors(<-answers)
		}
		return Colors{}
	}
}

// parseColors reads the colors out of a terminal's answers to colorQuery.
func parseColors(answers string) Colors {
	var c Colors
	for _, m := range colorReport.FindAllStringSubmatch(answers, -1) {
		rgb := color.RGBA{component(m[3]), component(m[4]), component(m[5]), 0xff}
		switch m[1] {
		case "10":
			c.Fg = rgb
		case "11":
			c.Bg = rgb
		default:
			if i, err := strconv.Atoi(m[2]); err == nil && i < len(c.Palette) {
				c.Palette[i] = rgb
			}
		}
	}
	return c
}

// component reads one of the hex components of an XParseColor rgb: color, of
// one to four digits, as 8 bits.
func component(hex string) uint8 {
	v, err := strconv.ParseUint(hex, 16, 32)
	if err != nil || len(hex) > 4 {
		return 0
	}
	maxValue := uint64(1)<<(4*len(hex)) - 1
	return uint8((v*0xff + maxValue/2) / maxValue)
}
