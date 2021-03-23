package ptools

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"sync"
	"time"

	//"fmt"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type Tool struct {
	Name            string   //工具名
	Path            string   //工具路径，包含工具名，安装&更新时按该路径操作
	TakeOver        bool     //工具更新是否由这里接管，false->用户自行更新
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
	NonKeyWords     []string //下载的文件不包含的关键字
	Fetch           string   //在压缩包解压得到的文件中取得某文件作为工具可执行文件
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

//func CreateTool() *Tool {
//	return &Tool{
//		Name:            "",
//		Path:            "",
//		TakeOver:        false,
//		Version:         "",
//		VersionApi:      "",
//		VersionApiCDN:   "",
//		DownloadLink:    "",
//		DownloadLinkCDN: "",
//		VersionRegExp:   "",
//		GithubRepo:      "",
//		IsGitHub:        false,
//		IsCLI:           false,
//		KeyWords:        []string{},
//	}
//}

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
//      - 版本相等->先下载cdn源，失败->下载官方源 //同时下载直到某一个下载完成
//      - srcVer > cdnVer -> 下载官方源
//      - srcVer < cdnVer -> 返回error "cdn version is above source version"
// - 根据 format 安装下载好的文件 isCompressed
//    - 压缩包->解压到"dir/工具名/"
//    - 非压缩包->移动到"dir/工具名/"
func (t *Tool) Install() error {
	dir, _ := path.Split(t.Path)
	if t.CheckExist() {
		if !t.TakeOver {
			fmt.Println("请用户自行更新工具")
			return nil
		}
	} else {
		if !t.TakeOver {
			fmt.Println("用户自行更新但是工具不存在，下面尝试安装")
		}

		//检查安装位置
		if !IsFileExisted(dir) {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return err
			}
		}
	}

	var srcVer, cdnVer, srcUrl, cdnUrl string
	var srcOK, cdnOK = false, false
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		//读取官方源版本号和下载地址
		if t.IsGitHub {
			if tVer, tUrl, err := t.ParseGithubLatestRelease(); err != nil {
				srcOK = false
			} else {
				srcVer = tVer
				srcUrl = tUrl
				srcOK = true
			}
		} else {
			//利用版本号API获得官方源版本
			if data, err := GetHttpData(t.VersionApi); err != nil {
				//log.Println(err)
				srcOK = false
			} else {
				srcVer = data
				srcUrl = t.DownloadLink
				srcOK = true
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		//获取CDN源版本和下载地址
		if data, err := GetHttpData(t.VersionApiCDN); err != nil {
			//log.Println(err)
			cdnOK = false
		} else {
			cdnVer = data
			cdnUrl = t.DownloadLinkCDN
			cdnOK = true
		}
	}()

	wg.Wait()

	//决定如何下载
	var tempDir, tempVer, filename string
	if srcOK && cdnOK && CompareVersion(srcVer, cdnVer) == 0 {
		//判断是否要更新
		if CompareVersion(t.Version, srcVer) == 0 {
			return nil
		}

		tempVer = srcVer
		//优先下载cdn源
		if _, err := GrabDownload("./temp/"+t.Name+"/cdn/", cdnUrl); err != nil {
			fmt.Println("cdn源下载失败，再尝试一次...")
			time.Sleep(time.Second * 2)
			if _, err := GrabDownload("./temp/"+t.Name+"/cdn/", cdnUrl); err != nil {
				fmt.Println("cdn源下载失败，正在下载src源")
				if _, err := GrabDownload("./temp/"+t.Name+"/src/", srcUrl); err != nil {
					return err
				} else {
					tempDir = FormatPath("./temp/" + t.Name + "/src/")
					_, filename = path.Split(srcUrl)
				}
			} else {
				tempDir = FormatPath("./temp/" + t.Name + "/cdn/")
				_, filename = path.Split(cdnUrl)
			}
		} else {
			tempDir = FormatPath("./temp/" + t.Name + "/cdn/")
			_, filename = path.Split(cdnUrl)
		}
	} else {
		var url string
		if !srcOK && !cdnOK {
			return errors.New("install/update failed on src and cdn mirror")
		} else if srcOK && !cdnOK {
			//判断是否要更新
			if CompareVersion(t.Version, srcVer) == 0 {
				return nil
			}
			//下载src
			tempDir = FormatPath("./temp/" + t.Name + "/src/")
			tempVer = srcVer
			url = srcUrl
		} else if !srcOK && cdnOK {
			//判断是否要更新
			if CompareVersion(t.Version, cdnVer) == 0 {
				return nil
			}
			//下载cdn
			tempDir = FormatPath("./temp/" + t.Name + "/cdn/")
			tempVer = cdnVer
			url = cdnUrl
		} else {
			//两个源都OK 判断版本号
			switch CompareVersion(srcVer, cdnVer) {
			case 1:
				//判断是否要更新
				if CompareVersion(t.Version, srcVer) == 0 {
					return nil
				}
				//下载src
				tempDir = FormatPath("./temp/" + t.Name + "/src/")
				tempVer = srcVer
				url = srcUrl
			case -1:
				//报错
				return errors.New("cdn version is above src version")
			}
		}

		//单线下载
		var err error
		if filename, err = GrabDownload(tempDir, url); err != nil {
			return err
		}
	}

	//fmt.Println(tempDir, tempVer, filename)

	//判断文件类型
	if IsCompressed(filename) {
		//解压
		if err := SafeDecompress(tempDir+filename, tempDir+"temp"); err != nil {
			fmt.Println(tempDir+filename, " -> ", tempDir+"temp")
			return err
		}

		//移除根目录
		ok, topDir := CheckTopDir(tempDir+"temp")
		fmt.Println("ok=", ok, "topDir=", topDir)
		if ok {
			_ = os.Rename(topDir, tempDir+"temp_swap")
			_ = os.RemoveAll(tempDir+"temp")
			_ = os.Rename(tempDir+"temp_swap", tempDir+"temp")
		}

		//根据Fetch从临时文件夹中取文件
		if t.Fetch == "" {
			//转移文件
			if err := XCopy(tempDir+"/temp", dir); err != nil {
				return err
			}
		} else {
			filepath := GetFilePathFromDir(tempDir+"/temp", t.Fetch)
			if filepath == "" {
				return errors.New("cannot find the file to fetch")
			}

			if err := XCopy(filepath, t.Path); err != nil {
				return err
			}
		}
	} else {
		//直接转移
		if err := XCopy(tempDir+filename, t.Path); err != nil {
			return err
		}
	}

	t.Version = tempVer
	return os.RemoveAll("./temp/" + t.Name)
	//return nil
}

