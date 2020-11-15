# NVTool

> NVENC video encoding tool based on NVEncC

#### Build

```powershell
packr2 
rsrc -manifest nvtool.exe.manifest -ico ./assets/icon.ico -arch amd64 -o rsrc.syso
go build -ldflags="-s -w -H windowsgui -linkmode external -extldflags -static" .
rm rsrc.syso
```

#### Screenshots

![图片](https://images-cdn.shimo.im/XUvfrHzbyPOKQn9c__original.png)

![图片](https://images-cdn.shimo.im/dLiWypVO9fbgAXPb__original.png)


Video encoding tool based on rigaya's NVEncC.

Check out the link below to see if your graphics card supports NVENC.

https://developer.nvidia.com/video-encode-and-decode-gpu-support-matrix-new

### Options Reference:

#### Preset

Encode quality preset.

- P1 (= performance)
- P2
- P3
- P4 (= default)
- P5
- P6
- P7 (= quality)

#### Quality

Set target quality when using VBR mode. (0.0-51.0, 0 = automatic)

#### Bitrate

Set bitrate in kbps.

To enable constant quality mode, set the value to 0.

#### AQ

- temporal - Enable adaptive quantization between frames.

- spatial - Enable adaptive quantization in frame.

#### Strength

Specify the AQ strength. 1 = weak, 15 = strong, 0 = auto

#### HEVC

Specify the output codec as HEVC, default is h264.

#### Resize

Set output resolution. When it is different from the input resolution, HW/GPU resizer (Lanczos 4x4) will be activated automatically.

If not specified, it will be same as the input resolution. (no resize)

*Special Values*

- 0 ... Will be same as input.
- One of width or height as negative value
  Will be resized keeping aspect ratio, and a value which could be divided by the negative value will be chosen.

Example

input  1280x720
Resize 1024x576 -> normal
Resize 960x0    -> resize to 960x720 (0 will be replaced to 720, same as input)
Resize 1920x-2  -> resize to 1920x1080 (calculated to keep aspect ratio)

### Filter Reference:

#### KNN

Strong noise reduction filter.

Parameters

- radius=<int> (default=3, 1-5)
  radius of filter.
- strength=<float> (default=0.08, 0.0 - 1.0)
  Strength of the filter.
- lerp=<float> (default=0.2, 0.0 - 1.0)
  The degree of blending of the original pixel to the noise reduction pixel.
- th_lerp=<float> (default=0.8, 0.0 - 1.0)
  Threshold of edge detection.

#### PMD

Rather weak noise reduction by modified pmd method, aimed to preserve edge while noise reduction.

Parameters

- apply_count=<int> (default=2, 1- )
  Number of times to apply the filter.
- strength=<float> (default=100, 0-100)
  Strength of the filter.
- threshold=<float> (default=100, 0-255)
  Threshold for edge detection. The smaller the value is, more will be detected as edge, which will be preserved.

#### Unsharp

unsharp filter, for edge and detail enhancement.

Parameters

- radius=<int> (default=3, 1-9)
  radius of edge / detail detection.
- weight=<float> (default=0.5, 0-10)
  Strength of edge and detail emphasis. Larger value will result stronger effect.
- threshold=<float> (default=10.0, 0-255)
  Threshold for edge and detail detection.

#### EdgeLevel

Edge level adjustment filter, for edge sharpening.

Parameters

- strength=<float> (default=5.0, -31 - 31)
  Strength of edge sharpening. Larger value will result stronger edge sharpening.
- threshold=<float> (default=20.0, 0 - 255)
  Noise threshold to avoid enhancing noise. Larger value will treat larger luminance change as noise.
- black=<float> (default=0.0, 0-31)
  strength to enhance dark part of edges.
- white=<float> (default=0.0, 0-31)
  strength to enhance bright part of edges.