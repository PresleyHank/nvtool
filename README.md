# NVTool

> NVENC video encoding tool based on FFMpeg

Make sure that the graphics card supports NVENC encoders and that FFMpeg and MediaInfo are properly installed before using it.
If not, we recommend using the windows package management tool [scoop](https://scoop.sh) to install it:

#### Install requirements

```powershell
Set-ExecutionPolicy RemoteSigned -scope CurrentUser
iwr -useb get.scoop.sh | iex
scoop install ffmpeg mediainfo
```

#### Build

```sh
go build -ldflags='-s -w -H windowsgui -linkmode external -extldflags -static' .
```

#### Screenshots

![图片](https://images-cdn.shimo.im/nXrLjx6Sm8rufqat__original.png)

![图片](https://images-cdn.shimo.im/rQnQ5c9NOyEv9PPv__original.png)

---

#### Release

- Windows [NVTool.zip](https://attachments-cdn.shimo.im/XErIszdTagxUpyqr.zip?attname=nvtool.zip)
