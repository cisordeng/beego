package xenon

import (
	"io"
	"mime/multipart"
	"os"
)

func SaveToFile(file multipart.File, tofile string) error {
	defer file.Close()
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