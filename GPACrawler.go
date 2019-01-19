package tyut_osc

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client = &http.Client{
		Jar:     jar,
		Timeout: time.Second * 15,
	}

	// 初始化Cookie
	_, er := client.Get(crawler.config.BaseLocationGPA[0])
	if er != nil {
		uid, _ := uuid.NewUUID()
		uids := strings.Split(uid.String(), "-")[0]
		logger.Error("无法连接GPA系统主机", zap.String("detail", er.Error()), zap.String("errid", uids))
		err = fmt.Errorf("无法连接GPA教务系统,错误id: %s", uids)
		return
	}

	// 登录验证
	formValues := url.Values{}
	formValues.Add("u", stuid)
	formValues.Add("p", stuPassword)
	formValues.Add("r", "on")
	resp, er := client.PostForm(crawler.config.BaseLocationGPA[0]+"/Hander/LoginAjax.ashx",
		formValues)
	if er != nil {

		return
	}
	defer resp.Body.Close()

	bodyData, _ := ioutil.ReadAll(resp.Body)
	bodyJson := map[string]interface{}{}
	err = json.Unmarshal(bodyData, &bodyJson)

	if err != nil || bodyJson["Code"] != 1.0 {
		if err == nil {
			logger.Warn("Cannot Login. "+"gpa系统返回值:"+fmt.Sprintf("%s", bodyJson), zap.String("stuid", stuid), zap.Time("time", time.Now()))
			err = fmt.Errorf("登陆失败,提示:%s", fmt.Sprint(bodyData))
		} else {
			uid, _ := uuid.NewUUID()
			uids := strings.Split(uid.String(), "-")[0]
			err = fmt.Errorf("未知异常. 错误id:%s", uids)
			return
		}

	} else {
		logger.Info("GPA 系统登陆成功", zap.String("stuid", stuid))
		err = nil
	}
	return
}

//FetchGpaJson 接受一个已经准备好并通过登录认证的http.Client指针, 返回gpa教务系统的原生json(已处理掉数组表达)
func (crawler *GpaCrawler) fetchGpaJson(stuid string, client *http.Client) (string, error) {
	result := ""

	uid, _ := uuid.NewUUID()
	uids := strings.Split(uid.String(), "-")[0]

	resp, err := client.PostForm(crawler.config.BaseLocationGPA[0]+"/Hander/Cj/CjAjax.ashx?rnd=0.04993201044579343",
		url.Values{"limit": {"40"}, "offset": {"0"}, "order": {"asc"}, "sort": {"jqzypm%2Cxh"},
			"do": {"xsgrcj"}, "xh": {stuid}})

	if err != nil {
		logger.Error("获取GPA原始信息失败", zap.String("stuid", stuid), zap.Time("time", time.Now()), zap.String("detail", err.Error()),
			zap.String("errid", uids))
		return "", fmt.Errorf("获取GPA教务系统信息失败.错误id:%s", uids)
	}
	defer resp.Body.Close()

	rawJsonBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("无法读取gpa原始信息响应体", zap.String("stuid", stuid), zap.Time("time", time.Now()))
		return "", fmt.Errorf("未知错误. 错误id:%s", uids)
	}

	// 去掉返回值中的数组表达
	result = strings.Trim(string(rawJsonBytes), "[]")
	return result, nil
}

func (crawler *GpaCrawler) GetGpaInfo(stuid string, stuPassword string, targetStuid string) (*GpaInfo, error) {
	uid, _ := uuid.NewUUID()
	uids := strings.Split(uid.String(), "-")[0]

	client, err := crawler.createClientAndLogin(stuid, stuPassword)
	if err != nil {
		return nil, err
	}

	jsonText, err := crawler.fetchGpaJson(targetStuid, client)
	if err != nil {
		return nil, err
	}

	// 解析Json
	gpainfo := GpaInfo{}
	err = json.Unmarshal([]byte(jsonText), &gpainfo)
	if err != nil {
		logger.Error("无法解析GPA JSON", zap.Time("time", time.Now()), zap.String("detail", err.Error()))
		return nil, fmt.Errorf("未知错误. 错误id:%s", uids)
	}
	return &gpainfo, nil
}
