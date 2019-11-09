package scheduleq

import (
	"sync/atomic"
	"time"
)

var count int64

type Timer struct {
	id     int64
	OnTime func()
}

func NewTimer(f func()) Timer {
	id := atomic.AddInt64(&count, 1) //И действительно, спасибо!
	return Timer{id: id, OnTime: f}
}

func EmptyTimer() Timer {
	return Timer{}
}

//события
type timerData struct {
	timer Timer
	time  time.Time
	index int
}

//implements heap interface
type timers []*timerData

func (t timers) Len() int {
	return len(t)
}

func (t timers) Less(i, j int) bool {
	return t[i].time.Before(t[j].time)
}

func (t timers) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
	t[i].index, t[j].index = t[j].index, t[i].index
}

func (t *timers) Push(x interface{}) {
	idx := len(*t)
	item := x.(*timerData)
	item.index = idx
	*t = append(*t, item)
}

func (t *timers) Pop() interface{} {
	old := *t
	n := len(old)
	item := old[n-1]
	item.index = -1
	*t = old[0 : n-1]
	return item
}
