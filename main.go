package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	tool "github.com/One-Studio/ptools/pkg"
	"github.com/cavaliercoder/grab"
)

func test1() {
	//ffcommand := "ffmpeg -i ~/movies/演示找A卡驱动.mp4 -c:v libx264 -crf 26 -preset 3 -c:a copy ~/测试.mp4 -y"
	ffcommand := "E:/测试/ffmpeg.exe -i E:/测试/测试ff.mp4 -c:v libx264 -crf 26 -preset 3 -c:a copy E:/测试/结果.mp4 -y"
	suspend := "C:\\Users\\Purp1e\\go\\src\\github.com\\One-Studio\\ptools\\temp\\pssuspend.exe"
	//x264command := "E:/测试/x264.exe E:/测试/测试ff.mp4 --crf 26 --preset slow -output E:/测试/结果.mp4"
	a := make(chan rune)
	defer close(a)
	go func() {
		time.Sleep(time.Second * 1)
		fmt.Println("触发暂停")
		tool.Pause(a)
		time.Sleep(time.Second * 2)
		fmt.Println("触发继续")
		tool.Resume(a)
		time.Sleep(time.Second * 2)
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

func testGrab() {
	//https://cdn.jsdelivr.net/gh/One-Studio/FFmpeg-Win64@master/download_link
	resp, err := grab.Get(".", "https://cdn.jsdelivr.net/gh/One-Studio/FFmpeg-Win64@master/download_link")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Download saved to", resp.Filename)
}

func testChan() {
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

	fmt.Println("测试结束", <-c)
	time.Sleep(time.Second)
	fmt.Println("测试结束", <-c)

}

func testWG() {
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

func testCompareVersion() {
	fmt.Println(tool.CompareVersion("v1.0.1-alpha", "v1.0.2-aa"))
	fmt.Println(tool.CompareVersion("a1.0.1", "v2.0.0-alpha.47"))
	fmt.Println(tool.CompareVersion("z1.0.1", "v2.0.0-alpha.47"))
	fmt.Println(tool.CompareVersion("v1.2.3", "v1.2.3"))
	fmt.Println(tool.CompareVersion("v1.2.3", "v1.2"))
}

func testTool() {
	//var t = tool.CreateTool()
	var t = tool.Tool{
		Name:            "hlae",
		Path:            "./bin/hlae/hlae.exe",
		TakeOver:        true,
		Version:         "",
		VersionApi:      "",
		VersionApiCDN:   "https://cdn.jsdelivr.net/gh/One-Studio/HLAE-Archive@master/version",
		DownloadLink:    "",
		DownloadLinkCDN: "https://cdn.jsdelivr.net/gh/One-Studio/HLAE-Archive@master/dist/hlae.zip",
		VersionRegExp:   "",
		GithubRepo:      "advancedfx/advancedfx",
		IsGitHub:        true,
		IsCLI:           false,
		KeyWords:        []string{"hlae", "zip"},
		NonKeyWords:     []string{".asc"},
	}

	fmt.Println(t.Install())
	fmt.Println("=====\n当前参数\n", t)
}

func testTool1() {
	var t = tool.Tool{
		Name:            "ffmpeg",
		Path:            "./bin/ffmpeg/ffmpeg.exe",
		TakeOver:        true,
		Version:         "",
		VersionApi:      "https://www.gyan.dev/ffmpeg/builds/release-version",
		VersionApiCDN:   "https://cdn.jsdelivr.net/gh/One-Studio/FFmpeg-Win64@master/version",
		DownloadLink:    "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.7z",
		DownloadLinkCDN: "https://cdn.jsdelivr.net/gh/One-Studio/FFmpeg-Win64@master/dist/ffmpeg-release-essentials.7z",
		VersionRegExp:   "ffmpeg version (\\S+)-essentials_build-www.gyan.dev",
		GithubRepo:      "",
		IsGitHub:        false,
		IsCLI:           true,
		Fetch:           "",
		//Fetch:           "ffmpeg.exe",
	}

	if err := t.Install(); err != nil {
		fmt.Println(err)
	}
}

func testXMove() {
	fmt.Println(tool.XMove("C:\\Users\\Purp1e\\go\\src\\github.com\\One-Studio\\ptools\\bin\\ffmpeg\\ffmpeg-release-essentials\\README.txt",
		"C:\\Users\\Purp1e\\Desktop\\"))
}

func testDecomp1() {
	//去除顶层文件夹
	//_ = os.MkdirAll("./temp/ffmpeg", os.ModePerm)
	tempDir := "./bin/"
	filename := "你好.7z"

	to := "./temp/ffmpeg"

	err := tool.SafeDecompress(tempDir+filename, to)
	if err != nil {
		log.Println(err)
	}

	//TODO 核心操作

	//fmt.Println("得到的exe路径:", tool.GetFilePathFromDir("./temp/hlae", "AfxHook.dat"))
	////temp\ffmpeg\AfxHook.dat
	//_ = tool.XMove("temp\\hlae\\AfxHook.dat", "D:\\afx.dat")
}

func testTopDir() {
	//fetch binary
	ok, path := tool.CheckTopDir("./bin/ffmpeg")
	if ok {
		fmt.Println(path)
	}
}

func testConfigDir() {
	fmt.Println(tool.ConfigDir())
}

//TODO 测试完这个，win平台的几个工具就能用了
func testFFmpeg() {
	command := []string{
		"bin\\ffmpeg.exe",
		"-i",
		"C:\\Users\\Purp1e\\Videos\\sb战术.mp4",
		"-vcodec",
		"libx264",
		"-crf",
		"20",
		"-preset",
		"slow",
		"C:\\Users\\Purp1e\\Videos\\sb战术_encode.mp4",
	}

	if out, err := tool.ExecArgs(command); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(out)
	}
	//	func(line string) {
	//		fmt.Println(line)
	//	}); err != nil {
	//	fmt.Println(err)
	//}
}

func main() {

	tool.Test()

	//testConfigDir()
	//testFFmpeg()
	//testChan()
	//testWG()
	//testTool()
	//testTool1()
	//testXMove()
	//testDecomp1()
	//testTopDir()
}
