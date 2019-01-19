package tyut_osc

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
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

	client = &http.Client{
		Jar:     jar,
		Timeout: time.Second * 15,
	}

	fmt.Println("config.BaseLocationGPA: ", crawler.config.BaseLocationGPA)

	// 初始化Cookie
	res, err := client.Get(crawler.config.BaseLocationGPA[0])
	fmt.Println("first get status: ", res.Status, " length: ", res.ContentLength)

	// 登录验证
	formValues := url.Values{}
	formValues.Add("u", stuid)
	formValues.Add("p", stuPassword)
	formValues.Add("r", "on")
	resp, err := client.PostForm(crawler.config.BaseLocationGPA[0]+"/Hander/LoginAjax.ashx",
		formValues)
	fmt.Println("login response length: ", resp.ContentLength, " status: ", resp.Status)
	if err != nil {
		errorWithTime("GPA 系统请求失败" + " - 正在尝试学生:" + stuid)
		return
	}
	defer resp.Body.Close()

	bodyData, _ := ioutil.ReadAll(resp.Body)
	bodyJson := map[string]interface{}{}
	json.Unmarshal(bodyData, &bodyJson)
	if bodyJson["Code"] != 1.0 {
		logger.Warn("Cannot Login. "+"gpa系统返回值:"+fmt.Sprintf("%s", bodyJson), zap.String("stuid", stuid), zap.Time("time", time.Now()))
		err = fmt.Errorf("登陆失败,提示:%s", fmt.Sprint(bodyData))
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

	rawJsonBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("无法读取gpa原始信息响应体", zap.String("stuid", stuid), zap.Time("time", time.Now()))
	}

	// 去掉返回值中的数组表达
	result = strings.Trim(string(rawJsonBytes), "[]")
	return result, nil
}

func (crawler *GpaCrawler) GetGpaInfo(stuid string, stuPassword string, targetStuid string) (string, error) {
	client, err := crawler.createClientAndLogin(stuid, stuPassword)
	if err != nil {
		return "", err
	}
	jsonText, err := crawler.fetchGpaJson(targetStuid, client)

	if err != nil {
		return "", err
	}

	return jsonText, nil
}
