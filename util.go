package tyut_osc

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"image"
	"io/ioutil"
	"log"
	"net/http"
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

func GenerateCaptchaGrainingSet() {
	resp, _ := http.Get("http://202.207.247.51:8065/validateCodeAction.do")
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	img = BinPic(img)
	imgbytes := Image2ByteArray(img)
	ioutil.WriteFile("temp.jpeg", imgbytes, 0644)
}
