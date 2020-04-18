package mediainfo

import (
	"bufio"
	"os/exec"
)

var (
	mediainfoBinary = "mediainfo"
)

func GetMediaInfo(mediaFile string) (mediainfo []string, err error) {
	cmd := exec.Command(mediainfoBinary, mediaFile)
	stdout, err := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		mediainfo = append(mediainfo, line)
	}

	if err = scanner.Err(); err != nil {
		return
	}
	return mediainfo, nil
}
