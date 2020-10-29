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

	execute "github.com/Nicify/nvtool/execute"
)

// Progress ...
type Progress struct {
	Percent         float64
	FramesProcessed int
	FPS             float64
	Bitrate         float64
	Remain          string
	GPU             int
	VE              int
	VD              int
	EstOutSize      string
}

type VPPKNNParam struct {
	Radius   int32
	Strength float32
	Lerp     float32
	ThLerp   float32
}

type VPPPMDParam struct {
	ApplyCount int32
	Strength   int32
	Threshold  int32
}

type VPPUnSharpParam struct {
	Radius    int32
	Weight    float32
	Threshold float32
}

type VPPEdgeLevelParam struct {
	Strength  float32
	Threshold float32
	Black     float32
	White     float32
}

type VPPSmoothParam struct {
	Quality int32
	QP      int32
	Prec    string
}

type VPPColorSpaceParam struct {
	HDR2SDR    string
	SourcePeak float32
	LdrNits    float32
}

var DefaultVPPKNNParam = VPPKNNParam{
	Radius:   3,
	Strength: 0.08,
	Lerp:     0.2,
	ThLerp:   0.8,
}

var DefaultVPPPMDParam = VPPPMDParam{
	ApplyCount: 2,
	Strength:   100,
	Threshold:  100,
}

var DefaultVPPUnSharpParam = VPPUnSharpParam{
	Radius:    3,
	Weight:    0.5,
	Threshold: 10.0,
}

var DefaultVPPEdgeLevelParam = VPPEdgeLevelParam{
	Strength:  10.0,
	Threshold: 20.0,
	Black:     0,
	White:     0,
}

var DefaultVPPSmoothParam = VPPSmoothParam{
	Quality: 6,
	QP:      12,
	Prec:    "fp32",
}

var DefaultVPPColorSpaceParam = VPPColorSpaceParam{
	HDR2SDR:    "hdr2sdr=hable",
	SourcePeak: 1000.0,
	LdrNits:    100.0,
}

var (
	binary string

	PresetOptions       = []string{"P1", "P2", "P3", "P4", "P5", "P6", "P7"}
	AQOptions           = []string{"aq-temporal", "aq"}
	AQOptionsForPreview = []string{"temporal", "spatial"}
)

func progress(stream io.ReadCloser, out chan Progress) {
	scanner := bufio.NewScanner(stream)
	scanner.Split(execute.Spit)

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
					Progress.GPU, _ = strconv.Atoi(strings.Split(filed, " ")[1])
					continue
				}
				if strings.HasPrefix(filed, "VE") {
					Progress.VE, _ = strconv.Atoi(strings.Split(filed, " ")[1])
					continue
				}
				if strings.HasPrefix(filed, "VD") {
					Progress.VD, _ = strconv.Atoi(strings.Split(filed, " ")[1])
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

func CheckDevice() (name string, err error) {
	stdout, _, err := execute.ExecSync(".", binary, "--check-device")
	if err != nil {
		return
	}
	r := regexp.MustCompile(`DeviceId #\d+: `)
	name = strings.Trim(r.ReplaceAllString(string(stdout), ""), "\n")
	return
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
