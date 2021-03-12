package app

import (
	"fmt"
	"time"
	"unsafe"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	c "github.com/Nicify/customwidget"
	"github.com/Nicify/nvtool/helper"
	"github.com/Nicify/nvtool/hooks"
	"github.com/Nicify/nvtool/nvenc"
	"github.com/Nicify/nvtool/preset"
	"github.com/Nicify/nvtool/win"
	"github.com/Nicify/theme"
)

const (
	windowPadding = 8
	contentWidth  = 734
	buttonWidth   = 68
	buttonHeight  = 24
)

func (app *Application) Render() {
	hooks.UseMounted(&app.mounted, app.OnMounted)
	hooks.UseWindowMove(app.Window.GLFWWindow, app.Window.MWDragArea, app.Window.MWMoveState)
	isEncoding := app.NVENC.IsEncoding()
	inputDisableFlag := app.ShouldDisableInput(isEncoding)
	useLayoutFlat := theme.UseLayoutFlat()
	useStyleButtonDark := theme.UseStyleButtonDark()
	defer useLayoutFlat.Pop()
	useLayoutFlat.Push()
	g.SingleWindow("NVTool").Layout(
		g.Group().Layout(
			g.Line(
				g.Image(app.Textures.texLogo).Size(18, 18),
				g.Label("NVENC Video Toolbox"),
				g.Dummy(-83, 0),
				g.Custom(useStyleButtonDark.Push),
				g.Button(".").Size(20, 20),
				g.Button("_").Size(20, 20).OnClick(func() {
					win.ShowWindow(win.HWND(unsafe.Pointer(app.Window.GLFWWindow.GetWin32Window())), win.SW_FORCEMINIMIZE)
				}),
				g.ImageButton(app.Textures.texButtonClose).Size(20, 20).OnClick(func() {
					app.Window.GLFWWindow.SetShouldClose(true)
				}),
				g.Custom(useStyleButtonDark.Pop),
			),
		),
		g.TabBar("maintab").Layout(
			g.TabItem("Encode").Layout(
				g.Child("control").Border(false).Flags(inputDisableFlag).Size(contentWidth, 92).Layout(
					g.Spacing(),
					g.Line(
						g.InputText("##video", &app.State.InputPath).Size(-((windowPadding+buttonWidth)/imgui.DPIScale)),
						c.WithHiDPIFont(app.Fonts.fontIosevka, app.Fonts.fontTamzenb, g.Layout{g.Button("Video").Size(buttonWidth, buttonHeight).OnClick(app.OnInputClick)}),
					),

					g.Spacing(),
					g.Line(
						g.InputText("##output", &app.State.OutputPath).Size(-((windowPadding+buttonWidth)/imgui.DPIScale)),
						c.WithHiDPIFont(app.Fonts.fontIosevka, app.Fonts.fontTamzenb, g.Layout{g.Button("Output").Size(buttonWidth, buttonHeight).OnClick(app.OnOutputClick)}),
					),

					g.Spacing(),
					g.Line(
						g.Label("Preset"),
						g.Combo("##preset", nvenc.PresetOptions[app.EncodingPreset.Preset], nvenc.PresetOptions, &app.EncodingPreset.Preset).Size(50),

						g.Label("Quality"),
						g.InputInt("##quality", &app.EncodingPreset.Quality).Size(24/imgui.DPIScale),

						g.Label("Bitrate"),
						g.InputInt("##bitrate", &app.EncodingPreset.Bitrate).Size(60/imgui.DPIScale),

						g.Label("AQ"),
						g.Combo("##aq", nvenc.AQOptionsForPreview[app.EncodingPreset.AQ], nvenc.AQOptionsForPreview, &app.EncodingPreset.AQ).Size(92),
						g.Label("-"),
						g.InputInt("##strength", &app.EncodingPreset.AQStrength).Size(24/imgui.DPIScale).OnChange(func() {
							app.EncodingPreset.AQStrength = helper.LimitValue(app.EncodingPreset.AQStrength, 0, 15)
						}),

						g.Checkbox("HEVC", &app.EncodingPreset.HEVC),

						g.Checkbox("Resize", &app.EncodingPreset.Resize),
						g.InputText("##outputRes", &app.EncodingPreset.OutputRes).Size(80/imgui.DPIScale).Flags(g.InputTextFlagsCallbackAlways).OnChange(func() {
							app.EncodingPreset.OutputRes = helper.LimitResValue(app.EncodingPreset.OutputRes)
						}),
					),
				),

				g.Spacing(),
				g.InputTextMultiline("##NVENCLog", &app.State.NVENCLog).Size(contentWidth, 200).Flags(g.InputTextFlagsReadOnly),
				g.Custom(func() {
					if isEncoding && time.Now().Second()%2 == 0 {
						imgui.BeginChild("##NVENCLog")
						imgui.SetScrollHereY(1)
						imgui.EndChild()
					}
				}),

				g.Spacing(),
				g.ProgressBar(app.State.Percent).Size(contentWidth, 20).Overlay(""),

				g.Line(
					g.Dummy(0, 5),
				),
				g.Line(
					g.Condition(app.State.GPUName != "",
						g.Layout{
							g.Line(
								g.Image(app.Textures.texGraphicsCard).Size(18, 18),
								g.Label(fmt.Sprintf("%s | GPU:%v%% VE:%v%% VD:%v%%", app.State.GPUName, app.Usage.GPU, app.Usage.VE, app.Usage.VD))),
						},
						g.Layout{
							g.Line(g.Dummy(18, 18)),
						},
					),

					g.Dummy(-(windowPadding+buttonWidth), 24),
					c.WithHiDPIFont(app.Fonts.fontIosevka, app.Fonts.fontTamzenb, g.Layout{g.Condition(isEncoding,
						g.Layout{g.Button("Cancel").Size(buttonWidth, buttonHeight).OnClick(app.NVENC.Stop)},
						g.Layout{g.Button("Run").Size(buttonWidth, buttonHeight).OnClick(app.OnRunClick)},
					)}),
				),
			),

			g.TabItem("Filter").Layout(
				g.Dummy(contentWidth, 5),
				g.Child("FilterContent").Border(false).Size(contentWidth, 0).Layout(

					g.Label("NoiseReduce"),
					g.Child("NoiseReduce").Border(false).Size(contentWidth, 150).Layout(
						g.Custom(func() {
							imgui.PushStyleColor(imgui.StyleColorChildBg, imgui.Vec4{X: 0.12, Y: 0.12, Z: 0.12, W: 0.99})
						}),
						g.Line(
							g.Child("KNN").Border(true).Size((contentWidth-8)*0.5, 0).Flags(inputDisableFlag).Layout(
								g.Line(g.Checkbox("KNN", &app.VPPSwitches.VPPKNN), g.Dummy(-52, 0), g.Button("Reset##ResetKNN").OnClick(func() {
									app.VPPParams.VPPKNNParam = preset.DefaultVPPParams.VPPKNNParam
								})),
								g.SliderInt("radius", &app.VPPParams.VPPKNNParam.Radius, 0, 5).Format("%.0f"),
								g.SliderFloat("strength", &app.VPPParams.VPPKNNParam.Strength, 0, 1).Format("%.2f"),
								g.SliderFloat("lerp", &app.VPPParams.VPPKNNParam.Lerp, 0, 1).Format("%.2f"),
								g.SliderFloat("th_lerp", &app.VPPParams.VPPKNNParam.ThLerp, 0, 1).Format("%.2f"),
							),
							g.Child("PMD").Border(true).Size((contentWidth-8)*0.5, 0).Flags(inputDisableFlag).Layout(
								g.Line(g.Checkbox("PMD", &app.VPPSwitches.VPPPMD), g.Dummy(-52, 0), g.Button("Reset##ResetPMD").OnClick(func() {
									app.VPPParams.VPPPMDParam = preset.DefaultVPPParams.VPPPMDParam
								})),
								g.SliderInt("applyCount", &app.VPPParams.VPPPMDParam.ApplyCount, 1, 100).Format("%.0f"),
								g.SliderInt("strength", &app.VPPParams.VPPPMDParam.Strength, 0, 100).Format("%.0f"),
								g.SliderInt("threshold", &app.VPPParams.VPPPMDParam.Threshold, 0, 255).Format("%.0f"),
							),
						),
						g.Custom(func() {
							imgui.PopStyleColorV(1)
						}),
					),

					g.Dummy(contentWidth, 5),
					g.Label("Sharpen"),
					g.Child("Sharpen").Border(false).Size(contentWidth, 150).Layout(
						g.Custom(func() {
							imgui.PushStyleColor(imgui.StyleColorChildBg, imgui.Vec4{X: 0.12, Y: 0.12, Z: 0.12, W: 0.99})
						}),
						g.Line(
							g.Child("UnSharp").Border(true).Size((contentWidth-8)*0.5, 0).Flags(inputDisableFlag).Layout(
								g.Line(g.Checkbox("UnSharp", &app.VPPSwitches.VPPUnSharp), g.Dummy(-52, 0), g.Button("Reset##ResetUnSharp").OnClick(func() {
									app.VPPParams.VPPUnSharpParam = preset.DefaultVPPParams.VPPUnSharpParam
								})),
								g.SliderInt("radius", &app.VPPParams.VPPUnSharpParam.Radius, 1, 9).Format("%.0f"),
								g.SliderFloat("weight", &app.VPPParams.VPPUnSharpParam.Weight, 0, 10).Format("%.2f"),
								g.SliderFloat("threshold", &app.VPPParams.VPPUnSharpParam.Threshold, 0, 255).Format("%.0f"),
							),
							g.Child("EdgeLevel").Border(true).Size((contentWidth-8)*0.5, 0).Flags(inputDisableFlag).Layout(
								g.Line(g.Checkbox("EdgeLevel", &app.VPPSwitches.VPPEdgeLevel), g.Dummy(-52, 0), g.Button("Reset##ResetEdgeLevel").OnClick(func() {
									app.VPPParams.VPPEdgeLevelParam = preset.DefaultVPPParams.VPPEdgeLevelParam
								})),
								g.SliderFloat("strength", &app.VPPParams.VPPEdgeLevelParam.Strength, -31, 31).Format("%.2f"),
								g.SliderFloat("threshold", &app.VPPParams.VPPEdgeLevelParam.Threshold, 0, 255).Format("%.2f"),
								g.SliderFloat("black", &app.VPPParams.VPPEdgeLevelParam.Black, 0, 31).Format("%.2f"),
								g.SliderFloat("white", &app.VPPParams.VPPEdgeLevelParam.White, 0, 31).Format("%.2f"),
							),
						),
						g.Custom(func() {
							imgui.PopStyleColorV(1)
						}),
					),
				),
			),

			g.TabItem("MediaInfo").Layout(
				g.Spacing(),
				g.InputTextMultiline("##mediaInfoLog", &app.State.MediaInfoLog).Size(contentWidth, 362.5).Flags(g.InputTextFlagsReadOnly),
			),

			g.TabItem("About").Layout(
				g.Spacing(),
				g.InputTextMultiline("##aboutText", &aboutText).Size(contentWidth, 362.5).Flags(g.InputTextFlagsReadOnly),
			),
		),
	)
}
