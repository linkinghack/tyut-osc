package tyut_osc

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/linkinghack/tyut-osc/DataModel"
	"go.uber.org/zap"
	"golang.org/x/net/publicsuffix"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type UrpCrawler struct {
	config *Configuration
}

func (u *UrpCrawler) SetConfiguration(conf *Configuration) {
	u.config = conf
}

// NewUrpCraler 加载配置文件并创建一个UrpCrawler实例，返回指针
func NewUrpCrawler() *UrpCrawler {
	uc := UrpCrawler{
		config: loadConfigFromFile("/tyuter/configs/tyut-osc-CrawlerConfig.json"),
	}

	return &uc
}

func (urp *UrpCrawler) CreateClientAndLogin(stuid string, stuPassword string) (client *http.Client, activateUrlIdx int, err error) {
	uid, _ := uuid.NewUUID()
	uids := strings.Split(uid.String(), "-")[0]
	err = fmt.Errorf("登录过程异常,错误id: %s", uids)
	defer logger.Sync()

	// 准备http.Client

	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client = &http.Client{
		Jar:     jar,
		Timeout: time.Second * 2,
	}

	// 1. 探测可用的url并初始化cookie
	// TODO: 使用redis存储可用URL,避免每次都探测
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
	for attempt := 0; !ok && attempt < urp.config.UrpLoginAttempt; attempt++ {
		message := fmt.Sprintf("尝试第%d次登录", (attempt + 1))
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
	ok = false // 登录结果,返回值
	defer logger.Sync()

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
		//ioutil.WriteFile(uids+"loginpage.html", loginPageBytes, 0644)
		logger.Info("登录失败", zap.String("stuid", stuid), zap.Time("time", time.Now()), zap.String("detail", uids+"loginpage.html"))
		return false, nil
	}
}

