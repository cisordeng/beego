// +build linux

package xos

import (
	"strconv"
)

type SystemProxyConfig struct{
	Mode string
	HttpHost string
	HttpPort int
	HttpsHost string
	HttpsPort int
	FTPHost string
	FTPPort int
	SocksHost string
	SocksPort int
	AutoconfigUrl string
}

func SetSystemProxy(httpPort int, socksPort int) error {
	commands := [][]string{
		{"gsettings", "set", "org.gnome.system.proxy.http", "host", "127.0.0.1"},
		{"gsettings", "set", "org.gnome.system.proxy.http", "port", strconv.Itoa(httpPort)},
		{"gsettings", "set", "org.gnome.system.proxy.https", "host", "127.0.0.1"},
		{"gsettings", "set", "org.gnome.system.proxy.https", "port", strconv.Itoa(httpPort)},
		{"gsettings", "set", "org.gnome.system.proxy.socks", "host", "127.0.0.1"},
		{"gsettings", "set", "org.gnome.system.proxy.socks", "port", strconv.Itoa(socksPort)},
		{"gsettings", "set", "org.gnome.system.proxy", "mode", "manual"},
	}
	for _, command := range commands {
		if err := NewCommand(command[0], command[1:]...).Run(); err != nil {
			return err
		}
	}
	return nil
}

func ClearSystemProxy() error {
	return NewCommand("gsettings", "set", "org.gnome.system.proxy", "mode", "none").Run()
}

func GetSystemProxy() *SystemProxyConfig {
	var config SystemProxyConfig
	commands := [][]string{
		{"gsettings", "get", "org.gnome.system.proxy", "mode"},
		{"gsettings", "get", "org.gnome.system.proxy.http", "host"},
		{"gsettings", "get", "org.gnome.system.proxy.http", "port"},
		{"gsettings", "get", "org.gnome.system.proxy.https", "host"},
		{"gsettings", "get", "org.gnome.system.proxy.https", "port"},
		{"gsettings", "get", "org.gnome.system.proxy.socks", "host"},
		{"gsettings", "get", "org.gnome.system.proxy.socks", "port"},
		{"gsettings", "get", "org.gnome.system.proxy.ftp", "host"},
		{"gsettings", "get", "org.gnome.system.proxy.ftp", "port"},
		{"gsettings", "get", "org.gnome.system.proxy", "autoconfig-url"},
	}

	config.Mode, _ = NewCommand(commands[0][0], commands[0][1:]...).GetRunOutput()
	config.HttpHost, _ = NewCommand(commands[1][0], commands[1][1:]...).GetRunOutput()
	data, _ := NewCommand(commands[2][0], commands[2][1:]...).GetRunOutput()
	config.HttpPort, _ = strconv.Atoi(data)
	config.HttpsHost, _ = NewCommand(commands[3][0], commands[3][1:]...).GetRunOutput()
	data, _ = NewCommand(commands[4][0], commands[4][1:]...).GetRunOutput()
	config.HttpsPort, _ = strconv.Atoi(data)
	config.FTPHost, _ = NewCommand(commands[5][0], commands[5][1:]...).GetRunOutput()
	data, _ = NewCommand(commands[6][0], commands[6][1:]...).GetRunOutput()
	config.FTPPort, _ = strconv.Atoi(data)
	config.SocksHost, _ = NewCommand(commands[7][0], commands[7][1:]...).GetRunOutput()
	data, _ = NewCommand(commands[8][0], commands[8][1:]...).GetRunOutput()
	config.SocksPort, _ = strconv.Atoi(data)
	config.AutoconfigUrl, _ = NewCommand(commands[9][0], commands[9][1:]...).GetRunOutput()
	return &config
}