//go:build darwin

package store

import (
	"golang.org/x/sys/unix"
)

// availableMemory returns total system memory in bytes on macOS.
// Returns 0 if it cannot be determined.
func availableMemory() (uint64, error) {
	mem, err := unix.SysctlUint64("hw.memsize")
	if err != nil {
		return 0, err
	}
	return mem, nil
}
