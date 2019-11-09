package scheduleq

import (
	"fmt"
	"testing"
	"time"
)

func TestScheduleQueue(t *testing.T) {
	var test = 0
	q := NewQueue()
	q.Plan(func(){
		test = test + 1
	}, 100)
	q.Plan(func(){
		test = test + 2
	}, 100)
	q.Plan(func(){
		test = test + 3
	}, 100)

	time.Sleep(time.Second)
	fmt.Printf("test is: %d", test)
}