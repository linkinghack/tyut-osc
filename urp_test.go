package tyut_osc

import (
	"fmt"
	"testing"
)

func Test_Channel(t *testing.T) {
	resource := make(chan int, 20)
	resource <- 2
	fmt.Println(len(resource))
}

func Test_UrpLogin(t *testing.T) {
	urpcrawler := NewUrpCrawler()
	_, idx, err := urpcrawler.createClientAndLogin("2015005973", "lolipop8974.")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("avaliable idx: ", idx)

}
