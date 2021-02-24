package ptools

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
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

