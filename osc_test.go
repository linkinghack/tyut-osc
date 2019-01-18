package tyut_osc

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_GpaMarshal(t *testing.T) {
	rawGpaJson := `{"xh":"2015005973","xm":"刘磊","bjh":"软件1516","bm":"软件1516","zyh":"160101","zym":"软件工程","xsh":"16","xsm":"软件学院","njdm":"2015","yqzxf":"188","yxzzsjxf":"8.32","zxf":"159.50","yxzxf":"159.50","cbjgxf":"0","sbjgxf":"0","pjxfjd":"3.80","gpabjpm":"4","gpazypm":"61","pjcj":"85.30","pjcjbjpm":"3","pjcjzypm":"69","jqxfcj":"84.93","jqbjpm":"4","jqzypm":"65","tsjqxfcj":"84.93","tjsj":"2019-01-17 01:00:04","bjrs":"30","zyrs":"968","dlrs":"","gpadlpm":"1148"}`
	gpa := &GpaInfo{}
	err := json.Unmarshal([]byte(rawGpaJson),gpa)
	if(err != nil){
		fmt.Println(err)
		t.Fail()
	}
	//decoder := json.NewDecoder(strings.NewReader(rawGpaJson))
	fmt.Println(gpa)
	fmt.Println((*gpa).AvgScore)
}
