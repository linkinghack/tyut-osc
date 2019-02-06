package DataModel

type SelectedCourse struct {
	TrainingScheme       string
	CourseId             string
	CourseSequenceNumber string
	CourseName           string
	CourseCredit         float64
	SelectionType        string //选修-必修
	ExamType             string
	TeacherName          string
	WayOfStudy           string // 修读方式
	SelectionStatus      string // 选课状态 置入
	TimeLocs             []TimeLocation
}

type TimeLocation struct {
	Weeks []int
	//WeekContinuous	bool // 是否连续星期  true: len(Weeks) = 2 代表开始和结束周; false: Weeks中保存上课的周数,长度不定
	Day      int    // Day of week
	Start    int    // 起始节次
	Length   int    // 课程节数
	Campus   string // 校区
	Building string // 教学楼
	Room     string // 教室号
}
