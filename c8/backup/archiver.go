package backup

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// Archiver represents type capable of archiving and
// restoring files.
type Archiver interface {
	DestFmt() string
	Archive(src, dest string) error
}

type zipper struct{}

// Zip is an Archiver that zips and unzips files.
var ZIP Archiver = (*zipper)(nil)

func (z *zipper) DestFmt() string {
	return "%d.zip"
}

func (z *zipper) Archive(src, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0777); err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	w := zip.NewWriter(out) /// create new zip.Writer type
	defer w.Close()

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil // skip
		}
		if err != nil {
			return err
		}
		in, err := os.Open(path) // open for reading
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
