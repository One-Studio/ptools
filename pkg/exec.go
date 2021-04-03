package ptools

import (
	"fmt"
	"strings"
)

//执行一次command指令 直接调用
func Exec(command string) (output string, err error) {
	cmdArgs := strings.Fields(command)

	return ExecArgs(cmdArgs)
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
