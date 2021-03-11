package mediainfo

import (
	"errors"
	"path/filepath"

	"github.com/Nicify/nvtool/helper"
)

type MediaInfo struct {
	binaryPath string
}

func New(binaryPath string) *MediaInfo {
	return &MediaInfo{binaryPath: binaryPath}
}

func (m *MediaInfo) GetMediaInfo(mediaFile string) (string, error) {
	abspath, err := filepath.Abs(mediaFile)
	if err != nil {
		return "", errors.New("file not found.")
	}
	stdout, _, _ := helper.ExecSync(".", m.binaryPath, abspath)
	mediainfo := string(stdout)
	return mediainfo, nil
}
