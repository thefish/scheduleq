package scheduleq

import (
	"container/heap"
	"errors"
	"time"
)

type Queue struct {
	heap timers
	//Указатели на кучу
	table        map[int64]*taskData
	Period       time.Duration
	MaxPerPeriod int
	MaxRetries   int
}

func NewQueue(period time.Duration, maxPerPeriod, maxRetries int) Queue {
	return Queue{
		table:        make(map[int64]*taskData),
		Period:       period,
		MaxPerPeriod: maxPerPeriod,
		MaxRetries:   maxRetries,
	}
}

func (q *Queue) Len() int {
	return len(q.heap)
}

// Schedule schedules a task for execution at time tm. If the
// task was already scheduled, it is rescheduled.
func (q *Queue) Schedule(t Task, tm time.Time) {
	if data, ok := q.table[t.id]; !ok {
		data = &taskData{t, tm, 0}
		heap.Push(&q.heap, data)
		q.table[t.id] = data
	} else {
		data.time = tm
		heap.Fix(&q.heap, data.index)
	}
}

// Unschedule unschedules a task's execution.
func (q *Queue) Unschedule(t Task) {
	if data, ok := q.table[t.id]; ok {
		heap.Remove(&q.heap, data.index)
		delete(q.table, t.id)
	}
}

// GetTime returns the time at which the task is scheduled.
// If the task isn't currently scheduled, an error is returned.
func (q *Queue) GetTime(t Task) (tm time.Time, err error) {
	if data, ok := q.table[t.id]; ok {
		return data.time, nil
	}
	return time.Time{}, errors.New("timerqueue: task not scheduled")
}

// IsScheduled returns true if the task is currently scheduled.
func (q *Queue) IsScheduled(t Task) bool {
	_, ok := q.table[t.id]
	return ok
}

// Clear unschedules all currently scheduled timers.
func (q *Queue) Clear() {
	q.heap, q.table = nil, make(map[int64]*taskData)
}

// PopFirst removes and returns the next task to be scheduled and
// the time at which it is scheduled to run.
func (q *Queue) PopFirst() (t Task, tm time.Time) {
	if len(q.heap) > 0 {
		data := heap.Pop(&q.heap).(*taskData)
		delete(q.table, data.task.id)
		return data.task, data.time
	}
	return EmptyTask(), time.Time{}
}

// PeekFirst returns the next task to be scheduled and the time
// at which it is scheduled to run. It does not modify the contents
// of the task queue.
func (q *Queue) PeekFirst() (t Task, tm time.Time) {
	if len(q.heap) > 0 {
		return q.heap[0].task, q.heap[0].time
	}
	return EmptyTask(), time.Time{}
}

// Advance executes OnTimer callbacks for all timers scheduled to be
// run before the time 'tm'. Executed timers are removed from the
// task queue.
func (q *Queue) Advance(tm time.Time) {
	for len(q.heap) > 0 && !tm.Before(q.heap[0].time) {
		data := q.heap[0]
		err := data.task.OnTime()
		if err != nil {
			newRetryCount := data.task.retries + 1
			if newRetryCount < q.MaxRetries {
				//reschedule
				data.task.retries = newRetryCount
				q.retry(data.task)
			} else {
				if data.task.OnRetryFail != nil {
					data.task.OnRetryFail()
				}
			}
		}
		heap.Remove(&q.heap, data.index)
		delete(q.table, data.task.id)
	}
}

// Plan schedules new task without respect to throttling, just on time
func (q *Queue) Plan(t *Task, after time.Time) {
	q.Schedule(*t, after)
}

// PlanWithThrottle schedules new task in throttled manner
func (q *Queue) PlanWithThrottle(t *Task) {
	q.Schedule(*t, q.getDelay())
}

func (q *Queue) retry(t Task) {
	q.Schedule(t, q.getDelay())
}

func (q *Queue) getDelay() time.Time {
	if q.MaxPerPeriod < 1 || q.Period == 0 {
		return time.Now()
	}
	periods := int(q.Len() / q.MaxPerPeriod)
	delay := time.Now().Add(time.Duration(periods) * q.Period)
	return delay
}
