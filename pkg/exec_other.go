// +build !windows

package ptools

import (
	"os/exec"
	"strings"
	"syscall"
)

//执行一次command指令 经过cmd
func CMDArgs(args []string) (output string, err error) {
	return ExecArgs(append([]string{"/bin/bash", "-c"}, args...))
}
func CMDRealtimeArgs(args []string, method func(line string)) error {
	return ExecRealtimeArgs(append([]string{"/bin/bash", "-c"}, args...), method)
}

//查找（环境变量+当前位置）可执行文件的位置
func GetBinaryPath(binary string) (string, error) {
	dir, err := CMD("which " + binary)
	dir = strings.TrimSpace(dir)
	return dir, err
}

func realtimeControl(cmd *exec.Cmd, signal chan<- rune) (err error) {
	for control := range signal {
		switch control {
		case 'p':
			//暂停
			_ = cmd.Process.Signal(syscall.SIGTSTP) //win下不可用
		case 'r':
			//继续
			_ = cmd.Process.Signal(syscall.SIGCONT) //win下不可用
		case 'q':
			//中止
			err = cmd.Process.Kill()
			break
		}
	}

	return
}

func doHideWindow(cmd *exec.Cmd) {
}
