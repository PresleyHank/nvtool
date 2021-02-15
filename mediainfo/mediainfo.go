package mediainfo

import (
	"errors"
	"path/filepath"

	"github.com/Nicify/nvtool/execute"
)

var (
	Binary string
)

func GetMediaInfo(mediaFile string) (string, error) {
	abspath, err := filepath.Abs(mediaFile)
	if err != nil {
		return "", errors.New("file not found.")
	}
	stdout, _, _ := execute.ExecSync(".", Binary, abspath)
	mediainfo := string(stdout)
	return mediainfo, nil
}

func init() {
	path, err := filepath.Abs("./core/MediaInfo.exe")
	if err == nil {
		Binary = path
	}
}
