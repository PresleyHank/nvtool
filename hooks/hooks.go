package hooks

import (
	"image"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type MWMoveState struct {
	moveable   bool
	prevMouseX int
	prevMouseY int
}

func UseMounted(mounted *bool, fn func()) {
	if !*mounted {
		*mounted = true
		fn()
	}
}

func UseWindowMove(glfwWindow *glfw.Window, dragArea image.Rectangle, mwMoveState *MWMoveState) {
	mouseX, mouseY := glfwWindow.GetCursorPos()
	prevPosX, prevPosY := glfwWindow.GetPos()

	if g.IsMouseClicked(0) {
		clickPoint := image.Point{X: int(mouseX / float64(imgui.DPIScale)), Y: int(mouseY / float64(imgui.DPIScale))}
		mwMoveState.moveable = clickPoint.In(dragArea)
		mwMoveState.prevMouseX = int(mouseX)
		mwMoveState.prevMouseY = int(mouseY)
	}
	if mwMoveState.moveable && g.IsMouseDown(0) {
		offsetX := int(mouseX) - mwMoveState.prevMouseX
		offsetY := int(mouseY) - mwMoveState.prevMouseY
		glfwWindow.SetPos(prevPosX+int(offsetX), prevPosY+int(offsetY))
	}
}
