// +build linux

package xos

import (
	"bytes"
	"context"
	"errors"
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
		Pdeathsig: syscall.SIGKILL,
	}
	return cmd
}

func NewCommandContext(ctx context.Context, name string, arg ...string) *Command {
	cmd := &Command{exec.CommandContext(ctx, name, arg...)}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGKILL,
	}
	return cmd
}

func (this *Command) GetRunOutput() (string, error) {
	var outBytes bytes.Buffer
	var errBytes bytes.Buffer
	this.Stdout = &outBytes
	this.Stderr = &errBytes
	err := this.Run()
	if errBytes.Len() > 0 {
		return "", errors.New(errBytes.String())
	}
	if err != nil {
		return "", err
	}
	return strings.Trim(strings.TrimSpace(outBytes.String()), "'"), nil
}
