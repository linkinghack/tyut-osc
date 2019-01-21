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
