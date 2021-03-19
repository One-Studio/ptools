package main

import (
	"fmt"
	tool "github.com/One-Studio/ptools/pkg"
	"github.com/cavaliercoder/grab"
	"log"
	"sync"
	"time"
)

func test1()  {
	//ffcommand := "ffmpeg -i ~/movies/演示找A卡驱动.mp4 -c:v libx264 -crf 26 -preset 3 -c:a copy ~/测试.mp4 -y"
	ffcommand := "E:/测试/ffmpeg.exe -i E:/测试/测试ff.mp4 -c:v libx264 -crf 26 -preset 3 -c:a copy E:/测试/结果.mp4 -y"
	suspend := "C:\\Users\\Purp1e\\go\\src\\github.com\\One-Studio\\ptools\\temp\\pssuspend.exe"
	//x264command := "E:/测试/x264.exe E:/测试/测试ff.mp4 --crf 26 --preset slow -output E:/测试/结果.mp4"
	a := make(chan rune)
	defer close(a)
	go func() {
		time.Sleep(time.Second *1)
		fmt.Println("触发暂停")
		tool.Pause(a)
		time.Sleep(time.Second *2)
		fmt.Println("触发继续")
		tool.Resume(a)
		time.Sleep(time.Second *2)
		fmt.Println("触发结束")
		tool.Quit(a)
	}()
	err := tool.ExecRealtimeControl(ffcommand, func(line string) {
		fmt.Println(line)
	}, a, suspend)
	close(a)

	if err != nil {
		log.Println(err)
	}

}

func testGrab()  {
	//https://cdn.jsdelivr.net/gh/One-Studio/FFmpeg-Win64@master/download_link
	resp, err := grab.Get(".", "https://cdn.jsdelivr.net/gh/One-Studio/FFmpeg-Win64@master/download_link")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Download saved to", resp.Filename)
}

func testChan()  {
	value := ""
	var c = make(chan string)
	defer close(c)
	go func() {
		fmt.Println("第一个go routine")
		value = "第一个"
		c <- "1"
	}()

	go func() {
		fmt.Println("第二个go routine")
		value = "第二个"
		c <- "2"
	}()


	fmt.Println("测试结束", <- c)
	time.Sleep(time.Second)
	fmt.Println("测试结束", <- c)

}

func testWG()  {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("第一件事做完")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("第二件事做完")
	}()

	wg.Wait()
	time.Sleep(time.Second)
}

func testCompareVersion()  {
	fmt.Println(tool.CompareVersion("v1.0.1-alpha", "v1.0.2-aa"))
	fmt.Println(tool.CompareVersion("a1.0.1", "v2.0.0-alpha.47"))
	fmt.Println(tool.CompareVersion("z1.0.1", "v2.0.0-alpha.47"))
	fmt.Println(tool.CompareVersion("v1.2.3", "v1.2.3"))
	fmt.Println(tool.CompareVersion("v1.2.3", "v1.2"))
}

func testTool()  {
	//var t = tool.CreateTool()
	var t = tool.Tool{
		Name: "hlae",
		Path: "./bin/hlae/hlae.exe",
		TakeOver: true,
		Version: "",
		VersionApi: "",
		VersionApiCDN: "https://cdn.jsdelivr.net/gh/One-Studio/HLAE-Archive@master/version",
		DownloadLink: "",
		DownloadLinkCDN: "https://cdn.jsdelivr.net/gh/One-Studio/HLAE-Archive@master/dist/hlae.zip",
		VersionRegExp: "",
		GithubRepo: "advancedfx/advancedfx",
		IsGitHub: true,
		IsCLI: false,
		KeyWords: []string{"hlae", "zip"},
		NonKeyWords: []string{".asc"},
	}

	fmt.Println(t.Install())
	fmt.Println("=====\n当前参数\n", t)
}

func testTool1()  {
	var t = tool.Tool{
		Name: "ffmpeg",
		Path: "./bin/ffmpeg/ffmpeg.exe",
		TakeOver: true,
		Version: "",
		VersionApi: "https://www.gyan.dev/ffmpeg/builds/release-version",
		VersionApiCDN: "https://cdn.jsdelivr.net/gh/One-Studio/FFmpeg-Win64@master/version",
		DownloadLink: "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.7z",
		DownloadLinkCDN: "https://cdn.jsdelivr.net/gh/One-Studio/FFmpeg-Win64@master/dist/ffmpeg-release-essentials.7z",
		VersionRegExp: "ffmpeg version (\\S+)-essentials_build-www.gyan.dev",
		GithubRepo: "advancedfx/advancedfx",
		IsGitHub: false,
		IsCLI: true,
		KeyWords: []string{},
		NonKeyWords: []string{},
	}

	fmt.Println(t.Install())
	fmt.Println("=====\n当前参数\n", t)
	fmt.Println(t.CheckExist())
	fmt.Println(t.GetCliVersion())
	fmt.Println(t.Update())
	fmt.Println("=====\n当前参数\n", t)
}


func main() {

	tool.Test()

	//testChan()
	//testWG()
	testTool()
	testTool1()
}
