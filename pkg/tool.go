package ptools

type tool struct {
	name            string //工具名
	path            string //完整调用路径
	version         string //版本号
	srcBias         bool   //true=偏向CDN源，false=偏向官方源
	verLink         string //获得版本号的链接
	dLink           string //获得下载地址的链接
	verLinkCDN      string //获得版本号的CDN链接
	dLinkCDN        string //获得下载地址的CDN链接
	isGitHub        bool   //是否为GitHub地址
	isCLI           bool   //是否为命令行程序
	Re2parseVersion string //从命令行程序解析版本号的正则表达式
}

//检查更新
func (t tool) CheckUpdate() {

}

//检查是否存在
func (t tool) CheckExist() bool {

	return false
}

//获取命令行工具的版本号
func (t tool) GetCliVersion() {

}

//安装工具
func (t tool) Install() {

}

//解析Github链接
func (t tool) ParseGithubSrc() {

}
