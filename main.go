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
	controlFlas      g.WindowFlags
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
	presetItems = []string{"slow", "bd"}
	rcItems     = []string{"vbr_hq", "vbr"}
	aqItems     = []string{"spatial", "temporal"}
)

var (
	preset int32
	rc     int32
	aq     int32
	cq     int32 = 26
	qmin   int32 = 16
	qmax   int32 = 26
	// bitrate       int32 = 6000
	// maxrate       int32 = 24000
	aqStrength int32 = 15
)

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
	progress = 0
	ffmpegLog = ""
	controlFlas = g.WindowFlagsNoInputs
	go ffmpeg.RunEncode(inputPath, outputPath, []string{
		"-c:a", "copy",
		// "-c:v", "libx264",
		"-c:v", "h264_nvenc",
		"-preset", "slow",
		"-profile:v", "high",
		"-level", "5.1",
		"-rc:v", "vbr_hq",
		"-cq", fmt.Sprint(cq),
		"-qmin", fmt.Sprint(qmin),
		"-qmax", fmt.Sprint(qmax),
		"-temporal-aq", "1",
		"-aq-strength:v", fmt.Sprint(aqStrength),
		// "-rc-lookahead:v", "32",
		// "-refs:v", "16",
		// "-bf:v", "3",
		"-coder:v", "cabac",
		// "-b:v", fmt.Sprintf("%dk", bitrate),
		// "-maxrate", fmt.Sprintf("%dk", maxrate),
		"-map", "0:0",
		"-f", "mp4",
	}, &progress, &ffmpegLog, &isEncoding, g.Update)
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
	ffmpegCmd := ffmpeg.GetFFMpegCmd()
	var err error
	if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "process", "where", "name='ffmpeg.exe'", "delete")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err = cmd.Run()
	} else {
		err = ffmpegCmd.Process.Kill()
	}
	if err != nil {
		fmt.Print(err)
		return
	}
	isEncoding = false
	controlFlas = g.WindowFlagsNone
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
					g.Child("control", false, 724, 90, controlFlas, g.Layout{
						g.Spacing(),
						g.Line(
							g.InputTextV("##video", -55, &inputPath, 0, nil, nil),
							g.ButtonV("video", 60, 22, handleInputClick),
						),
						g.Spacing(),
						g.Line(
							g.InputTextV("##output", -55, &outputPath, 0, nil, nil),
							g.ButtonV("output", 60, 22, handleOutputClick),
						),
						g.Spacing(),
						g.Line(
							g.Label("Preset"),
							g.Combo("##preset", presetItems[preset], presetItems, &preset, 85, 0, nil),

							g.Label("RC"),
							g.Combo("##rc", rcItems[rc], rcItems, &rc, 85, 0, nil),

							g.Label("CQ"),
							g.InputIntV("##cq", 40, &cq, 0, nil),

							g.Label("QMin"),
							g.InputIntV("##qmin", 40, &qmin, 0, nil),

							g.Label("QMax"),
							g.InputIntV("##qmax", 40, &qmax, 0, nil),

							g.Label("AQ"),
							g.Combo("##aq", aqItems[aq], aqItems, &aq, 85, 0, nil),

							g.Label("AQStrength"),
							g.InputIntV("##aqstrength", 40, &aqStrength, 0, validateAQStrength),

							// g.Label("Bitrate"),
							// g.InputIntV("k##bitrate", 60, &bitrate, 0, nil),

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

	if g.Context.GetPlatform().ShouldStop() {
		handleCancelClick()
	}
}

func main() {
	mw := g.NewMasterWindow("NVENC Video Toolbox", 740, 415, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsTransparent, loadFont)
	mw.SetBgColor(color.RGBA{0, 0, 0, 0})
	mw.SetDropCallback(handleDrop)
	mw.Main(loop)
}
