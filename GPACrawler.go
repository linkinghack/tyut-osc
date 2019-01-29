package tyut_osc

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/linkinghack/tyut-osc/DataModel"
	"go.uber.org/zap"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// GpaCrawler is the object representing the crawl engine
// Generally, just use the DefaultGpaCrawler is OK. The difference between different
//GpaCrawlers is just the Configuration.
type GpaCrawler struct {
	config *Configuration
}

func (e *GpaCrawler) SetConfiguration(conf *Configuration) {
	e.config = conf
}

// Create an instance of GpaCrawler.
// @returns The pointer of the GpaCrawler{} just created.
func NewGpaCrawler() *GpaCrawler {
	defer logger.Sync()
	// 初始化gpa教务系统配置
	defaultConfig := &Configuration{}
	DefaultGpaCrawler := &GpaCrawler{}

	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		logger.Warn("无法加载GPACrawler配置文件: config.json")
		defaultConfig = loadDefaultConfiguration()
		logger.Info("已使用默认配置创建GPACrawler")
	} else {
		err = json.Unmarshal(configFile, defaultConfig)
		if err != nil {
			defaultConfig = loadDefaultConfiguration()
		}
		// 配置文件错误格式不正确
		if defaultConfig.BaseLocationGPA == nil || defaultConfig.BaseLocationURP == nil {
			logger.Error("config.json 中无法读取所需信息。请正确定义BaseLocationURP:[]string 和 BaseLocationGPA:[]string")
			defaultConfig = loadDefaultConfiguration()
			logger.Info("使用默认配置创建GPACrawler", zap.Time("time", time.Now()))
		}
	}
	DefaultGpaCrawler.SetConfiguration(defaultConfig)
	logger.Info("Crawler init done.")
	return DefaultGpaCrawler
}

// createClientAndLogin 接受gpa教务系统学号和密码, 返回一个登陆状态ok的http.Client
func (crawler *GpaCrawler) createClientAndLogin(stuid string, stuPassword string) (client *http.Client, err error) {
	uid, _ := uuid.NewUUID()
	uids := strings.Split(uid.String(), "-")[0]
	defer logger.Sync()

	// 准备HttpClient
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client = &http.Client{
		Jar:     jar,
		Timeout: time.Second * 15,
	}

	// 初始化Cookie
	_, er := client.Get(crawler.config.BaseLocationGPA[0])
	if er != nil {
		logger.Error("无法连接GPA系统主机 可能超时错误", zap.String("errid", uids), zap.String("stuid", stuid), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
		err = fmt.Errorf("连接到GPA教务系统超时,错误id: %s", uids)
		client = nil
		return
	}

	// 登录验证
	formValues := url.Values{}
	formValues.Add("u", stuid)
	formValues.Add("p", stuPassword)
	formValues.Add("r", "on")
	resp, er := client.PostForm(crawler.config.BaseLocationGPA[0]+"/Hander/LoginAjax.ashx",
		formValues)
	if er != nil { // 一般不会发生，超时错误已在初始化cookie 阶段处理
		logger.Info("GPA教务系统登录失败", zap.String("stuid", stuid), zap.String("errid", uids), zap.Time("time", time.Now()))
		fmt.Errorf("GPA教务系统登录失败,错误id: %s", uids)
		return
	}
	defer resp.Body.Close()

	// 分析登录结果
	bodyData, _ := ioutil.ReadAll(resp.Body)
	bodyJson := map[string]interface{}{}
	err = json.Unmarshal(bodyData, &bodyJson)

	if err != nil || bodyJson["Code"] != 1.0 {
		if err == nil {
			logger.Warn("Cannot Login. "+"gpa系统返回值:"+fmt.Sprintf("%s", bodyJson), zap.String("stuid", stuid), zap.Time("time", time.Now()))
			err = fmt.Errorf("登陆失败,提示:%s", fmt.Sprint(bodyJson["Msg"]))
		} else {
			err = fmt.Errorf("未知异常. 错误id:%s", uids)
			return
		}

	} else {
		logger.Info("GPA 系统登陆成功", zap.String("stuid", stuid), zap.Time("time", time.Now()))
		err = nil
	}
	return
}

//FetchGpaJson 接受一个已经准备好并通过登录认证的http.Client指针, 返回gpa教务系统的原生json(已处理掉数组表达)
func (crawler *GpaCrawler) fetchGpaJson(stuid string, client *http.Client) (string, error) {
	defer logger.Sync()
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

func (crawler *GpaCrawler) GetGpaRank(stuid string, stuPassword string, targetStuid string) (*DataModel.GpaRank, error) {
	uid, _ := uuid.NewUUID()
	uids := strings.Split(uid.String(), "-")[0]
	defer logger.Sync()

	client, err := crawler.createClientAndLogin(stuid, stuPassword)
	if err != nil {
		return nil, err
	}

	jsonText, err := crawler.fetchGpaJson(targetStuid, client)
	if err != nil {
		return nil, err
	}

	// 解析Json
	gpainfo := DataModel.GpaInfo{}
	err = json.Unmarshal([]byte(jsonText), &gpainfo)
	if err != nil {
		logger.Error("无法解析GPA JSON", zap.Time("time", time.Now()), zap.String("detail", err.Error()))
		return nil, fmt.Errorf("未知错误,错误id:%s", uids)
	}
	gparank := DataModel.GpaRank(gpainfo)
	return &gparank, nil
}
