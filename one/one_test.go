package test

import (
	"log"
	"sync"
	"sync/atomic"
	"testing"
)

// 读写锁   和  互斥锁的测试

// 互斥锁测试
func TestOne(t *testing.T) {
	a := map[string]int{
		"one": 1,
	}
	var mu sync.Mutex
	var counter int64
	var wg sync.WaitGroup
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				mu.Lock()
				// 写测试
				a["one"] = i
				mu.Unlock()

				// 这里是原子操作
				atomic.AddInt64(&counter, 1)
			} else {
				mu.Lock()
				a1 := a["one"]
				a1 = a1
				mu.Unlock()

				atomic.AddInt64(&counter, 1)
			}
		}(i)
	}

	wg.Wait()
	log.Println(a)
	log.Println(counter)

	//=== RUN   TestOne
	//2019/11/12 19:21:56 map[one:999998]
	//2019/11/12 19:21:56 1000000
	//--- PASS: TestOne (0.36s)
}

// 读取锁测试
func TestTwo(t *testing.T) {
	a := map[string]int{
		"one": 1,
	}
	var rw sync.RWMutex
	var counter int64
	var wg sync.WaitGroup
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				rw.Lock()
				// 写测试
				a["one"] = i
				rw.Unlock()

				// 这里是原子操作
				atomic.AddInt64(&counter, 1)
			} else {
				rw.RLock()
				a1 := a["one"]
				a1 = a1
				rw.RUnlock()

				atomic.AddInt64(&counter, 1)
			}
		}(i)
	}

	wg.Wait()
	log.Println(a)
	log.Println(counter)

	//=== RUN   TestTwo
	//2019/11/12 19:25:04 map[one:999990]
	//2019/11/12 19:25:04 1000000
	//--- PASS: TestTwo (0.62s)
}

// 这是为什么  读取锁比互斥锁还慢?
