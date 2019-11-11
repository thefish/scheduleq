package scheduleq

import (
	"github.com/thefish/scheduleq"
	"testing"
	"time"
)

type stateHolder struct {
	counter int
}

func TestScheduleQueue(t *testing.T) {

	sh := &stateHolder{counter: 0}
	cnt := 0
	q := scheduleq.NewQueue()

	for i := 1; i < 4; i++ {
		increment := i 	//because _current_ i would be 3 in local scope at time of closure execution, 
				//so LAST value of i  will be used, and you'll get 9 as result
		
		cnt = cnt      	//you can do this, looks sick imo, but nevertheless binds to outer scope
		
		q.Plan(func() {
			func(holder *stateHolder) {
				cnt = cnt + increment
				holder.counter = holder.counter + increment
				t.Log("Planning to increase by ", increment)
			}(sh)
		}, 100)
	}
	time.Sleep(time.Second)
	q.Advance(time.Now())
	t.Log("holder struct = ", sh.counter)
	t.Log("local variable = ", cnt)
	if sh.counter != 6 || cnt != 6 {
		t.Fail()
	}
}
