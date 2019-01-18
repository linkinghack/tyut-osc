package tyut_osc

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// createClientAndLogin 接受gpa教务系统学号和密码, 返回一个登陆状态ok的http.Client
func (crawler *GpaCrawler) createClientAndLogin(stuid string, stuPassword string) (client *http.Client, err error) {
	// 准备HttpClient
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		//logger.Error("无法新建CookieJar",zap.Time("time",time.Now()))
		errorWithTime("无法新建CookieJar")
		return
	}

	duration, _ := time.ParseDuration("15s")
	client = &http.Client{
		Jar:     jar,
		Timeout: duration,
	}

	// 登录验证
	url_idx := 0 // 准备尝试的url下标
	resp, err := client.PostForm(crawler.config.BaseLocationGPA[url_idx]+"/Hander/LoginAjax.ashx",
		url.Values{"u": {stuid}, "p": {stuPassword}, "r": {"on"}})
	if err != nil {
		errorWithTime("GPA 系统请求失败" + " - 正在尝试学生:" + stuid)
		return
	}
	defer resp.Body.Close()

	bodyData := []byte{}
	bodyJson := map[string]interface{}{}
	resp.Body.Read(bodyData)
	json.Unmarshal(bodyData, &bodyJson)
	if bodyJson["Code"] != 1.0 {
		logger.Warn("Cannot Login. "+"gpa系统返回值:"+fmt.Sprintf("%s", bodyJson), zap.String("stuid", stuid), zap.Time("time", time.Now()))
		err = fmt.Errorf("登陆失败,提示:%s", bodyJson["Msg"])
	} else {
		logger.Info("GPA 系统登陆成功", zap.String("stuid", stuid))
		err = nil
	}
	return
}

//FetchGpaJson 接受一个已经准备好并通过登录认证的http.Client指针, 返回gpa教务系统的原生json(已处理掉数组表达)
func (crawler *GpaCrawler) fetchGpaJson(stuid string, client *http.Client) (string, error) {
	result := "[]"

	resp, err := client.PostForm(crawler.config.BaseLocationGPA[0]+"/Hander/Cj/CjAjax.ashx?rnd=0.04993201044579343",
		url.Values{"limit": {"40"}, "offset": {"0"}, "order": {"asc"}, "sort": {"jqzypm%2Cxh"},
			"do": {"xsgrcj"}, "xh": {stuid}})

	if err != nil {
		logger.Error("获取GPA原始信息失败", zap.String("stuid", stuid), zap.Time("time", time.Now()))
		return "", err
	}
	defer resp.Body.Close()

	rawJsonBytes := []byte{}
	_, err = resp.Body.Read(rawJsonBytes)
	if err != nil {
		logger.Error("无法读取gpa原始信息响应体", zap.String("stuid", stuid), zap.Time("time", time.Now()))
	}

	// 去掉返回值中的数组表达
	result = strings.Trim(string(rawJsonBytes), "[]")
	return result, nil
}
