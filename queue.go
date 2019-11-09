package scheduleq

import (
	"container/heap"
	"errors"
	"time"
)

type Queue struct {
	heap timers
	//Указатели на кучу
	table map[int64]*timerData
}

func NewQueue() Queue {
	return Queue{
		table: make(map[int64]*timerData),
	}
}

func (q *Queue) Len() int {
	return len(q.heap)
}

// Schedule schedules a timer for execution at time tm. If the
// timer was already scheduled, it is rescheduled.
func (q *Queue) Schedule(t Timer, tm time.Time) {
	if data, ok := q.table[t.id]; !ok {
		data = &timerData{t, tm, 0}
		heap.Push(&q.heap, data)
		q.table[t.id] = data
	} else {
		data.time = tm
		heap.Fix(&q.heap, data.index)
	}
}

// Unschedule unschedules a timer's execution.
func (q *Queue) Unschedule(t Timer) {
	if data, ok := q.table[t.id]; ok {
		heap.Remove(&q.heap, data.index)
		delete(q.table, t.id)
	}
}

// GetTime returns the time at which the timer is scheduled.
// If the timer isn't currently scheduled, an error is returned.
func (q *Queue) GetTime(t Timer) (tm time.Time, err error) {
	if data, ok := q.table[t.id]; ok {
		return data.time, nil
	}
	return time.Time{}, errors.New("timerqueue: timer not scheduled")
}

// IsScheduled returns true if the timer is currently scheduled.
func (q *Queue) IsScheduled(t Timer) bool {
	_, ok := q.table[t.id]
	return ok
}

// Clear unschedules all currently scheduled timers.
func (q *Queue) Clear() {
	q.heap, q.table = nil, make(map[int64]*timerData)
}

// PopFirst removes and returns the next timer to be scheduled and
// the time at which it is scheduled to run.
func (q *Queue) PopFirst() (t Timer, tm time.Time) {
	if len(q.heap) > 0 {
		data := heap.Pop(&q.heap).(*timerData)
		delete(q.table, data.timer.id)
		return data.timer, data.time
	}
	return EmptyTimer(), time.Time{}
}

// PeekFirst returns the next timer to be scheduled and the time
// at which it is scheduled to run. It does not modify the contents
// of the timer queue.
func (q *Queue) PeekFirst() (t Timer, tm time.Time) {
	if len(q.heap) > 0 {
		return q.heap[0].timer, q.heap[0].time
	}
	return EmptyTimer(), time.Time{}
}

// Advance executes OnTimer callbacks for all timers scheduled to be
// run before the time 'tm'. Executed timers are removed from the
// timer queue.
func (q *Queue) Advance(tm time.Time) {
	for len(q.heap) > 0 && !tm.Before(q.heap[0].time) {
		data := q.heap[0]
		heap.Remove(&q.heap, data.index)
		delete(q.table, data.timer.id)
		data.timer.OnTime()
	}
}

func (q *Queue) Plan(f func(), after int64) {
	q.Schedule(NewTimer(f), time.Now().Add(time.Millisecond*time.Duration(after)))
}
