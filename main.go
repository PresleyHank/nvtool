package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	ffmpeg "github.com/Nicify/nvtool/ffmpeg"
	"github.com/Nicify/nvtool/gpu"
	mediainfo "github.com/Nicify/nvtool/mediainfo"
	theme "github.com/Nicify/nvtool/theme"
	win "github.com/Nicify/nvtool/win"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type encodingPresets struct {
	preset     int32
	rc         int32
	aq         int32
	cq         int32
	qmin       int32
	qmax       int32
	bitrate    int32
	maxrate    int32
	aqStrength int32
}

const (
	windowPadding = 8
	contentWidth  = 734
	buttonWidth   = 68
	buttonHeight  = 24
)

var (
	lockFile           = path.Join(os.TempDir(), "nvtool.lock")
	ffmpegCmd          *exec.Cmd
	texLogo            *g.Texture
	texButtonClose     *g.Texture
	texGraphicsCard    *g.Texture
	texLoadingfinished = make(chan bool)
	mw                 *g.MasterWindow
	glfwWindow         *glfw.Window
	mwMoveable         bool
	prevMouseX         int
	prevMouseY         int
	selectedTebIndex   int
	inputPath          string
	outputPath         string
	fullDuration       uint
	isEncoding         bool
	progress           float32
	gpuName            string
	ffmpegLog          string
	mediaInfoLog       string = "Drag and drop media files here"
)

var defaultPreset = encodingPresets{
	qmin:       16,
	qmax:       24,
	bitrate:    19850,
	maxrate:    59850,
	aqStrength: 15,
}

func cleanOutput() {
	progress = 0
	ffmpegLog = ""
}

func onInputClick() {
	filePath := selectInputPath()
	if len(filePath) > 1 {
		progress = 0
		ffmpegLog = ""
		inputPath = filePath
		fileExt := path.Ext(inputPath)
		outputPath = strings.Replace(inputPath, fileExt, "_x264.mp4", 1)
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
	if isEncoding || invalidPath(inputPath, outputPath) {
		return
	}
	cleanOutput()
	go func() {
		isEncoding = true
		command := fmt.Sprintf(
			"-c:a copy -c:v h264_nvenc -preset %s -profile:v high -rc:v %s -qmin %d -qmax %d -strict_gop 1 -%s-aq 1 -aq-strength:v %d -b:v %dk -maxrate:v %dk -map 0 -f mp4",
			ffmpeg.PresetOptions[defaultPreset.preset],
			ffmpeg.RCOptions[defaultPreset.rc],
			defaultPreset.qmin,
			defaultPreset.qmax,
			ffmpeg.AQOptions[defaultPreset.aq],
			defaultPreset.aqStrength,
			defaultPreset.bitrate,
			defaultPreset.maxrate,
		)
		cmd, progress, _ := ffmpeg.RunEncode(inputPath, outputPath, strings.Split(command, " "))
		ffmpegCmd = cmd
		for msg := range progress {
			log.Printf("%+v", msg)
		}
		isEncoding = false
	}()
}

func setMediaInfo(inputPath string) {
	info, err := mediainfo.GetMediaInfo(inputPath)
	if err != nil {
		mediaInfoLog = fmt.Sprintf("Error: %s", err)
		return
	}
	mediaInfoLog = strings.Join(info, "\n")
	g.Update()
}

func onDrop(dropItem []string) {
	if isEncoding {
		return
	}
	inputPath = dropItem[0]
	fileExt := path.Ext(inputPath)
	outputPath = strings.Replace(inputPath, fileExt, "_x264.mp4", 1)
	go setMediaInfo(inputPath)
}

func dispose() {
	if ffmpegCmd == nil {
		return
	}
	stdin, err := ffmpegCmd.StdinPipe()
	if err == nil {
		stdin.Write([]byte("q\n"))
	}
	// ffmpegCmd.Process.Signal(os.Interrupt)
	if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "process", "where", "name='ffmpeg.exe'", "delete")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd.Run()
		return
	}
	go ffmpegCmd.Wait()
	isEncoding = false
}

