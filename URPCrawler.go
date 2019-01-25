package tyut_osc

import (
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/net/publicsuffix"
	_ "image/jpeg"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"
)

type UrpCrawler struct {
	config *Configuration
}

func (u *UrpCrawler) SetConfiguration(conf *Configuration) {
	u.config = conf
}

func NewUrpCrawler() *UrpCrawler {
	uc := UrpCrawler{
		config: loadConfigFromFile("config.json"),
	}

	return &uc
}

func (urp *UrpCrawler) createClientAndLogin(stuid string, stuPassword string) (client *http.Client, activateUrlIdx int, err error) {
	uid, _ := uuid.NewUUID()
	uids := strings.Split(uid.String(), "-")[0]

	// 准备http.Client
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client = &http.Client{
		Jar:     jar,
		Timeout: time.Second * 2,
	}

	// 1. 探测可用的url并初始化cookie
	foundActive := false
	var er error
	for i, v := range urp.config.BaseLocationURP {
		_, er := client.Get(v)
		if er == nil {
			foundActive = true
			activateUrlIdx = i
			break
		}
		fmt.Println(er.Error())
	}
	// 没有可以连通的urp系统url
	if !foundActive {
		client = nil
		activateUrlIdx = -1
		logger.Warn("无法连接教务系统,可能原因:超时", zap.String("stuid", stuid), zap.Time("time", time.Now()), zap.String("errid", uids), zap.String("detail", er.Error()))
		err = fmt.Errorf("连接教务系统超时, 错误id: %s", uids)
		return
	}

	// 2. 请求验证码
	getCaptcha, _ := http.NewRequest("GET", urp.config.BaseLocationURP[activateUrlIdx]+"/validateCodeAction.do?random=0.7616917022636875", nil)
	getCaptcha.Header.Set("random", string(rand.Int63()))
	getCaptcha.Header.Set("Accept-Encoding", "gzip, deflate")
	captchaResp, er := client.Do(getCaptcha)
	if er != nil {
		logger.Warn("获取URP系统验证码出错", zap.String("stuid", stuid), zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
		client = nil
		err = fmt.Errorf("登录教务系统出错,错误id: %s", uids)
		return
	}
	defer captchaResp.Body.Close()

	captchaPic, er := ioutil.ReadAll(captchaResp.Body)
	if er != nil {
		logger.Error("读取验证码错误", zap.Time("time", time.Now()))
		err = fmt.Errorf("登录过程异常, 错误id: %s", uids)
		client = nil
		return
	}

	// captcha pic temp file
	tmpfile, er := ioutil.TempFile(urp.config.TempDir, "captcha")
	tmpfile.Write(captchaPic)
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	// 3. 图片二值化处理

	err = nil
	return
}
