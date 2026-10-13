package textsize

import (
	"io"
	"os"
	"regexp"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/muesli/cancelreader"
	"golang.org/x/term"
)

// probeTimeout bounds how long startup waits on a terminal that never answers
// a cursor position request.
const probeTimeout = 500 * time.Millisecond

// envOverride names the environment variable that turns text sizing on or off
// without asking the terminal.
const envOverride = "KYMA_TEXT_SIZING"

var supported atomic.Bool

// probe draws a space two cells wide, then a space scaled to two by two cells,
// asking for the cursor position before and after each. A terminal that only
// implements widths, such as foot, moves the cursor for the first but not the
// second, and one that implements neither moves it for neither.
const probe = "\r" +
	ansi.RequestCursorPositionReport +
	"\x1b]66;w=2; \a" +
	ansi.RequestCursorPositionReport +
	"\x1b]66;s=2; \a" +
	ansi.RequestCursorPositionReport

var cursorReport = regexp.MustCompile(`\x1b\[(\d+);(\d+)R`)

// Supported reports whether [Detect] found a terminal that can scale text.
func Supported() bool {
	return supported.Load()
}

// Detect asks the terminal whether it implements the text sizing protocol and
// remembers the answer for [Supported]. It has to run before the TUI starts,
// since the answer arrives on stdin, which bubbletea reads from once it is
// running. Multiplexers such as tmux answer the query themselves without
// passing the escape codes through, so they correctly read as unsupported.
//
// Setting KYMA_TEXT_SIZING to a boolean skips the query and forces the answer.
func Detect() {
	if force, err := strconv.ParseBool(os.Getenv(envOverride)); err == nil {
		supported.Store(force)
		return
	}
	supported.Store(detect(os.Stdin, os.Stdout))
}

func detect(in, out *os.File) bool {
	if !term.IsTerminal(int(in.Fd())) || !term.IsTerminal(int(out.Fd())) {
		return false
	}

	state, err := term.MakeRaw(int(in.Fd()))
	if err != nil {
		return false
	}
	defer func() { _ = term.Restore(int(in.Fd()), state) }()

	r, err := cancelreader.NewReader(in)
	if err != nil {
		return false
	}
	defer func() { _ = r.Close() }()

	// The probe draws on the alternate screen, which bubbletea clears when it
	// takes over, so the test glyphs never end up in the scrollback.
	if _, err := io.WriteString(out, ansi.SetModeAltScreenSaveCursor+ansi.CursorHomePosition); err != nil {
		return false
	}
	defer func() { _, _ = io.WriteString(out, ansi.EraseEntireScreen+ansi.ResetModeAltScreenSaveCursor) }()

	return query(r, out, probeTimeout, r.Cancel)
}

// query writes the probe to w and reads the terminal's three cursor position
// reports from r. The column has to advance by two cells after each probe
// glyph. When the reports do not all arrive within timeout, cancel is called
// to unblock the pending read and the terminal counts as unsupported.
func query(r io.Reader, w io.Writer, timeout time.Duration, cancel func() bool) bool {
	if _, err := io.WriteString(w, probe); err != nil {
		return false
	}

	reports := make(chan []int, 1)
	go func() {
		reports <- readColumns(r, 3)
	}()

	var cols []int
	select {
	case cols = <-reports:
	case <-time.After(timeout):
		// A reader that cannot be cancelled stays blocked until the next
		// keypress, so only wait for it to give up when it says it will.
		if cancel() {
			<-reports
		}
		return false
	}

	return len(cols) == 3 && cols[1]-cols[0] == 2 && cols[2]-cols[1] == 2
}

// readColumns reads from r until it has seen n cursor position reports, or the
// reader fails, and returns the column of each report it saw.
func readColumns(r io.Reader, n int) []int {
	var (
		buf  []byte
		cols []int
	)

	chunk := make([]byte, 64)
	for {
		read, err := r.Read(chunk)
		buf = append(buf, chunk[:read]...)

		matches := cursorReport.FindAllSubmatch(buf, -1)
		if len(matches) >= n || err != nil {
			for _, m := range matches[:min(len(matches), n)] {
				col, _ := strconv.Atoi(string(m[2]))
				cols = append(cols, col)
			}
			return cols
		}
	}
}
