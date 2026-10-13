package chafa

var (
	// Gets a list of the platform-specific features this library was built with.
	GetBuiltinFeatures func() Features

	// Gets a list of the platform-specific features that are built in and usable on the runtime platform.
	GetSupportedFeatures func() Features

	// Takes a set of flags potentially returned from [GetBuiltinFeatures] or
	// [GetSupportedFeatures] and generates a human-readable ASCII string descriptor.
	DescribeFeatures func(features Features) string

	// Queries the maximum number of worker threads to use for parallel processing.
	GetNThreads func() int32

	// Sets the maximum number of worker threads to use for parallel processing,
	// or -1 to determine this automatically. The default is -1.
	//
	// Setting this to 0 or 1 will avoid using thread pools and instead perform
	// all processing in the main thread.
	SetNThreads func(n int32)

	// Queries the number of worker threads that will actually be used for
	// parallel processing.
	GetNActualThreads func() int32
)

type Features int32

const (
	CHAFA_FEATURE_MMX    Features = 0
	CHAFA_FEATURE_SSE41  Features = 1
	CHAFA_FEATURE_POPCNT Features = 2
	CHAFA_FEATURE_AVX2   Features = 3
)
