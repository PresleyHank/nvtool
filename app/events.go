package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/Nicify/nvtool/helper"
	"github.com/Nicify/nvtool/preset"
)

func (app *Application) OnRunClick() {
	if app.NVENC.IsEncoding() || helper.InvalidPath(app.State.InputPath, app.State.OutputPath) {
		return
	}
	go func() {
		defer app.Update()
		app.ResetState()
		args := preset.GetCommandLineArgs(app.EncodingPreset)
		progress, _ := app.NVENC.RunEncode(app.State.InputPath, app.State.OutputPath, args)
		for msg := range progress {
			app.State.Percent = float32(msg.Percent) / 100
			app.Usage.GPU = int32(msg.GPU)
			app.Usage.VE = int32(msg.VE)
			app.Usage.VD = int32(msg.VD)
			app.State.NVENCLog += fmt.Sprintf("\n%v frames: %.0f fps, %v kb/s, remain %s, est out size %s", msg.FramesProcessed, msg.FPS, msg.Bitrate, msg.Remain, msg.EstOutSize)
			app.State.NVENCLog = strings.Trim(app.State.NVENCLog, "\n")
			app.Update()
		}

		if app.NVENC.Cmd.ProcessState != nil && app.NVENC.Cmd.ProcessState.Success() {
			app.State.Percent = 1.0
		}
	}()
}

func (app *Application) OnMounted() {
	go func() {
		defer app.Window.GLFWWindow.SetOpacity(0.98)
		texLogo, _ := app.LoadTexture("embed/icon.png")
		texButtonClose, _ := app.LoadTexture("embed/close_white.png")
		texDropDown, _ := app.LoadTexture("embed/dropdown.png")
		texGraphicsCard, _ := app.LoadTexture("embed/graphics_card.png")
		app.Textures = &Textures{
			texLogo:         texLogo,
			texButtonClose:  texButtonClose,
			texDropDown:     texDropDown,
			texGraphicsCard: texGraphicsCard,
		}
	}()
}

func (app *Application) OnInputClick() {
	filePath := helper.SelectInputPath()
	if len(filePath) > 1 {
		app.State.InputPath = filePath
		fileExt := path.Ext(app.State.InputPath)
		app.State.OutputPath = strings.Replace(app.State.InputPath, fileExt, "_nvenc.mp4", 1)
		go app.GetMediaInfo(filePath)
	}
}

func (app *Application) OnOutputClick() {
	filePath := helper.SelectOutputPath()
	if len(filePath) > 1 {
		app.State.OutputPath = filePath
	}
}

func (app *Application) OnSecondInstance() {
	ioutil.WriteFile(app.LockFile, []byte("focus"), 0644)
	os.Exit(0)
}

func (app *Application) OnCommand(command string) {
	if command == "focus" {
		app.Window.GLFWWindow.Restore()
	}
}

func (app *Application) OnDrop(dropItem []string) {
	if app.NVENC.IsEncoding() {
		return
	}
	app.State.InputPath = dropItem[0]
	fileExt := path.Ext(app.State.InputPath)
	app.State.OutputPath = strings.Replace(app.State.InputPath, fileExt, "_nvenc.mp4", 1)
	go app.GetMediaInfo(app.State.InputPath)
}
