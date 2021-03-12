package preset

import (
	"fmt"
	"strings"

	nvenc "github.com/Nicify/nvtool/nvenc"
)

type EncodingPresets struct {
	HEVC       bool
	Preset     int32
	Quality    int32
	Bitrate    int32
	Maxrate    int32
	AQ         int32
	AQStrength int32
	Resize     bool
	OutputRes  string
	VPPSwitches
	VPPParams
}

type VPPSwitches struct {
	VPPKNN        bool
	VPPPMD        bool
	VPPUnSharp    bool
	VPPEdgeLevel  bool
	VPPSmooth     bool
	VPPColorSpace bool
}

type VPPParams struct {
	nvenc.VPPKNNParam
	nvenc.VPPPMDParam
	nvenc.VPPUnSharpParam
	nvenc.VPPEdgeLevelParam
}

var (
	DefaultVppSwitches = VPPSwitches{}

	DefaultVPPParams = VPPParams{
		nvenc.DefaultVPPKNNParam,
		nvenc.DefaultVPPPMDParam,
		nvenc.DefaultVPPUnSharpParam,
		nvenc.DefaultVPPEdgeLevelParam,
	}

	DefaultPreset = EncodingPresets{
		HEVC:        false,
		Preset:      6,
		Quality:     12,
		Bitrate:     19000,
		Maxrate:     59850,
		OutputRes:   "1920x1080",
		VPPSwitches: DefaultVppSwitches,
		VPPParams:   DefaultVPPParams,
	}
)

func GetCommandLineArgs(presets EncodingPresets) []string {
	codec := "h264"
	if presets.HEVC {
		codec = "hevc"
	}
	command := fmt.Sprintf("--codec %s --profile high --audio-copy --preset %s --vbr %v --vbr-quality %v --max-bitrate 60000 --bframes 4 --ref 16 --lookahead 32 --gop-len 250 --%s --aq-strength %v",
		codec,
		nvenc.PresetOptions[presets.Preset],
		presets.Bitrate,
		presets.Quality,
		nvenc.AQOptions[presets.AQ],
		presets.AQStrength,
	)
	args := strings.Split(command, " ")

	if presets.VPPSwitches.VPPKNN {
		param := presets.VPPKNNParam
		args = append(args, "--vpp-knn", fmt.Sprintf("radius=%v,strength=%.2f,lerp=%.1f,th_lerp=%.1f", param.Radius, param.Strength, param.Lerp, param.ThLerp))
	}

	if presets.VPPSwitches.VPPPMD {
		param := presets.VPPPMDParam
		args = append(args, "--vpp-pmd", fmt.Sprintf("apply_count=%v,strength=%v,threshold=%v", param.ApplyCount, param.Strength, param.Threshold))
	}

	if presets.VPPSwitches.VPPUnSharp {
		param := presets.VPPUnSharpParam
		args = append(args, "--vpp-unsharp", fmt.Sprintf("radius=%v,weight=%.1f,threshold=%.1f", param.Radius, param.Weight, param.Threshold))
	}

	if presets.VPPSwitches.VPPEdgeLevel {
		param := presets.VPPEdgeLevelParam
		args = append(args, "--vpp-edgelevel", fmt.Sprintf("strength=%v,threshold=%.1f,black=%v,white=%v", param.Strength, param.Threshold, param.Black, param.White))
	}

	if presets.Resize {
		args = append(args, "--vpp-resize", "lanczos2", "--output-res", presets.OutputRes)
	}
	return args
}
