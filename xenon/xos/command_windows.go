// +build windows

package xos

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"syscall"
)

type Command struct {
	*exec.Cmd
}

func NewCommand(name string, arg ...string) *Command {
	cmd := &Command{exec.Command(name, arg...)}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
		CreationFlags: 0x08000000,
	}
	return cmd
}

func NewCommandContext(ctx context.Context, name string, arg ...string) *Command {
	cmd := &Command{exec.CommandContext(ctx, name, arg...)}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
		CreationFlags: 0x08000000,
	}
	return cmd
}

func (this *Command) GetRunOutput() string {
	var outBytes bytes.Buffer
	var errBytes bytes.Buffer
	this.Stdout = &outBytes
	this.Stderr = &errBytes
	err := this.Run()
	if err != nil || errBytes.Len() > 0 {
		return ""
	}
	return strings.Trim(strings.TrimSpace(outBytes.String()), "'")
}
