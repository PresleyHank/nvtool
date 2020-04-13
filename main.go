package main

import (
	"strconv"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	ffmpeg "github.com/Nicify/enclite/ffmpeg"
)

var (
	selectedTebIndex int
	inputPath        string
	outputPath       string
	fullDuration     uint
	isEncoding       bool
	progress         float32
	log              string
)

var (
	presetItems = []string{"slow", "bd"}
	rcItems     = []string{"vbr_hq", "vbr"}
)

var (
	preset        int32
	rc            int32
	cq            int32 = 26
	qmin          int32 = 16
	qmax          int32 = 26
	useTemporalAq       = false
	aqStrength          = 15
)

func handleInputClick() {
	path := selectInputPath()
	if len(path) > 1 {
		inputPath = path
	}
}

func handleOutputClick() {
	path := selectOutputPath()
	if len(path) > 1 {
		outputPath = path
	}
}

func handleRunClick() {
	if isEncoding {
		return
	}
	ffmpeg.RunEncode(inputPath, outputPath, []string{
		"-c:a", "copy",
		"-c:v", "h264_nvenc",
		"-preset", "slow",
		"-profile:v", "high",
		"-level", "5.1",
		"-rc:v", "vbr_hq",
		"-cq", strconv.Itoa(int(cq)),
		"-qmin", strconv.Itoa(int(qmin)),
		"-qmax", strconv.Itoa(int(qmax)),
		"-temporal-aq", "1",
		"-aq-strength:v", strconv.Itoa(aqStrength),
		// "-rc-lookahead:v", "32",
		// "-refs:v", "16",
		// "-bf:v", "3",
		"-coder:v", "cabac",
		// "-b:v", "6000k",
		// "-maxrate", "24000k",
		"-map", "0:0",
		"-f", "mp4",
	}, &progress, &log, &isEncoding, g.Update)
}

func loop() {
	g.Context.IO().SetFontGlobalScale(1)
	g.SingleWindow("Overview",
		g.Layout{
			g.TabBar("##maintab", g.Layout{
				g.TabItem("Encode", g.Layout{
					g.Spacing(),
					g.Line(
						g.InputTextV("##video", -55, &inputPath, 0, nil, nil),
						g.ButtonV("video", 60, 24, handleInputClick),
					),
					g.Spacing(),
					g.Line(
						g.InputTextV("##output", -55, &outputPath, 0, nil, nil),
						g.ButtonV("output", 60, 24, handleOutputClick),
					),
					g.Spacing(),
					g.Line(
						g.Label("Preset"),
						g.Combo("##preset", presetItems[preset], presetItems, &preset, 80, 0, nil),

						g.Label("RC"),
						g.Combo("##rc", rcItems[rc], rcItems, &rc, 80, 0, nil),

						g.Label("CQ"),
						g.InputIntV("##cq", 50, &cq, 0, nil),

						g.Label("QMin"),
						g.InputIntV("##qmin", 50, &qmin, 0, nil),

						g.Label("QMax"),
						g.InputIntV("##qmax", 50, &qmax, 0, nil),
					),
					g.Spacing(),
					g.InputTextMultiline("", &log, 705, 200, 0, nil, func() {
						imgui.SetScrollHereY(1.0)
					}),
					g.Spacing(),
					g.ProgressBar(progress, 705, 20, ""),
					g.Line(
						g.Dummy(0, 5),
					),
					g.Line(
						g.Dummy(-67, 24),
						g.ButtonV("Run", 60, 24, func() {
							go handleRunClick()
						}),
					),
				},
				),
				g.TabItem("MediaInfo", g.Layout{
					g.Label("mediaInfo"),
				}),
			}),
		})
}

func main() {
	w := g.NewMasterWindow("Overview", 720, 420, g.MasterWindowFlagsNotResizable, loadFont)
	w.Main(loop)
}
