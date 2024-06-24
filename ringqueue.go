package ringqueue

import "sync/atomic"

type RingQueue[T any] struct {
	List    *[]T
	maxSize int
	head    int32
	tail    int32
	gap     int32
}

func NewRing[T any](maxSize ...int) *RingQueue[T] {
	var defaultSize = 1000
	if len(maxSize) >= 1 {
		defaultSize = maxSize[0]
	}
	list := make([]T, defaultSize)
	return &RingQueue[T]{
		List:    &list,
		maxSize: defaultSize,
		head:    0,
		tail:    0,
		gap:     0,
	}
}

func (queue *RingQueue[T]) Enqueue(value T) bool {
	for {
		tail := atomic.LoadInt32(&queue.tail)
		nextTail := (tail + 1) % int32(queue.maxSize)
		gap := atomic.LoadInt32(&queue.gap)
		if gap == int32(queue.maxSize) {
			return false // 队列已满
		}
		if atomic.CompareAndSwapInt32(&queue.tail, tail, nextTail) {
			(*queue.List)[tail] = value
			atomic.AddInt32(&queue.gap, 1)
			return true
		}
	}
}

func (queue *RingQueue[T]) Dequeue() (value T, ok bool) {
	for {
		if atomic.LoadInt32(&queue.gap) <= 0 {
			return
		}
		head := atomic.LoadInt32(&queue.head)
		nextHead := (head + 1) % int32(queue.maxSize)
		if atomic.CompareAndSwapInt32(&queue.head, head, nextHead) {
			value = (*queue.List)[head]
			atomic.AddInt32(&queue.gap, -1)
			ok = true
			return
		}
	}
}

func (queue *RingQueue[T]) Peek() (value T, ok bool) {
	if atomic.LoadInt32(&queue.gap) == 0 {
		return
	}
	head := atomic.LoadInt32(&queue.head)
	value = (*queue.List)[head]
	ok = true
	return
}

func (queue *RingQueue[T]) Empty() bool {
	return atomic.LoadInt32(&queue.gap) == 0
}

func (queue *RingQueue[T]) Size() int {
	return int(atomic.LoadInt32(&queue.gap))
}
