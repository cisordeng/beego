package xenon

import (
	"bytes"
	"os/exec"
	"strings"
)

func RunOSCommand(cmd *exec.Cmd) string {
	var outBytes bytes.Buffer
	var errBytes bytes.Buffer
	cmd.Stdout = &outBytes
	cmd.Stderr = &errBytes
	err := cmd.Run()
	if err != nil || errBytes.Len() > 0 {
		return ""
	}
	return strings.Trim(strings.TrimSpace(outBytes.String()), "'")
}
