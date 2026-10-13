package chafa

var (
	// Creates a new [Image]. The image is initially transparent
	// and dimensionless.
	ImageNew func() *Image

	// Adds a reference to image.
	ImageRef func(image *Image)

	// Removes a reference from image. When the reference count drops to zero,
	// the image is freed and can no longer be used.
	ImageUnref func(image *Image)

	// Assigns frame as the content for image. The image will keep its own
	// reference to the frame.
	ImageSetFrame func(image *Image, frame *Frame)
)

type Image struct {
	Refs  int32
	Frame *Frame
}
