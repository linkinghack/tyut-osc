package tyut_osc

import (
	"fmt"
	"go.uber.org/zap"
	"testing"
	"time"
)

func Test_Channel(t *testing.T) {
	resource := make(chan int, 20)
	resource <- 2
	fmt.Println(len(resource))
}

func Test_UrpLogin(t *testing.T) {

	for i := 0; i < 80; i++ {
		start := time.Now()
		urpcrawler := NewUrpCrawler()
		_, idx, err := urpcrawler.CreateClientAndLogin("2015005973", "lolipop8974.")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("avaliable idx: ", idx)

		end := time.Now()
		fmt.Println("耗时:", end.UnixNano()-start.UnixNano())
	}
}

func Test_PassedCourses(t *testing.T) {
	urp := NewUrpCrawler()
	client, idx, _ := urp.CreateClientAndLogin("2015005968", "304515")
	courses, _ := urp.GetPassedCourses(client, idx)
	fmt.Println(courses)
}

func Test_FailedCourses(t *testing.T) {
	urp := NewUrpCrawler()
	client, idx, _ := urp.CreateClientAndLogin("2015005968", "304515")
	fcourse, err := urp.GetFailedCourses(client, idx)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	} else {
		fmt.Println(fcourse)
	}
}

func Test_CourseList(t *testing.T) {
	urp := NewUrpCrawler()
	client, idx, _ := urp.CreateClientAndLogin("2016006359", "110628")
	list, err := urp.GetCourseList(client, idx)
	if err != nil {
		logger.Warn("课程列表获取失败", zap.String("detail", err.Error()))
		t.Fail()
	}
	fmt.Println(list)
}
