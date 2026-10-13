package chafa

import "runtime"

var (
	// Creates a new, blank [TermDb].
	TermDbNew func() *TermDb

	// Creates a new [TermDb] that's a copy of termDb.
	TermDbCopy func(termDb *TermDb) *TermDb

	// Adds a reference to termDb.
	TermDbRef func(termDb *TermDb)

	// Removes a reference from termDb.
	TermDbUnref func(termDb *TermDb)

	// Gets the global [TermDb]. This can normally be used safely in a read-only
	// capacity. The caller should not unref the returned object.
	TermDbGetDefault func() *TermDb

	// Builds a new [TermInfo] with capabilities implied by the provided
	// environment variables (principally the TERM variable, but also others).
	termDbDetect func(termDb *TermDb, envp **byte) *TermInfo

	// Builds a new [TermInfo] with fallback control sequences. This can be used
	// with unknown but presumably modern terminals, or to supplement missing
	// capabilities in a detected terminal.
	//
	// Fallback control sequences may cause unpredictable behavior and should
	// only be used as a last resort.
	TermDbGetFallbackInfo func(termDb *TermDb) *TermInfo
)

func TermDbDetect(termDb *TermDb, envp []string) *TermInfo {
	allocated := make([][]byte, len(envp))
	// envp is read by chafa as a NULL-terminated strv, so leave room for the
	// trailing nil.
	ptrs := make([]*byte, len(envp)+1)
	for i, s := range envp {
		allocated[i] = append([]byte(s), 0)
		ptrs[i] = &allocated[i][0]
	}

	info := termDbDetect(termDb, &ptrs[0])
	runtime.KeepAlive(allocated)
	runtime.KeepAlive(ptrs)
	return info
}

type TermDb struct {
	Refs int32
}
