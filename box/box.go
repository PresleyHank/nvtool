package box

import (
	"bytes"
	"io"

	"github.com/markbates/pkger"
)

func Find(path string) (data []byte, err error) {
	buf := bytes.NewBuffer(nil)
	f, err := pkger.Open(path)
	if err != nil {
		return data, err
	}
	io.Copy(buf, f)
	f.Close()
	data = buf.Bytes()
	return
}
