//go:build !unix

package img

// terminalCellSize returns the default cell size, since there is no portable
// way to ask the terminal for its size in pixels here.
func terminalCellSize() cellSize {
	return cellSize{defaultCellWidth, defaultCellHeight}
}
