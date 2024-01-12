package scheduler

import (
	"math"
	"sync/atomic"
)

type DEQueue interface {
	pushBottom(task *ImageTask)
	IsEmpty() bool //returns whether the queue is empty
	popTop() *ImageTask
	popBottom() *ImageTask
	Size() int
}

type bDeque struct {
	top    int32
	mask   int32
	bottom int32
	tasks  []*ImageTask
	size   int32
}

func NewBDEQueue() *bDeque {
	capacity := int32(math.Pow(2, 16))
	return &bDeque{
		top:    0,
		bottom: 0,
		tasks:  make([]*ImageTask, capacity),
		size:   0,
		mask:   capacity - 1,
	}
}

func (q *bDeque) popTop() *ImageTask {

	oldTopStamp := atomic.LoadInt32(&q.top)
	oldTop := oldTopStamp & q.mask
	oldStamp := (oldTopStamp >> 16) & 0xFFFF
	newTop := oldTop + 1
	newStamp := oldStamp + 1

	if q.bottom <= oldTop {
		return nil // Empty deque
	}
	index := oldTop & q.mask
	task := q.tasks[index]

	// task := q.tasks[oldTop]

	if atomic.CompareAndSwapInt32(&q.top, oldTopStamp, (newTop<<16)|newStamp) && task != nil {
		atomic.AddInt32(&q.size, -1)
		return task
	}

	return nil
}

func (q *bDeque) popBottom() *ImageTask {

	bottom := atomic.LoadInt32(&q.bottom)
	if bottom == 0 {

		return nil
	}

	newBottom := bottom - 1
	task := q.tasks[newBottom]

	oldTopStamp := atomic.LoadInt32(&q.top)
	oldTop := oldTopStamp & q.mask
	oldStamp := (oldTopStamp >> 16) & 0xFFFF
	newTop := int32(0)
	newStamp := oldStamp + 1

	if newBottom > oldTop {
		if atomic.CompareAndSwapInt32(&q.bottom, bottom, newBottom) {
			atomic.AddInt32(&q.size, -1)
			return task
		}
	}

	if newBottom == oldTop {
		atomic.StoreInt32(&q.bottom, 0)
		if atomic.CompareAndSwapInt32(&q.top, oldTopStamp, (newTop<<16)|newStamp) {
			atomic.AddInt32(&q.size, -1)
			return task
		}
	}

	atomic.StoreInt32(&q.top, (newTop<<16)|newStamp)
	atomic.StoreInt32(&q.bottom, 0)
	return nil
}

func (q *bDeque) pushBottom(task *ImageTask) {

	currentBottom := atomic.LoadInt32(&q.bottom)
	if currentBottom == int32(len(q.tasks)) {
		// Deque is full
		panic("Deque is full")
	}

	q.tasks[currentBottom] = task
	atomic.StoreInt32(&q.bottom, currentBottom+1)
	atomic.AddInt32(&q.size, 1)
}

func (q *bDeque) IsEmpty() bool {
	top := atomic.LoadInt32(&q.top) // Extract the top index
	bottom := atomic.LoadInt32(&q.bottom)
	return top >= bottom
}

func (q *bDeque) Size() int {
	return int(atomic.LoadInt32(&q.size))
}
