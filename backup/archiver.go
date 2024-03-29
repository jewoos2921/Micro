package backup

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type Archiver interface {
	DestFmt() string
	Archive(src, dest string) error
}
type zipper struct{}

// Zip은 파일을 압축하고 압축 해제하는 Archiver다
var ZIP Archiver = (*zipper)(nil)

func (z *zipper) Archive(src, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0777); err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	w := zip.NewWriter(out)
	defer w.Close()

	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil // 건너뜀
		}
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		f, err := w.Create(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, in)
		if err != nil {
			return err
		}
		return nil
	})
}
func (z *zipper) DestFmt() string {
	return "%d.zip"
}
