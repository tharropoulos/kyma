package chafa

var (
	// Creates a new [Placement] for the specified image and ID.
	// If id <= 0, an ID is assigned automatically.
	PlacementNew func(image *Image, id int32) *Placement

	// Adds a reference to placement.
	PlacementRef func(placement *Placement)

	// Removes a reference from placement. When remaining references drops to zero,
	// the placement is freed and can no longer be used.
	PlacementUnref func(placement *Placement)

	// Gets the tucking policy of placement. This describes how the image is
	// resized to fit placement's extents, and defaults to [CHAFA_TUCK_STRETCH].
	PlacementGetTuck func(placement *Placement) Tuck

	// Sets the tucking policy for placement to tuck . This describes how the
	// image is resized to fit placement 's extents, and defaults to [CHAFA_TUCK_STRETCH].
	PlacementSetTuck func(placement *Placement, tuck Tuck)

	// Gets the horizontal alignment of placement. This determines how any
	// padding added by the tucking policy is distributed, and defaults to [CHAFA_ALIGN_START].
	PlacementGetHAlign func(placement *Placement) Align

	// Sets the horizontal alignment of placement. This determines how any
	// padding added by the tucking policy is distributed, and defaults to [CHAFA_ALIGN_START].
	PlacementSetHAlign func(placement *Placement, align Align)

	// Gets the vertical alignment of placement. This determines how any padding
	// added by the tucking policy is distributed, and defaults to [CHAFA_ALIGN_START].
	PlacementGetVAlign func(placement *Placement) Align

	// Sets the vertical alignment of placement . This determines how any
	// padding added by the tucking policy is distributed.
	PlacementSetVAlign func(placement *Placement, align Align)
)

type Placement struct {
	Refs int32

	Image          *Image
	Id             int32
	Halign, Valign Align
	Tuck           Tuck
}

type Tuck int32

const (
	CHAFA_TUCK_STRETCH       Tuck = 0
	CHAFA_TUCK_FIT           Tuck = 1
	CHAFA_TUCK_SHRINK_TO_FIT Tuck = 2
	CHAFA_TUCK_MAX           Tuck = 3
)

type Align int32

const (
	CHAFA_ALIGN_START  Align = 0
	CHAFA_ALIGN_END    Align = 1
	CHAFA_ALIGN_CENTER Align = 2
	CHAFA_ALIGN_MAX    Align = 3
)
