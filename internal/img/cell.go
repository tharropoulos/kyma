package img

// Cell size assumed when the terminal does not report its size in pixels.
const (
	defaultCellWidth  = 18
	defaultCellHeight = 36
)

// cellSize is the size of one terminal cell in pixels.
type cellSize struct {
	width, height int32
}

// CellSize returns the size of one of the terminal's cells in pixels, or a
// default size when the terminal doesn't say.
func CellSize() (width, height int) {
	c := terminalCellSize()
	return int(c.width), int(c.height)
}

// ratio returns the cell's width divided by its height, which chafa calls the
// font ratio.
func (c cellSize) ratio() float32 {
	return float32(c.width) / float32(c.height)
}

// cellSizeFromWindow works out the size of a cell from the size of the window
// in cells and in pixels. Terminals that do not know their size in pixels
// report zero for it, which gets the default size instead.
func cellSizeFromWindow(cols, rows, xpixel, ypixel uint16) cellSize {
	if cols == 0 || rows == 0 || xpixel < cols || ypixel < rows {
		return cellSize{defaultCellWidth, defaultCellHeight}
	}
	return cellSize{int32(xpixel / cols), int32(ypixel / rows)}
}
