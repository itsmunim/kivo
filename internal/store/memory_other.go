//go:build !linux

package store

// availableMemory returns 0 on non-Linux platforms.
func availableMemory() (uint64, error) {
	return 0, nil
}
