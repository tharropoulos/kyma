package chafa

var (
	// Creates a new [Frame] containing a copy of the image data pointed to by data.
	FrameNew func(data []uint8, pixelType PixelType, width, height, rowstride int32) *Frame

	// Creates a new [Frame] embedding the data pointer. It's the caller's
	// responsibility to ensure the pointer remains valid for the lifetime of
	// the frame. The frame will not free the buffer when its reference count
	// drops to zero.
	//
	// THIS IS DANGEROUS API which should only be used when the life cycle of
	// the frame is short, stealing the buffer is impossible, and copying would
	// cause unacceptable performance degradation.
	//
	// Use [FrameNew] instead.
	FrameNewBorrow func(data []uint8, pixelType PixelType, width, height, rowstride int32) *Frame

	// Creates a new [Frame], which takes ownership of the data buffer. The
	// buffer will be freed with g_free() when the frame's reference count drops
	// to zero.
	FrameNewSteal func(data []uint8, pixelType PixelType, width, height, rowstride int32) *Frame

	// Adds a reference to frame.
	FrameRef func(frame *Frame)

	// Removes a reference from frame. When the reference count drops to zero,
	// the frame is freed and can no longer be used.
	FrameUnref func(frame *Frame)
)

type Frame struct {
	Refs                     int32
	PixelType                PixelType
	Width, Height, Rowstride int32

	Data []uint8

	DataIsOwned bool
}
