package export

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/gomonobolditalic"
	"golang.org/x/image/font/gofont/gomonoitalic"
	"golang.org/x/image/font/sfnt"
)

// A Family is a font in its regular, bold, italic and bold italic styles. Any
// but the regular one may be missing, and bold text is then drawn heavier by
// hand, and italic text upright.
type Family struct {
	Regular, Bold, Italic, BoldItalic *sfnt.Font
}

// style returns the font of the family to draw text in, and whether the font
// falls short of bold text, which then has to be thickened by hand.
func (f Family) style(bold, italic bool) (font *sfnt.Font, fakeBold bool) {
	switch {
	case bold && italic && f.BoldItalic != nil:
		return f.BoldItalic, false
	case bold && f.Bold != nil:
		return f.Bold, false
	case italic && f.Italic != nil:
		return f.Italic, bold
	}
	return f.Regular, bold
}

// LoadFamily reads a font file, TrueType, OpenType or the first font of a
// collection. When its name says it is the regular style, such as
// JetBrainsMono-Regular.ttf, the files beside it named for its other styles
// are read too.
func LoadFamily(path string) (Family, error) {
	regular, err := loadFont(path)
	if err != nil {
		return Family{}, err
	}
	f := Family{Regular: regular}

	dir, name := filepath.Split(path)
	if !strings.Contains(name, "Regular") {
		return f, nil
	}
	sibling := func(style string) *sfnt.Font {
		font, err := loadFont(filepath.Join(dir, strings.Replace(name, "Regular", style, 1)))
		if err != nil {
			return nil
		}
		return font
	}
	f.Bold, f.Italic, f.BoldItalic = sibling("Bold"), sibling("Italic"), sibling("BoldItalic")
	return f, nil
}

func loadFont(path string) (*sfnt.Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if font, err := sfnt.Parse(data); err == nil {
		return font, nil
	}
	collection, err := sfnt.ParseCollection(data)
	if err != nil {
		return nil, fmt.Errorf("%s is not a font: %w", path, err)
	}
	return collection.Font(0)
}

// goMono is the font that comes with kyma, drawn with when no font given has a
// glyph.
func goMono() Family {
	parse := func(data []byte) *sfnt.Font {
		font, err := sfnt.Parse(data)
		if err != nil {
			panic(err)
		}
		return font
	}
	return Family{
		Regular:    parse(gomono.TTF),
		Bold:       parse(gomonobold.TTF),
		Italic:     parse(gomonoitalic.TTF),
		BoldItalic: parse(gomonobolditalic.TTF),
	}
}

// systemFallbacks are fonts that come with the system, with many of the
// symbols Go Mono lacks, which are drawn with when they are installed. The
// terminal falls back on the system's fonts much the same way.
var systemFallbacks = []string{
	"/System/Library/Fonts/Menlo.ttc",
	"/System/Library/Fonts/Apple Symbols.ttf",
	"/System/Library/Fonts/Supplemental/Arial Unicode.ttf",
	"/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
	"/usr/share/fonts/TTF/DejaVuSansMono.ttf",
	"/usr/share/fonts/dejavu/DejaVuSansMono.ttf",
	"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
	"/usr/share/fonts/TTF/DejaVuSans.ttf",
}

// fallbacks returns Go Mono, then the system's fallback fonts installed here.
func fallbacks() []Family {
	families := []Family{goMono()}
	for _, path := range systemFallbacks {
		if font, err := loadFont(path); err == nil {
			families = append(families, Family{Regular: font})
		}
	}
	return families
}
