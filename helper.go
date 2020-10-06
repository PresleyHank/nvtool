package main

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	gpu "github.com/Nicify/nvtool/gpu"
	"github.com/sqweek/dialog"
)

func getGpuNames() string {
	gpuList, err := gpu.GetGPUInfo()
	if err != nil {
		return "Error getting GPU info"
	}
	return strings.Join(gpuList, " ")
}

func loadFont() {
	fonts := g.Context.IO().Fonts()
	font, _ := box.Find("iosevka.ttf")
	fonts.AddFontFromMemoryTTFV(font, 18, imgui.DefaultFontConfig, fonts.GlyphRangesChineseFull())
}

func loadImageFromMemory(imageData []byte) (imageRGBA *image.RGBA, err error) {
	r := bytes.NewReader(imageData)
	img, err := png.Decode(r)
	if err != nil {
		return nil, err
	}

	switch trueImg := img.(type) {
	case *image.RGBA:
		return trueImg, nil
	default:
		rgba := image.NewRGBA(trueImg.Bounds())
		draw.Draw(rgba, trueImg.Bounds(), trueImg, image.Pt(0, 0), draw.Src)
		return rgba, nil
	}
}

func selectInputPath() string {
	path, _ := dialog.File().Filter("video file", "mp4", "mkv", "mov", "flv", "avi").Load()
	return path
}

func selectOutputPath() string {
	path, _ := dialog.File().Filter("video file", "mp4", "mkv", "mov", "flv", "avi").Save()
	return path
}

func limitValue(val int32, min int32, max int32) int32 {
	if val > max {
		return max
	}
	if val < min {
		return min
	}
	return val
}

func invalidPath(inputPath string, outputPath string) bool {
	return inputPath == outputPath || inputPath == "" || outputPath == ""
}
