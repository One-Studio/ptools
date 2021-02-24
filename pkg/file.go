package ptools

import (
	"github.com/otiai10/copy"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

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

//获取当前程序路径 最好直接用 os.Getwd
func Getwd() (wd string) {
	wd, _ = os.Getwd()
	return
}

//获取指定路径下的所有文件，只搜索当前路径，不进入下一级目录，可匹配后缀过滤（suffix为空则不过滤）TODO
func ListDir(path2List, suffix string) (files []string, err error) {
	files = []string{}

	dir, err := ioutil.ReadDir(path2List)
	if err != nil {
		return nil, err
	}

	suffix = strings.ToLower(suffix) //匹配后缀

	for _, v := range dir {
		//if v.IsDir() {
		//	continue //忽略目录
		//}
		if len(suffix) == 0 || strings.HasSuffix(strings.ToLower(v.Name()), suffix) {
			//文件后缀匹配
			files = append(files, path.Join(path2List, v.Name()))
		}
	}

	return files, nil
}

//复制文件夹或者文件
func XCopy(from, to string) error {
	return copy.Copy(from, to)
	//from = FormatPath(from)
	//to = FormatPath(to)
	//
	////确保目标路径存在，否则复制报错exit status 4 || WTF cmd用就行 这里用就不行
	//exist := IsFileExisted(to)
	//if exist == false {
	//	err := os.Mkdir(to, os.ModePerm)
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//var command string
	//if runtime.GOOS == "windows" {
	//	command = "xcopy /I /E /Y /R " + strconv.Quote(from) + " " + strconv.Quote(to)
	//} else {
	//	command = "cp -R " + strconv.Quote(from) + " " + strconv.Quote(to)
	//}
	//
	//fmt.Println(command)
	//_, err := Exec(command)
	//return err
}

//移动文件夹或者文件
func XMove(from, to string) error {
	if err := copy.Copy(from, to); err != nil {
		return err
	}

	return os.RemoveAll(from)
}

