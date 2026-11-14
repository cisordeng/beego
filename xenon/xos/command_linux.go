// +build linux

package xos

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
)

type Command struct {
	*exec.Cmd
}

func NewCommand(name string, arg ...string) *Command {
	return &Command{exec.Command(name, arg...)}
}

func NewCommandContext(ctx context.Context, name string, arg ...string) *Command {
	return &Command{exec.CommandContext(ctx, name, arg...)}
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
