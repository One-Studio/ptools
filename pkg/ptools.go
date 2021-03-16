package ptools

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

//简单测试 打个招呼
func Test() {
	fmt.Println("Hello, world! It's ptools")
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

//利用HTTP Get请求获得数据的字节切片
func GetHttpDataByteSlice(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()

	return data, nil
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

//自定义Scanner分割的方式，\n和\r都分割
func ScanCRandLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	//优先分割换行\n
	if i := bytes.IndexAny(data, "\n"); i >= 0 {
		return i + 1, data[0:i], nil
	}

	//然后分割行首\r
	if i := bytes.IndexAny(data, "\r"); i >= 0 {
		return i + 1, data[0 : len(data)-1], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

//转换编码解决chcp936的中文乱码问题
func ConvertString(s string) string {
	return mahonia.NewDecoder("GBK").ConvertString(s)
}

//比较版本号 1: v1>v2  -1: v1<v2  0: v1=v2
func CompareVersion(v1, v2 string) int {
	re, err := regexp.Compile("\\d+|\\D+")
	if err != nil {
		log.Fatal(err)
	}

	s1 := re.FindAllString(v1, -1)
	s2 := re.FindAllString(v2, -1)

	var n1, n2 int64
	for i, isNum1, isNum2 := 0, false, false; i < len(s1) && i < len(s2); i++ {
		//排除版本号里的.
		if s1[i] == "." {
			continue
		}

		if isNum1, err = regexp.MatchString("\\d", s1[i]); err != nil {
			log.Fatal(err)
		} else if isNum1 {
			if n1, err = strconv.ParseInt(s1[i], 10, 32); err != nil {
				log.Fatal(err)
			}
		}

		if isNum2, err = regexp.MatchString("\\d", s2[i]); err != nil {
			log.Fatal(err)
		} else if isNum2 {
			if n2, err = strconv.ParseInt(s2[i], 10, 32); err != nil {
				log.Fatal(err)
			}
		}

		if isNum1 != isNum2 {
			//版本号格式不匹配
			log.Println("version formats of 2 strings are not correspond, using simple string compare method.")
			return strings.Compare(v1, v2)
		} else if isNum1 {
			if n1 > n2 {
				return 1
			} else if n1 < n2 {
				return -1
			}
		} else {
			if res := strings.Compare(s1[i], s2[i]); res != 0 {
				return res
			}
		}
	}

	//共有部分全部一致则根据长度决定
	if len(v1) > len(v2) {
		return 1
	} else if len(v1) < len(v2) {
		return -1
	} else {
		return 0
	}
}
