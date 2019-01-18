package tyut_osc

import "github.com/otiai10/gosseract"

var ocr *gosseract.Client

func init() {
	ocr = gosseract.NewClient()
	ocr.Languages = []string{"rnd"}
}
