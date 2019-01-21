package tyut_osc

import (
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
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
