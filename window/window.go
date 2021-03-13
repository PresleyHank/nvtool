package window

import (
	"image"
	"unsafe"

	"github.com/Nicify/nvtool/win"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type CompositionAttribute struct {
	AccentState int
	Flags       int
	Color       int
	AnimationID int
}

type DragConfig struct {
	DragArea   image.Rectangle
	Moveable   bool
	PrevMouseX int
	PrevMouseY int
}

type GLFWWindowConfig struct {
	Icon48px             *image.RGBA
	CompositionAttribute *CompositionAttribute
	DragConfig           *DragConfig
	FocusCallback        func(focused bool)
}

func ApplyWindowConfig(window *glfw.Window, config *GLFWWindowConfig) {
	dpi, _ := window.GetContentScale()
	window.SetIcon([]image.Image{config.Icon48px})
	hwnd := win.HWND(unsafe.Pointer(window.GetWin32Window()))
	attr := config.CompositionAttribute
	win.SetWindowCompositionAttribute(hwnd, attr.AccentState, attr.Flags, attr.Color, attr.AnimationID)
	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft && action == glfw.Press {
			xpos, ypos := w.GetCursorPos()
			clickPoint := image.Point{X: int(xpos / float64(dpi)), Y: int(ypos / float64(dpi))}
			config.DragConfig.Moveable = clickPoint.In(config.DragConfig.DragArea)
			config.DragConfig.PrevMouseX = int(xpos)
			config.DragConfig.PrevMouseY = int(ypos)
		}
	})
	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		prevPosX, prevPosY := w.GetPos()
		if config.DragConfig.Moveable && window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
			offsetX := int(xpos) - config.DragConfig.PrevMouseX
			offsetY := int(ypos) - config.DragConfig.PrevMouseY
			w.SetPos(prevPosX+int(offsetX), prevPosY+int(offsetY))
		}
	})
	window.SetFocusCallback(func(w *glfw.Window, focused bool) {
		if config.FocusCallback != nil {
			config.FocusCallback(focused)
		}
	})
}
