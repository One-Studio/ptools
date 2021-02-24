/*
 * MIT License
 *
 * Copyright (c) 2021. Purp1e
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
package ptools

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func Test() {
	fmt.Println("say hello to ptools~")
}

//快速打开文件和读内容
func ReadAll(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	//对内容的操作
	//ReadFile返回的是[]byte字节切片，要用string()方法转变成字符串
	//去除内容结尾的换行符
	str := strings.TrimRight(string(content), "\n")
	return str, nil
}

//快速文件写先清空再写入
func WriteFast(filePath string, content string) error {
	dir, _ := path.Split(filePath)
	exist := IsFileExisted(dir)
	if exist == false {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	err := ioutil.WriteFile(filePath, []byte(content), 0666)
	if err != nil {
		return err
	} else {
		return nil
	}
}

//判断文件/文件夹是否存在
func IsFileExisted(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if os.IsNotExist(err) {
		return false
	}
	return true
}

//获取指定路径下的所有文件，只搜索当前路径，不进入下一级目录，可匹配后缀过滤（suffix为空则不过滤）
func ListDir(dir, suffix string) (files []string, err error) {
	files = []string{}

	_dir, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	suffix = strings.ToLower(suffix) //匹配后缀

	for _, _file := range _dir {
		if _file.IsDir() {
			continue //忽略目录
		}
		if len(suffix) == 0 || strings.HasSuffix(strings.ToLower(_file.Name()), suffix) {
			//文件后缀匹配
			files = append(files, path.Join(dir, _file.Name()))
		}
	}

	return files, nil
}

//TODO 复制文件
func CopyFile() {
	//打开原始文件
	originalFile, err := os.Open("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer originalFile.Close()
	//创建新的文件作为目标文件
	newFile, err := os.Create("test_copy.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()
	//从源中复制字节到目标文件
	bytesWritten, err := io.Copy(newFile, originalFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Copied %d bytes.", bytesWritten)
	//将文件内容flush到硬盘中
	err = newFile.Sync()
	if err != nil {
		log.Fatal(err)
	}

}

//TODO 移动文件
func MoveFile() {

}

//利用HTTP Get请求获得数据
func GetHttpData(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	_ = resp.Body.Close()

	return string(data), nil
}

//下载文件 (下载地址，存放位置)
func DownloadFile(url string, location string) error {
	//利用HTTP下载文件并读取内容给data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		errorInfo := "http failed, check if file exists, HTTP Status Code:" + strconv.Itoa(resp.StatusCode)
		return errors.New(errorInfo)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	//确保下载位置存在
	_, fileName := path.Split(url)
	ok := IsFileExisted(location)
	if ok == false {
		err := os.Mkdir(location, os.ModePerm)
		if err != nil {
			return err
		}
	}
	//文件写入 先清空再写入 利用ioutil
	err = ioutil.WriteFile(location+"/"+fileName, data, 0666)
	if err != nil {
		return err
	} else {
		return nil
	}
}

//判断是不是non-ASCII
func IsNonASCII(str string) bool {
	re := regexp.MustCompile("[[:^ascii:]]")
	return re.MatchString(str)
}

//获取当前程序路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//规格化路径
func FormatPath(s string) string {
	return strings.TrimRight(strings.Replace(s, "\\", "/", -1), "\\")
}

//规格化到绝对路径
func FormatAbsPath(s string) string {
	if strings.HasPrefix(s, ".") {
		s = strings.Replace(s, ".", GetCurrentDirectory(), 1)
	}
	return strings.TrimRight(strings.Replace(s, "\\", "/", -1), "\\")
}

//复制文件夹
func CopyDir(from string, to string) error {
	from = FormatPath(from)
	to = FormatPath(to)

	//确保目标路径存在，否则复制报错exit status 4
	exist := IsFileExisted(to)
	if exist == false {
		err := os.Mkdir(to, os.ModePerm)
		if err != nil {
			return err
		}
	}

	var command string
	if runtime.GOOS == "windows" {
		command = "xcopy " + from + " " + to + " /I /E /Y /R"
	} else {
		command = "cp -R " + from + " " + to
	}

	out, err := Exec(command)
	if err != nil {
		log.Println(out, err)
	}
	return err
}

//执行一次command指令 跨平台兼容
func Exec(command string) (string, error) {
	cmdArgs := strings.Fields(command)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	//隐藏黑框
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	out, err := cmd.CombinedOutput()
	return string(out), err
	//var out []byte
	//var err error
	//var cmd *exec.Cmd
	//if runtime.GOOS == "windows" {
	//	cmd = exec.Command("cmd.exe", "/c", command)
	//} else {
	//	cmd = exec.Command("/bin/bash", "-c", command)
	//}
	////隐藏黑框
	//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	//
	//out, err = cmd.CombinedOutput()
	//return string(out), err
}

//执行一次command指令并实时输出每行结果 跨平台兼容
func ExecRealtime(command string) error {
	cmdArgs := strings.Fields(command)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	//隐藏黑框
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	in := bufio.NewScanner(stdout)
	for in.Scan() {
		cmdRe:=ConvertByte2String(in.Bytes(),"GB18030")
		fmt.Println(cmdRe)
	}

	err = cmd.Wait()
	return err
}

//转换编码解决乱码问题
func ConvertByte2String(byte []byte, charset string) string {
	var str string
	switch charset {
	case "GB18030":
		var decodeBytes,_=simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str= string(decodeBytes)
	case "UTF8":
		fallthrough
	default:
		str = string(byte)
	}
	return str
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

//TODO 暂停目前只能linux macos 使用
func FFmpegPauseWhenRunning(command string) {
	//command := "./x264 /Users/purp1e/vd/in.mp4 --crf 18 --preset 3 -o /Users/purp1e/vd/outx264.mp4"
	//command := "ffmpeg -i /Users/purp1e/vd/in.mp4 -vcodec libx264 -crf 18 -preset 1 -acodec copy /Users/purp1e/vd/out3.mp4 -y"
	cmdArgs := strings.Fields(command)
	fmt.Println(cmdArgs)
	//cmd := exec.Command(command)//"/bin/bash", "-c",
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	stderr, _ := cmd.StderrPipe()
	defer stderr.Close()
	//stdin, _ := cmd.StdinPipe()
	//defer stdin.Close()
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
	ww := false
	go func() {
		//匿名函数实现暂停恢复
		fmt.Println("进入匿名函数")
		fmt.Println("go func暂停3s")
		time.Sleep(time.Second * 3)
		fmt.Println("go func暂停结束")
		cmd.Process.Signal(syscall.SIGTRAP)
		//cmd.Process.Signal(syscall.SIGTSTP)	//win下不可用
		//cmd.
		//w <- true
		ww = true
		//sc.WriteString("^C")
		fmt.Println("信号量已设置")
		fmt.Println("go func暂停10s")
		time.Sleep(time.Second * 10)
		fmt.Println("go func暂停结束")
		//cmd.Process.Signal(syscall.SIGCONT)	//win下不可用
		ww = false
		//sc.WriteByte(0x72)
		fmt.Println("信号量已复原")

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

