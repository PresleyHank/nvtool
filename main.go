package main

import (
	"fmt"
	"image/color"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"syscall"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	ffmpeg "github.com/Nicify/nvtool/ffmpeg"
	gpu "github.com/Nicify/nvtool/gpu"
	mediainfo "github.com/Nicify/nvtool/mediainfo"
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

var (
	font             = imgui.Font(0)
	selectedTebIndex int
	inputPath        string
	outputPath       string
	fullDuration     uint
	isEncoding       bool
	progress         float32
	gpuInfo          string
	ffmpegLog        string
	mediaInfoLog     string = "Drag and drop media files here"
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
		ffmpeg.RunEncode(inputPath, outputPath, strings.Split(command, " "), &progress, &ffmpegLog, g.Update)
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
	isEncoding = false
	proc := ffmpeg.GetProcess()
	if proc == nil {
		return
	}
	stdin, err := proc.StdinPipe()
	if err == nil && stdin != nil {
		stdin.Write([]byte("q\n"))
		return
	}
	if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "process", "where", "name='ffmpeg.exe'", "delete")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd.Run()
		return
	}
	proc.Process.Kill()
}

func shouldDisableInput(b bool) (flag g.WindowFlags) {
	if b {
		return g.WindowFlagsNoInputs
	}
	return
}

func loop() {
	imgui.PushStyleVarFloat(imgui.StyleVarWindowBorderSize, 0)
	imgui.PushStyleVarFloat(imgui.StyleVarFrameBorderSize, 0)
	imgui.PushStyleVarFloat(imgui.StyleVarChildBorderSize, 0)
	imgui.PushStyleVarFloat(imgui.StyleVarFrameRounding, 3)
	g.Context.IO().SetFontGlobalScale(1)
	g.PushColorWindowBg(color.RGBA{50, 50, 50, 250})
	g.PushColorFrameBg(color.RGBA{10, 10, 10, 240})
	g.PushColorButton(color.RGBA{100, 100, 100, 255})
	g.PushColorButtonHovered(color.RGBA{120, 120, 120, 240})
	g.PushColorButtonActive(color.RGBA{80, 80, 80, 245})
	g.SingleWindow("NVENC Video Toolbox",
		g.Layout{
			g.TabBar("maintab", g.Layout{
				g.TabItem("Encode", g.Layout{
					g.Child("control", false, 734, 92, shouldDisableInput(isEncoding), g.Layout{
						g.Spacing(),
						g.Line(
							g.InputTextV("##video", -68/imgui.DPIScale, &inputPath, 0, nil, nil),
							g.ButtonV("video", 60, 22, onInputClick),
						),
						g.Spacing(),
						g.Line(
							g.InputTextV("##output", -68/imgui.DPIScale, &outputPath, 0, nil, nil),
							g.ButtonV("output", 60, 22, onOutputClick),
						),
						g.Spacing(),
						g.Line(
							g.Label("preset"),
							g.Combo("##preset", ffmpeg.PresetOptions[defaultPreset.preset], ffmpeg.PresetOptions, &defaultPreset.preset, 80, 0, nil),

							g.Label("rc"),
							g.Combo("##rc", ffmpeg.RCOptions[defaultPreset.rc], ffmpeg.RCOptions, &defaultPreset.rc, 80, 0, nil),

							// g.Label("cq"),
							// g.InputIntV("##cq", 25, &defaultPreset.cq, 0, nil),

							g.Label("qmin"),
							g.InputIntV("##qmin", 35, &defaultPreset.qmin, 0, nil),

							g.Label("qmax"),
							g.InputIntV("##qmax", 35, &defaultPreset.qmax, 0, nil),

							g.Label("aq"),
							g.Combo("##aq", ffmpeg.AQOptions[defaultPreset.aq], ffmpeg.AQOptions, &defaultPreset.aq, 95, 0, nil),

							// g.Label("aq-strength"),
							g.InputIntV("##aqstrength", 35, &defaultPreset.aqStrength, 0, func() {
								defaultPreset.aqStrength = limitValue(defaultPreset.aqStrength, 0, 15)
							}),

							g.Label("bitrate"),
							g.InputIntV("k##bitrate", 70, &defaultPreset.bitrate, 0, nil),

							// g.Label("Maxrate"),
							// g.InputIntV("k##maxrate", 65, &defaultPreset.maxrate, 0, nil),
						),
					}),
					g.Spacing(),
					g.InputTextMultiline("", &ffmpegLog, 734, 200, 0, nil, func() {
						imgui.SetScrollHereY(1.0)
					}),
					g.Spacing(),
					g.ProgressBar(progress, 734, 20, ""),
					g.Line(
						g.Dummy(0, 5),
					),
					g.Line(
						g.Label(gpuInfo),
						g.Dummy(-68, 24),
						g.Condition(isEncoding,
							g.Layout{g.ButtonV("Cancel", 60, 24, dispose)},
							g.Layout{g.ButtonV("Run", 60, 24, onRunClick)},
						),
					),
				},
				),
				g.TabItem("MediaInfo", g.Layout{
					g.Spacing(),
					g.InputTextMultiline("mediainfo", &mediaInfoLog, 734, 374, g.InputTextFlagsReadOnly, nil, nil),
				}),
			}),
		})
	g.PopStyleColorV(5)
	imgui.PopStyleVarV(4)
}

func main() {
	gpuList, err := gpu.GetGPUInfo()
	if err != nil {
		fmt.Printf("Error getting GPU info: %v", err)
	}
	gpuInfo = strings.Join(gpuList, " ")
	mw := g.NewMasterWindow("NVENC Video Toolbox 1.2", 750, 420, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsTransparent, loadFont)
	mw.SetBgColor(color.RGBA{0, 0, 0, 0})
	mw.SetDropCallback(onDrop)
	mw.Main(loop)
	dispose()
}
