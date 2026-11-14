// +build linux

package xos

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

func SetAutoStart(desktopFileName string, desktopBytes ...[]byte) error {
	if !strings.HasSuffix(desktopFileName, ".desktop") {
		desktopFileName = fmt.Sprintf("%s.desktop", desktopFileName)
	}
	usr, err := user.Current()
	if err != nil {
		return err
	}
	autoStartDir := fmt.Sprintf("%s/.config/autostart", usr.HomeDir)
	autoStartDesktopPath := fmt.Sprintf("%s/%s", autoStartDir, desktopFileName)
	err = ioutil.WriteFile(autoStartDesktopPath, desktopBytes[0], 0644)
	return err
}

func ClearAutoStart(desktopFileName string) error {
	if !strings.HasSuffix(desktopFileName, ".desktop") {
		desktopFileName = fmt.Sprintf("%s.desktop", desktopFileName)
	}
	usr, err := user.Current()
	if err != nil {
		return err
	}
	autoStartDir := fmt.Sprintf("%s/.config/autostart", usr.HomeDir)
	autoStartDesktopPath := fmt.Sprintf("%s/%s", autoStartDir, desktopFileName)
	err = os.Remove(autoStartDesktopPath)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}