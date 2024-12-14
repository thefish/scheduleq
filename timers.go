package scheduleq

import (
	"sync/atomic"
	"time"
)

var count int64

type Task struct {
	id      int64
	OnTime  func() error
	retries int
	OnFail  func()
}

func NewTaskWithRetry(f func() error) Task {
	id := atomic.AddInt64(&count, 1)
	return Task{
		id:      id,
		OnTime:  f,
		retries: 0,
	}
}

func EmptyTask() Task {
	return Task{}
}

// события
type taskData struct {
	timer Task
	time  time.Time
	index int
}

// implements heap interface
type timers []*taskData

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
	item := x.(*taskData)
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
