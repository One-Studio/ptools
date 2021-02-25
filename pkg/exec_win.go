// +build windows

package ptools

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

//执行一次command指令
func Exec(command string) (output string, err error) {
	cmd := exec.Command("cmd.exe", "/c", command)

	//隐藏黑框 !仅win下用
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	out, err := cmd.CombinedOutput()
	return string(out), err
}

//执行一次command指令且自定义方法处理每行结果
func ExecRealtime(command string, method func(line string)) error {
	//cmd/bash传参是为了使用二者自带的命令，直接exec无法使用这些命令
	cmd := exec.Command("cmd.exe", "/c", command)

	//隐藏黑框 !仅win下用
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

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
		fmt.Println("debug...结束")
	}()

	return cmd.Wait()
}

//执行一次command指令且实时输出每行结果
func ExecRealtimePrint(command string) error {
	return ExecRealtime(command, func(line string) {
		fmt.Println(line)
	})
}

//执行时实时输出每行并解决cmd chcp 936 输出乱码问题
func ExecRealtimePrintGBK(command string) error {
	return ExecRealtime(command, func(line string) {
		fmt.Println(ConvertString(line))
	})
}

//查找（环境变量+当前位置）可执行文件的位置
func GetBinaryPath(binary string) (string, error) {
	dir, err := Exec("where " + binary)
	dir = strings.TrimSpace(dir)
	return dir, err
}

//windows要用winPssuspend.exe 需指定其路径
func ExecRealtimeControl(command string, method func(line string), signal chan rune, winPssuspend string) error {
	if exist := IsFileExisted(winPssuspend); !exist {
		return errors.New("pssuspend.exe does not exist. check path string")
	}

	//cmd := exec.Command("cmd.exe", "/c", command)
	//实时控制要直接执行程序，不然获取的是cmd.exe，没法挂起
	cmdArgs := strings.Fields(command)
	fmt.Println(cmdArgs)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	//隐藏黑框 !仅win下用
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

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
				fmt.Println(FormatPath(winPssuspend) + " " + strconv.Itoa(cmd.Process.Pid))
				if _, err := Exec(FormatPath(winPssuspend) + " " + strconv.Itoa(cmd.Process.Pid)); err != nil {
					//log.Println(out)
					log.Println(err)
				}
			case 'r':
				//继续
				fmt.Println(FormatPath(winPssuspend) + " -r " + strconv.Itoa(cmd.Process.Pid))
				if _, err := Exec(FormatPath(winPssuspend) + " -r " + strconv.Itoa(cmd.Process.Pid)); err != nil {
					//log.Println(out)
					log.Println(err)
				}
			case 'q':
				//中止
				_ = cmd.Process.Kill()
			}
		}
	}()

	return cmd.Wait()
}

//实时控制的时候暂停
func ExecPause(a chan rune)  {
	a <- 'p'
}

//实时控制的时候继续
func ExecResume(a chan rune)  {
	a <- 'r'
}

//实时控制的时候结束
func ExecQuit(a chan rune)  {
	a <- 'q'
}