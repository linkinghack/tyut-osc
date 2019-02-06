package tyut_osc

import (
	"bytes"
	"github.com/google/uuid"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// DecodeGBK 接受GBK编码的字节数组，返回转换后的UTF-8编码数组
func DecodeGBK(src []byte) ([]byte, error) {
	in := bytes.NewReader(src)
	out := transform.NewReader(in, simplifiedchinese.GBK.NewDecoder())
	decoded, err := ioutil.ReadAll(out)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func GenerateCaptchaGrainingSet(n int) {

	for i := 0; i < n; i++ {
		uu, _ := uuid.NewRandom()
		uis := strings.Split(uu.String(), "-")[0]
		filename := "/home/linking/captchapics/" + uis + ".jpeg"
		resp, _ := http.Get("http://202.207.247.51:8065/validateCodeAction.do")
		img, _, err := image.Decode(resp.Body)
		if err != nil {
			log.Fatal(err.Error())
		}
		img = BinPic(img)
		imgbytes := Image2ByteArray(img)
		ioutil.WriteFile(filename, imgbytes, 0644)
	}

}

// ParseCourseWeeks 解析教务系统课程安排中的周数表达,返回所有需要上课的周次
//
// * 周次表达式是一个字符串，类似于  1-3周上 |  6-10,12-16周上 | 2,4,6,8周上
func ParseCourseWeeks(weekstr string) (weeks []int) {
	reg, _ := regexp.Compile(`[^\d-,]`)
	tmpstr := reg.ReplaceAllString(weekstr, "")
	weekpart := strings.Split(tmpstr, ",")
	for _, part := range weekpart {
		if strings.Index(part, "-") != -1 {
			fe := strings.Split(part, "-")
			from, _ := strconv.Atoi(fe[0])
			to, _ := strconv.Atoi(fe[1])
			for ; from <= to; from++ {
				weeks = append(weeks, from)
			}
		} else {
			w, _ := strconv.Atoi(part)
			weeks = append(weeks, w)
		}
	}
	return
}

// ParseCourseStartTime 解析课程起始节次,返回整数起始节次
//
// @param strexp 原生表达  '三小'
// @return int 3 | -1 cannot parse
func ParseCourseStartTime(strexp string) int {
	switch strexp {
	case "一小":
		return 1
	case "二小":
		return 2
	case "三小":
		return 3
	case "四小":
		return 4
	case "五小":
		return 5
	case "六小":
		return 6
	case "七小":
		return 7
	case "八小":
		return 8
	case "九小":
		return 9
	case "十小":
		return 10
	case "十一小":
		return 11
	case "十二小":
		return 12
	}
	return 0
}
