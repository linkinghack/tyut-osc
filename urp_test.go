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
	start := time.Now()

	urpcrawler := NewUrpCrawler()
	_, idx, err := urpcrawler.createClientAndLogin("2015005973", "lolipop8974.")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("avaliable idx: ", idx)

	end := time.Now()
	fmt.Println("耗时:",end.UnixNano() - start.UnixNano() )
}
