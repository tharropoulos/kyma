package img

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/ploMP4/chafa-go"

	"github.com/museslabs/kyma/docs"
	"github.com/museslabs/kyma/internal/img/placeholder"
)

// docsBackend is an image backend tailored for documentation rendering.
// It loads images from an embedded filesystem instead of the user's local FS.
type docsBackend struct {
	cache *chafaCache
}

// NewDocsBackend creates a new instance of docsBackend with an initialized cache.
func NewDocsBackend() *docsBackend {
	return &docsBackend{
		cache: &chafaCache{
			cache: map[string]string{},
		},
	}
}

func (b *docsBackend) SymbolsOnly() bool {
	return b.detectTerminal().pixelMode == chafa.CHAFA_PIXEL_MODE_SYMBOLS
}

func (b *docsBackend) Render(path string, width, height int, symbols bool) (string, error) {
	var err error

	defer func() {
		if r := recover(); r != nil {
			if recErr, ok := r.(error); ok {
				err = recErr
			} else {
				err = fmt.Errorf("panic: %v", r)
			}
		}
	}()

	cell := terminalCellSize()

	c := b.cache.Get(path, width, height, cell, symbols)
	if c != "" {
		return c, nil
	}

	out, err := b.render(path, int32(width), int32(height), cell, symbols)
	b.cache.Save(out, path, width, height, cell, symbols)

	return out, err
}

func (b docsBackend) render(path string, width, height int32, cell cellSize, symbols bool) (string, error) {
	pixels, pixelWidth, pixelHeight, err := b.load(path)
	if err != nil {
		return "", err
	}

	// Images drawn with placeholders are uploaded in kitty's protocol whatever
	// chafa makes of the terminal, see [placeholder.Draw].
	placeholders := !symbols && placeholder.Enabled()

	// Shrink the box the slide gave the image to the largest size that keeps
	// the image's proportions, since the canvas stretches it to fill the box.
	chafa.CalcCanvasGeometry(pixelWidth, pixelHeight, &width, &height, cell.ratio(), true, false)

	capabilities := b.detectTerminal()
	defer chafa.SymbolMapUnref(capabilities.symbolMap)
	defer chafa.TermInfoUnref(capabilities.termInfo)

	config := chafa.CanvasConfigNew()
	defer chafa.CanvasConfigUnref(config)

	chafa.CanvasConfigSetCanvasMode(config, capabilities.canvasMode)
	chafa.CanvasConfigSetGeometry(config, width, height)
	chafa.CanvasConfigSetPassthrough(config, capabilities.passthrough)
	chafa.CanvasConfigSetSymbolMap(config, capabilities.symbolMap)
	chafa.CanvasConfigSetCellGeometry(config, cell.width, cell.height)

	switch {
	case symbols:
		chafa.CanvasConfigSetPixelMode(config, chafa.CHAFA_PIXEL_MODE_SYMBOLS)
	case placeholders:
		chafa.CanvasConfigSetPixelMode(config, chafa.CHAFA_PIXEL_MODE_KITTY)
		chafa.CanvasConfigSetPassthrough(config, chafa.CHAFA_PASSTHROUGH_NONE)
	default:
		chafa.CanvasConfigSetPixelMode(config, capabilities.pixelMode)
	}

	canvas := chafa.CanvasNew(config)
	defer chafa.CanvasUnRef(canvas)

	chafa.CanvasDrawAllPixels(
		canvas,
		chafa.CHAFA_PIXEL_RGBA8_UNASSOCIATED,
		pixels,
		pixelWidth,
		pixelHeight,
		pixelWidth*nChannels,
	)
	printable := chafa.CanvasPrint(canvas, nil)

	if placeholders {
		pixels, err := kittyPixels(printable.String())
		if err != nil {
			return "", err
		}
		return placeholder.Draw(pixels, int(width), int(height))
	}
	return printable.String(), nil
}

func (b docsBackend) detectTerminal() chafaTerminalCapabilities {
	termInfo := chafa.TermDbDetect(chafa.TermDbGetDefault(), os.Environ())

	mode := chafa.TermInfoGetBestCanvasMode(termInfo)
	pixelMode := chafa.TermInfoGetBestPixelMode(termInfo)

	passthrough := chafa.CHAFA_PASSTHROUGH_NONE
	if chafa.TermInfoGetIsPixelPassthroughNeeded(termInfo, pixelMode) {
		passthrough = chafa.TermInfoGetPassthroughType(termInfo)
	}

	symbolMap := chafa.SymbolMapNew()
	chafa.SymbolMapAddByTags(symbolMap, chafa.TermInfoGetSafeSymbolTags(termInfo))

	return chafaTerminalCapabilities{
		termInfo:    termInfo,
		canvasMode:  mode,
		pixelMode:   pixelMode,
		passthrough: passthrough,
		symbolMap:   symbolMap,
	}
}

func (b docsBackend) load(path string) (pixels []uint8, width, height int32, err error) {
	file, err := docs.FS.Open(path)
	if err != nil {
		return nil, 0, 0, err
	}
	defer file.Close()

	var img image.Image

	switch filepath.Ext(path) {
	case "png":
		img, err = png.Decode(file)
		if err != nil {
			return nil, 0, 0, err
		}
	case "jpg", "jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			return nil, 0, 0, err
		}
	default:
		img, _, err = image.Decode(file)
		if err != nil {
			return nil, 0, 0, err
		}
	}

	bounds := img.Bounds()
	width = int32(bounds.Dx())
	height = int32(bounds.Dy())

	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, img, bounds.Min, draw.Src)

	return rgbaImg.Pix, width, height, nil
}
