package nvenc

import (
	"bytes"
	"io"
	"os/exec"

	"syscall"
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