func (urp *UrpCrawler) getCaptcha(stuid string, stuPassword string, client *http.Client, activateUrlIdx int) (captcha string, err error) {
	uid, _ := uuid.NewUUID()
	uids := strings.Split(uid.String(), "-")[0]
	err = fmt.Errorf("无法登录教务系统,错误id: %s", uids)
	defer logger.Sync()

	captcha = ""
	ocr := OcrPool.Get()
	defer OcrPool.Put(ocr)

	for len(captcha) != 4 {
		// 2. 请求验证码
		getCaptcha, _ := http.NewRequest("GET", urp.config.BaseLocationURP[activateUrlIdx]+"/validateCodeAction.do", nil)

		captchaResp, er := client.Do(getCaptcha)
		if er != nil {
			logger.Warn("获取URP系统验证码出错", zap.String("stuid", stuid), zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
			return
		}
		defer captchaResp.Body.Close()

		// 3. 图片二值化处理
		img, _, er := image.Decode(captchaResp.Body)
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
	}
	return captcha, nil
}

// GetPassedCourses 返回所有已出成绩学期列表,包含每学期成绩
func (urp *UrpCrawler) GetPassedCourses(client *http.Client, activateUrlIdx int) (terms []DataModel.Term, err error) {
	uid, _ := uuid.NewUUID()
	uids := strings.Split(uid.String(), "-")[0]
	err = fmt.Errorf("无法获取已通过成绩,错误id: %s", uids)
	defer logger.Sync()

	resp, er := client.Get(urp.config.BaseLocationURP[activateUrlIdx] + "/gradeLnAllAction.do?type=ln&oper=qbinfo&lnxndm=")
	if er != nil {
		logger.Warn("无法请求学期成绩页面", zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	passedCoursesHtmlBytes, _ := ioutil.ReadAll(resp.Body)
	passedCoursesHtmlBytes, er = DecodeGBK(passedCoursesHtmlBytes)
	if er != nil {
		logger.Warn("GBK编码转换错误", zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", err.Error()))
		return nil, err
	}

	//ioutil.WriteFile("grade.html",passedCoursesHtmlBytes,0644)

	// 分析html结构
	doc, er := goquery.NewDocumentFromReader(bytes.NewReader(passedCoursesHtmlBytes))
	if er != nil {
		logger.Warn("goquery无法解析成绩页面", zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", err.Error()))
		return nil, err
	}

	// 1. 学期列表
	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		term := DataModel.Term{}
		termname, _ := selection.Attr("name") //2015-2016学年秋(两学期)
		term.TermDescription = termname

		temprunes := []rune(termname)
		if temprunes[11] == '秋' {
			term.TermOrder = 1
			term.TermYear, _ = strconv.Atoi(string(temprunes[0:4]))
		} else if temprunes[11] == '春' {
			term.TermOrder = 0
			term.TermYear, _ = strconv.Atoi(string(temprunes[5:9]))
		}

		terms = append(terms, term)
	})

	// 2. 课程成绩列表
	// 每个.titleTop2 对应一个学期
	doc.Find(".titleTop2").Each(func(i int, termhtml *goquery.Selection) {
		passedcourses := []DataModel.PassedCourse{}

		// 每条课程成绩信息存在一个.odd中
		termhtml.Find(".odd").Each(func(i int, coursehtml *goquery.Selection) {
			course := DataModel.PassedCourse{}

			coursehtml.Find("td").Each(func(tdidx int, field *goquery.Selection) {

				switch tdidx {
				case 0:
					{
						course.CourseId = strings.TrimSpace(field.Text())
					}
				case 1:
					{
						course.CourseSequenceNumber = strings.TrimSpace(field.Text())
					}
				case 2:
					{
						course.CourseName = strings.TrimSpace(field.Text())
					}
				case 3:
					{
						course.EnglishCourseName = strings.TrimSpace(field.Text())
					}
				case 4:
					{
						course.CourseCredit, _ = strconv.ParseFloat(strings.TrimSpace(field.Text()), 64)
					}
				case 5:
					{
						course.SelectionProperty = strings.TrimSpace(field.Text())
					}
				case 6:
					{
						// 分中文成绩和数字成绩两种解决
						scoreStr := strings.TrimSpace(field.Text())
						reg, _ := regexp.Compile(`[[^(0-9)+.(0-9)+]]`)
						scoreFiltered := reg.ReplaceAllString(scoreStr, "")
						if len(scoreFiltered) > 0 {
							course.Score, _ = strconv.ParseFloat(scoreStr, 64)
						} else {
							course.ChScore = scoreStr
							course.Score = -1 //标记使用中文成绩
						}
					}
				}
			})

			passedcourses = append(passedcourses, course)
		})

		terms[i].PassedCourses = passedcourses
	})

	//fmt.Println(terms)

	return terms, nil
}

/**
GetFailedCourses 返回挂科成绩列表,包含曾挂科和现挂科
*/
func (urp *UrpCrawler) GetFailedCourses(client *http.Client, activateUrlIdx int) (fcourses []DataModel.FailedCourse, err error) {
	uid, _ := uuid.NewUUID()
	uids := strings.Split(uid.String(), "-")[0]
	err = fmt.Errorf("无法获取不及格科目,错误id: %s", uids)
	defer logger.Sync()

	resp, er := client.Get(urp.config.BaseLocationURP[activateUrlIdx] + "/gradeLnAllAction.do?type=ln&oper=bjg")
	if er != nil {
		logger.Warn("无法请求不及格成绩页面", zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	failedCourseHtmlBytes, er := ioutil.ReadAll(resp.Body)

	if er != nil {
		logger.Warn("不及格成绩页面响应读取错误", zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
		return nil, err
	}
	failedCourseHtmlBytes, er = DecodeGBK(failedCourseHtmlBytes)
	if er != nil {
		logger.Warn("GBK解码错误", zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
		return nil, err
	}

	// goquery解析页面
	doc, er := goquery.NewDocumentFromReader(bytes.NewReader(failedCourseHtmlBytes))
	if er != nil {
		logger.Warn("不及格页面goquery解析错误", zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
		return nil, err
	}

	// 分析页面
	doc.Find(".displayTag").Each(func(dispIdx int, dispTag *goquery.Selection) {
		dispTag.Find(".odd").Each(func(i int, fc *goquery.Selection) {
			failed := DataModel.FailedCourse{}
			fc.Find("td").Each(func(i int, failedCourseProperty *goquery.Selection) {
				switch i {
				case 0:
					failed.CourseId = strings.TrimSpace(failedCourseProperty.Text())
				case 1:
					failed.CourseSequenceNumber = strings.TrimSpace(failedCourseProperty.Text())
				case 2:
					failed.CourseName = strings.TrimSpace(failedCourseProperty.Text())
				case 3:
					failed.EnglishCourseName = strings.TrimSpace(failedCourseProperty.Text())
				case 4:
					{
						failed.CourseCredit, _ = strconv.ParseFloat(strings.TrimSpace(failedCourseProperty.Text()), 64)
					}
				case 5:
					failed.SelectionProperty = strings.TrimSpace(failedCourseProperty.Text())
				case 6:
					failed.Score, _ = strconv.ParseFloat(strings.TrimSpace(failedCourseProperty.Text()), 64)
				case 7:
					failed.ExamTime, _ = time.Parse("20060102", strings.TrimSpace(failedCourseProperty.Text()))
				case 8:
					failed.Reason = strings.TrimSpace(failedCourseProperty.Text())

				} //switch
			})

			// 判断是否仍在挂
			if dispIdx == 0 {
				failed.StillFail = true
			} else {
				failed.StillFail = false
			}

			fcourses = append(fcourses, failed) // 加入到最终返回列表
		})

	})
	return fcourses, nil
}

// GetCourseList 获取课程表页面并解析
func (urp *UrpCrawler) GetCourseList(client *http.Client, activeUrlIdx int) (seletects []DataModel.SelectedCourse, err error) {
	uid, _ := uuid.NewUUID()
	uids := strings.Split(uid.String(), "-")[0]
	err = fmt.Errorf("无法获取课程表信息,错误id: %s", uids)
	defer logger.Sync()

	resp, er := client.Get(urp.config.BaseLocationURP[activeUrlIdx] + "/xkAction.do?actionType=6")
	if er != nil {
		logger.Warn("请求课程表页面出错", zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	// 解码GBK
	pageHtmlBytes, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		logger.Warn("课程表页面响应读取错误", zap.String("errid", uids), zap.Time("time", time.Now()), zap.String("detail", er.Error()))
		return nil, err
	}
	pageDecodedBytes, _ := DecodeGBK(pageHtmlBytes)

	var tmpCourse *DataModel.SelectedCourse
	// 解析页面
	doc, er := goquery.NewDocumentFromReader(bytes.NewReader(pageDecodedBytes))
	doc.Find(".odd").Each(func(i int, coursesTR *goquery.Selection) {
		if coursesTR.Find("td").Length() > 7 {
			tmpCourse = new(DataModel.SelectedCourse)
			var timeloc = DataModel.TimeLocation{} // 第一个上课时间地点安排

			coursesTR.Find("td").Each(func(i int, courseProperty *goquery.Selection) {
				switch i {
				case 0:
					tmpCourse.TrainingScheme = strings.TrimSpace(courseProperty.Text())
				case 1:
					tmpCourse.CourseId = strings.TrimSpace(courseProperty.Text())
				case 2:
					tmpCourse.CourseName = strings.TrimSpace(courseProperty.Text())
				case 3:
					tmpCourse.CourseSequenceNumber = strings.TrimSpace(courseProperty.Text())
				case 4:
					tmpCourse.CourseCredit, _ = strconv.ParseFloat(strings.TrimSpace(courseProperty.Text()), 64)
				case 5:
					tmpCourse.SelectionType = strings.TrimSpace(courseProperty.Text())
				case 6:
					tmpCourse.ExamType = strings.TrimSpace(courseProperty.Text())
				case 7:
					tmpCourse.TeacherName = strings.TrimSpace(courseProperty.Text())
				case 9:
					tmpCourse.WayOfStudy = strings.TrimSpace(courseProperty.Text())
				case 10:
					tmpCourse.SelectionStatus = strings.TrimSpace(courseProperty.Text())
				case 11:
					timeloc.Weeks = ParseCourseWeeks(strings.TrimSpace(courseProperty.Text()))
				case 12:
					timeloc.Day, _ = strconv.Atoi(strings.TrimSpace(courseProperty.Text()))
				case 13:
					timeloc.Start = ParseCourseStartTime(strings.TrimSpace(courseProperty.Text()))
				case 14:
					timeloc.Length, _ = strconv.Atoi(strings.TrimSpace(courseProperty.Text()))
				case 15:
					timeloc.Campus = strings.TrimSpace(courseProperty.Text())
				case 16:
					timeloc.Building = strings.TrimSpace(courseProperty.Text())
				case 17:
					timeloc.Room = strings.TrimSpace(courseProperty.Text())
				}
			})
			tmpCourse.TimeLocs = append(tmpCourse.TimeLocs, timeloc) // 添加第一个上课时间地点

		} else { // td < 7
			fmt.Println("td counts less than 7")
			// 上一个课程行的附加上课时间地点信息
			otherTimeLoc := DataModel.TimeLocation{}
			coursesTR.Find("td").Each(func(i int, courseProperty *goquery.Selection) {
				switch i {
				case 0:
					otherTimeLoc.Weeks = ParseCourseWeeks(strings.TrimSpace(courseProperty.Text()))
				case 1:
					otherTimeLoc.Day, _ = strconv.Atoi(strings.TrimSpace(courseProperty.Text()))
				case 2:
					otherTimeLoc.Start = ParseCourseStartTime(strings.TrimSpace(courseProperty.Text()))
				case 3:
					otherTimeLoc.Length, _ = strconv.Atoi(strings.TrimSpace(courseProperty.Text()))
				case 4:
					otherTimeLoc.Campus = strings.TrimSpace(courseProperty.Text())
				case 5:
					otherTimeLoc.Building = strings.TrimSpace(courseProperty.Text())
				case 6:
					otherTimeLoc.Room = strings.TrimSpace(courseProperty.Text())
				}
			})
			fmt.Println("timeloc: ", otherTimeLoc)
			tmpCourse.TimeLocs = append(tmpCourse.TimeLocs, otherTimeLoc)
		}

		// 下一个tr包含课程名信息或已经是最后一个节点，则已经准备好一个Course
		if coursesTR.Next().Find("td").Length() > 7 || coursesTR.Next().Size() < 1 {
			seletects = append(seletects, *tmpCourse)
		}

	})

	return seletects, nil
}
