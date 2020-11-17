package mediainfo

import (
	"errors"
	"path/filepath"

	"github.com/Nicify/nvtool/execute"
)

var (
	binary string
)

func GetMediaInfo(mediaFile string) (string, error) {
	abspath, err := filepath.Abs(mediaFile)
	if err != nil {
		return "", errors.New("file not found.")
	}
	stdout, _, _ := execute.ExecSync(".", binary, abspath)
	mediainfo := string(stdout)
	return mediainfo, nil
}

func init() {
	path, _ := filepath.Abs("./core/mediainfo.exe")
	binary = path
}
