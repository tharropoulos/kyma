package img

import (
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"strconv"
	"strings"
)

// kittyPixels reads an image back out of the escape codes chafa draws it with
// in kitty's graphics protocol: the keys in the first code, then the pixels as
// 32-bit RGBA, in chunks spread over as many codes as it takes, each chunk in
// base64 of its own. Chafa scales the image to fit its cells before it writes
// it out, which is the size it is drawn with placeholders at, see
// [placeholder.Draw].
func kittyPixels(codes string) (*image.NRGBA, error) {
	var (
		width, height int
		pixels        []byte
	)
	for first := true; codes != ""; first = false {
		code, ok := strings.CutPrefix(codes, "\x1b_G")
		end := strings.Index(code, "\x1b\\")
		if !ok || end < 0 {
			return nil, errors.New("not an image in kitty's graphics protocol")
		}
		codes = code[end+len("\x1b\\"):]

		keys, payload, _ := strings.Cut(code[:end], ";")
		chunk, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return nil, err
		}
		pixels = append(pixels, chunk...)
		if !first {
			continue
		}
		for kv := range strings.SplitSeq(keys, ",") {
			key, value, _ := strings.Cut(kv, "=")
			switch key {
			case "f":
				if value != "32" {
					return nil, fmt.Errorf("pixels in format %s rather than RGBA", value)
				}
			case "s":
				width, _ = strconv.Atoi(value)
			case "v":
				height, _ = strconv.Atoi(value)
			}
		}
	}

	if width <= 0 || height <= 0 || len(pixels) != 4*width*height {
		return nil, fmt.Errorf("%d bytes of pixels for a %dx%d image", len(pixels), width, height)
	}
	return &image.NRGBA{Pix: pixels, Stride: 4 * width, Rect: image.Rect(0, 0, width, height)}, nil
}
