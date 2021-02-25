package main

import (
	"fmt"
	tool "github.com/One-Studio/ptools/pkg"
	"log"
	"time"
)

func main() {

	tool.Test()

	ffcommand := "ffmpeg -i ~/movies/演示找A卡驱动.mp4 -c:v libx264 -crf 26 -preset 3 -c:a copy ~/测试.mp4 -y"
	//ffcommand := "E:/测试/ffmpeg.exe -i E:/测试/测试ff.mp4 -c:v libx264 -crf 26 -preset 3 -c:a copy E:/测试/结果.mp4 -y"
	suspend := "C:\\Users\\Purp1e\\go\\src\\github.com\\One-Studio\\ptools\\temp\\pssuspend.exe"
	//x264command := "E:/测试/x264.exe E:/测试/测试ff.mp4 --crf 26 --preset slow -output E:/测试/结果.mp4"
	a := make(chan rune)
	go func() {
		time.Sleep(time.Second *1)
		fmt.Println("触发暂停")
		a <- 'p'
		time.Sleep(time.Second *2)
		fmt.Println("触发继续")
		a <- 'r'
		time.Sleep(time.Second *2)
		fmt.Println("触发结束")
		a <- 'q'
	}()
	err := tool.ExecRealtimeControl(ffcommand, func(line string) {
		fmt.Println(line)
	}, a, suspend)

	if err != nil {
		log.Println(err)
	}

}
