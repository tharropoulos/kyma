package chafa

// Calculates an optimal geometry for a [Canvas] given the width and height
// of an input image, maximum width and height of the canvas, font ratio, zoom and
// stretch preferences.
//
// srcWidth and srcHeight must both be zero or greater.
//
// destWidthInout and destHeightInout must point to integers containing the
// maximum dimensions of the canvas in character cells. These will be replaced
// by the calculated values, which may be zero if one of the input dimensions is
// zero. If one or both of the input parameters is negative, they will be treated
// as unspecified and calculated based on the remaining parameters and aspect ratio.
//
// fontRatio is the font's width divided by its height. 0.5 is a typical value.
var CalcCanvasGeometry func(
	srcWidth, srcHeight int32,
	destWidthInout, destHeightInout *int32,
	fontRatio float32,
	zoom, stretch bool,
)
