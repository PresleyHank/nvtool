package ffmpeg

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

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
	process      *exec.Cmd
	encodingTime uint
)

func GetProcess() *exec.Cmd {
	return process
}

func GetVideoMeta(inputPath string) (uint, []byte, error) {
	args := append(prefix, "-hide_banner", "-ss", "3", "-skip_frame", "nokey", "-i", inputPath, "-vf", "thumbnail=10", "-frames:v", "1", "-vsync", "0", "-f", "image2", "-")
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

func RunEncode(inputPath string, outputPath string, args []string, progress *float32, log *string, next func()) {
	fullDuration, _, err := GetVideoMeta(inputPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	args = append([]string{"-i", inputPath}, args...)
	args = append(prefix, args...)
	args = append(args, outputPath)
	process = exec.Command(binary, args...)
	process.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	stderr, err := process.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}
	err = process.Start()
	if err != nil {
		fmt.Println(err)
	}
	var isEncoding bool
	var logList []string
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		line := scanner.Text()
		if isEncoding {
			logList = append(logList, line)
		}
		matches := regexp.MustCompile(encodingTimeRegexString).FindStringSubmatch(line)
		if len(matches) == 5 {
			encodingTime = getDurationFromTimeParams(matches)
			*progress = float32(encodingTime) / float32(fullDuration)
			next()
		}
		if regexp.MustCompile(encodingSpeedRegexString).MatchString(line) {
			isEncoding = true
			logChunk := strings.Join(logList, " ")
			logChunk = strings.ReplaceAll(logChunk, "frame=", "\nframe=")
			logChunk = strings.ReplaceAll(logChunk, "= ", "=")
			*log += logChunk
			logList = logList[:0]
			next()
		}
	}
}
