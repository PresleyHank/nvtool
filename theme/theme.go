package theme

import (
	"image/color"

	"github.com/AllenDang/giu/imgui"
)

type EffectMethod struct {
	Push func()
	Pop  func()
}

type Palette struct {
	Tint   color.RGBA
	Hover  color.RGBA
	Active color.RGBA
	BG     color.RGBA
}

func RGBAToVec4(rgba color.RGBA) imgui.Vec4 {
	return imgui.Vec4{
		X: float32(rgba.R) / 255,
		Y: float32(rgba.G) / 255,
		Z: float32(rgba.B) / 255,
		W: float32(rgba.A) / 255,
	}
}

func SetThemeDark(style *imgui.Style) {
	style.SetColor(imgui.StyleColorText, RGBAToVec4(color.RGBA{230, 230, 230, 230}))
	style.SetColor(imgui.StyleColorWindowBg, RGBAToVec4(color.RGBA{50, 50, 50, 250}))
	style.SetColor(imgui.StyleColorFrameBg, RGBAToVec4(color.RGBA{10, 10, 10, 240}))
	style.SetColor(imgui.StyleColorButton, RGBAToVec4(color.RGBA{100, 100, 100, 255}))
	style.SetColor(imgui.StyleColorButtonHovered, RGBAToVec4(color.RGBA{120, 120, 120, 240}))
	style.SetColor(imgui.StyleColorButtonActive, RGBAToVec4(color.RGBA{80, 80, 80, 245}))
	style.SetColor(imgui.StyleColorTab, RGBAToVec4(color.RGBA{0, 0, 0, 0}))
	style.SetColor(imgui.StyleColorTabActive, RGBAToVec4(color.RGBA{18, 150, 219, 255}))
	style.SetColor(imgui.StyleColorTabHovered, RGBAToVec4(color.RGBA{18, 150, 219, 255}))
	style.SetColor(imgui.StyleColorScrollbarBg, RGBAToVec4(color.RGBA{27, 27, 27, 250}))
	style.SetColor(imgui.StyleColorScrollbarGrab, RGBAToVec4(color.RGBA{73, 73, 73, 255}))
	style.SetColor(imgui.StyleColorScrollbarGrabActive, RGBAToVec4(color.RGBA{164, 164, 164, 255}))
	style.SetColor(imgui.StyleColorScrollbarGrabHovered, RGBAToVec4(color.RGBA{140, 140, 140, 250}))
}

func UseLayoutFlat() EffectMethod {
	return EffectMethod{
		func() {
			imgui.PushStyleVarFloat(imgui.StyleVarWindowBorderSize, 0)
			imgui.PushStyleVarFloat(imgui.StyleVarFrameBorderSize, 0)
			imgui.PushStyleVarFloat(imgui.StyleVarChildBorderSize, 0)
			imgui.PushStyleVarFloat(imgui.StyleVarFrameRounding, 0)
			imgui.PushStyleVarFloat(imgui.StyleVarChildRounding, 0)
			imgui.PushStyleVarVec2(imgui.StyleVarWindowPadding, imgui.Vec2{X: 8, Y: 6})
		},
		func() { imgui.PopStyleVarV(6) },
	}
}

func UseStyleDarkButton() EffectMethod {
	return EffectMethod{
		func() {
			imgui.PushStyleVarFloat(imgui.StyleVarFrameRounding, 0)
			imgui.PushStyleVarVec2(imgui.StyleVarFramePadding, imgui.Vec2{X: 0, Y: 0})
			imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 0.125, Y: 0.125, Z: 0.125, W: 1})
			imgui.PushStyleColor(imgui.StyleColorButtonHovered, imgui.Vec4{X: 0.125, Y: 0.125, Z: 0.125, W: 0.90})
			imgui.PushStyleColor(imgui.StyleColorButtonActive, imgui.Vec4{X: 0.125, Y: 0.125, Z: 0.125, W: 0.95})
		},
		func() {
			imgui.PopStyleVarV(2)
			imgui.PopStyleColorV(3)
		},
	}
}

func UseFont(font imgui.Font) EffectMethod {
	return EffectMethod{
		func() {
			imgui.PushFont(font)
		},
		func() {
			imgui.PopFont()
		},
	}
}
