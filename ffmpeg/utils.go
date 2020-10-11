package ffmpeg

import (
	"bytes"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"syscall"
)

const durationRegexString = `(\d{2}):(\d{2}):(\d{2})\.(\d{2})`

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
	if err != nil {
		return nil, bufErr.Bytes(), err
	}
	return buf.Bytes(), bufErr.Bytes(), err
}

func DurationToSec(time string) uint {
	matches := regexp.MustCompile(durationRegexString).FindStringSubmatch(time)
	var (
		hour     uint64
		min      uint64
		sec      uint64
		ms       uint64
		duration uint
	)
	hour, _ = strconv.ParseUint(matches[1], 10, 32)
	min, _ = strconv.ParseUint(matches[2], 10, 32)
	sec, _ = strconv.ParseUint(matches[3], 10, 32)
	ms, _ = strconv.ParseUint(matches[4], 10, 32)
	duration = uint(hour*60*60*1000 + min*60*1000 + sec*1000 + ms*10)

	return duration
}
