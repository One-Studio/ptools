// +build windows

package ptools

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

//执行一次command指令 经过cmd
func CMDArgs(args []string) (output string, err error) {
	return ExecArgs(append([]string{"cmd.exe", "/c"}, args...))
}

func CMDRealtimeArgs(args []string, method func(line string)) error {
	return ExecRealtimeArgs(append([]string{"cmd.exe", "/c"}, args...), method)
}

func CMDRealtimeControlArgs(args []string, method func(line string), signal chan rune, winPssuspend string) error {
	return ExecRealtimeControlArgs(append([]string{"cmd.exe", "/c"}, args...), method, signal, winPssuspend)
}

func realtimeControl(cmd *exec.Cmd, signal chan rune, winPssuspend string) {
	for control := range signal {
		switch control {
		case 'p':
			//暂停
			fmt.Println(FormatPath(winPssuspend) + " " + strconv.Itoa(cmd.Process.Pid))
			if _, err := CMD(FormatPath(winPssuspend) + " " + strconv.Itoa(cmd.Process.Pid)); err != nil {
				//log.Println(out)
				log.Println(err)
			}
		case 'r':
			//继续
			fmt.Println(FormatPath(winPssuspend) + " -r " + strconv.Itoa(cmd.Process.Pid))
			if _, err := CMD(FormatPath(winPssuspend) + " -r " + strconv.Itoa(cmd.Process.Pid)); err != nil {
				//log.Println(out)
				log.Println(err)
			}
		case 'q':
			//中止
			_ = cmd.Process.Kill()
		}
	}
}

//查找（环境变量+当前位置）可执行文件的位置
func GetBinaryPath(binary string) (string, error) {
	dir, err := CMD("where " + binary)
	dir = strings.TrimSpace(dir)
	return dir, err
}

func doHideWindow(cmd *exec.Cmd) {
	//隐藏黑框 !仅win下用
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

// func realtimeControl(cmd *exec.Cmd, signal chan rune) (err error) {
// 	for control := range signal {
// 		switch control {
// 		case 'p':
// 			//暂停
// 			fmt.Println(FormatPath(winPssuspend) + " " + strconv.Itoa(cmd.Process.Pid))
// 			if _, err := CMD(FormatPath(winPssuspend) + " " + strconv.Itoa(cmd.Process.Pid)); err != nil {
// 				//log.Println(out)
// 				log.Println(err)
// 			}
// 		case 'r':
// 			//继续
// 			fmt.Println(FormatPath(winPssuspend) + " -r " + strconv.Itoa(cmd.Process.Pid))
// 			if _, err := CMD(FormatPath(winPssuspend) + " -r " + strconv.Itoa(cmd.Process.Pid)); err != nil {
// 				//log.Println(out)
// 				log.Println(err)
// 			}
// 		case 'q':
// 			//中止
// 			_ = cmd.Process.Kill()
// 		}
// 	}

// 	return
// }
