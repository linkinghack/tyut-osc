package tyut_osc

import (
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type UrpCrawler struct {
	config *Configuration
}

func (u *UrpCrawler) SetConfiguration(conf *Configuration) {
	u.config = conf
}

func (urp *UrpCrawler) createClientAndLogin(stuid string, stuPassword string) (client *http.Client, activateUrlIdx int, err error) {
	uid, _ := uuid.NewUUID()
	uids := strings.Split(uid.String(), "-")[0]

	// 准备http.Client
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client = &http.Client{
		Jar:     jar,
		Timeout: time.Second * 10,
	}

	// 1. 探测可用的url并初始化cookie
	foundActive := false
	for i, v := range urp.config.BaseLocationURP {
		_, er := client.Get(v)
		if er == nil {
			foundActive = true
			activateUrlIdx = i
			break
		}
	}
	// 没有可以连通的urp系统url
	if !foundActive {
		client = nil
		activateUrlIdx = -1
		err = fmt.Errorf("连接教务系统超时, 错误id: %s", uids)
		return
	}

	// 2. 请求验证码
	getCaptcha, _ := http.NewRequest("GET", urp.config.BaseLocationURP[activateUrlIdx], nil)
	getCaptcha.Header.Set("random", string(rand.Int63()))
	captchaResp, er := client.Do(getCaptcha)
	if er != nil {
		logger.Warn("获取URP系统验证码出错", zap.String("stuid", stuid), zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
		client = nil
		err = fmt.Errorf("登录教务系统出错,错误id: %s", uids)
	}
	captchaPic, er := ioutil.ReadAll(captchaResp.Body)
	if er != nil {
		logger.Error("读取验证码错误", zap.Time("time", time.Now()))
		err = fmt.Errorf("登录过程异常, 错误id: %s", uids)
		client = nil
		return
	}
	ioutil.WriteFile("tempcaptcha.jpg", captchaPic, 0644)

	err = nil
	return
}
