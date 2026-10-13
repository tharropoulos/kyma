package textsize

import (
	"fmt"
	"io"
	"testing"
	"time"
)

// fakeTerminal answers the probe the way a terminal would, reporting the
// cursor at each column in turn once it has read the probe.
func fakeTerminal(t *testing.T, cols ...int) (io.Reader, io.Writer, func() bool) {
	t.Helper()

	probeR, probeW := io.Pipe()
	replyR, replyW := io.Pipe()

	go func() {
		buf := make([]byte, len(probe))
		if _, err := io.ReadFull(probeR, buf); err != nil {
			return
		}
		for _, col := range cols {
			if _, err := fmt.Fprintf(replyW, "\x1b[1;%dR", col); err != nil {
				return
			}
		}
	}()

	t.Cleanup(func() {
		_ = probeR.Close()
		_ = replyW.Close()
	})

	cancel := func() bool {
		_ = replyR.CloseWithError(io.ErrClosedPipe)
		return true
	}

	return replyR, probeW, cancel
}

func TestQuery(t *testing.T) {
	tests := []struct {
		name string
		cols []int
		want bool
	}{
		{name: "width and scale", cols: []int{1, 3, 5}, want: true},
		{name: "width only", cols: []int{1, 3, 3}, want: false},
		{name: "neither", cols: []int{1, 1, 1}, want: false},
		{name: "scale without width", cols: []int{1, 2, 4}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, cancel := fakeTerminal(t, tt.cols...)

			if got := query(r, w, time.Second, cancel); got != tt.want {
				t.Errorf("query() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestQuery_TimesOut(t *testing.T) {
	// The terminal reads the probe but only ever answers one of the requests.
	r, w, cancel := fakeTerminal(t, 1)

	start := time.Now()
	if query(r, w, 50*time.Millisecond, cancel) {
		t.Error("query() = true for a terminal that never answered")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("query() took %v to give up", elapsed)
	}
}

func TestDetect_Override(t *testing.T) {
	t.Cleanup(func() { supported.Store(false) })

	t.Setenv(envOverride, "true")
	Detect()
	if !Supported() {
		t.Errorf("Supported() = false with %s=true", envOverride)
	}

	t.Setenv(envOverride, "0")
	Detect()
	if Supported() {
		t.Errorf("Supported() = true with %s=0", envOverride)
	}
}
