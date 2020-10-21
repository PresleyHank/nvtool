package nvenc

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

// Progress ...
type Progress struct {
	Percent         float64
	FramesProcessed int
	FPS             float64
	Bitrate         float64
	Remain          string
	GPUUsage        int
	VEUsage         int
	VDUsage         int
	EstOutSize      string
}

var (
	binary string

	PresetOptions = []string{"default", "performance", "quality"}
	AQOptions     = []string{"temporal", "spatial"}
)

func progress(stream io.ReadCloser, out chan Progress) {
	scanner := bufio.NewScanner(stream)
	scanner.Split(spit)

	buf := make([]byte, 2)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)

	for scanner.Scan() {
		Progress := new(Progress)
		line := scanner.Text()
		if strings.Contains(line, "frames:") && strings.Contains(line, "remain") {
			r := regexp.MustCompile(`\[(\d+.\d+)%]`)
			st := r.ReplaceAllString(line, "$1 percent,")
			st = strings.ReplaceAll(st, "%", "")
			st = strings.ReplaceAll(st, "frames:", "frames,")
			st = strings.ReplaceAll(st, "est out size", "EST")

			for _, filed := range strings.Split(st, ", ") {
				if strings.HasSuffix(filed, "percent") {
					Progress.Percent, _ = strconv.ParseFloat(strings.Split(filed, " ")[0], 64)
					continue
				}
				if strings.HasSuffix(filed, "frames") {
					Progress.FramesProcessed, _ = strconv.Atoi(strings.Split(filed, " ")[0])
					continue
				}
				if strings.HasSuffix(filed, "fps") {
					Progress.FPS, _ = strconv.ParseFloat(strings.Split(filed, " ")[0], 64)
					continue
				}
				if strings.HasSuffix(filed, "kb/s") {
					Progress.Bitrate, _ = strconv.ParseFloat(strings.Split(filed, " ")[0], 64)
					continue
				}
				if strings.HasPrefix(filed, "remain") {
					Progress.Remain = strings.Split(filed, " ")[1]
					continue
				}
				if strings.HasPrefix(filed, "GPU") {
					Progress.GPUUsage, _ = strconv.Atoi(strings.Split(filed, " ")[1])
					continue
				}
				if strings.HasPrefix(filed, "VE") {
					Progress.VEUsage, _ = strconv.Atoi(strings.Split(filed, " ")[1])
					continue
				}
				if strings.HasPrefix(filed, "VD") {
					Progress.VDUsage, _ = strconv.Atoi(strings.Split(filed, " ")[1])
					continue
				}
				if strings.HasPrefix(filed, "EST") {
					Progress.EstOutSize = strings.Split(filed, " ")[1]
					continue
				}

			}

			out <- *Progress
		}
	}
}

// RunEncode ...
func RunEncode(inputPath string, outputPath string, args []string) (*exec.Cmd, <-chan Progress, error) {
	out := make(chan Progress)
	args = append([]string{"-i", inputPath}, args...)
	args = append(args, "-o", outputPath)
	fmt.Println(args)
	cmd := exec.Command(binary, args...)
	cmd.Path = binary
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	go func() {
		io.Copy(os.Stdout, stdout)
	}()
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		progress(stderr, out)
	}()

	go func() {
		defer close(out)
		err = cmd.Wait()
	}()

	return cmd, out, nil
}

func init() {
	path, err := filepath.Abs("./bin/NVEncC64.exe")
	if err != nil {
		panic("NVEncC64.exe not found!")
	}
	binary = path
}
