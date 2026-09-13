//go:build !linux && !darwin

package store

// availableMemory returns 0 on non-Linux, non-macOS platforms.
func availableMemory() (uint64, error) {
	return 0, nil
}