//检查更新
//@param 空
//@return error->错误
func (t *Tool) Update() error {
	return t.Install()
}

//检查工具是否存在
//@param 空
//@return bool->是否存在
func (t *Tool) CheckExist() bool {
	return IsFileExisted(t.Path)
}

//检查环境变量，有就设置Path和TakeOver=False，尝试获取Version
func (t *Tool) CheckEnvPath() bool {
	if _, err := Exec(t.Name); err != nil {
		return false
	}

	t.TakeOver = false
	t.Path = t.Name
	_ = t.SetCliVersion()

	return true
}

//获取命令行工具的版本号
//@param 空
//@return string ver->版本号, error->错误
//算法：
// - 判断工具存在 false -> 返回error
// - isCLI==false -> 返回error
// - 调用工具但不加参数，获得输出
// - 利用VersionRegExp获取版本号，获取失败则返回error
func (t *Tool) GetCliVersion() (ver string, err error) {
	if t.CheckExist() == false {
		return "", errors.New("tool does not exist")
	}

	if t.IsCLI == false {
		return "", errors.New("it is not a cli program")
	}

	output, _ := Exec(t.Path)
	re, err := regexp.Compile(t.VersionRegExp)
	if err != nil {
		return "", err
	}
	str := re.FindStringSubmatch(output)
	if len(str) < 2 {
		return "", errors.New("failed to match version string of tool: " + t.Name)
	}
	return str[1], nil
}

//设置命令行版本号
func (t *Tool) SetCliVersion() error {
	ver, err := t.GetCliVersion()
	if err != nil {
		return err
	}

	t.Version = ver
	return nil
}

//解析从Github的API得到的json数据，获得版本号和下载地址
//@param []byte json数据
//@return ver->版本号, url->下载链接, error->错误
//说明：string类型数据要转换成byte切片
func (t *Tool) ParseGithubApiData(jsonData []byte) (ver, url string, err error) {
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

				//过滤关键字
				for _, keyword := range t.KeyWords {
					if !strings.Contains(file.Name, keyword) {
						ok = false
						break
					}
				}

				//过滤非关键字
				if ok {
					for _, nonkeyword := range t.NonKeyWords {
						if strings.Contains(file.Name, nonkeyword) {
							ok = false
							break
						}
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
func (t *Tool) ParseGithubApi(api string) (ver, url string, err error) {
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
func (t *Tool) ParseGithubLatestRelease() (ver, url string, err error) {
	jsonData, err := GetHttpDataByteSlice("https://api.github.com/repos/" + t.GithubRepo + "/releases/latest")
	if err != nil {
		return "", "", err
	}

	return t.ParseGithubApiData(jsonData)
}
