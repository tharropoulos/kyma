package chafa

import "unsafe"

var (
	// Creates a new [SymbolMap] representing a set of Unicode symbols.
	// The symbol map starts out empty.
	SymbolMapNew func() *SymbolMap

	// Creates a new [SymbolMap] that's a copy of symbolMap.
	SymbolMapCopy func(symbolMap *SymbolMap) *SymbolMap

	// Adds a reference to symbolMap.
	SymbolMapRef func(symbolMap *SymbolMap)

	// Removes a reference from symbolMap. When remaining references drops to
	// zero, the symbol map is freed and can no longer be used.
	SymbolMapUnref func(symbolMap *SymbolMap)

	// Adds symbols matching the set of tags to symbolMap.
	SymbolMapAddByTags func(symbolMap *SymbolMap, tags SymbolTags)

	// Adds symbols in the code point range starting with first and ending
	// with last to symbolMap.
	SymbolMapAddByRange func(symbolMap *SymbolMap, first, last rune)

	// Removes symbols matching the set of tags from symbolMap.
	SymbolMapRemoveByTags func(symbolMap *SymbolMap, tags SymbolTags)

	// Removes symbols in the code point range starting with first and ending
	// with last from symbolMap.
	SymbolMapRemoveByRange func(symbolMap *SymbolMap, first, last rune)

	// Parses a string consisting of symbol tags separated by [+-,] and applies
	// the pattern to symbolMap . If the string begins with + or -, it's
	// understood to be relative to the current set in symbolMap, otherwise the
	// map is cleared first.
	//
	// The symbol tags are string versions of [SymbolTags], i.e. [all, none,
	// space, solid, stipple, block, border, diagonal, dot, quad, half, hhalf,
	// vhalf, braille, technical, geometric, ascii, extra].
	//
	// Examples: "block,border" sets map to contain symbols matching either of
	// those tags. "+block,border-dot,stipple" adds block and border symbols then
	// removes dot and stipple symbols.
	//
	// If there is a parse error, none of the changes are applied.
	SymbolMapApplySelectors func(symbolMap *SymbolMap, selectors string) bool

	// Queries whether a symbol map is allowed to use built-in glyphs for symbol
	// selection. This can be turned off if you want to use your own glyphs
	// exclusively (see [SymbolMapAddGlyph]).
	//
	// Defaults to TRUE.
	SymbolMapGetAllowBuiltinGlyphs func(symbolMap *SymbolMap) bool

	// Controls whether a symbol map is allowed to use built-in glyphs for symbol
	// selection. This can be turned off if you want to use your own glyphs
	// exclusively (see [SymbolMapAddGlyph]).
	//
	// Defaults to TRUE.
	SymbolMapSetAllowBuiltinGlyphs func(symbolMap *SymbolMap, allow bool)

	// Returns data for the glyph corresponding to codePoint stored in symbolMap.
	// Any of pixelsOut , widthOut , heightOut and rowstrideOut can be nil,
	// in which case the corresponding data is not retrieved.
	//
	// If pixelsOut is not nil, a pointer to freshly allocated memory containing
	// height * rowstride bytes in the pixel format specified by pixelFormat will
	// be stored at this address. It must be freed using g_free() when you're
	// done with it.
	//
	// Monochrome glyphs (the only kind currently supported) will be rendered
	// as opaque white on a transparent black background
	// (0xffffffff for inked pixels and 0x00000000 for uninked).
	SymbolMapGetGlyph func(
		symbolMap *SymbolMap,
		codePoint rune,
		pixelFormat PixelType,
		pixelsOut **byte,
		widthOut, heightOut, rowstrideOut *int32,
	) bool

	// Assigns a rendered glyph to a Unicode code point. This tells Chafa what
	// the glyph looks like so the corresponding symbol can be used appropriately
	// in output.
	//
	// Assigned glyphs override built-in glyphs and any earlier glyph that may
	// have been assigned to the same code point.
	//
	// If the input is in a format with an alpha channel, the alpha channel will
	// be used for the shape. If not, an average of the color channels will be used.
	SymbolMapAddGlyph func(
		symbolMap *SymbolMap,
		codePoint rune,
		pixelFormat PixelType,
		pixels unsafe.Pointer,
		width, height, rowstride int32,
	)
)

