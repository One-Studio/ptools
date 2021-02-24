package ptools

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/axgle/mahonia"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"
)

//自定义Scanner分割的方式，\n和\r都分割
func ScanCRorLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	//优先分割换行\n
	if i := bytes.IndexAny(data, "\n"); i >= 0 {
		return i + 1, data[0:i], nil
	}

	//然后分割行首\r
	if i := bytes.IndexAny(data, "\r"); i >= 0 {
		return i + 1, data[0:len(data)-1], nil
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
		scanner.Split(ScanCRorLF)
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

//TODO 暂时只能linux macos使用 windows要用pstools.exe
func ExecRealtimePause(command string, method func(line string)) {
	cmdArgs := strings.Fields(command)
	fmt.Println(cmdArgs)
	//cmd := exec.Command(command)//"/bin/bash", "-c",
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	stderr, _ := cmd.StderrPipe()
	defer stderr.Close()
	in := bytes.Buffer{}
	cmd.Stdin = &in
	//stdout, _ := cmd.StdoutPipe()
	//defer stdout.Close()
	//stdout, _ := cmd.StdoutPipe()
	_ = cmd.Start()
	//原方法
	scanner := bufio.NewScanner(stderr)
	//scanner1 := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanRunes)
	//sc := bufio.NewWriter(stdin)
	//go func() {
	//	for scanner1.Scan() {
	//		fmt.Println(scanner1.Text())
	//	}
	//}()

	//不行，只能暂停scanner，实际还在运行
	//var w = make(chan bool, 1)
	//ww := false
	go func() {
		//TODO 解决windows下没法暂停的问题
		//     思路是，inputpipe中传入按键PauseBreak和Enter键
		//     暂用外部工具pssuspend.exe代替

		fmt.Println("go func暂停2s")
		time.Sleep(time.Second * 2)
		fmt.Println("go func暂停结束")
		//fmt.Println("go func暂停1s")
		//time.Sleep(time.Second * 1)
		//fmt.Println("go func暂停结束")
		//cmd.Process.Signal(syscall.SIGALRM)
		//
		//fmt.Println("go func暂停1s")
		//time.Sleep(time.Second * 1)
		//fmt.Println("go func暂停结束")
		//cmd.Process.Signal(syscall.SIGFPE)
		//
		//fmt.Println("go func暂停1s")
		//time.Sleep(time.Second * 1)
		//fmt.Println("go func暂停结束")
		//cmd.Process.Signal(syscall.SIGILL)

		////匿名函数实现暂停恢复
		//fmt.Println("进入匿名函数")
		//fmt.Println("go func暂停3s")
		//time.Sleep(time.Second * 3)
		//fmt.Println("go func暂停结束")
		//cmd.Process.Signal(syscall.SIGTRAP)
		////cmd.Process.Signal(syscall.SIGTSTP)	//win下不可用
		////cmd.
		////w <- true
		//ww = true
		////sc.WriteString("^C")
		//fmt.Println("信号量已设置")
		//fmt.Println("go func暂停10s")
		//time.Sleep(time.Second * 10)
		//fmt.Println("go func暂停结束")
		////cmd.Process.Signal(syscall.SIGCONT)	//win下不可用
		//ww = false
		////sc.WriteByte(0x72)
		//fmt.Println("信号量已复原")

		//直接杀掉进程：
		//cmd.Process.Kill()
		//w <- false
	}()

	var line, t string
	r := regexp.MustCompile("(frame=\\s*(\\d+) fps=\\s*[\\d.]+ q=\\s*[\\d.-]+ (L?)size=\\s*[\\d\\S]+ time=[\\d:.]+ bitrate=[\\d.]*\\S?bits/s speed=\\s*[\\d.]+x)")
	for scanner.Scan() {
		t = scanner.Text()
		line += t
		if t == "\n" {
			fmt.Print("获得: " + line)
			line = ""
		} else {
			res := r.FindString(line)
			if len(res) != 0 {
				fmt.Println("#", line)
				line = ""
			}
		}
		//处理暂停恢复
		//<- w
		//if ww == true {
		//	for ww == true {
		//		fmt.Println("scanner阻塞")
		//		time.Sleep(time.Second)
		//	}
		//}
	}

	_ = cmd.Wait()
}

