package ffmpeg

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

// Progress ...
type Progress struct {
	FramesProcessed string
	CurrentTime     string
	CurrentBitrate  string
	Progress        float64
	Speed           string
}

const (
	PresetSlow = iota
	PresetMedium
	PresetFast
	PresetBD
)

const (
	CQVbrHQ = iota
	CQCbrHQ
)

const (
	AQTemporal = iota
	AQSpatial
)

var (
	PresetOptions = []string{"slow", "medium", "fast", "bd"}
	RCOptions     = []string{"vbr_hq", "cbr_hq"}
	AQOptions     = []string{"temporal", "spatial"}
)

var (
	binary       = "ffmpeg"
	prefix       = []string{"-y", "-hide_banner"}
	encodingTime uint
)

func spit(data []byte, atEOF bool) (advance int, token []byte, spliterror error) {
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

func GetVideoMeta(inputPath string) (uint, []byte, error) {
	args := append(prefix, "-ss", "3", "-skip_frame", "nokey", "-i", inputPath, "-vf", "thumbnail=10", "-frames:v", "1", "-vsync", "0", "-f", "image2", "-")
	preview, output, err := execSync(".", binary, args...)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		clean := strings.TrimSpace(line)
		if strings.HasPrefix(clean, "Duration:") {
			matches := regexp.MustCompile(durationRegexString).FindStringSubmatch(clean)
			duration := getDurationFromTimeParams(matches)
			return duration, nil, nil
		}
	}
	return 0, preview, err
}

func progress(stream io.ReadCloser, dursec float64, out chan Progress) {
	scanner := bufio.NewScanner(stream)
	scanner.Split(spit)

	buf := make([]byte, 2)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)

	for scanner.Scan() {
		Progress := new(Progress)
		line := scanner.Text()

		if strings.Contains(line, "frame=") && strings.Contains(line, "time=") && strings.Contains(line, "bitrate=") {
			var re = regexp.MustCompile(`=\s+`)
			st := re.ReplaceAllString(line, `=`)

			f := strings.Fields(st)

			var framesProcessed string
			var currentTime string
			var currentBitrate string
			var currentSpeed string

			for j := 0; j < len(f); j++ {
				field := f[j]
				fieldSplit := strings.Split(field, "=")

				if len(fieldSplit) > 1 {
					fieldname := strings.Split(field, "=")[0]
					fieldvalue := strings.Split(field, "=")[1]

					if fieldname == "frame" {
						framesProcessed = fieldvalue
					}

					if fieldname == "time" {
						currentTime = fieldvalue
					}

					if fieldname == "bitrate" {
						currentBitrate = fieldvalue
					}
					if fieldname == "speed" {
						currentSpeed = fieldvalue
					}
				}
			}

			timesec := DurToSec(currentTime)

			progress := (timesec * 100) / float64(dursec)
			Progress.Progress = progress

			Progress.CurrentBitrate = currentBitrate
			Progress.FramesProcessed = framesProcessed
			Progress.CurrentTime = currentTime
			Progress.Speed = currentSpeed

			out <- *Progress
		}
	}
}

// RunEncode ...
func RunEncode(inputPath string, outputPath string, args []string) (*exec.Cmd, <-chan Progress, error) {
	dursec, _, err := GetVideoMeta(inputPath)
	out := make(chan Progress)
	args = append([]string{"-i", inputPath}, args...)
	args = append(prefix, args...)
	args = append(args, outputPath)
	cmd := exec.Command(binary, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	stderr, _ := cmd.StderrPipe()
	err = cmd.Start()
	if err != nil {
	}

	go func() {
		progress(stderr, float64(dursec), out)
	}()

	go func() {
		defer close(out)
		err = cmd.Wait()
	}()

	return cmd, out, nil
}
