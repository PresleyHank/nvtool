package gpu

import (
	"github.com/jaypipes/ghw"
)

// GetGPUInfo get the name of the installed GPUs
func GetGPUInfo() (gpuList []string, err error) {
	gpu, err := ghw.GPU()
	for _, card := range gpu.GraphicsCards {
		gpuList = append(gpuList, card.DeviceInfo.Product.Name)
	}
	return
}
