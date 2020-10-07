package main

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	gpu "github.com/Nicify/nvtool/gpu"
	"github.com/gobuffalo/packr/v2"
	"github.com/sqweek/dialog"
)

type hwnd uintptr

type accentpolicy struct {
	nAccentState int
	nFlags       int
	nColor       int
	nAnimationID int
}

type wincompattrdata struct {
	nAttribute int
	pData      *accentpolicy
	ulDataSize uintptr
}

var (
	box                               = packr.New("assets", "./assets")
	mod                               = windows.NewLazyDLL("user32.dll")
	procSetWindowCompositionAttribute = mod.NewProc("SetWindowCompositionAttribute")
)

func setWindowCompositionAttribute(hwnd hwnd) {
	accent := accentpolicy{3, 0, 0, 0}
	data := wincompattrdata{19, &accent, unsafe.Sizeof(accent)}
	procSetWindowCompositionAttribute.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&data)),
	)
}

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

func imageToTexture(filename string) (*g.Texture, error) {
	imageByte, err := box.Find(filename)
	if err != nil {
		return nil, err
	}
	imageRGBA, _ := loadImageFromMemory(imageByte)
	textureID, err := g.NewTextureFromRgba(imageRGBA)
	return textureID, err
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
