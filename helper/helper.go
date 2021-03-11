package helper

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"syscall"

	"github.com/melbahja/got"
	"github.com/saracen/go7z"
	"github.com/sqweek/dialog"
)

func Spit(data []byte, atEOF bool) (advance int, token []byte, spliterror error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}
	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		// We have a cr terminated line
		return i + 1, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func ExecSync(pwd string, command string, args ...string) ([]byte, []byte, error) {
	cmd := exec.Command(command, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Dir = pwd

	buf := &bytes.Buffer{}
	bufErr := &bytes.Buffer{}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	go io.Copy(buf, stdout)
	go io.Copy(bufErr, stderr)
	err := cmd.Run()
	if err != nil {
		return nil, bufErr.Bytes(), err
	}
	return buf.Bytes(), bufErr.Bytes(), err
}

func ByteCountDecimal(b int64) string {
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

func ByteCountBinary(b int64) string {
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

func SelectInputPath() string {
	path, _ := dialog.File().Filter("video file", "mp4", "mkv", "mov", "flv", "avi").Load()
	return path
}

func SelectOutputPath() string {
	path, _ := dialog.File().Filter("video file", "mp4", "mkv", "mov", "flv", "avi").Save()
	return path
}

func LimitValue(val int32, min int32, max int32) int32 {
	if val > max {
		return max
	}
	if val < min {
		return min
	}
	return val
}

func LimitResValue(val string) string {
	r := regexp.MustCompile(`[^0-9x-]`)
	return r.ReplaceAllString(val, "")
}

func InvalidPath(inputPath string, outputPath string) bool {
	return inputPath == outputPath || inputPath == "" || outputPath == ""
}

func GetErrorLevel(processState *os.ProcessState) (int, bool) {
	if processState.Success() {
		return 0, true
	} else if t, ok := processState.Sys().(syscall.WaitStatus); ok {
		return t.ExitStatus(), true
	} else {
		return 255, false
	}
}

func LoadImageFromMemory(imageData []byte) (imageRGBA *image.RGBA, err error) {
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

func Extract7z(file string, dist string) (files []string, err error) {
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

func Download(url string, dest string, onProgress func(float32)) error {
	d := &got.Download{
		URL:  url,
		Dest: dest,
	}

	if err := d.Init(); err != nil {
		return err
	}

	go d.RunProgress(func(d *got.Download) {
		progress := float32(d.Size()) / float32(d.TotalSize())
		onProgress(progress)
	})
	err := d.Start()
	d.StopProgress = true
	return err
}
