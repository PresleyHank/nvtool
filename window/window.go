package window

import (
	"image"
	"unsafe"

	g "github.com/AllenDang/giu"
	"github.com/Nicify/nvtool/win"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type CompositionAttribute struct {
	AccentState int
	Flags       int
	Color       int
	AnimationID int
}

type GLFWWindowConfig struct {
	Icon48px             *image.RGBA
	TPS                  int
	CompositionAttribute *CompositionAttribute
	FocusCallback        func(focused bool)
}

func ApplyWindowConfig(window *glfw.Window, config *GLFWWindowConfig) {
	window.SetIcon([]image.Image{config.Icon48px})
	hwnd := win.HWND(unsafe.Pointer(window.GetWin32Window()))
	attr := config.CompositionAttribute
	win.SetWindowCompositionAttribute(hwnd, attr.AccentState, attr.Flags, attr.Color, attr.AnimationID)
	g.Context.GetPlatform().SetTPS(config.TPS)
	window.SetFocusCallback(func(w *glfw.Window, focused bool) {
		if config.FocusCallback != nil {
			config.FocusCallback(focused)
		}
	})
}
