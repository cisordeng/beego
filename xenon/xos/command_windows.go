// +build windows

package xos

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

// load kernel32.dll
var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	procCreateJobObjectW = kernel32.NewProc("CreateJobObjectW")
	procSetInformationJobObject = kernel32.NewProc("SetInformationJobObject")
	procOpenProcess = kernel32.NewProc("OpenProcess")
	procAssignProcessToJobObject = kernel32.NewProc("AssignProcessToJobObject")
)

// windows define vars
const (
	JobObjectExtendedLimitInformation = 9
	JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE = 0x2000
	PROCESS_SET_QUOTA                  = 0x0100
	PROCESS_TERMINATE                  = 0x0001
)

// windows define struct
type JOBOBJECT_BASIC_LIMIT_INFORMATION struct {
	PerProcessUserTimeLimit int64
	PerJobUserTimeLimit int64
	LimitFlags uint32
	MinimumWorkingSetSize uintptr
	MaximumWorkingSetSize uintptr
	ActiveProcessLimit uint32
	Affinity uintptr
	PriorityClass uint32
	SchedulingClass uint32
}

type IO_COUNTERS struct {
	ReadOperationCount uint64
	WriteOperationCount uint64
	OtherOperationCount uint64
	ReadTransferCount uint64
	WriteTransferCount uint64
	OtherTransferCount uint64
}

type JOBOBJECT_EXTENDED_LIMIT_INFORMATION struct {
	BasicLimitInformation JOBOBJECT_BASIC_LIMIT_INFORMATION
	IoInfo IO_COUNTERS
	ProcessMemoryLimit uintptr
	JobMemoryLimit uintptr
	PeakProcessMemoryUsed uintptr
	PeakJobMemoryUsed uintptr
}

type Command struct {
	*exec.Cmd
	windowsJob uintptr
}

func NewCommand(name string, arg ...string) *Command {
	cmd := &Command{
		Cmd: exec.Command(name, arg...),
		windowsJob: 0,
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
		CreationFlags: 0x08000000,
	}

	// create windows job
	cmd.windowsJob, _, _ = procCreateJobObjectW.Call(0, 0)
	if cmd.windowsJob == 0 {
		return cmd
	}
	// set cmd parent exit to kill cmd process
	info := JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	info.BasicLimitInformation.LimitFlags = JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE

	ret, _, _ := procSetInformationJobObject.Call(
		cmd.windowsJob,
		uintptr(JobObjectExtendedLimitInformation),
		uintptr(unsafe.Pointer(&info)),
		unsafe.Sizeof(info),
	)
	if ret == 0 {
		syscall.CloseHandle(syscall.Handle(cmd.windowsJob))
		return cmd
	}

	return cmd
}

func NewCommandContext(ctx context.Context, name string, arg ...string) *Command {
	cmd := &Command{
		Cmd: exec.CommandContext(ctx, name, arg...),
		windowsJob: 0,
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
		CreationFlags: 0x08000000,
	}

	// create windows job
	cmd.windowsJob, _, _ = procCreateJobObjectW.Call(0, 0)
	if cmd.windowsJob == 0 {
		return cmd
	}
	// set cmd parent exit to kill cmd process
	info := JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	info.BasicLimitInformation.LimitFlags = JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE

	ret, _, _ := procSetInformationJobObject.Call(
		cmd.windowsJob,
		uintptr(JobObjectExtendedLimitInformation),
		uintptr(unsafe.Pointer(&info)),
		unsafe.Sizeof(info),
	)
	if ret == 0 {
		syscall.CloseHandle(syscall.Handle(cmd.windowsJob))
		return cmd
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

func (this *Command) Start() error {
	err := this.Cmd.Start()
	if this.windowsJob == 0 {
		return err
	}
	if err != nil {
		syscall.CloseHandle(syscall.Handle(this.windowsJob))
		return err
	}

	// get cmd process
	pid := this.Process.Pid
	process, _, err := procOpenProcess.Call(
		uintptr(PROCESS_SET_QUOTA|PROCESS_TERMINATE),
		0,
		uintptr(pid),
	)
	if process == 0 {
		syscall.CloseHandle(syscall.Handle(this.windowsJob))
		_ = this.Process.Kill()
		return err
	}
	// close temp process
	defer syscall.CloseHandle(syscall.Handle(process))

	// cmd process to job
	ret, _, err := procAssignProcessToJobObject.Call(this.windowsJob, process)
	if ret == 0 {
		syscall.CloseHandle(syscall.Handle(this.windowsJob))
		_ = this.Process.Kill()
		return err
	}
	return nil
}