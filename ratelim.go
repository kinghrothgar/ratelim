package ratelim

import (
	"container/list"
	"sync"
	"time"
)

type Queue struct {
	limit  int
	period time.Duration
	times  *list.List
	sync.RWMutex
}

func (rlq *Queue) Action() (bool, time.Duration) {
	rlq.Lock()
	defer rlq.Unlock()
	if rlq.times.Len() < rlq.limit {
		rlq.times.PushBack(time.Now())
		return true, 0
	}
	te := rlq.times.Front()
	now := time.Now()
	diff := now.Sub(te.Value.(time.Time))
	// TODO should this include =
	if diff >= rlq.period {
		rlq.times.Remove(te)
		rlq.times.PushBack(now)
		return true, 0
	}
	return false, rlq.period - diff
}

func New(limit int, period time.Duration) *Queue {
	return &Queue{
		limit:  limit,
		period: period,
		times:  list.New(),
	}
}
