package tyut_osc

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/linkinghack/tyut-osc/DataModel"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func Test_GpaMarshal(t *testing.T) {
	rawGpaJson := `{"xh":"2015005973","xm":"刘磊","bjh":"软件1516","bm":"软件1516","zyh":"160101","zym":"软件工程","xsh":"16","xsm":"软件学院","njdm":"2015","yqzxf":"188","yxzzsjxf":"8.32","zxf":"159.50","yxzxf":"159.50","cbjgxf":"0","sbjgxf":"0","pjxfjd":"3.80","gpabjpm":"4","gpazypm":"61","pjcj":"85.30","pjcjbjpm":"3","pjcjzypm":"69","jqxfcj":"84.93","jqbjpm":"4","jqzypm":"65","tsjqxfcj":"84.93","tjsj":"2019-01-17 01:00:04","bjrs":"30","zyrs":"968","dlrs":"","gpadlpm":"1148"}`
	gpa := DataModel.GpaInfo{}
	err := json.Unmarshal([]byte(rawGpaJson), &gpa)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	//decoder := json.NewDecoder(strings.NewReader(rawGpaJson))
	fmt.Println(gpa)
	fmt.Println(gpa.AvgScore)
}

func Test_JsonParse(t *testing.T) {
	gpaLoginStatus := `{"Result":false,"Code":1,"Msg":"登陆成功","Msg1":"/Default.aspx","Msg2":null,"Msg3":null,"Msg4":null,"Msg5":null}`
	bodyData := []byte(gpaLoginStatus)
	bodyJson := map[string]interface{}{}
	err := json.Unmarshal(bodyData, &bodyJson)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(bodyJson)
	fmt.Println(reflect.TypeOf(bodyJson["Code"]))     // float64
	fmt.Println(int(bodyJson["Code"].(float64)) == 1) //true
}

func TestGpaCrawler_GetGpaInfo(t *testing.T) {
	gpacrawler := NewGpaCrawler()
	text, err := gpacrawler.GetGpaRank("2015005973", "lolipop8974", "2015005973")
	if err != nil {
		t.Fail()
		panic(err)
	}
	fmt.Println(text)
}

func Test_HttpClient(t *testing.T) {
	formValues := url.Values{}
	formValues.Add("u", "2015005973")
	formValues.Add("p", "lolipop8974.")
	formValues.Add("r", "on")

	requestValue := "u=2015005973&p=lolipop8974.&r=on"
	request, _ := http.NewRequest("POST", "http://202.207.247.60/Hander/LoginAjax.ashx", strings.NewReader(requestValue))
	//request.Header.Set("User-Agent","Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:64.0) Gecko/20100101 Firefox/64.0")
	//request.Header.Set("Host","202.207.247.60")
	request.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	request.Header.Set("Accept-Encoding", "gzip, deflate")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	//request.Header.Set("Cookie","ASP.NET_SessionId=wuox4n44t1gvloxsodzyhyjk; ValidateCode=pfrag")
	//request.Header.Set("Cache-Control","max-age=0")

	//request.Header.Set("Connection","keep-alive")

	client := http.Client{}
	resp, _ := client.Do(request)
	fmt.Println(resp.Status)

	fmt.Println(resp.ContentLength)
	//resp,_ := http.PostForm("http://202.207.247.60/Hander/LoginAjax.ashx",formValues)

	//res,_ := http.Get("http://202.207.247.60")
	//fmt.Println(res)

	// 重要是这里，在首部添加Accept-Encoding 属性为gzip后，response中的body将不再自动进行gzip解码
	// 需要手动进行
	gzipdata, _ := gzip.NewReader(resp.Body)
	data, _ := ioutil.ReadAll(gzipdata)

	defer resp.Body.Close()

	fmt.Println(string(data))
	ioutil.WriteFile("GPALoginResponse.txt", data, 0644)
}
