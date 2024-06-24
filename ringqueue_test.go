package ringqueue

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRingQueue_cq(t *testing.T) {
	q := NewRing[int]()

	var wg sync.WaitGroup
	var enqueueCount int32
	var dequeueCount int32

	// 多协程并发写入测试
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				q.Enqueue(id*100 + j)
				atomic.AddInt32(&enqueueCount, 1)
				time.Sleep(time.Millisecond * 10)
			}
		}(i)
	}

	// 多协程并发读取测试
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				value, ok := q.Dequeue()
				if ok {
					t.Logf("Goroutine %d dequeued %d\n", id, value)
					atomic.AddInt32(&dequeueCount, 1)
				} else {
					t.Logf("Goroutine %d failed to dequeue\n", id)
				}
				time.Sleep(time.Millisecond * 15)
			}
		}(i)
	}

	wg.Wait()

	// 验证队列是否为空
	if !q.Empty() {
		t.Errorf("Queue is not empty. Size: %d", q.Size())
	}

	// 验证写入和读取的数量是否匹配
	if enqueueCount != dequeueCount {
		t.Errorf("Mismatch between enqueue count (%d) and dequeue count (%d)", enqueueCount, dequeueCount)
	}
}

func BenchmarkRingQueueEnqueueDequeue(b *testing.B) {
	q := NewRing[int]()
	b.ResetTimer() // 重置计时器，开始计时
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q.Enqueue(1)
			q.Dequeue()
		}
	})
}
