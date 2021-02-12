package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"
	"unsafe"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	c "github.com/Nicify/customwidget"
	"github.com/Nicify/nvtool/hooks"
	mediainfo "github.com/Nicify/nvtool/mediainfo"
	nvenc "github.com/Nicify/nvtool/nvenc"
	win "github.com/Nicify/nvtool/win"
	theme "github.com/Nicify/theme"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type encodingPresets struct {
	hevc       bool
	preset     int32
	quality    int32
	bitrate    int32
	maxrate    int32
	aq         int32
	aqStrength int32
	resize     bool
	outputRes  string
	vppSwitches
	vppParams
}

type vppSwitches struct {
	vppKNN        bool
	vppPMD        bool
	vppUnSharp    bool
	vppEdgeLevel  bool
	vppSmooth     bool
	vppColorSpace bool
}

type vppParams struct {
	nvenc.VPPKNNParam
	nvenc.VPPPMDParam
	nvenc.VPPUnSharpParam
	nvenc.VPPEdgeLevelParam
	nvenc.VPPColorSpaceParam
}

type usage struct {
	GPU int32
	VE  int32
	VD  int32
}

const (
	windowPadding = 8
	contentWidth  = 734
	buttonWidth   = 68
	buttonHeight  = 24
)

var (
	lockFile = path.Join(os.TempDir(), "nvtool.lock")
	mounted  bool

	fontTamzenr imgui.Font
	fontTamzenb imgui.Font
	fontIosevka imgui.Font

	texLogo         *g.Texture
	texButtonClose  *g.Texture
	texDropDown     *g.Texture
	texGraphicsCard *g.Texture

	mw          *g.MasterWindow
	glfwWindow  *glfw.Window
	mwMoveState = hooks.MWMoveState{}
	mwDragArea  = image.Rectangle{image.ZP, image.Point{750, 30}}

	nvencCmd   *exec.Cmd
	inputPath  string
	outputPath string
	percent    float32
	gpuName    string

	nvencLog     string
	mediaInfoLog string = "Drag and drop media files here"
)

var defaultVppSwitches = vppSwitches{}

var defaultVppParams = vppParams{
	nvenc.DefaultVPPKNNParam,
	nvenc.DefaultVPPPMDParam,
	nvenc.DefaultVPPUnSharpParam,
	nvenc.DefaultVPPEdgeLevelParam,
	nvenc.DefaultVPPColorSpaceParam,
}

var defaultPreset = encodingPresets{
	hevc:        false,
	preset:      6,
	quality:     12,
	bitrate:     19000,
	maxrate:     59850,
	outputRes:   "1920x1080",
	vppSwitches: defaultVppSwitches,
	vppParams:   defaultVppParams,
}

var utilization = usage{
	GPU: 0,
	VD:  0,
	VE:  0,
}

func isEncoding() bool {
	if nvencCmd == nil || (nvencCmd.ProcessState != nil && nvencCmd.ProcessState.Exited()) {
		return false
	}
	return true
}

func resetState() {
	percent = 0
	nvencLog = ""
	g.Update()
}

func onInputClick() {
	filePath := selectInputPath()
	if len(filePath) > 1 {
		// percent = 0
		// nvencLog = ""
		inputPath = filePath
		fileExt := path.Ext(inputPath)
		outputPath = strings.Replace(inputPath, fileExt, "_nvenc.mp4", 1)
		go setMediaInfo(filePath)
	}
}

func onOutputClick() {
	filePath := selectOutputPath()
	if len(filePath) > 1 {
		outputPath = filePath
	}
}

