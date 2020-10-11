package ffmpeg

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

// Progress ...
type Progress struct {
	FramesProcessed int
	CurrentSize     string
	CurrentTime     string
	CurrentBitrate  string
	Progress        float64
	Q               float64
	FPS             int
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

const durationRegexString = `(\d{2}):(\d{2}):(\d{2})\.(\d{2})`

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
			duration := DurationToSec(matches)
			return duration, nil, nil
		}
	}
	return 0, preview, err
}

func progress(stream io.ReadCloser, durationInMs uint, out chan Progress) {
	scanner := bufio.NewScanner(stream)
	scanner.Split(spit)

	buf := make([]byte, 2)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)

	for scanner.Scan() {
		Progress := new(Progress)
		line := scanner.Text()
		line = strings.ReplaceAll(line, "frame=", "\nframe=")
		line = strings.ReplaceAll(line, "= ", "=")
		if strings.Contains(line, "frame=") && strings.Contains(line, "time=") && strings.Contains(line, "bitrate=") {
			var re = regexp.MustCompile(`=\s+`)
			st := re.ReplaceAllString(line, `=`)

			f := strings.Fields(st)

			var framesProcessed int
			var currentSize string
			var currentTime string
			var currentBitrate string
			var progress float64
			var q float64
			var fps int
			var speed string

			for j := 0; j < len(f); j++ {
				field := f[j]
				fieldSplit := strings.Split(field, "=")

				if len(fieldSplit) > 1 {
					fieldname := strings.Split(field, "=")[0]
					fieldvalue := strings.Split(field, "=")[1]

					switch fieldname {
					case "frame":
						framesProcessed, _ = strconv.Atoi(fieldvalue)
					case "fps":
						fps, _ = strconv.Atoi(fieldvalue)
					case "q":
						q, _ = strconv.ParseFloat(fieldvalue, 64)
					case "size":
						currentSize = fieldvalue
					case "Lsize":
						currentSize = fieldvalue
					case "time":
						currentTime = fieldvalue
					case "bitrate":
						currentBitrate = fieldvalue
					case "speed":
						speed = fieldvalue
					default:
						log.Printf("%s: %v", fieldname, fieldvalue)
					}
				}
			}

			matches := regexp.MustCompile(durationRegexString).FindStringSubmatch(currentTime)
			timesec := DurationToSec(matches)
			progress = float64(timesec) / float64(durationInMs)
			Progress.Progress = progress

			Progress.FramesProcessed = framesProcessed
			Progress.FPS = fps
			Progress.Q = q
			Progress.CurrentSize = currentSize
			Progress.CurrentTime = currentTime
			Progress.CurrentBitrate = currentBitrate
			Progress.Speed = speed

			out <- *Progress
		}
	}
}

// RunEncode ...
func RunEncode(inputPath string, outputPath string, args []string) (*exec.Cmd, <-chan Progress, error) {
	durationInMs, _, err := GetVideoMeta(inputPath)
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
		progress(stderr, durationInMs, out)
	}()

	go func() {
		defer close(out)
		err = cmd.Wait()
	}()

	return cmd, out, nil
}
