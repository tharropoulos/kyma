package main

import (
	"fmt"
	"os"

	"github.com/ploMP4/chafa-go"
)

const (
	FONT_WIDTH  = 18
	FONT_HEIGHT = 36
	N_CHANNELS  = 4
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	pixels, width, height, err := chafa.Load(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	termInfo := chafa.TermDbDetect(chafa.TermDbGetDefault(), os.Environ())
	defer chafa.TermInfoUnref(termInfo)

	mode := chafa.TermInfoGetBestCanvasMode(termInfo)
	pixelMode := chafa.TermInfoGetBestPixelMode(termInfo)

	passthrough := chafa.CHAFA_PASSTHROUGH_NONE
	if chafa.TermInfoGetIsPixelPassthroughNeeded(termInfo, pixelMode) {
		passthrough = chafa.TermInfoGetPassthroughType(termInfo)
	}

	config := chafa.CanvasConfigNew()
	defer chafa.CanvasConfigUnref(config)

	chafa.CanvasConfigSetGeometry(config, 40, 20)
	chafa.CanvasConfigSetCellGeometry(config, FONT_WIDTH, FONT_HEIGHT)
	chafa.CanvasConfigSetCanvasMode(config, mode)
	chafa.CanvasConfigSetPixelMode(config, pixelMode)
	chafa.CanvasConfigSetPassthrough(config, passthrough)

	chafa.CalcCanvasGeometry(
		width, height,
		&config.Width, &config.Height,
		float32(FONT_WIDTH)/float32(FONT_HEIGHT),
		false, false,
	)
	chafa.CanvasConfigSetGeometry(config, config.Width, config.Height)

	canvas := chafa.CanvasNew(config)
	defer chafa.CanvasUnRef(canvas)

	frame := chafa.FrameNew(
		pixels,
		chafa.CHAFA_PIXEL_RGBA8_UNASSOCIATED,
		width,
		height,
		width*N_CHANNELS,
	)
	defer chafa.FrameUnref(frame)

	img := chafa.ImageNew()
	defer chafa.ImageUnref(img)

	chafa.ImageSetFrame(img, frame)

	placement := chafa.PlacementNew(img, 1)
	defer chafa.PlacementUnref(placement)

	chafa.CanvasSetPlacement(canvas, placement)

	gs := chafa.CanvasPrint(canvas, termInfo)

	fmt.Println(gs)
}
