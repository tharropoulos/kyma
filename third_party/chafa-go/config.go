package chafa

var (
	// Creates a new [CanvasConfig] with default settings.
	// This object can later be used in the creation of a [Canvas].
	CanvasConfigNew func() *CanvasConfig

	// Creates a new [CanvasConfig] that's a copy of config.
	CanvasConfigCopy func(config *CanvasConfig) *CanvasConfig

	// Adds a reference to config.
	CanvasConfigRef func(config *CanvasConfig)

	// Removes a reference from config.
	CanvasConfigUnref func(config *CanvasConfig)

	// Returns config's width and height in character cells in the provided output locations.
	CanvasConfigGetGeometry func(config *CanvasConfig, widthOut, heightOut *int32)

	// Sets config's width and height in character cells to width x height.
	CanvasConfigSetGeometry func(config *CanvasConfig, width, height int32)

	// Returns config's cell width and height in pixels in the provided output locations.
	CanvasConfigGetCellGeometry func(config *CanvasConfig, cellWidthOut, cellHeightOut *int32)

	// Sets config's cell width and height in pixels to cellWidth x cellHeight.
	CanvasConfigSetCellGeometry func(config *CanvasConfig, cellWidth, cellHeight int32)

	// Returns config's [PixelMode].
	CanvasConfigGetPixelMode func(config *CanvasConfig) PixelMode

	// Sets config's stored [PixelMode] to pixelMode.
	// This determines how pixel graphics are rendered in the output.
	CanvasConfigSetPixelMode func(config *CanvasConfig, pixelMode PixelMode)

	// Returns config's [CanvasMode]. This determines how colors (and color control codes)
	// are used in the output.
	CanvasConfigGetCanvasMode func(config *CanvasConfig) CanvasMode

	// Sets config's stored [CanvasMode] to mode.
	// This determines how colors (and color control codes) are used in the output.
	CanvasConfigSetCanvasMode func(config *CanvasConfig, mode CanvasMode)

	// Returns config's [ColorExtractor]. This determines how colors
	// are approximated in character symbol output.
	CanvasConfigGetColorExtractor func(config *CanvasConfig) ColorExtractor

	// Sets config's stored [ColorExtractor] to colorExtractor.
	// This determines how colors are approximated in character symbol output.
	CanvasConfigSetColorExtractor func(config *CanvasConfig, colorExtractor ColorExtractor)

	// Returns config's [ColorSpace].
	CanvasConfigGetColorSpace func(config *CanvasConfig) ColorSpace

	// Sets config's stored [ColorSpace] to colorSpace.
	CanvasConfigSetColorSpace func(config *CanvasConfig, colorSpace ColorSpace)

	// Queries whether automatic image preprocessing is enabled. This allows Chafa
	// to boost contrast and saturation in an attempt to improve legibility.
	// The type of preprocessing applied (if any) depends on the canvas mode.
	CanvasConfigGetPreprocessingEnabled func(config *CanvasConfig) bool

	// Indicates whether automatic image preprocessing should be enabled. This
	// allows Chafa to boost contrast and saturation in an attempt to improve legibility.
	// The type of preprocessing applied (if any) depends on the canvas mode.
	CanvasConfigSetPreprocessingEnabled func(config *CanvasConfig, preprocessing_enabled bool)

	// Returns a pointer to the symbol map belonging to config. This can be
	// inspected using the [SymbolMap] getter functions, but not changed.
	CanvasConfigPeekSymbolMap func(config *CanvasConfig) *SymbolMap

	// Assigns a copy of symbolMap to config.
	CanvasConfigSetSymbolMap func(config *CanvasConfig, symbolMap *SymbolMap)

	// Returns a pointer to the fill symbol map belonging to config. This can
	// be inspected using the [SymbolMap] getter functions, but not changed.
	//
	// Fill symbols are assigned according to their overall foreground to
	// background coverage, disregarding shape.
	CanvasConfigPeekFillSymbolMap func(config *CanvasConfig) *SymbolMap

	// Returns the threshold above which full transparency will be used.
	CanvasConfigGetTransparencyThreshold func(config *CanvasConfig) float32

	// Sets the threshold above which full transparency will be used.
	CanvasConfigSetTransparencyThreshold func(config *CanvasConfig, alphaThreshold float32)

	// Queries whether to use foreground colors only, leaving the background
	// unmodified in the canvas output. This is relevant only when the [PixelMode]
	// is set to [CHAFA_PIXEL_MODE_SYMBOLS].
	//
	// When this is set, the canvas will emit escape codes to set the foreground color only.
	CanvasConfigGetFgOnlyEnabled func(config *CanvasConfig) bool

	// Indicates whether to use foreground colors only, leaving the background
	// unmodified in the canvas output. This is relevant only when the [PixelMode]
	// is set to [CHAFA_PIXEL_MODE_SYMBOLS].
	//
	//When this is set, the canvas will emit escape codes to set the foreground color only.
	CanvasConfigSetFgOnlyEnabled func(config *CanvasConfig, fgOnlyEnabled bool)

	// Gets the assumed foreground color of the output device. This is used to
	// determine how to apply the foreground pen in FGBG modes.
	CanvasConfigGetFgColor func(config *CanvasConfig) uint32

	// Sets the assumed foreground color of the output device. This is used to
	// determine how to apply the foreground pen in FGBG modes.
	CanvasConfigSetFgColor func(config *CanvasConfig, fgColorPackedRGB uint32)

	// Gets the assumed background color of the output device. This is used to
	// determine how to apply the background pen in FGBG modes.
	CanvasConfigGetBgColor func(config *CanvasConfig) uint32

	// Sets the assumed background color of the output device. This is used to
	// determine how to apply the background and transparency pens in FGBG modes,
	// and will also be substituted for partial transparency.
	CanvasConfigSetBgColor func(config *CanvasConfig, bgColorPackedRGB uint32)

	// Gets the work/quality tradeoff factor. A higher value means more time
	// and memory will be spent towards a higher quality output.
	CanvasConfigGetWorkFactor func(config *CanvasConfig) float32

	// Sets the work/quality tradeoff factor. A higher value means more time
	// and memory will be spent towards a higher quality output.
	CanvasConfigSetWorkFactor func(config *CanvasConfig, workFactor float32)

	// Returns config's [DitherMode].
	CanvasConfigGetDitherMode func(config *CanvasConfig) DitherMode

	// Sets config's stored [DitherMode] to ditherMode.
	CanvasConfigSetDitherMode func(config *CanvasConfig, ditherMode DitherMode)

	// Returns config's dither grain size in widthOut and heightOut.
	CanvasConfigGetDitherGrainSize func(config *CanvasConfig, widthOut, heightOut *int32)

	// Sets config's stored dither grain size to width by height pixels.
	// These values can be 1, 2, 4 or 8. 8 corresponds to the size of an entire
	// character cell. The default is 4 pixels by 4 pixels.
	CanvasConfigSetDitherGrainSize func(config *CanvasConfig, width, height int32)

	// Returns the relative intensity of the dithering pattern applied during
	// image conversion. 1.0 is the default, corresponding to a moderate intensity.
	CanvasConfigGetDitherIntensity func(config *CanvasConfig) float32

	// Sets config's stored relative intensity of the dithering pattern applied
	// during image conversion. 1.0 is the default, corresponding to a moderate
	// intensity. Possible values range from 0.0 to infinity, but in practice,
	// values above 10.0 are rarely useful.
	CanvasConfigSetDitherIntensity func(config *CanvasConfig, intensity float32)

	// Returns config's optimization flags. When enabled, these may produce
	// more compact output at the cost of reduced compatibility and increased
	// CPU use. Output quality is unaffected.
	CanvasConfigGetOptimizations func(config *CanvasConfig) Optimizations

	// Sets config's stored optimization flags. When enabled, these may produce
	// more compact output at the cost of reduced compatibility and increased
	// CPU use. Output quality is unaffected.
	CanvasConfigSetOptimizations func(config *CanvasConfig, optimizations Optimizations)

	// Returns config's [Passthrough] setting. This defaults to [CHAFA_PASSTHROUGH_NONE].
	//
	// Passthrough is needed to transmit certain escape codes to the outer terminal
	// when running in an inner terminal environment like tmux. When enabled,
	// this will happen automatically as needed, dictated by information
	// contained in a [TermInfo].
	CanvasConfigGetPassthrough func(config *CanvasConfig) Passthrough

	// Indicates which passthrough mode to use. This defaults to [CHAFA_PASSTHROUGH_NONE].
	//
	// [Passthrough] is needed to transmit certain escape codes to the outer
	// terminal when running in an inner terminal environment like tmux. When
	// enabled, this will happen automatically as needed, dictated by information
	// contained in a [TermInfo].
	CanvasConfigSetPassthrough func(config *CanvasConfig, passthrough Passthrough)
)

type CanvasConfig struct {
	Refs int32

	Width, Height                       int32
	CellWidth, CellHeight               int32
	CanvasMode                          CanvasMode
	ColorSpace                          ColorSpace
	DitherMode                          DitherMode
	ColorExtractor                      ColorExtractor
	PixelMode                           PixelMode
	DitherGrainWidth, DitherGrainHeight int32
	DitherIntensity                     float32
	FgColorPackedRgb                    uint32
	BgColorPackedRgb                    uint32
	AlphaThreshold                      int32 // 0-255. 255 = no alpha in output
	WorkFactor                          float32
	SymbolMap                           SymbolMap
	FillSymbolMap                       SymbolMap
	PreprocessingEnabled                bool
	FgOnlyEnabled                       bool
	Optimizations                       Optimizations
	Passthrough                         Passthrough
}