func shouldDisableInput(b bool) (flag g.WindowFlags) {
	if b {
		return g.WindowFlagsNoInputs
	}
	return
}

func shouldWindowMove() {
	mousePos := g.GetMousePos()
	prevPosX, prevPosY := glfwWindow.GetPos()
	if g.IsMouseClicked(0) {
		mwMoveable = float32(mousePos.Y) < 50*imgui.DPIScale
		prevMouseX = mousePos.X
		prevMouseY = mousePos.Y
	}
	if mwMoveable && g.IsMouseDown(0) {
		offsetX := mousePos.X - prevMouseX
		offsetY := mousePos.Y - prevMouseY
		glfwWindow.SetPos(prevPosX+int(offsetX), prevPosY+int(offsetY))
	}
}

func loop() {
	shouldWindowMove()
	useLayoutFlat := theme.UseLayoutFlat()
	useStyleButtonDark := theme.UseStyleButtonDark()
	defer useLayoutFlat.Pop()
	useLayoutFlat.Push()
	g.SingleWindow("NVTool",
		g.Layout{
			g.Line(
				g.Image(texLogo, 18, 18),
				g.Label("NVTool 1.5"),
				g.Dummy(-83, 0),
				g.Custom(useStyleButtonDark.Push),
				g.ButtonV(".", 20, 20, func() {}),
				g.ButtonV("_", 20, 20, func() {
					win.ShowWindow(win.HWND(unsafe.Pointer(glfwWindow.GetWin32Window())), win.SW_FORCEMINIMIZE)
				}),
				g.ImageButton(texButtonClose, 20, 20, func() {
					glfwWindow.SetShouldClose(true)
				}),
				g.Custom(useStyleButtonDark.Pop),
			),
			g.TabBar("maintab", g.Layout{
				g.TabItem("Encode", g.Layout{
					g.Child("control", false, contentWidth, 92, shouldDisableInput(isEncoding), g.Layout{
						g.Spacing(),
						g.Line(
							g.InputTextV("##video", -((windowPadding+buttonWidth)/imgui.DPIScale), &inputPath, 0, nil, nil),
							g.ButtonV("Video", buttonWidth, buttonHeight, onInputClick),
						),

						g.Spacing(),
						g.Line(
							g.InputTextV("##output", -((windowPadding+buttonWidth)/imgui.DPIScale), &outputPath, g.InputTextFlagsReadOnly, nil, nil),
							g.ButtonV("Output", buttonWidth, buttonHeight, onOutputClick),
						),

						g.Spacing(),
						g.Line(
							g.Label("Preset"),
							g.Combo("##preset", ffmpeg.PresetOptions[defaultPreset.preset], ffmpeg.PresetOptions, &defaultPreset.preset, 80, 0, nil),

							g.Label("RC"),
							g.Combo("##rc", ffmpeg.RCOptions[defaultPreset.rc], ffmpeg.RCOptions, &defaultPreset.rc, 80, 0, nil),

							g.Label("QMin"),
							g.InputIntV("##qmin", 40, &defaultPreset.qmin, 0, nil),

							g.Label("QMax"),
							g.InputIntV("##qmax", 40, &defaultPreset.qmax, 0, nil),

							g.Label("AQ"),
							g.Combo("##aq", ffmpeg.AQOptions[defaultPreset.aq], ffmpeg.AQOptions, &defaultPreset.aq, 95, 0, nil),

							g.InputIntV("##aqstrength", 40, &defaultPreset.aqStrength, 0, func() {
								defaultPreset.aqStrength = limitValue(defaultPreset.aqStrength, 0, 15)
							}),

							g.Label("Bitrate"),
							g.InputIntV("##bitrate", 70, &defaultPreset.bitrate, 0, nil),

							// g.Label("Maxrate"),
							// g.InputIntV("k##maxrate", 65, &defaultPreset.maxrate, 0, nil),
						),
					}),

					g.Spacing(),
					g.InputTextMultiline("##ffmpegLog", &ffmpegLog, contentWidth, 200, g.InputTextFlagsReadOnly, nil, func() {
						imgui.SetScrollHereY(1.0)
					}),

					g.Spacing(),
					g.ProgressBar(progress, contentWidth, 20, ""),

					g.Line(
						g.Dummy(0, 5),
					),
					g.Line(
						g.Condition(gpuName != "", g.Layout{
							g.Line(
								g.Image(texGraphicsCard, 18, 18),
								g.Label(gpuName),
							),
						}, nil),
						g.Dummy(-(windowPadding+buttonWidth), 24),
						g.Condition(isEncoding,
							g.Layout{g.ButtonV("Cancel", buttonWidth, buttonHeight, dispose)},
							g.Layout{g.ButtonV("Run", buttonWidth, buttonHeight, onRunClick)},
						),
					),
				},
				),

				g.TabItem("MediaInfo", g.Layout{
					g.Spacing(),
					g.InputTextMultiline("##mediaInfoLog", &mediaInfoLog, contentWidth, 360, g.InputTextFlagsReadOnly, nil, nil),
				}),

				// g.TabItem("Settings", g.Layout{
				// 	g.Custom(func() {
				// 		imgui.PushStyleColor(imgui.StyleColorChildBg, imgui.Vec4{X: 0.12, Y: 0.12, Z: 0.12, W: 0.99})
				// 	}),

				// 	g.Spacing(),
				// 	g.Label("Interface"),
				// 	g.Child("Interface", true, contentWidth, 95, g.WindowFlagsAlwaysUseWindowPadding, g.Layout{}),

				// 	g.Spacing(),
				// 	g.Label("Encoding"),
				// 	g.Child("Encoding", true, contentWidth, 95, g.WindowFlagsAlwaysUseWindowPadding, g.Layout{}),

				// 	g.Spacing(),
				// 	g.Label("Binary"),
				// 	g.Child("Binary", true, contentWidth, 95, g.WindowFlagsAlwaysUseWindowPadding, g.Layout{}),

				// 	g.Custom(func() {
				// 		imgui.PopStyleColorV(1)
				// 	}),
				// }),
			}),
		})
}

