// +build !windows

package ptools

import (
	"log"
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

func CMDRealtimeControlArgs(args []string, method func(line string), signal chan rune, winPssuspend string) error {
	return CMDRealtimeControlArgs(append([]string{"/bin/bash", "-c"}), method, signal, winPssuspend)
}

func realtimeControl(cmd *exec.Cmd, signal chan rune, winPssuspend string) {
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
			err := cmd.Process.Kill()
			if err != nil {
				log.Println(err)
			}

			break
		}
	}
}

//查找（环境变量+当前位置）可执行文件的位置
func GetBinaryPath(binary string) (string, error) {
	dir, err := CMD("which " + binary)
	dir = strings.TrimSpace(dir)
	return dir, err
}


func doHideWindow(cmd *exec.Cmd) {
}
