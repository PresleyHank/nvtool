package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"syscall"

	g "github.com/AllenDang/giu"
	"github.com/fsnotify/fsnotify"
	"github.com/melbahja/got"
	"github.com/saracen/go7z"
	"github.com/sqweek/dialog"
)

func nTrue(b ...bool) int {
	n := 0
	for _, v := range b {
		if v {
			n++
		}
	}
	return n
}

func byteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func byteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
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
	imageByte, err := assets.ReadFile(filename)
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

func limitResValue(val string) string {
	r := regexp.MustCompile(`[^0-9x-]`)
	return r.ReplaceAllString(val, "")
}

func invalidPath(inputPath string, outputPath string) bool {
	return inputPath == outputPath || inputPath == "" || outputPath == ""
}

func getErrorLevel(processState *os.ProcessState) (int, bool) {
	if processState.Success() {
		return 0, true
	} else if t, ok := processState.Sys().(syscall.WaitStatus); ok {
		return t.ExitStatus(), true
	} else {
		return 255, false
	}
}

func initSingleInstanceLock(lockFile string, onSecondInstance func(), onCommand func(command string)) (unlock func()) {
	if err := os.Remove(lockFile); err != nil && !os.IsNotExist(err) {
		onSecondInstance()
	}
	f, _ := os.Create(lockFile)
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
					command, _ := ioutil.ReadFile(lockFile)
					onCommand(string(command))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(lockFile)
	if err != nil {
		log.Fatal(err)
	}
	return func() {
		f.Close()
		watcher.Close()
	}
}

func extract7z(file string, dist string) (files []string, err error) {
	sz, err := go7z.OpenReader(file)
	if err != nil {
		return
	}
	defer sz.Close()

	for {
		hdr, err := sz.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if hdr.IsEmptyStream && !hdr.IsEmptyFile {
			if err := os.MkdirAll(hdr.Name, os.ModePerm); err != nil {
				return nil, err
			}
			continue
		}

		f, err := os.Create(path.Join(dist, hdr.Name))
		if err != nil {
			return nil, err
		}
		defer f.Close()

		if _, err := io.Copy(f, sz); err != nil {
			return files, err
		}
		files = append(files, f.Name())
	}
	return
}

func download(url string, dest string, onProgress func(float32)) ([]string, error) {
	tmp := path.Join(os.TempDir(), "nvtool_download.zip")
	defer os.Remove(tmp)
	d := &got.Download{
		URL:  url,
		Dest: tmp,
	}

	if err := d.Init(); err != nil {
		fmt.Print(err)
	}

	go d.RunProgress(func(d *got.Download) {
		progress := float32(d.Size()) / float32(d.TotalSize())
		onProgress(progress)
	})
	err := d.Start()
	d.StopProgress = true
	if err != nil {
		return nil, err
	}
	return extract7z(d.Name(), dest)
}