func applyWindowProperties(window *glfw.Window) {
	data, err := box.Find("icon_48px.png")
	if err != nil {
		log.Fatal("icon_48px.png read failed.")
	}
	icon48px, _ := loadImageFromMemory(data)
	glfwWindow.SetIcon([]image.Image{icon48px})
	hwnd := win.HWND(unsafe.Pointer(glfwWindow.GetWin32Window()))
	win.SetWindowCompositionAttribute(hwnd, 3, 0, 0, 0)
	glfwWindow.SetFocusCallback(func(w *glfw.Window, focused bool) {
		if focused {
			glfwWindow.SetOpacity(0.98)
			return
		}
		glfwWindow.SetOpacity(1)
	})
}

func loadTexture() {
	texButtonClose, _ = imageToTexture("close_white.png")
	texLogo, _ = imageToTexture("icon.png")
	texGraphicsCard, _ = imageToTexture("graphics_card.png")
}

func init() {
	runtime.LockOSThread()

	if err := os.Remove(lockFile); err != nil && !os.IsNotExist(err) {
		ioutil.WriteFile(lockFile, []byte("focus"), 0644)
		os.Exit(0)
	}
}

func main() {
	defer dispose()
	unlock := initSingleInstanceLock()
	defer unlock()
	go loadTexture()
	gpuList, _ := gpu.GetGPUList()
	gpuName = gpuList[0]
	mw = g.NewMasterWindow("NVTool", 750, 435, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFrameless|g.MasterWindowFlagsTransparent, loadFont)
	currentStyle := imgui.CurrentStyle()
	theme.SetThemeDark(&currentStyle)
	platform := g.Context.GetPlatform().(*imgui.GLFW)
	glfwWindow = platform.GetWindow()
	applyWindowProperties(glfwWindow)
	mw.SetBgColor(color.RGBA{0, 0, 0, 0})
	mw.SetDropCallback(onDrop)
	mw.Main(loop)
}
