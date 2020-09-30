package ffmpeg

import (
	"bytes"
	"io"
	"os/exec"
	"strconv"
	"syscall"
)

const (
	durationRegexString      = `Duration: (\d{2}):(\d{2}):(\d{2})\.(\d{2})`
	encodingTimeRegexString  = `time=(\d{2}):(\d{2}):(\d{2})\.(\d{2})`
	encodingSpeedRegexString = `speed=\d+\.\d+x`
)

func execSync(pwd string, command string, args ...string) ([]byte, []byte, error) {
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
	if err := cmd.Run(); err != nil {
		return nil, bufErr.Bytes(), err
	}
	return buf.Bytes(), bufErr.Bytes(), err
}

func getDurationFromTimeParams(time []string) uint {
	var (
		hour     uint64
		min      uint64
		sec      uint64
		ms       uint64
		duration uint
	)
	hour, err := strconv.ParseUint(time[1], 10, 32)
	if err == nil {
		min, err = strconv.ParseUint(time[2], 10, 32)
	}
	if err == nil {
		sec, err = strconv.ParseUint(time[3], 10, 32)
	}
	if err == nil {
		ms, err = strconv.ParseUint(time[4], 10, 32)
	}
	if err == nil {
		duration = uint(hour*60*60*1000 + min*60*1000 + sec*1000 + ms*10)
	}
	return duration
}
