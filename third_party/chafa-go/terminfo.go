package chafa

var (
	// Creates a new, blank [TermInfo].
	TermInfoNew func() *TermInfo

	// Creates a new [TermInfo] that's a copy of termInfo.
	TermInfoCopy func(termInfo *TermInfo) *TermInfo

	// Adds a reference to termInfo.
	TermInfoRef func(termInfo *TermInfo)

	// Removes a reference from termInfo.
	TermInfoUnref func(termInfo *TermInfo)

	// Gets the string equivalent of seq stored in termInfo.
	TermInfoGetSeq func(termInfo *TermInfo, seq TermSeq) string

	// Sets the control sequence string equivalent of seq stored in termInfo to str.
	//
	// The string may contain argument indexes to be substituted with integers on
	// formatting. The indexes are preceded by a percentage character and
	// start at 1, i.e. %1, %2, %3, etc.
	//
	// The string's length after formatting must not exceed [CHAFA_TERM_SEQ_LENGTH_MAX] bytes.
	// Each argument can add up to four digits, or three for those specified as 8-bit integers.
	// If the string could potentially exceed this length when formatted,
	// [TermInfoSetSeq] will return FALSE.
	//
	// If parsing fails or str is too long, any previously existing sequence will be left untouched.
	//
	// Passing NULL for str clears the corresponding control sequence.
	TermInfoSetSeq func(termInfo *TermInfo, seq TermSeq, str string, err **GError) bool // TODO:

	// Checks if termInfo can emit seq.
	TermInfoHaveSeq func(termInfo *TermInfo, seq TermSeq) bool

	// Formats the terminal sequence seq, inserting positional arguments.
	// The seq's number of arguments must be supplied exactly.
	//
	// The argument list must be terminated by -1, or undefined behavior will result.
	//
	// If the wrong number of arguments is supplied, or an argument is out of range,
	// this function will return NULL. Otherwise, it returns a zero-terminated string
	// that must be freed with g_free().
	//
	// If you want compile-time validation of arguments, consider using one of
	// the specific chafa_term_info_emit_*() functions. They are also faster,
	// but require you to allocate at least [CHAFA_TERM_SEQ_LENGTH_MAX] bytes up front.
	TermInfoEmitSeq func(termInfo *TermInfo, seq TermSeq, args ...any) string

	// Attempts to parse a terminal sequence from an input data array.
	// If successful, [CHAFA_PARSE_SUCCESS] will be returned, the input pointer
	// will be advanced and the parsed length will be subtracted from inputLen.
	TermInfoParseSeq func(
		term_info *TermInfo,
		seq TermSeq,
		input []string,
		inputLen *int32,
		argsOut *uint32,
	) ParseResult

	// Supplements missing sequences in termInfo with ones copied from source.
	TermInfoSupplement func(termInfo, source *TermInfo)

	// Prints the control sequence for [CHAFA_TERM_SEQ_RESET_TERMINAL_SOFT].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitResetTerminalSoft func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_RESET_TERMINAL_HARD].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitResetTerminalHard func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_RESET_ATTRIBUTES].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitResetAttributes func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_CLEAR].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitClear func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_CURSOR_TO_POS].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitCursorToPos func(termInfo *TermInfo, dest *string, x, y uint32) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_CURSOR_TO_TOP_LEFT].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitCursorToTopLeft func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_CURSOR_TO_BOTTOM_LEFT].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitCursorToBottomLeft func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_CURSOR_UP].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitCursorUp func(termInfo *TermInfo, dest *string, n uint32) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_CURSOR_DOWN].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitCursorDown func(termInfo *TermInfo, dest *string, n uint32) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_CURSOR_LEFT].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitCursorLeft func(termInfo *TermInfo, dest *string, n uint32) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_CURSOR_RIGHT].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitCursorRight func(termInfo *TermInfo, dest *string, n uint32) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_CURSOR_UP_1].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitCursorUp1 func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_CURSOR_DOWN_1].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitCursorDown1 func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_CURSOR_LEFT_1].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitCursorLeft1 func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_CURSOR_RIGHT_1].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitCursorRight1 func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_CURSOR_UP_SCROLL].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitCursorUpScroll func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_CURSOR_DOWN_SCROLL].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitCursorDownScroll func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_INSERT_CELLS].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitInsertCells func(termInfo *TermInfo, dest *string, n uint32) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_DELETE_CELLS].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitDeleteCells func(termInfo *TermInfo, dest *string, n uint32) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_INSERT_ROWS].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitInsertRows func(termInfo *TermInfo, dest *string, n uint32) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_DELETE_ROWS].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitDeleteRows func(termInfo *TermInfo, dest *string, n uint32) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_ENABLE_CURSOR].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitEnableCursor func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_DISABLE_CURSOR].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitDisableCursor func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_ENABLE_ECHO].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitEnableEcho func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_DISABLE_ECHO].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitDisableEcho func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_ENABLE_INSERT].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitEnableInsert func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_DISABLE_INSERT].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitDisableInsert func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_ENABLE_WRAP].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitEnableWrap func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_DISABLE_WRAP].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitDisableWrap func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_ENABLE_BOLD].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitEnableBold func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_INVERT_COLORS].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitInvertColors func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_COLOR_BG_8].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetColorBg8 func(termInfo *TermInfo, dest *string, pen uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_COLOR_FG_8].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetColorFg8 func(termInfo *TermInfo, dest *string, pen uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_COLOR_FGBG_8].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetColorFgbg8 func(termInfo *TermInfo, dest *string, fgPen, bgPen uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_COLOR_FG_16].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetColorFg16 func(termInfo *TermInfo, dest *string, pen uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_COLOR_BG_16].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetColorBg16 func(termInfo *TermInfo, dest *string, pen uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_COLOR_FGBG_16].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetColorFgbg16 func(termInfo *TermInfo, dest *string, fgPen, bgPen uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_COLOR_FG_256].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetColorFg256 func(termInfo *TermInfo, dest *string, pen uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_COLOR_BG_256].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetColorBg256 func(termInfo *TermInfo, dest *string, pen uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_COLOR_FGBG_256].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetColorFgbg256 func(termInfo *TermInfo, dest *string, fgPen, bgPen uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_COLOR_FG_DIRECT].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetColorFgDirect func(termInfo *TermInfo, dest *string, r, g, b uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_COLOR_BG_DIRECT].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetColorBgDirect func(termInfo *TermInfo, dest *string, r, g, b uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_COLOR_FGBG_DIRECT].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetColorFgbgDirect func(termInfo *TermInfo, dest *string, fgR, fgG, fgB, bgR, bgG, bgB uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_RESET_COLOR_FG].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitResetColorFg func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_RESET_COLOR_BG].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitResetColorBg func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_RESET_COLOR_FGBG].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitResetColorFgbg func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_DEFAULT_FG].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetDefaultFg func(termInfo *TermInfo, dest *string, r, g, b uint16) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_DEFAULT_BG].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetDefaultBg func(termInfo *TermInfo, dest *string, r, g, b uint16) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_RESET_DEFAULT_FG].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitResetDefaultFg func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_RESET_DEFAULT_BG].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitResetDefaultBg func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_QUERY_DEFAULT_FG].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitQueryDefaultFg func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_QUERY_DEFAULT_BG].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitQueryDefaultBg func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_REPEAT_CHAR].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitRepeatChar func(termInfo *TermInfo, dest *string, n uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_SCROLLING_ROWS].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetScrollingRows func(termInfo *TermInfo, dest *string, top, bottom uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_RESET_SCROLLING_ROWS].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitResetScrollingRows func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SAVE_CURSOR_POS].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSaveCursorPos func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_RESTORE_CURSOR_POS].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitRestoreCursorPos func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_BEGIN_SIXELS].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	//
	// All three parameters (p1 , p2 and p3 ) can normally be set to 0.
	TermInfoEmitBeginSixels func(termInfo *TermInfo, dest *string, p1, p2, p3 uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_END_SIXELS].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitEndSixels func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_ENABLE_SIXEL_SCROLLING].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitEnableSixelScrolling func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_DISABLE_SIXEL_SCROLLING].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitDisableSixelScrolling func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_SIXEL_ADVANCE_DOWN].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetSixelAdvanceDown func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_SET_SIXEL_ADVANCE_RIGHT].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitSetSixelAdvanceRight func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_BEGIN_KITTY_IMMEDIATE_IMAGE_V1].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	//
	// bpp must be set to either 24 for RGB data, 32 for RGBA, or 100 to embed a PNG file.
	//
	// This sequence must be followed by zero or more paired sequences of type
	// [CHAFA_TERM_SEQ_BEGIN_KITTY_IMAGE_CHUNK] and [CHAFA_TERM_SEQ_END_KITTY_IMAGE_CHUNK]
	// with base-64 encoded image data between them.
	//
	// When the image data has been transferred, [CHAFA_TERM_SEQ_END_KITTY_IMAGE] must be emitted.
	TermInfoEmitBeginKittyImmediateImageV1 func(
		termInfo *TermInfo,
		dest *string,
		bpp, widthPixels, heightPixels, widthCells, heightCells uint8,
	) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_BEGIN_KITTY_IMMEDIATE_IMAGE_V1].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	//
	// bpp must be set to either 24 for RGB data, 32 for RGBA, or 100 to embed a PNG file.
	//
	// This sequence must be followed by zero or more paired sequences of type
	// [CHAFA_TERM_SEQ_BEGIN_KITTY_IMAGE_CHUNK] and [CHAFA_TERM_SEQ_END_KITTY_IMAGE_CHUNK]
	// with base-64 encoded image data between them.
	//
	// When the image data has been transferred, [CHAFA_TERM_SEQ_END_KITTY_IMAGE] must be emitted.
	TermInfoEmitBeginKittyImmediateVirtImageV1 func(
		termInfo *TermInfo,
		dest *string,
		bpp, widthPixels, heightPixels, widthCells, heightCells, id uint8,
	) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_END_KITTY_IMAGE].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitEndKittyImage func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_BEGIN_KITTY_IMAGE_CHUNK].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitBeginKittyImageChunk func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_END_KITTY_IMAGE_CHUNK].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitEndKittyImageChunk func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_BEGIN_ITERM2_IMAGE].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	//
	// This sequence must be followed by base64-encoded image file data. The image
	// can be any format supported by MacOS, e.g. PNG, JPEG, TIFF, GIF. When the
	// image data has been transferred, [CHAFA_TERM_SEQ_END_ITERM2_IMAGE] must be emitted.
	TermInfoEmitBeginIterm2Image func(termInfo *TermInfo, dest *string, width, height uint8) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_END_ITERM2_IMAGE].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitEndIterm2Image func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_BEGIN_SCREEN_PASSTHROUGH].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	//
	// Any control sequences between the beginning and end passthrough seqs must
	// be escaped by turning \033 into \033\033.
	TermInfoEmitBeginScreenPassthrough func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_END_SCREEN_PASSTHROUGH].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	//
	// Any control sequences between the beginning and end passthrough seqs must
	// be escaped by turning \033 into \033\033.
	TermInfoEmitEndScreenPassthrough func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_ENABLE_ALT_SCREEN].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitEnableAltScreen func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_DISABLE_ALT_SCREEN].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitDisableAltScreen func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_BEGIN_TMUX_PASSTHROUGH].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	//
	// Any control sequences between the beginning and end passthrough seqs must
	// be escaped by turning \033 into \033\033.
	TermInfoEmitBeginTmuxPassthrough func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_END_TMUX_PASSTHROUGH].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	//
	// Any control sequences between the beginning and end passthrough seqs must
	// be escaped by turning \033 into \033\033.
	TermInfoEmitEndTmuxPassthrough func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_RETURN_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitReturnKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_BACKSPACE_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitBackspaceKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_DELETE_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitDeleteKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_DELETE_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitDeleteCtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_DELETE_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitDeleteShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_INSERT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitInsertKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_INSERT_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitInsertCtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_INSERT_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitInsertShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_HOME_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitHomeKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_HOME_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitHomeCtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_HOME_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitHomeShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_END_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitEndKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_END_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitEndCtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_END_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitEndShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_UP_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitUpKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_UP_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitUpCtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_UP_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitUpShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_DOWN_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitDownKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_DOWN_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitDownCtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_DOWN_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitDownShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_LEFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitLeftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_LEFT_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitLeftCtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_LEFT_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitLeftShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_RIGHT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitRightKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_RIGHT_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitRightCtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_RIGHT_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitRightShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_PAGE_UP_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitPageUpKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_PAGE_UP_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitPageUpCtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_PAGE_UP_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitPageUpShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_PAGE_DOWN_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitPageDownKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_PAGE_DOWN_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitPageDownCtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_PAGE_DOWN_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitPageDownShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_TAB_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitTabKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_TAB_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitTabShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F1_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF1Key func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F1_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF1CtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F1_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF1ShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F2_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF2Key func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F2_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF2CtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F2_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF2ShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F3_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF3Key func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F3_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF3CtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F3_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF3ShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F4_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF4Key func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F4_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF4CtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F4_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF4ShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F5_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF5Key func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F5_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF5CtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F5_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF5ShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F6_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF6Key func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F6_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF6CtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F6_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF6ShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F7_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF7Key func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F7_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF7CtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F7_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF7ShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F8_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF8Key func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F8_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF8CtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F8_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF8ShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F9_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF9Key func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F9_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF9CtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F9_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF9ShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F10_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF10Key func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F10_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF10CtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F10_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF10ShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F11_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF11Key func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F11_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF11CtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F11_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF11ShiftKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F12_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF12Key func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F12_CTRL_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF12CtrlKey func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_F12_SHIFT_KEY].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitF12ShiftKey func(termInfo *TermInfo, dest *string) string

	// Terminal emulators and applications are often nested, with the inner
	// application's capabilities limiting, extending or modifying the outer's.
	//
	// Examples are terminal multiplexers like Screen and tmux, or terminal
	// emulators running inside editors like Emacs and vi.
	//
	// This merges the outer and inner sequences into a single [ChafaTermInfo]
	// according to the following rules.
	//
	// For sequences marked as inherited in inner:
	//
	//  - If either inner or outer sequence is nil, pick the outer sequence.
	//  - Otherwise, pick inner sequence.
	//
	// For sequences not marked as inherited, always use the inner sequence.
	//
	// This allows for using the inner term's sequences while clearing them if
	// the outer term does not support the sequence at all. This is useful for
	// muxers (e.g. fbterm supports 256 colors, but with private seqs; we want to
	// use the inner mux' corresponding seqs).
	//
	// The merged [TermInfo] is a new instance, with the initial reference owned
	// by the caller.
	//
	// This function can be used repeatedly to create chains that're arbitrarily long,
	// but is unlikely to be useful beyond three levels
	// (terminal emulator, multiplexer, application).
	TermInfoChain func(outer, inner *TermInfo) *TermInfo

	// Prints the control sequence for [CHAFA_TERM_SEQ_CELL_SIZE_PX].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitCellSizePx func(termInfo *TermInfo, dest *string, heightPx, widthPx uint32) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_PRIMARY_DEVICE_ATTRIBUTES].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitPrimaryDeviceAttributes func(termInfo *TermInfo, dest *string, args *uint32, nArgs int32) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_QUERY_CELL_SIZE_PX].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitQueryCellSizePx func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_QUERY_PRIMARY_DEVICE_ATTRIBUTES].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitQueryPrimaryDeviceAttributes func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_QUERY_TEXT_AREA_SIZE_CELLS].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX bytes],
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitQueryTextAreaSizeCells func(termInfo *TermInfo, dest *string) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_QUERY_TEXT_AREA_SIZE_PX].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitQueryTextAreaSizePx func(termInfo *TermInfo, dest *string) string

	// Behaves like [TermInfoEmitSeq]
	TermInfoEmitSeqValist func(termInfo *TermInfo, seq TermSeq, args ...any) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_TEXT_AREA_SIZE_CELLS].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX bytes],
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitTextAreaSizeCells func(termInfo *TermInfo, dest *string, heightCells, widthCells uint32) string

	// Prints the control sequence for [CHAFA_TERM_SEQ_TEXT_AREA_SIZE_PX].
	//
	// dest must have enough space to hold [CHAFA_TERM_SEQ_LENGTH_MAX] bytes,
	// even if the emitted sequence is shorter. The output will not be zero-terminated.
	TermInfoEmitTextAreaSizePx func(termInfo *TermInfo, dest *string, heightPx, widthPx uint32) string

	// Gets the optimal [CanvasMode] supported by termInfo.
	TermInfoGetBestCanvasMode func(termInfo *TermInfo) CanvasMode

	// Gets the optimal [PixelMode] supported by termInfo.
	TermInfoGetBestPixelMode func(termInfo *TermInfo) PixelMode

	// Gets whether seq can be inherited from the outer [TermInfo] when chaining
	// with [TermInfoChain].
	TermInfoGetInheritSeq func(termInfo *TermInfo, seq TermSeq) bool

	// Gets if passthrough should be used to convey images in pixel_mode.
	// If needed, the passthrough type can be gotten from [TermInfoGetPassthroughType].
	TermInfoGetIsPixelPassthroughNeeded func(termInfo *TermInfo, pixelMode PixelMode) bool

	// Gets the name associated with the termInfo. This may be nil. The returned
	// string belongs to termInfo, and is only valid until the next operation on
	// this data structure.
	TermInfoGetName func(termInfo *TermInfo) string

	// Gets the passthrough mode supported by termInfo.
	TermInfoGetPassthroughType func(termInfo *TermInfo) Passthrough

	// Gets the quirks associated with termInfo.
	TermInfoGetQuirks func(termInfo *TermInfo) TermQuirks

	// Gets the [SymbolTags] that are likely safe to use with termInfo.
	// The [SymbolTags] are a bitwise OR of flags from the enum.
	TermInfoGetSafeSymbolTags func(termInfo *TermInfo) SymbolTags

	// Checks if termInfo has terminal sequences that support canvasMode.
	TermInfoIsCanvasModeSupported func(termInfo *TermInfo, canvasMode CanvasMode) bool

	// Checks if termInfo has terminal sequences that support pixelMode.
	TermInfoIsPixelModeSupported func(termInfo *TermInfo, pixelMode PixelMode) bool

	// Attempts to parse a terminal sequence from an input data array. If successful,
	// [CHAFA_PARSE_SUCCESS] will be returned, the input pointer will be advanced
	// and the parsed length will be subtracted from inputLen.
	//
	// Any numeric parsed arguments are returned as an array starting at argsOut ,
	// which must have room for up to [CHAFA_TERM_SEQ_ARGS_MAX] elements.
	//
	// The number of parsed arguments is returned in nArgsOut. This is useful for
	// seqs with a variable number of arguments, like
	// [CHAFA_TERM_SEQ_PRIMARY_DEVICE_ATTRIBUTES].
	//
	// Either or both of argsOut and nArgsOut can be NULL, in which case nothing
	// is returned for that parameter.
	TermInfoParseSeqVarargs func(
		term_info *TermInfo,
		seq TermSeq,
		input []string,
		inputLen *int32,
		argsOut *uint32,
		nArgsOut *int32,
	) ParseResult

	// Sets whether seq can be inherited from the outer [TermInfo] when chaining
	// with [TermInfoChain].
	TermInfoSetInheritSeq func(termInfo *TermInfo, seq TermSeq, inherit bool)

	// Specifies if passthrough should be used to convey images in pixelMode.
	TermInfoSetIsPixelPassthroughNeeded func(termInfo *TermInfo, pixelMode PixelMode, pixelPassthroughNeeded bool)

	// Assigns a new name to termInfo. The name should be a short lowercase ASCII
	// string that uniquely identifies the terminal or program described by termInfo.
	TermInfoSetName func(termInfo *TermInfo, name string)

	// Assigns a set of quirks to termInfo.
	TermInfoSetQuirks func(termInfo *TermInfo, quirks TermQuirks)

	// Sets the [SymbolTags] that are likely safe to use with termInfo.
	// The tags are a bitwise OR of flags from the enum.
	TermInfoSetSafeSymbolTags func(termInfo *TermInfo, tags SymbolTags)
)

