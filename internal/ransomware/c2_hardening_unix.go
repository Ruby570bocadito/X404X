//go:build !windows

package ransomware

import (
	"os"
	"os/exec"
)

func isUserActiveNative() bool {
	if data, err := os.ReadFile("/sys/class/tty/tty0/active"); err == nil {
		if len(data) > 3 {
			return true
		}
	}
	cmd := exec.Command("w", "-h")
	out, err := cmd.Output()
	if err == nil && len(out) > 0 {
		return true
	}
	return false
}
