package DataModel

// GpaInfo 代表从GPA查询系统中查询GPA排名获取到的数据结构
// GPA 原始信息中是此对象的数组,注意处理
type GpaInfo struct {
	StudentId     string `json:"xh"`          //xh 学号
	Name          string `json:"xm"`          //xm 姓名
	ClassId       string `json:"bjh"`         //bjh 班级号
	ClassName     string `json:"bm"`          //bm 班级名
	MajorId       string `json:"zyh"`         //zyh 专业号
	MajorName     string `json:"zym"`         //zym 专业名
	CollegeId     string `json:"xsh"`         //xsh 学院号
	CollegeName   string `json:"xsm"`         //xsm  学院名
	Grade         int    `json:"njdm,string"` //njdm 年级代码
	NumOfClassStu int    `json:"bjrs,string"` //bjrs 班级人数
	NumOfMajorStu int    `json:"zyrs,string"` //zyrs 专业人数

	RequiredCredit   float64 `json:"yqzxf,string"`    //yqzxf 要求总学分
	ActivityCredit   float64 `json："yxzzsjxf,string"` //yxzzsjxf 已修自主实践学分
	TotalCredit      float64 `json:"zxf,string"`      //zxf 总学分
	GottedCredit     float64 `json:"yxzxf,string"`    //yxzxf 已修总学分
	OnceFailedCredit float64 `json:"cbjgxf,string"`   //cbjgxf 曾不及格学分
	NowFailedCredit  float64 `json:"sbjgxf,string"`   //sbjgxf 尚不及格学分

	Gpa              float64 `json:"pjxfjd,string"`  //pjxfjd 平均学分绩点 GPA
	GpaClassRank     int     `json:"gpabjpm,string"` //gpabjpm GPA班级排名
	GpaMajorRank     int     `json:"gpazypm,string"` //gpazypm GPA专业排名
	GpaMajorCateRank int     `json:"gpadlpm,string"` //gpadlpm GPA大类排名
	//NumOfMajorCate int `json:"dlrs,string"`  // 大类人数

	AvgScore          float64 `json:"pjcj,string"`     //pjcj 平均成绩
	AvgScoreClassRank int     `json:"pjcjbjpm,string"` //pjcjbjpm 平均成绩班级排名
	AvgScoreMajorRank int     `json:"pjcjzypm,string"` //pjcjzypm 平均成绩专业排名
	WeightedScore     float64 `json:"jqxfcj,string"`   //jqxfcj 加权学分成绩
	WeightedClassRank int     `json:"jqbjpm,string"`   //jqbjpm 加权班级排名
	WeightedMajorRank int     `json:"jqzypm,string"`   //jqzypm 加权专业排名

	TsWeightedScore float64 `json:"tsjqxfcj,string"` //tsjqxfcj
	CountDate       string  `json:"tjsj"`            //tjsj 统计时间 %Y-%m-%d %H:%i:%S
}

/**
[{"xh":"2015005973","xm":"刘磊","bjh":"软件1516","bm":"软件1516","zyh":"160101","zym":"软件工程","xsh":"16","xsm":"软件学院","njdm":"2015","yqzxf":"188","yxzzsjxf":"8.32","zxf":"159.50","yxzxf":"159.50","cbjgxf":"0","sbjgxf":"0","pjxfjd":"3.80","gpabjpm":"4","gpazypm":"61","pjcj":"85.30","pjcjbjpm":"3","pjcjzypm":"69","jqxfcj":"84.93","jqbjpm":"4","jqzypm":"65","tsjqxfcj":"84.93","tjsj":"2019-01-17 01:00:04","bjrs":"30","zyrs":"968","dlrs":"","gpadlpm":"1148"}]
*/

// GpaDetail 与GpaInfo表达的信息一样,用于输出,其json tag已消除
// 在应用中可以直接由GpaInfo强制类型转换
type GpaDetail struct {
	StudentId     string //xh 学号
	Name          string //xm 姓名
	ClassId       string //bjh 班级号
	ClassName     string //bm 班级名
	MajorId       string //zyh 专业号
	MajorName     string //zym 专业名
	CollegeId     string //xsh 学院号
	CollegeName   string //xsm  学院名
	Grade         int    //njdm 年级代码
	NumOfClassStu int    //bjrs 班级人数
	NumOfMajorStu int    //zyrs 专业人数

	RequiredCredit   float64 //yqzxf 要求总学分
	ActivityCredit   float64 //yxzzsjxf 已修自主实践学分
	TotalCredit      float64 //zxf 总学分
	GottedCredit     float64 //yxzxf 已修总学分
	OnceFailedCredit float64 //cbjgxf 曾不及格学分
	NowFailedCredit  float64 //sbjgxf 尚不及格学分

	Gpa              float64 //pjxfjd 平均学分绩点 GPA
	GpaClassRank     int     //gpabjpm GPA班级排名
	GpaMajorRank     int     //gpazypm GPA专业排名
	GpaMajorCateRank int     //gpadlpm GPA大类排名
	//NumOfMajorCate int   // 大类人数

	AvgScore          float64 //pjcj 平均成绩
	AvgScoreClassRank int     //pjcjbjpm 平均成绩班级排名
	AvgScoreMajorRank int     //pjcjzypm 平均成绩专业排名
	WeightedScore     float64 //jqxfcj 加权学分成绩
	WeightedClassRank int     //jqbjpm 加权班级排名
	WeightedMajorRank int     //jqzypm 加权专业排名

	TsWeightedScore float64 //tsjqxfcj
	CountDate       string  //tjsj 统计时间 %Y-%m-%d %H:%i:%S

}
