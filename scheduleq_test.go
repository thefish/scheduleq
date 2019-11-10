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
		increment := i //because _current_ i would be out of scope, instead LAST value will be used, and you'll get 12 as result
		cnt = cnt      //or you can do this as a variant, looks sick imo, but nevertheless binds to outer scope
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
