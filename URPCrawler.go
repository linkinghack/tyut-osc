package tyut_osc

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/net/publicsuffix"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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
	err = fmt.Errorf("登录过程异常,错误id: %s", uids)

	// 准备http.Client
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client = &http.Client{
		Jar:     jar,
		Timeout: time.Second * 2,
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
		fmt.Println(er.Error())
	}
	// 没有可以连通的urp系统url
	if !foundActive {
		client = nil
		activateUrlIdx = -1
		logger.Warn("无法连接教务系统,可能原因:超时", zap.String("stuid", stuid), zap.Time("time", time.Now()), zap.String("errid", uids))
		err = fmt.Errorf("连接教务系统超时,错误id: %s", uids)
		return
	}

	ok := false
	for i := 0; !ok && i < 5; i++ {
		message := fmt.Sprintf("尝试第%d次登录", (i + 1))
		logger.Info(message, zap.String("stuid", stuid), zap.Time("time", time.Now()))

		// 2. 获取验证码
		captcha, er := urp.getCaptcha(stuid, stuPassword, client, activateUrlIdx)
		if er != nil {
			err = er
			client = nil
			return
		}

		// 3. 登录
		ok, er = urp.login(stuid, stuPassword, captcha, client, activateUrlIdx)
		if er != nil {
			err = er
			client = nil
			return
		}
	}

	// 多次尝试登录失败
	if !ok {
		logger.Info("登录失败", zap.String("stuid", stuid), zap.Time("time", time.Now()))
		err = fmt.Errorf("登录失败,教务系统密码是否正确？")
		client = nil
		return
	}

	err = nil
	return
}

// login 尝试一次登录,并返回登录结果
func (urp *UrpCrawler) login(stuid string, stuPassword string, captcha string, client *http.Client, activateUrlIdx int) (ok bool, err error) {
	uid, _ := uuid.NewUUID()
	uids := strings.Split(uid.String(), "-")[0]
	err = fmt.Errorf("登录过程异常,错误id: %s", uids)
	ok = false

	// 5. 登录
	loginformvalues := url.Values{}
	loginformvalues.Set("zjh", stuid)
	loginformvalues.Set("mm", stuPassword)
	loginformvalues.Set("v_yzm", captcha)
	loginResp, er := client.PostForm(urp.config.BaseLocationURP[activateUrlIdx]+"/loginAction.do", loginformvalues)
	if er != nil {
		logger.Warn("登录请求错误", zap.String("errid", uids), zap.String("stuid", stuid), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
		return
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != 200 {
		logger.Info("登录请求异常", zap.String("errid", uids), zap.String("stuid", stuid), zap.Time("time", time.Now()))
		return
	}

	loginPageBytes, _ := ioutil.ReadAll(loginResp.Body)
	loginPageBytes, er = DecodeGBK(loginPageBytes)
	if er != nil {
		logger.Warn("gbk编码转换错误", zap.String("errid", uids), zap.String("stuid", stuid), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
		return
	}
	ioutil.WriteFile("loginpage.html", loginPageBytes, 0644)

	// 分析登录结果
	doc, er := goquery.NewDocumentFromReader(bytes.NewReader(loginPageBytes))
	if er != nil {
		logger.Warn("Cannot parse the body with goquery", zap.String("detail", er.Error()))
		err = er
		return
	}

	framenodes := doc.Find("frame").Nodes
	if len(framenodes) > 0 {
		logger.Info("URP登录成功", zap.String("stuid", stuid), zap.Time("time", time.Now()))
		return true, nil
	} else {
		return false, nil
	}

}

func (urp *UrpCrawler) getCaptcha(stuid string, stuPassword string, client *http.Client, activateUrlIdx int) (captcha string, err error) {
	uid, _ := uuid.NewUUID()
	uids := strings.Split(uid.String(), "-")[0]
	err = fmt.Errorf("无法登录教务系统,错误id: %s", uids)

	captcha = ""
	ocr := OcrPool.Get()
	defer OcrPool.Put(ocr)

	for len(captcha) != 4 {
		//for len(captcha) != 4 && repeatCount < 5{  //尝试五次验证码
		// 2. 请求验证码
		getCaptcha, _ := http.NewRequest("GET", urp.config.BaseLocationURP[activateUrlIdx]+"/validateCodeAction.do?random=0.7616917022636875", nil)
		getCaptcha.Header.Set("random", string(rand.Int63()))
		//getCaptcha.Header.Set("Accept-Encoding", "gzip, deflate")
		captchaResp, er := client.Do(getCaptcha)
		if er != nil {
			logger.Warn("获取URP系统验证码出错", zap.String("stuid", stuid), zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
			return
		}
		defer captchaResp.Body.Close()

		captchaPic, er := ioutil.ReadAll(captchaResp.Body)
		if er != nil {
			logger.Error("读取验证码错误", zap.String("stuid", stuid), zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
			return
		}

		// captcha pic temp file
		tmpfile, er := ioutil.TempFile(urp.config.TempDir, "*captcha.jpeg")
		tmpfile.Write(captchaPic)
		tmpfile.Close()
		//defer os.Remove(tmpfile.Name())

		// 3. 图片二值化处理
		capfile, _ := os.Open(tmpfile.Name())
		defer capfile.Close()
		img, _, er := image.Decode(capfile)
		if er != nil {
			logger.Warn("captcha二值化:图片读取失败", zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
			return
		}
		imgbytes := Image2ByteArray(BinPic(img))

		// 4. 识别captcha图片
		er = ocr.SetImageFromBytes(imgbytes)
		if er != nil {
			logger.Warn("OCR engine 识别出错", zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
			return
		}

		captcha, er = ocr.Text()
		if er != nil {
			logger.Warn("OCR engine 识别出错", zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
			return
		}
		captcha = CaptchaTextFilt(captcha)

		fmt.Println("captcha: ", captcha, "len: ", len(captcha), " file: ", capfile.Name())
	}

	return captcha, nil
}
