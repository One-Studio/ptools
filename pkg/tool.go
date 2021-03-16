package ptools

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"log"
	"os"
	"strings"
)

type tool struct {
	Name            string   //工具名
	Path            string   //调用路径，包含工具名
	Version         string   //版本号
	VersionApi      string   //获得版本号的官方 API
	VersionApiCDN   string   //获得版本号的CDN API
	DownloadLink    string   //官方源的下载地址
	DownloadLinkCDN string   //CDN源的下载地址
	VersionRegExp   string   //从命令行程序解析版本号的正则表达式
	GithubRepo      string   //GitHub仓库的"用户名/仓库名"
	IsGitHub        bool     //是否为GitHub地址
	IsCLI           bool     //是否为命令行程序
	KeyWords        []string //下载的文件的关键字
	//TODO 用户自行更新还是这里自动更新
	//TODO 下载之后解压还是直接copy 根据格式 zip/rar等解压 此外
}

//Github Asset
type Asset struct {
	URL                string `json:"url"`
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	ContentType        string `json:"content_type"`
	State              string `json:"state"`
	Size               int    `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

//Github latest info
type GitHubLatest struct {
	URL     string  `json:"url"`
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Message string  `json:"message"`
	Assets  []Asset `json:"assets"`
}

//安装工具
func (t *tool) Install(dir string) error {
	var srcVersion, cdnVersion, srcFilename, cdnFilename string
	var srcOK, cdnOK  = false, false

	//检查安装位置
	if !IsFileExisted(dir) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return err
		}
	}

	//TODO 并发检查官方源和CDN源
	//读取官方源和CDN源的版本号和下载地址
	if t.IsGitHub {
		srcData, err := GetHttpDataByteSlice("https://api.github.com/repos/" + t.GithubRepo + "/releases/latest")
		if err != nil {
			log.Println("failed to get GitHub API response.", err)
		}

		var latestInst GitHubLatest
		var jsonx = jsoniter.ConfigCompatibleWithStandardLibrary
		err = jsonx.Unmarshal(srcData, &latestInst)
		if err != nil || latestInst.Message == "Not Found" || strings.Contains(latestInst.Message, "API rate limit") {
			log.Println("failed to parse GitHub API. " + latestInst.Message)
		} else {
			//设置官方源版本
			srcVersion = latestInst.TagName

			//根据关键词过滤得到下载文件的链接
			for _, file := range latestInst.Assets {
				if file.State == "uploaded" {
					ok := true
					for _, keyword := range t.KeyWords {
						if !strings.Contains(file.Name, keyword) {
							ok = false
							break
						}
					}

					if ok {
						t.DownloadLink = file.BrowserDownloadURL
						srcFilename = file.Name
						srcOK = true
						break
					}
				}
			}
		}
	} else {
		//利用版本号获得官方源版本
		if data, err := GetHttpData(t.VersionApi); err != nil {
			//TODO 错误处理
		} else {
			srcVersion = data
		}
	}

	//获取CDN源版本
	if data, err := GetHttpData(t.VersionApiCDN); err != nil {
		//TODO 错误处理
	} else {
		cdnVersion = data
	}

	//决定如何下载
	if srcVersion == cdnVersion {

	}

	//下载工具 TODO
	src, cdn, path := make(chan bool), make(chan bool), ""
	go func() {
		_ = DownloadFile(t.DownloadLink, "./temp/1/")	//TODO 不会重复的下载位置

		src <- true
	}()
	go func() {
		_ = DownloadFile(t.DownloadLinkCDN, "./temp/2/")

		cdn <- true
	}()


	select {
	case <-src:
		path = "./temp/1/" + srcFilename
	case <-cdn:
		path = "./temp/2/" + cdnFilename
	}

	//安装工具 TODO
	//1-解压

	//2-直接转移
	_ = XCopy(path, dir)

	return nil
}

//检查更新
func (t *tool) CheckUpdate() error {

	return nil
}

//检查是否存在
func (t *tool) CheckExist() bool {

	return false
}

//获取命令行工具的版本号
func (t *tool) GetCliVersion() {

}

//解析Github的API，获得版本号和下载地址
func (t *tool) ParseGithubApi() {

}
