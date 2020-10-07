package win

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type HWND uintptr

type AccentPolicy struct {
	nAccentState int
	nFlags       int
	nColor       int
	nAnimationID int
}

type WinCompAttrData struct {
	nAttribute int
	pData      *AccentPolicy
	ulDataSize uintptr
}

const (
	SW_HIDE            = 0
	SW_NORMAL          = 1
	SW_SHOWNORMAL      = 1
	SW_SHOWMINIMIZED   = 2
	SW_MAXIMIZE        = 3
	SW_SHOWMAXIMIZED   = 3
	SW_SHOWNOACTIVATE  = 4
	SW_SHOW            = 5
	SW_MINIMIZE        = 6
	SW_SHOWMINNOACTIVE = 7
	SW_SHOWNA          = 8
	SW_RESTORE         = 9
	SW_SHOWDEFAULT     = 10
	SW_FORCEMINIMIZE   = 11
)

var (
	libuser32                     *windows.LazyDLL
	showWindow                    *windows.LazyProc
	setWindowCompositionAttribute *windows.LazyProc
)

func init() {
	libuser32 = windows.NewLazySystemDLL("user32.dll")
	showWindow = libuser32.NewProc("ShowWindow")
	setWindowCompositionAttribute = libuser32.NewProc("SetWindowCompositionAttribute")
}

// SetWindowCompositionAttribute set the composition attribute of window
func SetWindowCompositionAttribute(hwnd HWND) (r1 uintptr, r2 uintptr, lastErr error) {
	accent := AccentPolicy{3, 0, 0, 0}
	data := WinCompAttrData{19, &accent, unsafe.Sizeof(accent)}
	return setWindowCompositionAttribute.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&data)),
	)
}

func ShowWindow(hWnd HWND, nCmdShow int32) (r1 uintptr, r2 uintptr, lastErr error) {
	return showWindow.Call(
		uintptr(hWnd),
		uintptr(nCmdShow),
	)
}
