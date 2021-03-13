package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	assets "github.com/Nicify/nvtool/assets"
	"github.com/Nicify/nvtool/helper"
	"github.com/fsnotify/fsnotify"
)

func (app *Application) GetMediaInfo(inputPath string) {
	info, err := app.MediaInfo.GetMediaInfo(inputPath)
	if err != nil {
		app.State.MediaInfoLog = fmt.Sprintf("Error: %s", err)
		return
	}
	app.State.MediaInfoLog = info
	app.Update()
}

func (app *Application) LoadFont() {
	fonts := g.Context.IO().Fonts()
	fontIosevkaTTF, _ := assets.EmbedFS.ReadFile("embed/iosevka.ttf")
	fontIosevka := fonts.AddFontFromMemoryTTFV(fontIosevkaTTF, 18, imgui.DefaultFontConfig, fonts.GlyphRangesChineseFull())
	fontTamzenbTTF, _ := assets.EmbedFS.ReadFile("embed/tamzen8x16b.ttf")
	fontTamzenb := fonts.AddFontFromMemoryTTFV(fontTamzenbTTF, 16, imgui.DefaultFontConfig, fonts.GlyphRangesChineseFull())
	fontTamzenrTTF, _ := assets.EmbedFS.ReadFile("embed/tamzen8x16r.ttf")
	fontTamzenr := fonts.AddFontFromMemoryTTFV(fontTamzenrTTF, 16, imgui.DefaultFontConfig, fonts.GlyphRangesChineseFull())
	app.Fonts = &CustomFonts{
		fontIosevka: fontIosevka,
		fontTamzenb: fontTamzenb,
		fontTamzenr: fontTamzenr,
	}
}

func (app *Application) NewTextureFromMemory(imageByte []byte) (*g.Texture, error) {
	imageRGBA, _ := helper.LoadImageFromMemory(imageByte)
	textureID, err := g.NewTextureFromRgba(imageRGBA)
	return textureID, err
}

func (app *Application) LoadTexture() {
	load := func(filename string, textureID **g.Texture) {
		imageByte, _ := assets.EmbedFS.ReadFile(filename)
		newTex, _ := app.NewTextureFromMemory(imageByte)
		*textureID = newTex
	}
	go func() {
		defer app.Window.GLFWWindow.SetOpacity(0.98)
		// defer func(t time.Time) { fmt.Printf("--- Time Elapsed: %v ---\n", time.Since(t)) }(time.Now())
		load("embed/icon.png", &app.Textures.texLogo)
		load("embed/close_white.png", &app.Textures.texButtonClose)
		// load("embed/dropdown.png", &app.Textures.texDropDown)
		load("embed/graphics_card.png", &app.Textures.texGraphicsCard)
	}()
}

func (app *Application) CheckCore() {
	if _, err := os.Stat("core"); os.IsNotExist(err) {
		os.Mkdir("core", 0777)
	}

	if _, err := os.Stat("./core/MediaInfo.exe"); os.IsNotExist(err) {
		bytes, err := assets.EmbedFS.ReadFile("embed/MediaInfo.7z")
		if err != nil {
			return
		}
		tmpPath := path.Join(os.TempDir(), "mediainfo.7z")
		ioutil.WriteFile(tmpPath, bytes, 0777)
		helper.Extract7z(tmpPath, "core")
	}

	if _, err := os.Stat("./core/NVEncC64.exe"); os.IsNotExist(err) {
		app.State.NVENCLog = "Downloading NVEncC...\n"
		app.Update()
		tmp := path.Join(os.TempDir(), "NVEncC.7z")
		err := helper.Download("https://hub.fastgit.org/rigaya/NVEnc/releases/download/5.29/NVEncC_5.29_x64.7z", tmp, func(progress float32) {
			app.State.Percent = progress * 0.95
			app.Update()
		})
		if err != nil {
			fmt.Printf("Download failed %s", err)
			return
		}
		files, _ := helper.Extract7z(tmp, "./core")
		app.State.Percent = 1
		app.State.NVENCLog += strings.Join(files, "\n") + "\nDownload completed."
		os.Remove(tmp)
		app.Update()
	}
}

func (app *Application) SingleInstanceLock() (unlock func()) {
	if err := os.Remove(app.LockFile); err != nil && !os.IsNotExist(err) {
		app.OnSecondInstance()
	}
	f, _ := os.Create(app.LockFile)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					command, _ := ioutil.ReadFile(app.LockFile)
					app.OnCommand(string(command))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(app.LockFile)
	if err != nil {
		log.Fatal(err)
	}
	return func() {
		f.Close()
		watcher.Close()
	}
}

func (app *Application) ShouldDisableInput(b bool) (flag g.WindowFlags) {
	if b {
		return g.WindowFlagsNoInputs
	}
	return g.WindowFlagsNone
}
