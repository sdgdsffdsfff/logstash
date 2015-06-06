package utils

import "container/list"
import "sync"

// FIFO
type Queue struct {
	index int
	list  *list.List
	mux   sync.Mutex
}

func (q *Queue) Push(v interface{}) {
	q.mux.Lock()
	q.list.PushBack(v)
	q.index++
	q.mux.Unlock()
}

func (q *Queue) Pop() (v interface{}) {
	q.mux.Lock()
	v = q.list.Back()
	q.index--
	q.mux.Unlock()

	return
}
