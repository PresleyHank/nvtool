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
	mediainfo "github.com/Nicify/nvtool/mediainfo"
)

var (
	font             = imgui.Font(0)
	selectedTebIndex int
	inputPath        string
	outputPath       string
	fullDuration     uint
	isEncoding       bool
	progress         float32
	ffmpegLog        string
	mediaInfoLog     string
)

var (
	presetItems = []string{"slow", "medium", "fast", "bd"}
	rcItems     = []string{"vbr_hq", "cbr_hq"}
	aqItems     = []string{"temporal", "spatial"}
)

var (
	preset     int32
	rc         int32
	aq         int32
	cq         int32
	qmin       int32 = 16
	qmax       int32 = 24
	bitrate    int32 = 19850
	maxrate    int32 = 60000
	aqStrength int32 = 15
)

func cleanOutput() {
	progress = 0
	ffmpegLog = ""
}

func handleInputClick() {
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

func handleOutputClick() {
	filePath := selectOutputPath()
	if len(filePath) > 1 {
		outputPath = filePath
	}
}

func handleRunClick() {
	if isEncoding || invalidPath(inputPath, outputPath) {
		return
	}
	cleanOutput()
	go func() {
		isEncoding = true
		command := fmt.Sprintf(
			"-c:a copy -c:v h264_nvenc -preset %s -profile:v high -level 5.1 -rc:v %s -qmin %d -qmax %d -strict_gop 1 -%s-aq 1 -aq-strength:v %d -b:v %dk -maxrate:v %dk -map 0 -f mp4",
			presetItems[preset],
			rcItems[rc],
			qmin,
			qmax,
			aqItems[aq],
			aqStrength,
			bitrate,
			maxrate,
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
}

func handleDrop(dropItem []string) {
	if isEncoding {
		return
	}
	inputPath = dropItem[0]
	fileExt := path.Ext(inputPath)
	outputPath = strings.Replace(inputPath, fileExt, "_x264.mp4", 1)
	go setMediaInfo(inputPath)
}

func handleCancelClick() {
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
					g.Child("control", false, 724, 90, shouldDisableInput(isEncoding), g.Layout{
						g.Spacing(),
						g.Line(
							g.InputTextV("##video", 655, &inputPath, 0, nil, nil),
							g.ButtonV("video", 60, 22, handleInputClick),
						),
						g.Spacing(),
						g.Line(
							g.InputTextV("##output", 655, &outputPath, 0, nil, nil),
							g.ButtonV("output", 60, 22, handleOutputClick),
						),
						g.Spacing(),
						g.Line(
							g.Label("preset"),
							g.Combo("##preset", presetItems[preset], presetItems, &preset, 85, 0, nil),

							g.Label("rc"),
							g.Combo("##rc", rcItems[rc], rcItems, &rc, 85, 0, nil),

							// g.Label("cq"),
							// g.InputIntV("##cq", 40, &cq, 0, nil),

							g.Label("qmin"),
							g.InputIntV("##qmin", 40, &qmin, 0, nil),

							g.Label("qmax"),
							g.InputIntV("##qmax", 40, &qmax, 0, nil),

							g.Label("aq"),
							g.Combo("##aq", aqItems[aq], aqItems, &aq, 85, 0, nil),

							// g.Label("aq-strength"),
							g.InputIntV("##aqstrength", 40, &aqStrength, 0, validateAQStrength),

							g.Label("bitrate"),
							g.InputIntV("k##bitrate", 70, &bitrate, 0, nil),

							// g.Label("Maxrate"),
							// g.InputIntV("k##maxrate", 60, &maxrate, 0, nil),
						),
					}),
					g.Spacing(),
					g.InputTextMultiline("", &ffmpegLog, 724, 200, 0, nil, func() {
						imgui.SetScrollHereY(1.0)
					}),
					g.Spacing(),
					g.ProgressBar(progress, 725, 20, ""),
					g.Line(
						g.Dummy(0, 5),
					),
					g.Line(
						g.Dummy(-67, 24),
						g.Condition(isEncoding,
							g.Layout{g.ButtonV("Cancel", 60, 24, handleCancelClick)},
							g.Layout{g.ButtonV("Run", 60, 24, handleRunClick)},
						),
					),
				},
				),
				g.TabItem("MediaInfo", g.Layout{
					g.Spacing(),
					g.InputTextMultiline("mediainfo", &mediaInfoLog, 724, 370, g.InputTextFlagsReadOnly, nil, nil),
				}),
			}),
		})
	g.PopStyleColorV(5)
	imgui.PopStyleVarV(4)
}

func main() {
	mw := g.NewMasterWindow("NVENC Video Toolbox 1.1", 740, 415, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsTransparent, loadFont)
	mw.SetBgColor(color.RGBA{0, 0, 0, 0})
	mw.SetDropCallback(handleDrop)
	mw.Main(loop)
	handleCancelClick()
}
