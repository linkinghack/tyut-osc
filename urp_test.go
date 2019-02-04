package tyut_osc

import (
	"fmt"
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
	client, idx, _ := urp.CreateClientAndLogin("2015005973", "lolipop8974.")
	urp.GetPassedCourses(client, idx)
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
