package tyut_osc

import (
	"bytes"
	"github.com/otiai10/gosseract"
	"go.uber.org/atomic"
	"image"
	"image/color"
	"image/jpeg"
	"strings"
	"sync"
	"unicode"
)

type OcrEnginePool struct {
	active *atomic.Int32
	size   *atomic.Int32
	p      sync.Pool
}

func (p *OcrEnginePool) SetSize(size int32) {
	p.size.Store(size)
}

func (p *OcrEnginePool) Get() *gosseract.Client {
	client := p.p.Get().(*gosseract.Client)
	p.active.Add(1)
	return client
}

func (p *OcrEnginePool) Put(c *gosseract.Client) {
	if p.active.Load() > p.size.Load() {
		c.Close()
		c = nil //GC
	} else {
		p.active.Add(-1)
		p.p.Put(c)
	}
}

func NewOcrEnginePool(size int32, initialActive int32) *OcrEnginePool {
	if size < 0 {
		return nil
	}

	sp := sync.Pool{
		New: func() interface{} {
			client := gosseract.NewClient()
			client.Languages = []string{"rnd"}
			return client
		},
	}

	ocrpool := OcrEnginePool{
		active: atomic.NewInt32(0),
		size:   atomic.NewInt32(size),
		p:      sp,
	}

	// Create initial engines
	if initialActive > 0 {
		for i := 0; i < int(initialActive); i++ {
			ocrpool.Put(gosseract.NewClient())
		}
	}

	return &ocrpool
}

// 提供一些图像二值化处理的方法

// BinPic 将一个RGB图片转为黑白图片,
// 为tesseract识别图片优化.
// rawPic 应该为一个指针
func BinPic(rawPic image.Image) *(image.Gray) {
	bound := rawPic.Bounds()
	newgraypic := image.NewGray(bound)
	for i := 0; i < bound.Dx(); i++ {
		for j := 0; j < bound.Dy(); j++ {
			r, g, b, _ := rawPic.At(i, j).RGBA()

			var pointColor color.Color // 新颜色
			if r+g+b > 99999 {         // 优化的阈值
				pointColor = color.White
			} else {
				pointColor = color.Black
			}

			newgraypic.Set(i, j, color.RGBAModel.Convert(pointColor))
		}
	}
	return newgraypic
}

// img应该为一个指针
func Image2ByteArray(img image.Image) []byte {
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)
	return buf.Bytes()
}

func CaptchaTextFilt(rawtext string) string {
	chars := []rune(rawtext)
	count := 0
	var result strings.Builder
	for _, v := range chars {

		if unicode.IsDigit(v) || unicode.IsLetter(v) {
			result.WriteRune(v)
		}
		count++
		if count >= 4 {
			break
		}
	}

	return result.String()
}
