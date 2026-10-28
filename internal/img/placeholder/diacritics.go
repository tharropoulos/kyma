package placeholder

// diacritics are the combining marks that number the rows and columns of an
// image's placeholder cells, the mark at index i standing for i. They are the
// first 35 of kitty's rowcolumn-diacritics.txt, the table chafa and kitty's own
// kitten carry. The marks after these include Hebrew and Arabic accents that
// go-runewidth, and reflow with it, measure as a cell wide, which would throw
// out every transition that cuts lines with them. Images are drawn in tiles
// small enough to be numbered with these instead.
var diacritics = [...]rune{
	0x0305, 0x030D, 0x030E, 0x0310, 0x0312, 0x033D, 0x033E,
	0x033F, 0x0346, 0x034A, 0x034B, 0x034C, 0x0350, 0x0351,
	0x0352, 0x0357, 0x035B, 0x0363, 0x0364, 0x0365, 0x0366,
	0x0367, 0x0368, 0x0369, 0x036A, 0x036B, 0x036C, 0x036D,
	0x036E, 0x036F, 0x0483, 0x0484, 0x0485, 0x0486, 0x0487,
}