func onRunClick() {
	if isEncoding() ||
		invalidPath(inputPath, outputPath) ||
		strings.HasSuffix(nvencLog, "Get input file information...") {
		return
	}
	go func() {
		defer g.Update()
		resetState()
		codec := "h264"
		if defaultPreset.hevc {
			codec = "hevc"
		}
		command := fmt.Sprintf("--codec %s --profile high --audio-copy --preset %s --vbr %v --vbr-quality %v --max-bitrate 60000 --bframes 4 --ref 16 --lookahead 32 --gop-len 250 --%s --aq-strength %v --vpp-resize lanczos2 --vpp-perf-monitor --ssim",
			codec,
			nvenc.PresetOptions[defaultPreset.preset],
			defaultPreset.bitrate,
			defaultPreset.quality,
			nvenc.AQOptions[defaultPreset.aq],
			defaultPreset.aqStrength,
		)
		args := strings.Split(command, " ")

		// if defaultPreset.vppPresets.vppColorSpace {
		// 	args = append(args, "--vpp-colorspace", defaultPreset.vppPresets.vppColorSpaceParam)
		// }

		if defaultPreset.vppSwitches.vppKNN {
			param := defaultPreset.VPPKNNParam
			args = append(args, "--vpp-knn", fmt.Sprintf("radius=%v,strength=%.2f,lerp=%.1f,th_lerp=%.1f", param.Radius, param.Strength, param.Lerp, param.ThLerp))
		}

		if defaultPreset.vppSwitches.vppPMD {
			param := defaultPreset.VPPPMDParam
			args = append(args, "--vpp-pmd", fmt.Sprintf("apply_count=%v,strength=%v,threshold=%v", param.ApplyCount, param.Strength, param.Threshold))
		}

		if defaultPreset.vppSwitches.vppUnSharp {
			param := defaultPreset.VPPUnSharpParam
			args = append(args, "--vpp-unsharp", fmt.Sprintf("radius=%v,weight=%.1f,threshold=%.1f", param.Radius, param.Weight, param.Threshold))
		}

		if defaultPreset.vppSwitches.vppEdgeLevel {
			param := defaultPreset.VPPEdgeLevelParam
			args = append(args, "--vpp-edgelevel", fmt.Sprintf("strength=%v,threshold=%.1f,black=%v,white=%v", param.Strength, param.Threshold, param.Black, param.White))
		}

		if defaultPreset.resize {
			args = append(args, "--output-res", defaultPreset.outputRes)
		}

		cmd, progress, _ := nvenc.RunEncode(inputPath, outputPath, args)
		nvencCmd = cmd
		for msg := range progress {
			percent = float32(msg.Percent) / 100
			utilization.GPU = int32(msg.GPU)
			utilization.VE = int32(msg.VE)
			utilization.VD = int32(msg.VD)
			nvencLog += fmt.Sprintf("\n%v frames: %.0f fps, %v kb/s, remain %s, est out size %s", msg.FramesProcessed, msg.FPS, msg.Bitrate, msg.Remain, msg.EstOutSize)
			nvencLog = strings.Trim(nvencLog, "\n")
			g.Update()
		}

		if nvencCmd.ProcessState != nil && nvencCmd.ProcessState.Success() {
			percent = 1.0
		}
	}()
}

func setMediaInfo(inputPath string) {
	info, err := mediainfo.GetMediaInfo(inputPath)
	if err != nil {
		mediaInfoLog = fmt.Sprintf("Error: %s", err)
		return
	}
	mediaInfoLog = info
	g.Update()
}

func onDrop(dropItem []string) {
	if isEncoding() {
		return
	}
	inputPath = dropItem[0]
	fileExt := path.Ext(inputPath)
	outputPath = strings.Replace(inputPath, fileExt, "_nvenc.mp4", 1)
	go setMediaInfo(inputPath)
}

func dispose() {
	if nvencCmd == nil {
		return
	}
	nvencCmd.Process.Kill()
	go nvencCmd.Wait()
}

func shouldDisableInput(b bool) (flag g.WindowFlags) {
	if b {
		return g.WindowFlagsNoInputs
	}
	return g.WindowFlagsNone
}

