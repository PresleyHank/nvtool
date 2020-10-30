packr2 
rsrc -manifest nvtool.exe.manifest -ico ./assets/icon.ico -arch amd64 -o rsrc.syso
go build -ldflags="-s -w -H windowsgui -linkmode external -extldflags -static" .
rm rsrc.syso