package tyut_osc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func Test_GpaMarshal(t *testing.T) {
	rawGpaJson := `{"xh":"2015005973","xm":"刘磊","bjh":"软件1516","bm":"软件1516","zyh":"160101","zym":"软件工程","xsh":"16","xsm":"软件学院","njdm":"2015","yqzxf":"188","yxzzsjxf":"8.32","zxf":"159.50","yxzxf":"159.50","cbjgxf":"0","sbjgxf":"0","pjxfjd":"3.80","gpabjpm":"4","gpazypm":"61","pjcj":"85.30","pjcjbjpm":"3","pjcjzypm":"69","jqxfcj":"84.93","jqbjpm":"4","jqzypm":"65","tsjqxfcj":"84.93","tjsj":"2019-01-17 01:00:04","bjrs":"30","zyrs":"968","dlrs":"","gpadlpm":"1148"}`
	gpa := &GpaInfo{}
	err := json.Unmarshal([]byte(rawGpaJson), gpa)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	//decoder := json.NewDecoder(strings.NewReader(rawGpaJson))
	fmt.Println(gpa)
	fmt.Println((*gpa).AvgScore)
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

/*
func Test_Logger(t *testing.T) {
	Logger.Info("works well",zap.Time("time",time.Now()))
}*/
