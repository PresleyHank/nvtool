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
bash build_win.sh
```

#### Screenshots

![图片](https://images-cdn.shimo.im/clz5A7fc57SDBHzu__original.png)

![图片](https://images-cdn.shimo.im/dLiWypVO9fbgAXPb__original.png)
