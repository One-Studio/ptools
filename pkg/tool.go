package ptools

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"sync"

	//"fmt"
	jsoniter "github.com/json-iterator/go"
	"os"
	"strings"
)

type tool struct {
	Name            string   //工具名
	Path            string   //工具路径，包含工具名，安装&更新时按该路径操作
	TakeOver		bool	 //工具更新是否由这里接管，false->用户自行更新
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

//Github Latest Info
type GitHubLatest struct {
	URL     string  `json:"url"`
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Message string  `json:"message"`
	Assets  []Asset `json:"assets"`
}

//安装/更新工具
//@param 空
//@return error->错误
//说明：工具是否在环境变量中不在此包中实现，判断之后设置tool的参数，则不论在何处都是适用的
//算法：
// - 检查工具 不存在 -> 尝试安装/更新, 存在 -> TakeOver==false -> 停止更新
// - 分割路径，检查安装位置dir，不存在则创建，创建失败返回error
// - 同时检查官方源和CDN源，根据结果设置srcOK/cdnOK为true/false
// - 根据 srcOK cdnOK srcVer cdnVer 决定安装方式
//    - srcOK/cdnOK 均为false->返回error
//    - srcOK/cdnOK true/false各一->直接下载
//    - srcOK/cdnOK 均为true->比较srcVer和cdnVer
//      - 版本相等->同时下载直到某一个下载完成
//      - srcVer > cdnVer -> 下载官方源
//      - srcVer < cdnVer -> 返回error "cdn version is above source version"
// - 根据 filename 和 format 安装下载好的文件
//    - 压缩包->解压到"dir/工具名/"
//    - 非压缩包->移动到"dir/工具名/"
func (t *tool) Install() error {
	dir, bin := path.Split(t.Path)

	if t.CheckExist() {
		if t.TakeOver == false {
			fmt.Println("请用户自行更新工具")
			return nil
		}
	} else {
		if t.TakeOver {
			fmt.Println("用户自行更新但是工具不存在，下面尝试安装")
		}

		//检查安装位置
		if !IsFileExisted(t.Path) {
			if err := os.Mkdir(dir, os.ModePerm); err != nil {
				return err
			}
		}
	}

	var srcVer, cdnVer string
	var srcOK, cdnOK  = false, false


	//TODO 并发检查官方源和CDN源
	//读取官方源和CDN源的版本号和下载地址
	if t.IsGitHub {
		err := nil
		srcVer, t.DownloadLink, err = t.ParseGithubLatestRelease()
		if err != nil {
			srcOK = false
		}

	} else {
		//利用版本号获得官方源版本
		if data, err := GetHttpData(t.VersionApi); err != nil {
			//TODO 错误处理
		} else {
			srcVer = data
		}
	}

	//获取CDN源版本
	if data, err := GetHttpData(t.VersionApiCDN); err != nil {
		//TODO 错误处理
	} else {
		cdnVer = data
	}

	//决定如何下载
	if srcVer == cdnVer {

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
//@param 空
//@return error->错误
func (t *tool) Update() error {
	return t.Install()
}

//检查工具是否存在
//@param 空
//@return bool->是否存在
func (t *tool) CheckExist() bool {
	return IsFileExisted(t.Path)
}

//获取命令行工具的版本号
//@param 空
//@return string ver->版本号, error->错误
//算法：
// - 判断工具存在 false -> 返回error
// - isCLI==false -> 返回error
// - 调用工具但不加参数，获得输出
// - 利用VersionRegExp获取版本号，获取失败则返回error
func (t *tool) GetCliVersion() (ver string, err error) {
	if t.CheckExist() == false {
		return "", errors.New("tool does not exist")
	}

	if t.IsCLI == false {
		return "", errors.New("it is not a cli program")
	}

	output, err := Exec(t.Path)
	if err != nil {
		return
	} else {
		re, err := regexp.Compile(t.VersionRegExp)
		if err != nil {
			return
		}

		return re.FindString(output), nil
	}
}

//解析从Github的API得到的json数据，获得版本号和下载地址
//@param []byte json数据
//@return ver->版本号, url->下载链接, error->错误
//说明：string类型数据要转换成byte切片
//算法：
// -
func (t *tool) ParseGithubApiData(jsonData []byte) (ver, url string, err error) {
	var latestInst GitHubLatest
	var jsonX = jsoniter.ConfigCompatibleWithStandardLibrary
	err = jsonX.Unmarshal(jsonData, &latestInst)
	if err != nil || latestInst.Message == "Not Found" || strings.Contains(latestInst.Message, "API rate limit") {
		err = errors.New("failed to parse GitHub API. " + latestInst.Message)
	} else {
		//设置官方源版本
		ver = latestInst.TagName

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
					url = file.BrowserDownloadURL
					break
				}
			}
		}

		if url == "" {
			err = errors.New("no keywords matched download link is found")
		}
	}

	return
}

//解析Github的API，获得版本号和下载地址
//@param string api 接口的完整链接
//@return ver->版本号, url->下载链接, error->错误
//算法：
// - 尝试获取切片格式的数据，出错则返回error
// - 调用ParseGithubApiData
func (t *tool) ParseGithubApi(api string) (ver, url string, err error) {
	jsonData, err := GetHttpDataByteSlice(api)
	if err != nil {
		return "", "", err
	}

	return t.ParseGithubApiData(jsonData)
}

//解析Github的API，获得版本号和下载地址
//@param 空
//@return ver->版本号, url->下载链接, error->错误
//算法：
// - 利用tool.GithubRepo的用户名/仓库名得到api的链接
// - 尝试获取切片格式的数据，出错则返回error
// - 调用ParseGithubApiData
func (t *tool) ParseGithubLatestRelease() (ver, url string, err error) {
	jsonData, err := GetHttpDataByteSlice("https://api.github.com/repos/" + t.GithubRepo + "/releases/latest")
	if err != nil {
		return "", "", err
	}

	return t.ParseGithubApiData(jsonData)
}
