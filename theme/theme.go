package theme

import (
	"image/color"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
)

// WithDarkTheme apply theme to widget builder
func WithDarkTheme(builder func()) {
	imgui.PushStyleVarFloat(imgui.StyleVarWindowBorderSize, 0)
	imgui.PushStyleVarFloat(imgui.StyleVarFrameBorderSize, 0)
	imgui.PushStyleVarFloat(imgui.StyleVarChildBorderSize, 0)
	imgui.PushStyleVarFloat(imgui.StyleVarFrameRounding, 0)
	imgui.PushStyleVarFloat(imgui.StyleVarChildRounding, 0)
	imgui.PushStyleVarVec2(imgui.StyleVarWindowPadding, imgui.Vec2{X: 8, Y: 6})
	imgui.PushStyleColor(imgui.StyleColorTab, imgui.Vec4{X: 0.125, Y: 0.125, Z: 0.125, W: 1.0})
	imgui.PushStyleColor(imgui.StyleColorTabActive, imgui.Vec4{X: 0.19, Y: 0.19, Z: 0.19, W: 0.941})
	// imgui.PushStyleColor(imgui.StyleColorTabActive, imgui.Vec4{X: 0.815, Y: 0.007, Z: 0.105, W: 1.0})
	// imgui.PushStyleColor(imgui.StyleColorTabActive, imgui.Vec4{X: 0.45, Y: 0.69, Z: 0.117, W: 1.0})
	imgui.PushStyleColor(imgui.StyleColorTabHovered, imgui.Vec4{X: 0.125, Y: 0.125, Z: 0.125, W: 0.941})
	imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 0.80, Y: 0.80, Z: 0.80, W: 1.0})
	g.PushColorWindowBg(color.RGBA{50, 50, 50, 250})
	g.PushColorFrameBg(color.RGBA{10, 10, 10, 240})
	g.PushColorButton(color.RGBA{100, 100, 100, 255})
	g.PushColorButtonHovered(color.RGBA{120, 120, 120, 240})
	g.PushColorButtonActive(color.RGBA{80, 80, 80, 245})
	builder()
	g.PopStyleColorV(9)
	imgui.PopStyleVarV(6)
}

// BeginStyleButtonDark push dark button style
func BeginStyleButtonDark() {
	imgui.PushStyleVarFloat(imgui.StyleVarFrameRounding, 0)
	imgui.PushStyleVarVec2(imgui.StyleVarFramePadding, imgui.Vec2{X: 0, Y: 0})
	imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 0.125, Y: 0.125, Z: 0.125, W: 1})
	imgui.PushStyleColor(imgui.StyleColorButtonHovered, imgui.Vec4{X: 0.125, Y: 0.125, Z: 0.125, W: 0.90})
	imgui.PushStyleColor(imgui.StyleColorButtonActive, imgui.Vec4{X: 0.125, Y: 0.125, Z: 0.125, W: 0.95})
}

// EndStyleButtonDark pop dark button style
func EndStyleButtonDark() {
	imgui.PopStyleVarV(2)
	imgui.PopStyleColorV(3)
}
