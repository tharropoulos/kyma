//go:build unix

package img

import (
	"os"

	"golang.org/x/sys/unix"
)

// terminalCellSize asks the terminal on stdout how big its cells are, falling
// back to the default size when stdout is not a terminal.
func terminalCellSize() cellSize {
	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return cellSize{defaultCellWidth, defaultCellHeight}
	}
	return cellSizeFromWindow(ws.Col, ws.Row, ws.Xpixel, ws.Ypixel)
}
