package web

import (
	"embed"
	"io/fs"
	"os"
)

//go:embed *
var web embed.FS

var FS = func(dir string) fs.FS {
	if d, _ := fs.ReadDir(web, "dist"); len(d) == 0 {
		return os.DirFS(dir)
	}
	f, _ := fs.Sub(web, "dist")
	if _, err := os.Stat(dir); nil != err {
		return f
	}
	return fileSystem{
		os.DirFS(dir),
		f,
	}
}

type fileSystem []fs.FS

func (f fileSystem) Open(name string) (file fs.File, err error) {
	for _, i := range f {
		if file, err = i.Open(name); err == nil {
			break
		}
	}
	return
}
