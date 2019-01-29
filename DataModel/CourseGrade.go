package DataModel

type PassedCourse struct {
	Id                   string
	CourseId             string
	CourseSequenceNumber int
	CourseName           string
	EnglishCourseName    string
	CourseCredit         float64
	SelectionProperty    string
	Score                float64
	Semester             string
	ChScore              string //中文成绩
}

type FailedCourse struct {
	CourseId             string  //课程号
	CourseSequenceNumber int     //课序号
	CourseName           string  // 课程名
	EnglishCourseName    string  //英文课程名
	CourseCredit         float64 //课程学分
	SelectionProperty    string  //课程属性 - 综合必修
	Score                float64 //得分
	ExamTime             string  //20170626
	StillFail            bool
}

type Term struct {
	TermDescription string
	TermYear        int
	TermOrder       int // 0-春季学期, 1-秋季学期
	PassedCourses   []PassedCourse
}
