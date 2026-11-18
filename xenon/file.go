package xenon

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveToFile(file multipart.File, tofile string) error {
	defer file.Seek(0, 0)
	f, err := os.OpenFile(tofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	defer f.Seek(0, 0)
	_, err = io.Copy(f, file)
	return err
}

func AbsolutePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	exePath, err := os.Executable()
	PanicNotNilError(err)
	exeDir := filepath.Dir(exePath)

	return filepath.Join(exeDir, filepath.FromSlash(path))
}

func ExeFilePath() string {
	exePath, err := os.Executable()
	PanicNotNilError(err)
	return exePath
}

func ExeDirPath() string {
	exePath, err := os.Executable()
	PanicNotNilError(err)
	return filepath.Dir(exePath)
}