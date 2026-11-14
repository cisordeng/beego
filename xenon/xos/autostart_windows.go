// +build windows

package xos

import (
	"fmt"
	"os"
)

func SetAutoStart(exeName string, bytes ...bytes) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	return NewCommand("reg", "add", regPath, "/v", exeName, "/t", "REG_SZ", "/d", fmt.Sprintf(`"%s"`, exePath), "/f").Run()
}

func ClearAutoStart(exeName string) error {
	_ = NewCommand("reg", "delete", regPath, "/v", exeName, "/f").Run()
	return nil
}