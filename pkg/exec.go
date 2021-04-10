package ptools

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

//// 1

//参数以切片形式存放
func ExecArgs(args []string) (output string, err error) {
	cmd := exec.Command(args[0], args[1:]...)

	//隐藏黑框，仅win下需要
	doHideWindow(cmd)

	out, err := cmd.CombinedOutput()
	return string(out), err
}

//执行一次command指令 直接调用
func Exec(command string) (output string, err error) {
	return ExecArgs(strings.Fields(command))
}

//// 2

//调用CMD或者bash执行指令，适用于终端指令
func CMD(command string) (output string, err error) {
	return CMDArgs(strings.Fields(command))
}

//// 3

//执行一次command指令且自定义方法处理每行结果
func ExecRealtimeArgs(args []string, method func(line string)) error {
	//cmd/bash传参是为了使用二者自带的命令，直接exec无法使用这些命令
	cmd := exec.Command(args[0], args[1:]...)

	//隐藏黑框，仅win下需要
	doHideWindow(cmd)

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

func ExecRealtime(command string, method func(line string)) error {
	return ExecRealtimeArgs(strings.Fields(command), method)
}

func CMDRealtime(command string, method func(line string)) error {
	return CMDRealtimeArgs(strings.Fields(command), method)
}

//// TODO:
// @ 指令[]string的方法带args后缀 指令string的在分割指令之后调用args
// 1. Exec 执行指令
// 2. CMD 根据系统执行命令，调用Exec
//    windows "cmd.exe", "/c", args...
//    other   "/bin/bash", "-c", args...
// 3. doRealtime 执行指令，实时对每行字符串进行操作
//    - ExecRealtime
//    - CMDRealtime
// 4. doRealtimeControl 执行指令，实时对每行字符串进行操作，且实时暂停/继续/结束
//    - ExecRealtimeControl
//    - CMDRealtimeControl

// 操作一：分割行
// go func() {
// 	scanner.Split(ScanCRandLF)
// 	for scanner.Scan() {
// 		对每一行的操作
// 		method(string(scanner.Bytes()))
// 	}
// }()
// 操作二：实时暂停/继续/结束 TODO: 解决实时控制的问题
// 根据 win or others
// realtimeControl

////TODO: END

//// 4

//执行指令 实时控制 winPssuspend留空
func ExecRealtimeControlArgs(args []string, method func(line string), signal chan rune, winPssuspend string) error {
	cmd := exec.Command(args[0], args[1:]...)

	//隐藏黑框，仅win下需要
	doHideWindow(cmd)

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

	//TODO: 实时控制
	go func() {
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
	}()
	// go func() {
	// 	for control := range signal {
	// 		switch control {
	// 		case 'p':
	// 			//暂停
	// 			cmd.Process.Signal(syscall.SIGTSTP) //win下不可用
	// 		case 'r':
	// 			//继续
	// 			cmd.Process.Signal(syscall.SIGCONT) //win下不可用
	// 		case 'q':
	// 			//中止
	// 			err = cmd.Process.Kill()
	// 			break
	// 		}
	// 	}
	// }()

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

//windows要用winPssuspend.exe 需指定其路径 其他系统留空
func ExecRealtimeControl(command string, method func(line string), signal chan rune, winPssuspend string) error {
	cmdArgs := strings.Fields(command)
	return ExecRealtimeControlArgs(cmdArgs, method, signal, winPssuspend)
}

//实时控制的时候暂停
func Pause(a chan rune) {
	a <- 'p'
}

//实时控制的时候继续
func Resume(a chan rune) {
	a <- 'r'
}

//实时控制的时候结束
func Quit(a chan rune) {
	a <- 'q'
}
