package ptools

import (
	"github.com/otiai10/copy"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

//规格化路径
func FormatPath(s string) string {
	if runtime.GOOS == "windows" {
		s = strings.Replace(s, "/", "\\", -1)
	} else {
		s = strings.Replace(s, "\\", "/", -1)
	}
	return s
}

//规格化到绝对路径
func FormatAbsPath(s string) string {
	if strings.HasPrefix(s, ".") {
		s = strings.Replace(s, ".", Getwd(), 1)
	}
	return FormatPath(s)
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
	filePath = FormatAbsPath(filePath)
	dir, _ := path.Split(filePath)
	exist := IsFileExisted(dir)
	if exist == false {
		_ = os.Mkdir(dir, os.ModePerm)
		//这里跳过检测，写文件的时候出错同样会报错
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

//去除顶层文件夹 TODO 借鉴ListDir只能获取一级目录 优化算法
func CheckTopDir(dir string) (bool, string) {
	res, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("failed to read dir: %v\n", err)
		return false, ""
	}
	
	if len(res) == 1 {
		if res[0].IsDir() {
			return true, dir + "/" + res[0].Name()
		}
	}

	return false, ""



	//var paths []string
	//var isDirs []bool
	//first := true
	//slashCount, t := 6657, 0
	//var splt string
	//if runtime.GOOS == "windows" {
	//	splt = "\\"
	//} else {
	//	splt = "/"
	//}
	//
	////获取最高层的文件/文件夹 O(n)
	//err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
	//	if f == nil {
	//		return err
	//	}
	//	
	//	if first {
	//		first = false
	//		return nil
	//	}
	//
	//	path = FormatPath(path)
	//
	//	if t = strings.Count(path, splt); t < slashCount {
	//		paths = nil
	//		isDirs = nil
	//		slashCount = t
	//	}
	//	if t == slashCount {
	//		paths = append(paths, path)
	//		if f.IsDir() {
	//			isDirs = append(isDirs, true)
	//		} else {
	//			isDirs = append(isDirs, false)
	//		}
	//	}
	//
	//	return nil
	//})
	//if err != nil {
	//	log.Printf("error when filepath.Walk(): %v\n", err)
	//}
	//
	////分析得出结果
	//if len(paths) == 1 {
	//	if isDirs[0] {
	//		return true, paths[0]
	//	}
	//}
	//
	//return false, ""
}

//遍历寻找某个文件
func GetFilePathFromDir(dir, name string) (result string) {
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if result != "" {
			return nil
		}

		if f.IsDir() {
			return nil
		}

		if f.Name() == name {
			result = path
		}

		return nil
	})
	if err != nil {
		log.Printf("error when filepath.Walk(): %v\n", err)
	}

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
}

//移动文件夹或者文件
func XMove(from, to string) error {
	if err := copy.Copy(from, to); err != nil {
		return err
	}

	return os.RemoveAll(from)
}