func loop() {
	hooks.UseMounted(&mounted, loadTexture)
	hooks.UseWindowMove(glfwWindow, mwDragArea, &mwMoveState)
	isEncoding := isEncoding()
	inputDisableFlag := shouldDisableInput(isEncoding)
	useLayoutFlat := theme.UseLayoutFlat()
	useStyleButtonDark := theme.UseStyleButtonDark()
	defer useLayoutFlat.Pop()
	useLayoutFlat.Push()
	g.SingleWindow("NVTool").Layout(
		g.Group().Layout(
			g.Line(
				g.Image(texLogo).Size(18, 18),
				g.Label("NVENC Video Toolbox 2.2"),
				g.Dummy(-83, 0),
				g.Custom(useStyleButtonDark.Push),
				g.Button(".").Size(20, 20),
				g.Button("_").Size(20, 20).OnClick(func() {
					win.ShowWindow(win.HWND(unsafe.Pointer(glfwWindow.GetWin32Window())), win.SW_FORCEMINIMIZE)
				}),
				g.ImageButton(texButtonClose).Size(20, 20).OnClick(func() {
					glfwWindow.SetShouldClose(true)
				}),
				g.Custom(useStyleButtonDark.Pop),
			),
		),
		g.TabBar("maintab").Layout(
			g.TabItem("Encode").Layout(
				g.Child("control").Border(false).Flags(inputDisableFlag).Size(contentWidth, 92).Layout(
					g.Spacing(),
					g.Line(
						g.InputText("##video", &inputPath).Size(-((windowPadding+buttonWidth)/imgui.DPIScale)),
						c.WithHiDPIFont(fontIosevka, fontTamzenb, g.Layout{g.Button("Video").Size(buttonWidth, buttonHeight).OnClick(onInputClick)}),
					),

					g.Spacing(),
					g.Line(
						g.InputText("##output", &outputPath).Size(-((windowPadding+buttonWidth)/imgui.DPIScale)),
						c.WithHiDPIFont(fontIosevka, fontTamzenb, g.Layout{g.Button("Output").Size(buttonWidth, buttonHeight).OnClick(onOutputClick)}),
					),

					g.Spacing(),
					g.Line(
						g.Label("Preset"),
						g.Combo("##preset", nvenc.PresetOptions[defaultPreset.preset], nvenc.PresetOptions, &defaultPreset.preset).Size(50),

						g.Label("Quality"),
						g.InputInt("##quality", &defaultPreset.quality).Size(24),

						g.Label("Bitrate"),
						g.InputInt("##bitrate", &defaultPreset.bitrate).Size(60),

						g.Label("AQ"),
						g.Combo("##aq", nvenc.AQOptionsForPreview[defaultPreset.aq], nvenc.AQOptionsForPreview, &defaultPreset.aq).Size(92),
						g.Label("-"),
						g.InputInt("##strength", &defaultPreset.aqStrength).Size(24).OnChange(func() {
							defaultPreset.aqStrength = limitValue(defaultPreset.aqStrength, 0, 15)
						}),

						g.Checkbox("HEVC", &defaultPreset.hevc),

						g.Checkbox("Resize", &defaultPreset.resize),
						g.InputText("##outputRes", &defaultPreset.outputRes).Size(80).Flags(g.InputTextFlagsCallbackAlways).OnChange(func() {
							defaultPreset.outputRes = limitResValue(defaultPreset.outputRes)
						}),
					),
				),

				g.Spacing(),
				g.InputTextMultiline("##nvencLog", &nvencLog).Size(contentWidth, 200).Flags(g.InputTextFlagsReadOnly),
				g.Custom(func() {
					if isEncoding && time.Now().Second()%2 == 0 {
						imgui.BeginChild("##nvencLog")
						imgui.SetScrollHereY(1)
						imgui.EndChild()
					}
				}),

				g.Spacing(),
				g.ProgressBar(percent).Size(contentWidth, 20).Overlay(""),

				g.Line(
					g.Dummy(0, 5),
				),
				g.Line(
					g.Condition(gpuName != "",
						g.Layout{
							g.Line(
								g.Image(texGraphicsCard).Size(18, 18),
								g.Label(fmt.Sprintf("%s | GPU:%v%% VE:%v%% VD:%v%%", gpuName, utilization.GPU, utilization.VE, utilization.VD))),
						},
						g.Layout{
							g.Line(g.Dummy(18, 18)),
						},
					),

					g.Dummy(-(windowPadding+buttonWidth), 24),
					c.WithHiDPIFont(fontIosevka, fontTamzenb, g.Layout{g.Condition(isEncoding,
						g.Layout{g.Button("Cancel").Size(buttonWidth, buttonHeight).OnClick(dispose)},
						g.Layout{g.Button("Run").Size(buttonWidth, buttonHeight).OnClick(onRunClick)},
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
								g.Line(g.Checkbox("KNN", &defaultPreset.vppSwitches.vppKNN), g.Dummy(-52, 0), g.Button("Reset##ResetKNN").OnClick(func() {
									defaultPreset.VPPKNNParam = defaultVppParams.VPPKNNParam
								})),
								g.SliderInt("radius", &defaultPreset.VPPKNNParam.Radius, 0, 5).Format("%.0f"),
								g.SliderFloat("strength", &defaultPreset.VPPKNNParam.Strength, 0, 1).Format("%.2f"),
								g.SliderFloat("lerp", &defaultPreset.VPPKNNParam.Lerp, 0, 1).Format("%.2f"),
								g.SliderFloat("th_lerp", &defaultPreset.VPPKNNParam.ThLerp, 0, 1).Format("%.2f"),
							),
							g.Child("PMD").Border(true).Size((contentWidth-8)*0.5, 0).Flags(inputDisableFlag).Layout(
								g.Line(g.Checkbox("PMD", &defaultPreset.vppSwitches.vppPMD), g.Dummy(-52, 0), g.Button("Reset##ResetPMD").OnClick(func() {
									defaultPreset.VPPPMDParam = defaultVppParams.VPPPMDParam
								})),
								g.SliderInt("applyCount", &defaultPreset.VPPPMDParam.ApplyCount, 1, 100).Format("%.0f"),
								g.SliderInt("strength", &defaultPreset.VPPPMDParam.Strength, 0, 100).Format("%.0f"),
								g.SliderInt("threshold", &defaultPreset.VPPPMDParam.Threshold, 0, 255).Format("%.0f"),
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
								g.Line(g.Checkbox("UnSharp", &defaultPreset.vppSwitches.vppUnSharp), g.Dummy(-52, 0), g.Button("Reset##ResetUnSharp").OnClick(func() {
									defaultPreset.VPPUnSharpParam = defaultVppParams.VPPUnSharpParam
								})),
								g.SliderInt("radius", &defaultPreset.VPPUnSharpParam.Radius, 1, 9).Format("%.0f"),
								g.SliderFloat("weight", &defaultPreset.VPPUnSharpParam.Weight, 0, 10).Format("%.2f"),
								g.SliderFloat("threshold", &defaultPreset.VPPUnSharpParam.Threshold, 0, 255).Format("%.0f"),
							),
							g.Child("EdgeLevel").Border(true).Size((contentWidth-8)*0.5, 0).Flags(inputDisableFlag).Layout(
								g.Line(g.Checkbox("EdgeLevel", &defaultPreset.vppSwitches.vppEdgeLevel), g.Dummy(-52, 0), g.Button("Reset##ResetEdgeLevel").OnClick(func() {
									defaultPreset.VPPEdgeLevelParam = defaultVppParams.VPPEdgeLevelParam
								})),
								g.SliderFloat("strength", &defaultPreset.VPPEdgeLevelParam.Strength, -31, 31).Format("%.2f"),
								g.SliderFloat("threshold", &defaultPreset.VPPEdgeLevelParam.Threshold, 0, 255).Format("%.2f"),
								g.SliderFloat("black", &defaultPreset.VPPEdgeLevelParam.Black, 0, 31).Format("%.2f"),
								g.SliderFloat("white", &defaultPreset.VPPEdgeLevelParam.White, 0, 31).Format("%.2f"),
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
				g.InputTextMultiline("##mediaInfoLog", &mediaInfoLog).Size(contentWidth, 362.5).Flags(g.InputTextFlagsReadOnly),
			),

			g.TabItem("About").Layout(
				g.Spacing(),
				g.InputTextMultiline("##aboutText", &aboutText).Size(contentWidth, 362.5).Flags(g.InputTextFlagsReadOnly),
			),
		),
	)
}

func applyWindowProperties(window *glfw.Window) {
	data, _ := box.Find("icon_48px.png")
	icon48px, _ := loadImageFromMemory(data)
	glfwWindow.SetIcon([]image.Image{icon48px})
	hwnd := win.HWND(unsafe.Pointer(glfwWindow.GetWin32Window()))
	win.SetWindowCompositionAttribute(hwnd, 3, 0, 0, 0)
	g.Context.GetPlatform().SetTPS(240)
	glfwWindow.SetFocusCallback(func(w *glfw.Window, focused bool) {
		if focused {
			glfwWindow.SetOpacity(0.98)
			return
		}
		glfwWindow.SetOpacity(1)
	})
}

func loadFont() {
	fonts := g.Context.IO().Fonts()
	fontIosevkaTTF, _ := box.Find("iosevka.ttf")
	fontIosevka = fonts.AddFontFromMemoryTTFV(fontIosevkaTTF, 18, imgui.DefaultFontConfig, fonts.GlyphRangesChineseFull())
	fontTamzenbTTF, _ := box.Find("tamzen8x16b.ttf")
	fontTamzenb = fonts.AddFontFromMemoryTTFV(fontTamzenbTTF, 16, imgui.DefaultFontConfig, fonts.GlyphRangesChineseFull())
	fontTamzenrTTF, _ := box.Find("tamzen8x16r.ttf")
	fontTamzenr = fonts.AddFontFromMemoryTTFV(fontTamzenrTTF, 16, imgui.DefaultFontConfig, fonts.GlyphRangesChineseFull())
}

func loadTexture() {
	go func() {
		texLogo, _ = imageToTexture("icon.png")
		texButtonClose, _ = imageToTexture("close_white.png")
		texDropDown, _ = imageToTexture("dropdown.png")
		texGraphicsCard, _ = imageToTexture("graphics_card.png")
	}()
}

func checkCore() {
	if _, err := os.Stat(nvenc.Binary); os.IsNotExist(err) {
		go func() {
			defer g.Update()
			nvencLog = "Downloading core.zip...\n"
			files, err := download("https://hub.fastgit.org/rigaya/NVEnc/releases/download/5.27/NVEncC_5.27_x64.7z", "./core", func(progress float32) {
				percent = progress
				g.Update()
			})
			if err != nil {
				fmt.Printf("Download failed %s", err)
			}
			percent = 1
			nvencLog += strings.Join(files, "\n") + "\nDownload completed."
		}()
	}
}

func onSecondInstance() {
	ioutil.WriteFile(lockFile, []byte("focus"), 0644)
	os.Exit(0)
}

func onCommand(command string) {
	if command == "focus" {
		glfwWindow.Restore()
	}
}

func init() {
	runtime.LockOSThread()

}

func main() {
	defer dispose()
	unlock := initSingleInstanceLock(lockFile, onSecondInstance, onCommand)
	defer unlock()

	gpuName, _ = nvenc.CheckDevice()
	mw = g.NewMasterWindow("NVTool", 750, 435, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFrameless|g.MasterWindowFlagsTransparent, loadFont)
	mw.SetBgColor(color.RGBA{0, 0, 0, 0})
	mw.SetDropCallback(onDrop)

	currentStyle := imgui.CurrentStyle()
	theme.SetThemeDark(&currentStyle)

	platform := g.Context.GetPlatform().(*imgui.GLFW)
	glfwWindow = platform.GetWindow()
	applyWindowProperties(glfwWindow)
	checkCore()
	mw.Run(loop)
}
