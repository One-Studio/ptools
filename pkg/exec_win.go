// +build windows

package ptools

import (
	"os/exec"
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

// //参数以切片形式存放
// func ExecArgs(args []string) (output string, err error) {
// 	cmd := exec.Command(args[0], args[1:]...)

// 	out, err := cmd.CombinedOutput()
// 	return string(out), err
// }

// //执行一次command指令且自定义方法处理每行结果
// func ExecRealtime(command string, method func(line string)) error {
// 	//cmd/bash传参是为了使用二者自带的命令，直接exec无法使用这些命令
// 	cmd := exec.Command("cmd.exe", "/c", command)

// 	//隐藏黑框 !仅win下用
// 	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

// 	//标准输出pipe
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		return err
// 	}

// 	//标准错误pipe->标准输出pipe
// 	cmd.Stderr = cmd.Stdout

// 	//开始运行
// 	if err = cmd.Start(); err != nil {
// 		return err
// 	}
// 	scanner := bufio.NewScanner(stdout)

// 	//实时处理输出
// 	go func() {
// 		scanner.Split(ScanCRandLF)
// 		for scanner.Scan() {
// 			//对每一行的操作
// 			method(string(scanner.Bytes()))
// 		}
// 		fmt.Println("debug...结束")
// 	}()

// 	return cmd.Wait()
// }

//查找（环境变量+当前位置）可执行文件的位置
func GetBinaryPath(binary string) (string, error) {
	dir, err := CMD("where " + binary)
	dir = strings.TrimSpace(dir)
	return dir, err
}

// //参数以切片形式存放
// func ExecRealtimeControlArgs(args []string, method func(line string), signal chan rune, winPssuspend string) error 	{
// 	if exist := IsFileExisted(winPssuspend); !exist {
// 		return errors.New("pssuspend.exe does not exist. check path string")
// 	}

// 	cmd := exec.Command(args[0], args[1:]...)

// 	//隐藏黑框 !仅win下用
// 	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

// 	//标准输出pipe
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		return err
// 	}

// 	//标准错误pipe->标准输出pipe
// 	cmd.Stderr = cmd.Stdout

// 	//开始运行
// 	if err = cmd.Start(); err != nil {
// 		return err
// 	}
// 	scanner := bufio.NewScanner(stdout)

// 	//实时处理输出
// 	go func() {
// 		scanner.Split(ScanCRandLF)
// 		for scanner.Scan() {
// 			//对每一行的操作
// 			method(string(scanner.Bytes()))
// 		}
// 	}()

// 	go func() {
// 		for control := range signal{
// 			switch control {
// 			case 'p':
// 				//暂停
// 				fmt.Println(FormatPath(winPssuspend) + " " + strconv.Itoa(cmd.Process.Pid))
// 				if _, err := CMD(FormatPath(winPssuspend) + " " + strconv.Itoa(cmd.Process.Pid)); err != nil {
// 					//log.Println(out)
// 					log.Println(err)
// 				}
// 			case 'r':
// 				//继续
// 				fmt.Println(FormatPath(winPssuspend) + " -r " + strconv.Itoa(cmd.Process.Pid))
// 				if _, err := CMD(FormatPath(winPssuspend) + " -r " + strconv.Itoa(cmd.Process.Pid)); err != nil {
// 					//log.Println(out)
// 					log.Println(err)
// 				}
// 			case 'q':
// 				//中止
// 				_ = cmd.Process.Kill()
// 			}
// 		}
// 	}()

// 	return cmd.Wait()
// }

func doHideWindow(cmd *exec.Cmd) {
	//隐藏黑框 !仅win下用
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

// 	go func() {

// 	}()

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
