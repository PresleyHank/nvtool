package nvenc

import (
	"bufio"
	"io"
	"os"
	"os/exec"
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

var (
	PresetOptions       = []string{"P1", "P2", "P3", "P4", "P5", "P6", "P7"}
	AQOptions           = []string{"aq-temporal", "aq"}
	AQOptionsForPreview = []string{"temporal", "spatial"}
)

type NVENC struct {
	binaryPath string
	Cmd        *exec.Cmd
}

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

func New(binaryPath string) *NVENC {
	return &NVENC{binaryPath: binaryPath}
}

func (n *NVENC) CheckDevice() (name string, err error) {
	stdout, _, err := execute.ExecSync(".", n.binaryPath, "--check-device")
	if err != nil {
		return
	}
	r := regexp.MustCompile(`DeviceId #\d+: `)
	name = strings.Trim(r.ReplaceAllString(string(stdout), ""), "\n")
	return
}

func (n *NVENC) RunEncode(inputPath string, outputPath string, args []string) (<-chan Progress, error) {
	out := make(chan Progress)
	args = append([]string{"-i", inputPath}, args...)
	args = append(args, "-o", outputPath)
	n.Cmd = exec.Command(n.binaryPath, args...)
	n.Cmd.Path = n.binaryPath
	n.Cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	stderr, _ := n.Cmd.StderrPipe()
	stdout, _ := n.Cmd.StdoutPipe()
	go func() {
		io.Copy(os.Stdout, stdout)
	}()

	err := n.Cmd.Start()
	if err != nil {
		return nil, err
	}

	go func() {
		progress(stderr, out)
	}()

	go func() {
		defer close(out)
		err = n.Cmd.Wait()
	}()

	return out, err
}

func (n *NVENC) IsEncoding() bool {
	if n.Cmd == nil || (n.Cmd.ProcessState != nil && n.Cmd.ProcessState.Exited()) {
		return false
	}
	return true
}

func (n *NVENC) Stop() {
	if n.Cmd == nil {
		return
	}
	n.Cmd.Process.Kill()
	go n.Cmd.Wait()
}
