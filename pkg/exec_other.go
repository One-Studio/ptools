// +build !windows

package ptools

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

//执行一次command指令 跨平台兼容
func CMD(command string) (output string, err error) {
	cmd := exec.Command("/bin/bash", "-c", command)

	out, err := cmd.CombinedOutput()
	return string(out), err
}

//执行一次command指令 直接调用
func Exec(command string) (output string, err error) {
	cmdArgs := strings.Fields(command)
	//fmt.Println(cmdArgs)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

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

//执行一次command指令且实时输出每行结果
func ExecRealtimePrint(command string) error {
	return ExecRealtime(command, func(line string) {
		fmt.Println(line)
	})
}

//执行时解决cmd chcp 936 乱码问题 P.S. !win好像不用
func ExecRealtimePrintGBK(command string) error {
	return ExecRealtime(command, func(line string) {
		fmt.Println(ConvertString(line))
	})
}

//查找（环境变量+当前位置）可执行文件的位置
func GetBinaryPath(binary string) (string, error) {
	dir, err := CMD("which " + binary)
	dir = strings.TrimSpace(dir)
	return dir, err
}

//winPssuspend留空
func ExecRealtimeControl(command string, method func(line string), signal chan rune, winPssuspend string) error {
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

//实时控制的时候暂停
func Pause(a chan rune)  {
	a <- 'p'
}

//实时控制的时候继续
func Resume(a chan rune)  {
	a <- 'r'
}

//实时控制的时候结束
func Quit(a chan rune)  {
	a <- 'q'
}