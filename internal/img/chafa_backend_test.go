package img

import (
	"math"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

// TestChafaBackendKeepsAspectRatio renders the 500x300 docs logo into boxes of
// different shapes and checks that it comes out with its own proportions,
// rather than stretched to fill the box.
func TestChafaBackendKeepsAspectRatio(t *testing.T) {
	const (
		path        = "../../docs/kyma_logo.png"
		imageAspect = 500.0 / 300.0
	)

	tests := []struct {
		name          string
		width, height int32
		cell          cellSize
	}{
		{name: "wide box", width: 43, height: 10, cell: cellSize{18, 36}},
		{name: "tall box", width: 43, height: 30, cell: cellSize{18, 36}},
		{name: "narrow cells", width: 43, height: 10, cell: cellSize{10, 25}},
		{name: "wide cells", width: 43, height: 10, cell: cellSize{16, 26}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := NewChafaBackend().render(path, tt.width, tt.height, tt.cell, true)
			if err != nil {
				t.Fatalf("render() error = %v", err)
			}

			lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
			rows := int32(len(lines))
			var cols int32
			for _, line := range lines {
				cols = max(cols, int32(ansi.StringWidth(line)))
			}

			if cols > tt.width || rows > tt.height {
				t.Fatalf("image is %dx%d cells, larger than its %dx%d box", cols, rows, tt.width, tt.height)
			}
			if cols != tt.width && rows != tt.height {
				t.Errorf("image is %dx%d cells, filling neither side of its %dx%d box", cols, rows, tt.width, tt.height)
			}

			aspect := float64(cols*tt.cell.width) / float64(rows*tt.cell.height)
			if math.Abs(aspect-imageAspect)/imageAspect > 0.05 {
				t.Errorf("image is drawn %.2f wide for every 1 tall, want %.2f", aspect, imageAspect)
			}
		})
	}
}
