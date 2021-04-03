// +build !windows

package ptools

import (
	"bufio"
	"os/exec"
	"strings"
	"syscall"
)

//执行一次command指令
func CMD(command string) (output string, err error) {
	cmd := exec.Command("/bin/bash", "-c", command)

	out, err := cmd.CombinedOutput()
	return string(out), err
}

//参数以切片形式存放
func ExecArgs(args []string) (output string, err error) {
	cmd := exec.Command(args[0], args[1:]...)

	out, err := cmd.CombinedOutput()
	return string(out), err
}

//执行一次command指令且自定义方法处理每行结果
func ExecRealtime(command string, method func(line string)) error {
	//cmd/bash传参是为了使用二者自带的命令，直接exec无法使用这些命令
	cmd := exec.Command("/bin/bash", "-c", command)

	//标准输出pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	//标准错误pipe->标准输出pipe
	cmd.Stderr = cmd.Stdout

	//开始运行
	if err = cmd.Start(); err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)

	//实时处理输出
	go func() {
		scanner.Split(ScanCRandLF)
		for scanner.Scan() {
			//对每一行的操作
			method(string(scanner.Bytes()))
		}
	}()

	return cmd.Wait()
}

//查找（环境变量+当前位置）可执行文件的位置
func GetBinaryPath(binary string) (string, error) {
	dir, err := CMD("which " + binary)
	dir = strings.TrimSpace(dir)
	return dir, err
}

//执行指令 实时控制 winPssuspend留空
func ExecRealtimeControlArgs(args []string, method func(line string), signal chan rune, winPssuspend string) error {
	//cmd/bash传参是为了使用二者自带的命令，直接exec无法使用这些命令
	args = append([]string{"-c"}, args...)
	cmd := exec.Command("/bin/bash", args...)

	//标准输出pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	//标准错误pipe->标准输出pipe
	cmd.Stderr = cmd.Stdout

	//开始运行
	if err = cmd.Start(); err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)

	//实时处理输出
	go func() {
		scanner.Split(ScanCRandLF)
		for scanner.Scan() {
			//对每一行的操作
			method(string(scanner.Bytes()))
		}
	}()

	go func() {
		for control := range signal{
			switch control {
			case 'p':
				//暂停
				cmd.Process.Signal(syscall.SIGTSTP)	//win下不可用
			case 'r':
				//继续
				cmd.Process.Signal(syscall.SIGCONT)	//win下不可用
			case 'q':
				//中止
				err = cmd.Process.Kill()
				break
			}
		}
	}()

	return cmd.Wait()
}