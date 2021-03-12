package app

import (
	"image"
	"os"
	"path"
	"path/filepath"

	"sync"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/Nicify/nvtool/hooks"
	"github.com/Nicify/nvtool/mediainfo"
	"github.com/Nicify/nvtool/nvenc"
	"github.com/Nicify/nvtool/preset"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type CustomFonts struct {
	fontTamzenr imgui.Font
	fontTamzenb imgui.Font
	fontIosevka imgui.Font
}

type Textures struct {
	texLogo         *g.Texture
	texButtonClose  *g.Texture
	texDropDown     *g.Texture
	texGraphicsCard *g.Texture
}

type Window struct {
	MW          *g.MasterWindow
	GLFWWindow  *glfw.Window
	MWMoveState *hooks.MWMoveState
	MWDragArea  image.Rectangle
}

type State struct {
	InputPath    string
	OutputPath   string
	Percent      float32
	GPUName      string
	NVENCLog     string
	MediaInfoLog string
}

type Usage struct {
	GPU int32
	VE  int32
	VD  int32
}

type Application struct {
	NVENC          *nvenc.NVENC
	MediaInfo      *mediainfo.MediaInfo
	LockFile       string
	mounted        bool
	Window         *Window
	Fonts          *CustomFonts
	Textures       *Textures
	State          *State
	Usage          *Usage
	EncodingPreset preset.EncodingPresets
	VPPSwitches    preset.VPPSwitches
	VPPParams      preset.VPPParams
	Update         func()
}

var (
	Version = 2.4
	app     *Application
	once    sync.Once
)

func GetInstance() *Application {
	once.Do(func() {
		basePath, _ := filepath.Abs(".")
		nPath := filepath.Join(basePath, "core", "NVEncC64.exe")
		mPath := filepath.Join(basePath, "core", "MediaInfo.exe")
		app = &Application{
			LockFile:  path.Join(os.TempDir(), "nvtool.lock"),
			NVENC:     nvenc.New(nPath),
			MediaInfo: mediainfo.New(mPath),
			Fonts:     &CustomFonts{},
			Textures:  &Textures{},
			State: &State{
				MediaInfoLog: "Drag and drop media files here",
			},
			Usage:          &Usage{},
			Update:         g.Update,
			EncodingPreset: preset.DefaultPreset,
			VPPSwitches:    preset.DefaultVppSwitches,
			VPPParams:      preset.DefaultVPPParams,
		}
	})
	return app
}

func (app *Application) InstallCore() {
	app.CheckCore()
	app.State.GPUName, _ = app.NVENC.CheckDevice()
	app.Update()
}

func (app *Application) Mount(window *Window) *Application {
	app.Window = window
	return app
}

func (app *Application) Run() {
	app.Window.MW.Run(app.Render)
}

func (app *Application) ResetState() {
	app.State.NVENCLog = ""
	app.State.Percent = 0
}
