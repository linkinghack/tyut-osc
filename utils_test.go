package tyut_osc

import (
	"fmt"
	"github.com/google/uuid"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func Test_HttpClient_Basic(t *testing.T) {
	resp, _ := http.Get("http://baidu.com")
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.Header)
	fmt.Println(string(data))
}

// google uuid
func Test_UUID(t *testing.T) {
	for i := 0; i < 100; i++ {
		uidd, _ := uuid.NewRandom()
		fmt.Println(uidd)
	}
}

func Test_Logger(t *testing.T) {
	logger.Info("is caller ok?")
}

func Test_PicHandle(t *testing.T) {
	fi, er := os.Open("temp2.jpeg")
	defer fi.Close()

	if er != nil {
		return
	}

	img, _, _ := image.Decode(fi)
	bound := img.Bounds()
	newgraypic := image.NewGray(bound)
	for i := 0; i < bound.Dx(); i++ {
		for j := 0; j < bound.Dy(); j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			//var nR, nG, nB uint8
			/*nR = uint8((float64(r) / 65535.0) * 255.0)
			nG = uint8((float64(g) / 65535.0) * 255.0)
			nB = uint8((float64(b) / 65535.0) * 255.0)
			*/
			var pointColor color.Color

			if r+g+b > 98888 {
				pointColor = color.White
			} else {
				pointColor = color.Black
			}

			newgraypic.Set(i, j, color.RGBAModel.Convert(pointColor))
		}
	}

	newfi, err := os.OpenFile("binarycaptcha.jpeg", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("open file error")
	}
	jpeg.Encode(newfi, newgraypic, nil)
	newfi.Close()

}

func TestFloatCount(t *testing.T) {
	fmt.Println(23 / 45.0)
}
