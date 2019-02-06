package DataModel

import "time"

type PassedCourse struct {
	CourseId             string
	CourseSequenceNumber string
	CourseName           string
	EnglishCourseName    string
	CourseCredit         float64
	SelectionProperty    string
	Score                float64
	Semester             string
	ChScore              string //中文成绩
}

type FailedCourse struct {
	CourseId             string    //课程号
	CourseSequenceNumber string    //课序号
	CourseName           string    // 课程名
	EnglishCourseName    string    //英文课程名
	CourseCredit         float64   //课程学分
	SelectionProperty    string    //课程属性 - 综合必修
	Score                float64   //得分
	ExamTime             time.Time //20170626
	StillFail            bool
	Reason               string
}

type Term struct {
	TermDescription string
	TermYear        int
	TermOrder       int // 0-春季学期, 1-秋季学期
	PassedCourses   []PassedCourse
}

type Wrapper struct {
	PassedCourse
	Id string
}
