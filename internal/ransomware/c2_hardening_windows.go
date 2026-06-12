//go:build windows

package ransomware

import (
	"syscall"
	"unsafe"
)

var (
	user32            = syscall.NewLazyDLL("user32.dll")
	getLastInputInfo  = user32.NewProc("GetLastInputInfo")
	getTickCount      = user32.NewProc("GetTickCount")
)

type lastInputInfo struct {
	cbSize uint32
	dwTime uint32
}

func isUserActiveNative() bool {
	var lii lastInputInfo
	lii.cbSize = uint32(unsafe.Sizeof(lii))
	ret, _, _ := getLastInputInfo.Call(uintptr(unsafe.Pointer(&lii)))
	if ret == 0 {
		return true
	}
	uptime, _, _ := getTickCount.Call()
	idleTime := uint32(uptime) - lii.dwTime
	return idleTime < 60000
}