type SymbolMap struct {
	Refs int32

	NeedRebuild      bool
	UseBuiltinGlyphs bool

	Glyphs    unsafe.Pointer
	Glyphs2   unsafe.Pointer // Wide glyphs with left/right bitmaps
	Selectors unsafe.Pointer

	// /* Remaining fields are populated by chafa_symbol_map_prepare () */

	// Narrow symbols
	Symbols       []Symbol
	NSymbols      int32
	PackedBitmaps *uint64

	// Wide symbols
	Symbols2       []Symbol2
	NSymbols2      int32
	PackedBitmaps2 *uint64
}

type Symbol struct {
	Sc                 SymbolTags
	C                  uint8
	Coverage           byte //gchar *
	MaskU32            *uint32
	FgWeight, BgWeight int32
	Bitmap             uint64
	Popcount           int32
}

type Symbol2 struct {
	Sym [2]Symbol
}

type SymbolTags int32

const (
	CHAFA_SYMBOL_TAG_NONE      SymbolTags = 0
	CHAFA_SYMBOL_TAG_SPACE     SymbolTags = (1 << 0)
	CHAFA_SYMBOL_TAG_SOLID     SymbolTags = (1 << 1)
	CHAFA_SYMBOL_TAG_STIPPLE   SymbolTags = (1 << 2)
	CHAFA_SYMBOL_TAG_BLOCK     SymbolTags = (1 << 3)
	CHAFA_SYMBOL_TAG_BORDER    SymbolTags = (1 << 4)
	CHAFA_SYMBOL_TAG_DIAGONAL  SymbolTags = (1 << 5)
	CHAFA_SYMBOL_TAG_DOT       SymbolTags = (1 << 6)
	CHAFA_SYMBOL_TAG_QUAD      SymbolTags = (1 << 7)
	CHAFA_SYMBOL_TAG_HHALF     SymbolTags = (1 << 8)
	CHAFA_SYMBOL_TAG_VHALF     SymbolTags = (1 << 9)
	CHAFA_SYMBOL_TAG_HALF      SymbolTags = ((CHAFA_SYMBOL_TAG_HHALF) | (CHAFA_SYMBOL_TAG_VHALF))
	CHAFA_SYMBOL_TAG_INVERTED  SymbolTags = (1 << 10)
	CHAFA_SYMBOL_TAG_BRAILLE   SymbolTags = (1 << 11)
	CHAFA_SYMBOL_TAG_TECHNICAL SymbolTags = (1 << 12)
	CHAFA_SYMBOL_TAG_GEOMETRIC SymbolTags = (1 << 13)
	CHAFA_SYMBOL_TAG_ASCII     SymbolTags = (1 << 14)
	CHAFA_SYMBOL_TAG_ALPHA     SymbolTags = (1 << 15)
	CHAFA_SYMBOL_TAG_DIGIT     SymbolTags = (1 << 16)
	CHAFA_SYMBOL_TAG_ALNUM     SymbolTags = CHAFA_SYMBOL_TAG_ALPHA | CHAFA_SYMBOL_TAG_DIGIT
	CHAFA_SYMBOL_TAG_NARROW    SymbolTags = (1 << 17)
	CHAFA_SYMBOL_TAG_WIDE      SymbolTags = (1 << 18)
	CHAFA_SYMBOL_TAG_AMBIGUOUS SymbolTags = (1 << 19)
	CHAFA_SYMBOL_TAG_UGLY      SymbolTags = (1 << 20)
	CHAFA_SYMBOL_TAG_LEGACY    SymbolTags = (1 << 21)
	CHAFA_SYMBOL_TAG_SEXTANT   SymbolTags = (1 << 22)
	CHAFA_SYMBOL_TAG_WEDGE     SymbolTags = (1 << 23)
	CHAFA_SYMBOL_TAG_LATIN     SymbolTags = (1 << 24)
	CHAFA_SYMBOL_TAG_IMPORTED  SymbolTags = (1 << 25)
	CHAFA_SYMBOL_TAG_OCTANT    SymbolTags = (1 << 26)
	CHAFA_SYMBOL_TAG_EXTRA     SymbolTags = (1 << 30)
	CHAFA_SYMBOL_TAG_BAD       SymbolTags = CHAFA_SYMBOL_TAG_AMBIGUOUS | CHAFA_SYMBOL_TAG_UGLY
	CHAFA_SYMBOL_TAG_ALL       SymbolTags = ^(CHAFA_SYMBOL_TAG_EXTRA | CHAFA_SYMBOL_TAG_BAD)
)
