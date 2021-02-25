package ptools

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/axgle/mahonia"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

//自定义Scanner分割的方式，\n和\r都分割
func ScanCRandLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	//优先分割换行\n
	if i := bytes.IndexAny(data, "\n"); i >= 0 {
		return i + 1, data[0:i], nil
	}

	//然后分割行首\r
	if i := bytes.IndexAny(data, "\r"); i >= 0 {
		return i + 1, data[0 : len(data)-1], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

//执行一次command指令 跨平台兼容
func Exec(command string) (output string, err error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/c", command)
	} else {
		cmd = exec.Command("/bin/bash", "-c", command)
	}
	//隐藏黑框
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	out, err := cmd.CombinedOutput()
	return string(out), err
}

//执行一次command指令且自定义方法处理每行结果 跨平台兼容
func ExecRealtime(command string, method func(line string)) error {
	//跨平台兼容，cmd/bash传参是为了使用二者自带的命令，直接exec无法使用这些命令
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/c", command)
	} else {
		cmd = exec.Command("/bin/bash", "-c", command)
	}

	//隐藏黑框
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

	//实时处理输出 TODO 可能有并发隐患导致内存溢出
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

//执行一次command指令且实时输出每行结果 跨平台兼容
func ExecRealtimePrint(command string) error {
	return ExecRealtime(command, func(line string) {
		fmt.Println(line)
	})
}

//执行时解决cmd chcp936的中文乱码问题
func ExecRealtimePrintGBK(command string) error {
	return ExecRealtime(command, func(line string) {
		fmt.Println(ConvertString(line))
	})
}

//转换编码解决乱码问题 字符串
func ConvertString(s string) string {
	return mahonia.NewDecoder("GBK").ConvertString(s)
}

//查找（环境变量+当前位置）可执行文件的位置 跨平台兼容
func GetBinaryPath(binary string) (string, error) {
	var command string
	if runtime.GOOS == "windows" {
		command = "where " + binary
	} else {
		command = "which " + binary
	}

	dir, err := Exec(command)
	dir = strings.TrimSpace(dir)
	return dir, err
}


//windows要用winPssuspend.exe 需指定其路径
//其他系统留空
func ExecRealtimePause(command string, method func(line string), a chan rune, winPssuspend string) error {
	//跨平台兼容，cmd/bash传参是为了使用二者自带的命令，直接exec无法使用这些命令
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/c", command)
	} else {
		cmd = exec.Command("/bin/bash", "-c", command)
	}

	//隐藏黑框
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

	//实时处理输出 TODO 可能有并发隐患导致内存溢出
	go func() {
		scanner.Split(ScanCRandLF)
		for scanner.Scan() {
			//对每一行的操作
			method(string(scanner.Bytes()))
		}
		fmt.Println("debug...结束")
	}()

	go func() {
		for {
			switch <-a {
			case 'p':
				//暂停
				if runtime.GOOS == "windows" {
					fmt.Println(FormatPath(winPssuspend) + " " + strconv.Itoa(cmd.Process.Pid))
					if out, err := Exec(FormatPath(winPssuspend) + " " + strconv.Itoa(cmd.Process.Pid)); err != nil {
						log.Println(out)
						log.Println(err)
					}
				} else {
					cmd.Process.Signal(syscall.SIGTSTP)	//win下不可用
				}
				a <- ' '
			case 'r':
				//继续
				if runtime.GOOS == "windows" {
					fmt.Println(FormatPath(winPssuspend) + " -r " + strconv.Itoa(cmd.Process.Pid))
					if out, err := Exec(FormatPath(winPssuspend) + " -r " + strconv.Itoa(cmd.Process.Pid)); err != nil {
						log.Println(out)
						log.Println(err)
					}
				} else {
					cmd.Process.Signal(syscall.SIGCONT)	//win下不可用
				}
				a <- ' '
			case 'q':
				//中止
				if err = cmd.Process.Kill(); err != nil {
					log.Println(err)
				}
				break
			}

			time.Sleep(time.Second * 1)
		}
	}()

	return cmd.Wait()
}
