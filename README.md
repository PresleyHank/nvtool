# NVTool

> 基于 FFMpeg 的 NVENC 视频压制工具

使用前请确保显卡支持 NVENC 编码器且 FFMpeg 与 MediaInfo 已正确安装，如未安装建议使用 windows 包管理工具 [scoop](https://scoop.sh)来安装:

install requirements

```powershell
Set-ExecutionPolicy RemoteSigned -scope CurrentUser
iwr -useb get.scoop.sh | iex
scoop install ffmpeg mediainfo
```

build

```sh
go build --ldflags '-s -w -extldflags "-static" -H=windowsgui' .
```

screenshots

![图片](https://uploader.shimo.im/f/l2Jc4yLrJSUdBEzW.png!thumbnail)

![图片](https://uploader.shimo.im/f/BjgbrHlAiuQe8TwM.png!thumbnail)

release

- Windows [NVTool.zip](https://attachments-cdn.shimo.im/sZZHbm7aVeceNHhK/NVTool.zip)
