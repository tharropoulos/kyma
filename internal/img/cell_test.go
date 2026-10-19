package img

import "testing"

func TestCellSizeFromWindow(t *testing.T) {
	tests := []struct {
		name                       string
		cols, rows, xpixel, ypixel uint16
		want                       cellSize
	}{
		{name: "reported", cols: 100, rows: 40, xpixel: 1600, ypixel: 1360, want: cellSize{16, 34}},
		{name: "no pixel size", cols: 100, rows: 40, want: cellSize{defaultCellWidth, defaultCellHeight}},
		{name: "no cells", xpixel: 1600, ypixel: 1360, want: cellSize{defaultCellWidth, defaultCellHeight}},
		{name: "fewer pixels than cells", cols: 100, rows: 40, xpixel: 50, ypixel: 1360, want: cellSize{defaultCellWidth, defaultCellHeight}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cellSizeFromWindow(tt.cols, tt.rows, tt.xpixel, tt.ypixel); got != tt.want {
				t.Errorf("cellSizeFromWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}
