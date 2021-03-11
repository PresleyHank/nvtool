package main

import (
	"image"
	"image/color"
	"runtime"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	applation "github.com/Nicify/nvtool/app"
	assets "github.com/Nicify/nvtool/assets"
	"github.com/Nicify/nvtool/helper"
	"github.com/Nicify/nvtool/hooks"
	window "github.com/Nicify/nvtool/window"
	"github.com/Nicify/theme"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	app := applation.GetInstance()
	defer app.NVENC.Stop()
	unlock := app.SingleInstanceLock()
	defer unlock()

	mw := g.NewMasterWindow("NVTool", 750, 435, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFrameless|g.MasterWindowFlagsTransparent, app.LoadFont)
	mw.SetBgColor(color.RGBA{0, 0, 0, 0})
	mw.SetDropCallback(app.OnDrop)

	currentStyle := imgui.CurrentStyle()
	theme.SetThemeDark(&currentStyle)

	platform := g.Context.GetPlatform().(*imgui.GLFW)
	glfwWindow := platform.GetWindow()
	glfwWindow.SetOpacity(0)

	data, _ := assets.EmbedFS.ReadFile("embed/icon_48px.png")
	icon48px, _ := helper.LoadImageFromMemory(data)
	window.ApplyWindowConfig(glfwWindow, &window.GLFWWindowConfig{
		Icon48px:             icon48px,
		TPS:                  240,
		CompositionAttribute: &window.CompositionAttribute{AccentState: 3, Flags: 0, Color: 0, AnimationID: 0},
		FocusCallback: func(w *glfw.Window, focused bool) {
			if focused {
				w.SetOpacity(0.98)
				return
			}
			w.SetOpacity(1)
		},
	})
	go app.InstallCore()
	app.Mount(&applation.Window{
		GLFWWindow:  glfwWindow,
		MW:          mw,
		MWMoveState: &hooks.MWMoveState{},
		MWDragArea:  image.Rectangle{image.Point{}, image.Point{750, 30}},
	}).Run()
}
