package scheduleq

import (
	"testing"
	"time"
)

type stateHolder struct {
	counter int
}

func TestScheduleQueue(t *testing.T) {

	sh := &stateHolder{counter: 0}
	cnt := 0
	q := NewQueue(0, 0, 0)

	for i := 1; i < 4; i++ {
		increment := i //because _current_ i would be 3 in local scope at time of closure execution,
		//so LAST value of i  will be used, and you'll get 9 as result

		cnt = cnt //you can do this, looks sick imo, but nevertheless binds to outer scope

		q.Plan(func() error {
			func(holder *stateHolder) {
				cnt = cnt + increment
				holder.counter = holder.counter + increment
				t.Log("Planning to increase by ", increment)
			}(sh)
			return nil
		}, time.Now().Add(100*time.Millisecond))
	}
	q.Advance(time.Now().Add(time.Second))
	t.Log("holder struct = ", sh.counter)
	t.Log("local variable = ", cnt)
	if sh.counter != 6 || cnt != 6 {
		t.Fail()
	}
}

func TestThrottledScheduleq(t *testing.T) {
	q := NewQueue(100*time.Millisecond, 3, 1)
	//three groups of tasks, 100ms after each other
	for i := 1; i < 8; i++ {
		increment := i //because _current_ i would be 3 in local scope at time of closure execution,
		q.PlanWithThrottle(func() error {
			t.Logf("Func %d running", increment)
			return nil
		})
	}
	t.Log("len = ", q.Len())
	if q.Len() != 7 {
		t.Errorf("q.Len() != 7, got %v instead", q.Len())
	}
	tt := time.Now()
	q.Advance(tt)
	q.Advance(tt) //calling twice gets no execution of later planned tasks
	t.Log("len = ", q.Len())
	if q.Len() != 4 {
		t.Errorf("q.Len() != 4, got %v instead", q.Len())
	}
	t.Log("len = ", q.Len())
	q.Advance(tt.Add(90 * time.Millisecond))
	if q.Len() != 4 {
		t.Errorf("q.Len() != 4, got %v instead", q.Len())
	}
	q.Advance(tt.Add(100 * time.Millisecond))
	t.Log("len = ", q.Len())
	if q.Len() != 1 {
		t.Errorf("q.Len() != 1, got %v instead", q.Len())
	}
	q.Advance(tt.Add(200 * time.Millisecond))
	t.Log("len = ", q.Len())
	if q.Len() != 0 {
		t.Errorf("q.Len() != 0, got %v instead", q.Len())
	}
}
