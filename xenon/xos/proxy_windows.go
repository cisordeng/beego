// +build windows

package xos

import (
	"fmt"
)

func SetSystemProxy(httpPort int, socksPort int) error {
	proxyAddr := fmt.Sprintf("127.0.0.1:%d", httpPort)
	regPath := `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`

	bypass := "localhost;127.*;10.*;172.16.*;172.17.*;172.18.*;172.19.*;172.20.*;172.21.*;172.22.*;172.23.*;172.24.*;172.25.*;172.26.*;172.27.*;172.28.*;172.29.*;172.30.*;172.31.*;192.168.*"

	commands := [][]string{
		{"reg", "add", regPath, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f"},
		{"reg", "add", regPath, "/v", "ProxyServer", "/t", "REG_SZ", "/d", proxyAddr, "/f"},
		{"reg", "add", regPath, "/v", "ProxyOverride", "/t", "REG_SZ", "/d", bypass, "/f"},
		{"reg", "add", regPath, "/v", "AutoConfigURL", "/t", "REG_SZ", "/d", "", "/f"},
	}

	for _, command := range commands {
		err := NewCommand(command[0], command[1:]...).Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func ClearSystemProxy() error {
	regPath := `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`

	commands := [][]string{
		{"reg", "add", regPath, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f"},
		{"reg", "add", regPath, "/v", "ProxyServer", "/t", "REG_SZ", "/d", "", "/f"},
		{"reg", "add", regPath, "/v", "AutoConfigURL", "/t", "REG_SZ", "/d", "", "/f"},
	}

	for _, command := range commands {
		err := NewCommand(command[0], command[1:]...).Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func GetSystemProxyWindows() *SystemProxyConfig {
	return nil
}