type TermInfo struct {
	Refs                   int32
	Name                   string
	SeqStr                 [CHAFA_TERM_SEQ_MAX][CHAFA_TERM_SEQ_LENGTH_MAX]byte
	SeqArgs                [CHAFA_TERM_SEQ_MAX][CHAFA_TERM_SEQ_ARGS_MAX]SeqArgInfo
	UnparsedStr            [CHAFA_TERM_SEQ_MAX]string
	PixelPassthroughNeeded [CHAFA_PIXEL_MODE_MAX]uint8
	InheritSeq             [CHAFA_TERM_SEQ_MAX]uint8
	Quirks                 TermQuirks
	SafeSymbolTags         SymbolTags
}

type SeqArgInfo struct {
	IsVarargs bool
	PreLen    uint8
	ArgIndex  uint8
}

type TermSeq int32

const (
	CHAFA_TERM_SEQ_RESET_TERMINAL_SOFT                 TermSeq = 0
	CHAFA_TERM_SEQ_RESET_TERMINAL_HARD                 TermSeq = 1
	CHAFA_TERM_SEQ_RESET_ATTRIBUTES                    TermSeq = 2
	CHAFA_TERM_SEQ_CLEAR                               TermSeq = 3
	CHAFA_TERM_SEQ_INVERT_COLORS                       TermSeq = 4
	CHAFA_TERM_SEQ_CURSOR_TO_TOP_LEFT                  TermSeq = 5
	CHAFA_TERM_SEQ_CURSOR_TO_BOTTOM_LEFT               TermSeq = 6
	CHAFA_TERM_SEQ_CURSOR_TO_POS                       TermSeq = 7
	CHAFA_TERM_SEQ_CURSOR_UP_1                         TermSeq = 8
	CHAFA_TERM_SEQ_CURSOR_UP                           TermSeq = 9
	CHAFA_TERM_SEQ_CURSOR_DOWN_1                       TermSeq = 10
	CHAFA_TERM_SEQ_CURSOR_DOWN                         TermSeq = 11
	CHAFA_TERM_SEQ_CURSOR_LEFT_1                       TermSeq = 12
	CHAFA_TERM_SEQ_CURSOR_LEFT                         TermSeq = 13
	CHAFA_TERM_SEQ_CURSOR_RIGHT_1                      TermSeq = 14
	CHAFA_TERM_SEQ_CURSOR_RIGHT                        TermSeq = 15
	CHAFA_TERM_SEQ_CURSOR_UP_SCROLL                    TermSeq = 16
	CHAFA_TERM_SEQ_CURSOR_DOWN_SCROLL                  TermSeq = 17
	CHAFA_TERM_SEQ_INSERT_CELLS                        TermSeq = 18
	CHAFA_TERM_SEQ_DELETE_CELLS                        TermSeq = 19
	CHAFA_TERM_SEQ_INSERT_ROWS                         TermSeq = 20
	CHAFA_TERM_SEQ_DELETE_ROWS                         TermSeq = 21
	CHAFA_TERM_SEQ_SET_SCROLLING_ROWS                  TermSeq = 22
	CHAFA_TERM_SEQ_ENABLE_INSERT                       TermSeq = 23
	CHAFA_TERM_SEQ_DISABLE_INSERT                      TermSeq = 24
	CHAFA_TERM_SEQ_ENABLE_CURSOR                       TermSeq = 25
	CHAFA_TERM_SEQ_DISABLE_CURSOR                      TermSeq = 26
	CHAFA_TERM_SEQ_ENABLE_ECHO                         TermSeq = 27
	CHAFA_TERM_SEQ_DISABLE_ECHO                        TermSeq = 28
	CHAFA_TERM_SEQ_ENABLE_WRAP                         TermSeq = 29
	CHAFA_TERM_SEQ_DISABLE_WRAP                        TermSeq = 30
	CHAFA_TERM_SEQ_SET_COLOR_FG_DIRECT                 TermSeq = 31
	CHAFA_TERM_SEQ_SET_COLOR_BG_DIRECT                 TermSeq = 32
	CHAFA_TERM_SEQ_SET_COLOR_FGBG_DIRECT               TermSeq = 33
	CHAFA_TERM_SEQ_SET_COLOR_FG_256                    TermSeq = 34
	CHAFA_TERM_SEQ_SET_COLOR_BG_256                    TermSeq = 35
	CHAFA_TERM_SEQ_SET_COLOR_FGBG_256                  TermSeq = 36
	CHAFA_TERM_SEQ_SET_COLOR_FG_16                     TermSeq = 37
	CHAFA_TERM_SEQ_SET_COLOR_BG_16                     TermSeq = 38
	CHAFA_TERM_SEQ_SET_COLOR_FGBG_16                   TermSeq = 39
	CHAFA_TERM_SEQ_BEGIN_SIXELS                        TermSeq = 40
	CHAFA_TERM_SEQ_END_SIXELS                          TermSeq = 41
	CHAFA_TERM_SEQ_REPEAT_CHAR                         TermSeq = 42
	CHAFA_TERM_SEQ_BEGIN_KITTY_IMMEDIATE_IMAGE_V1      TermSeq = 43
	CHAFA_TERM_SEQ_END_KITTY_IMAGE                     TermSeq = 44
	CHAFA_TERM_SEQ_BEGIN_KITTY_IMAGE_CHUNK             TermSeq = 45
	CHAFA_TERM_SEQ_END_KITTY_IMAGE_CHUNK               TermSeq = 46
	CHAFA_TERM_SEQ_BEGIN_ITERM2_IMAGE                  TermSeq = 47
	CHAFA_TERM_SEQ_END_ITERM2_IMAGE                    TermSeq = 48
	CHAFA_TERM_SEQ_ENABLE_SIXEL_SCROLLING              TermSeq = 49
	CHAFA_TERM_SEQ_DISABLE_SIXEL_SCROLLING             TermSeq = 50
	CHAFA_TERM_SEQ_ENABLE_BOLD                         TermSeq = 51
	CHAFA_TERM_SEQ_SET_COLOR_FG_8                      TermSeq = 52
	CHAFA_TERM_SEQ_SET_COLOR_BG_8                      TermSeq = 53
	CHAFA_TERM_SEQ_SET_COLOR_FGBG_8                    TermSeq = 54
	CHAFA_TERM_SEQ_RESET_DEFAULT_FG                    TermSeq = 55
	CHAFA_TERM_SEQ_SET_DEFAULT_FG                      TermSeq = 56
	CHAFA_TERM_SEQ_QUERY_DEFAULT_FG                    TermSeq = 57
	CHAFA_TERM_SEQ_RESET_DEFAULT_BG                    TermSeq = 58
	CHAFA_TERM_SEQ_SET_DEFAULT_BG                      TermSeq = 59
	CHAFA_TERM_SEQ_QUERY_DEFAULT_BG                    TermSeq = 60
	CHAFA_TERM_SEQ_RETURN_KEY                          TermSeq = 61
	CHAFA_TERM_SEQ_BACKSPACE_KEY                       TermSeq = 62
	CHAFA_TERM_SEQ_TAB_KEY                             TermSeq = 63
	CHAFA_TERM_SEQ_TAB_SHIFT_KEY                       TermSeq = 64
	CHAFA_TERM_SEQ_UP_KEY                              TermSeq = 65
	CHAFA_TERM_SEQ_UP_CTRL_KEY                         TermSeq = 66
	CHAFA_TERM_SEQ_UP_SHIFT_KEY                        TermSeq = 67
	CHAFA_TERM_SEQ_DOWN_KEY                            TermSeq = 68
	CHAFA_TERM_SEQ_DOWN_CTRL_KEY                       TermSeq = 69
	CHAFA_TERM_SEQ_DOWN_SHIFT_KEY                      TermSeq = 70
	CHAFA_TERM_SEQ_LEFT_KEY                            TermSeq = 71
	CHAFA_TERM_SEQ_LEFT_CTRL_KEY                       TermSeq = 72
	CHAFA_TERM_SEQ_LEFT_SHIFT_KEY                      TermSeq = 73
	CHAFA_TERM_SEQ_RIGHT_KEY                           TermSeq = 74
	CHAFA_TERM_SEQ_RIGHT_CTRL_KEY                      TermSeq = 75
	CHAFA_TERM_SEQ_RIGHT_SHIFT_KEY                     TermSeq = 76
	CHAFA_TERM_SEQ_PAGE_UP_KEY                         TermSeq = 77
	CHAFA_TERM_SEQ_PAGE_UP_CTRL_KEY                    TermSeq = 78
	CHAFA_TERM_SEQ_PAGE_UP_SHIFT_KEY                   TermSeq = 79
	CHAFA_TERM_SEQ_PAGE_DOWN_KEY                       TermSeq = 80
	CHAFA_TERM_SEQ_PAGE_DOWN_CTRL_KEY                  TermSeq = 81
	CHAFA_TERM_SEQ_PAGE_DOWN_SHIFT_KEY                 TermSeq = 82
	CHAFA_TERM_SEQ_HOME_KEY                            TermSeq = 83
	CHAFA_TERM_SEQ_HOME_CTRL_KEY                       TermSeq = 84
	CHAFA_TERM_SEQ_HOME_SHIFT_KEY                      TermSeq = 85
	CHAFA_TERM_SEQ_END_KEY                             TermSeq = 86
	CHAFA_TERM_SEQ_END_CTRL_KEY                        TermSeq = 87
	CHAFA_TERM_SEQ_END_SHIFT_KEY                       TermSeq = 88
	CHAFA_TERM_SEQ_INSERT_KEY                          TermSeq = 89
	CHAFA_TERM_SEQ_INSERT_CTRL_KEY                     TermSeq = 90
	CHAFA_TERM_SEQ_INSERT_SHIFT_KEY                    TermSeq = 91
	CHAFA_TERM_SEQ_DELETE_KEY                          TermSeq = 92
	CHAFA_TERM_SEQ_DELETE_CTRL_KEY                     TermSeq = 93
	CHAFA_TERM_SEQ_DELETE_SHIFT_KEY                    TermSeq = 94
	CHAFA_TERM_SEQ_F1_KEY                              TermSeq = 95
	CHAFA_TERM_SEQ_F1_CTRL_KEY                         TermSeq = 96
	CHAFA_TERM_SEQ_F1_SHIFT_KEY                        TermSeq = 97
	CHAFA_TERM_SEQ_F2_KEY                              TermSeq = 98
	CHAFA_TERM_SEQ_F2_CTRL_KEY                         TermSeq = 99
	CHAFA_TERM_SEQ_F2_SHIFT_KEY                        TermSeq = 100
	CHAFA_TERM_SEQ_F3_KEY                              TermSeq = 101
	CHAFA_TERM_SEQ_F3_CTRL_KEY                         TermSeq = 102
	CHAFA_TERM_SEQ_F3_SHIFT_KEY                        TermSeq = 103
	CHAFA_TERM_SEQ_F4_KEY                              TermSeq = 104
	CHAFA_TERM_SEQ_F4_CTRL_KEY                         TermSeq = 105
	CHAFA_TERM_SEQ_F4_SHIFT_KEY                        TermSeq = 106
	CHAFA_TERM_SEQ_F5_KEY                              TermSeq = 107
	CHAFA_TERM_SEQ_F5_CTRL_KEY                         TermSeq = 108
	CHAFA_TERM_SEQ_F5_SHIFT_KEY                        TermSeq = 109
	CHAFA_TERM_SEQ_F6_KEY                              TermSeq = 110
	CHAFA_TERM_SEQ_F6_CTRL_KEY                         TermSeq = 111
	CHAFA_TERM_SEQ_F6_SHIFT_KEY                        TermSeq = 112
	CHAFA_TERM_SEQ_F7_KEY                              TermSeq = 113
	CHAFA_TERM_SEQ_F7_CTRL_KEY                         TermSeq = 114
	CHAFA_TERM_SEQ_F7_SHIFT_KEY                        TermSeq = 115
	CHAFA_TERM_SEQ_F8_KEY                              TermSeq = 116
	CHAFA_TERM_SEQ_F8_CTRL_KEY                         TermSeq = 117
	CHAFA_TERM_SEQ_F8_SHIFT_KEY                        TermSeq = 118
	CHAFA_TERM_SEQ_F9_KEY                              TermSeq = 119
	CHAFA_TERM_SEQ_F9_CTRL_KEY                         TermSeq = 120
	CHAFA_TERM_SEQ_F9_SHIFT_KEY                        TermSeq = 121
	CHAFA_TERM_SEQ_F10_KEY                             TermSeq = 122
	CHAFA_TERM_SEQ_F10_CTRL_KEY                        TermSeq = 123
	CHAFA_TERM_SEQ_F10_SHIFT_KEY                       TermSeq = 124
	CHAFA_TERM_SEQ_F11_KEY                             TermSeq = 125
	CHAFA_TERM_SEQ_F11_CTRL_KEY                        TermSeq = 126
	CHAFA_TERM_SEQ_F11_SHIFT_KEY                       TermSeq = 127
	CHAFA_TERM_SEQ_F12_KEY                             TermSeq = 128
	CHAFA_TERM_SEQ_F12_CTRL_KEY                        TermSeq = 129
	CHAFA_TERM_SEQ_F12_SHIFT_KEY                       TermSeq = 130
	CHAFA_TERM_SEQ_RESET_COLOR_FG                      TermSeq = 131
	CHAFA_TERM_SEQ_RESET_COLOR_BG                      TermSeq = 132
	CHAFA_TERM_SEQ_RESET_COLOR_FGBG                    TermSeq = 133
	CHAFA_TERM_SEQ_RESET_SCROLLING_ROWS                TermSeq = 134
	CHAFA_TERM_SEQ_SAVE_CURSOR_POS                     TermSeq = 135
	CHAFA_TERM_SEQ_RESTORE_CURSOR_POS                  TermSeq = 136
	CHAFA_TERM_SEQ_SET_SIXEL_ADVANCE_DOWN              TermSeq = 137
	CHAFA_TERM_SEQ_SET_SIXEL_ADVANCE_RIGHT             TermSeq = 138
	CHAFA_TERM_SEQ_ENABLE_ALT_SCREEN                   TermSeq = 139
	CHAFA_TERM_SEQ_DISABLE_ALT_SCREEN                  TermSeq = 140
	CHAFA_TERM_SEQ_BEGIN_SCREEN_PASSTHROUGH            TermSeq = 141
	CHAFA_TERM_SEQ_END_SCREEN_PASSTHROUGH              TermSeq = 142
	CHAFA_TERM_SEQ_BEGIN_TMUX_PASSTHROUGH              TermSeq = 143
	CHAFA_TERM_SEQ_END_TMUX_PASSTHROUGH                TermSeq = 144
	CHAFA_TERM_SEQ_BEGIN_KITTY_IMMEDIATE_VIRT_IMAGE_V1 TermSeq = 145
	CHAFA_TERM_SEQ_QUERY_PRIMARY_DEVICE_ATTRIBUTES     TermSeq = 146
	CHAFA_TERM_SEQ_PRIMARY_DEVICE_ATTRIBUTES           TermSeq = 147
	CHAFA_TERM_SEQ_QUERY_TEXT_AREA_SIZE_CELLS          TermSeq = 148
	CHAFA_TERM_SEQ_TEXT_AREA_SIZE_CELLS                TermSeq = 149
	CHAFA_TERM_SEQ_QUERY_TEXT_AREA_SIZE_PX             TermSeq = 150
	CHAFA_TERM_SEQ_TEXT_AREA_SIZE_PX                   TermSeq = 151
	CHAFA_TERM_SEQ_QUERY_CELL_SIZE_PX                  TermSeq = 152
	CHAFA_TERM_SEQ_CELL_SIZE_PX                        TermSeq = 153
	CHAFA_TERM_SEQ_MAX                                 TermSeq = 154
)

const CHAFA_TERM_SEQ_LENGTH_MAX = 96

const CHAFA_TERM_SEQ_ARGS_MAX = 24

type TermQuirks int32

const (
	CHAFA_TERM_QUIRK_SIXEL_OVERSHOOT TermQuirks = (1 << 0)
)

type ParseResult int32

const (
	CHAFA_PARSE_SUCCESS ParseResult = 0
	CHAFA_PARSE_FAILURE ParseResult = 1
	CHAFA_PARSE_AGAIN   ParseResult = 2
)
