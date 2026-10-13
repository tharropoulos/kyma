package main

import (
	"fmt"

	"github.com/ploMP4/chafa-go"
)

const (
	PIX_WIDTH  = 3
	PIX_HEIGHT = 3
	N_CHANNELS = 4
)

func main() {
	pixels := []uint8{
		0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0xff,
		0x00, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0xff,
		0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0xff,
	}

	// Specify the symbols we want
	symbolMap := chafa.SymbolMapNew()
	defer chafa.SymbolMapUnref(symbolMap)

	chafa.SymbolMapAddByTags(symbolMap, chafa.CHAFA_SYMBOL_TAG_ALL)

	//  Set up a configuration with the symbols and the canvas size in characters
	config := chafa.CanvasConfigNew()
	defer chafa.CanvasConfigUnref(config)

	chafa.CanvasConfigSetGeometry(config, 40, 20)
	chafa.CanvasConfigSetSymbolMap(config, symbolMap)

	// Create canvas
	canvas := chafa.CanvasNew(config)
	defer chafa.CanvasUnRef(canvas)

	// Draw pixels to canvas
	chafa.CanvasDrawAllPixels(
		canvas,
		chafa.CHAFA_PIXEL_RGBA8_UNASSOCIATED,
		pixels,
		PIX_WIDTH,
		PIX_HEIGHT,
		PIX_WIDTH*N_CHANNELS,
	)

	// Generate a string that will show the canvas contents on a terminal
	gs := chafa.CanvasPrint(canvas, nil)

	fmt.Println(gs)
}
