package main

import (
	"fmt"
	"os"

	"github.com/ploMP4/chafa-go"
)

const (
	FONT_WIDTH  = 11
	FONT_HEIGHT = 24
	N_CHANNELS  = 4
)

func main() {
	pixels, width, height, err := chafa.Load("./examples/image.png")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	config := chafa.CanvasConfigNew()
	defer chafa.CanvasConfigUnref(config)

	chafa.CanvasConfigSetGeometry(config, 30, 30)
	chafa.CanvasConfigSetCellGeometry(config, FONT_WIDTH, FONT_HEIGHT)

	widthNew := config.Width
	heightNew := config.Height

	chafa.CalcCanvasGeometry(
		width, height,
		&widthNew, &heightNew,
		float32(FONT_WIDTH)/float32(FONT_HEIGHT),
		false, false,
	)
	chafa.CanvasConfigSetGeometry(config, widthNew, heightNew)

	canvas := chafa.CanvasNew(config)
	defer chafa.CanvasUnRef(canvas)

	chafa.CanvasDrawAllPixels(
		canvas,
		chafa.CHAFA_PIXEL_RGBA8_UNASSOCIATED,
		pixels,
		int32(width),
		int32(height),
		int32(width)*N_CHANNELS,
	)

	gs := chafa.CanvasPrint(canvas, nil)

	fmt.Println(gs)
}
