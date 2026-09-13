//go:build linux

package store

import (
	"fmt"
	"syscall"
)

// availableMemory returns total system memory in bytes on Linux.
// Returns 0 and logs if it cannot be determined.
func availableMemory() (uint64, error) {
	var info syscall.Sysinfo_t
	if err := syscall.Sysinfo(&info); err != nil {
		return 0, fmt.Errorf("sysinfo: %w", err)
	}
	// Totalram is in bytes on Linux.
	return info.Totalram, nil
}
