package tool

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// Compress
// files 可以是不同目录下的文件或目录
func Compress(files []*os.File, dest string) error {
	if err := MkdirIfNotExists(filepath.Dir(dest)); err != nil {
		return err
	}

	d, err := os.Create(dest)
	if err != nil {
		return err
	}

	w := zip.NewWriter(d)
	d.Close()
	defer w.Close()
	for _, file := range files {
		if err = compress(file, "", w); err != nil {
			return err
		}
	}

	return nil
}

func compress(file *os.File, prefix string, w *zip.Writer) error {
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.IsDir() {
		prefix += "/" + stat.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, info := range fileInfos {
			f, err := os.Open(file.Name() + "/" + info.Name())
			if err != nil {
				return err
			}
			if err = compress(f, prefix, w); err != nil {
				return err
			}
		}
		return nil
	}
	header, err := zip.FileInfoHeader(stat)
	if err != nil {
		return err
	}
	header.Name = prefix + "/" + header.Name
	writer, err := w.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}
	return nil
}
