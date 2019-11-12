package test

import (
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 原子操作


// 原子操作 和 互斥锁测试(原子)
func TestOne(t *testing.T) {
	var i int64
	go func() {
		for {
			go func() {
				atomic.AddInt64(&i,1)
			}()
		}
	}()

	time.Sleep(time.Second * 10)
	log.Println(i)

	//=== RUN   TestOne
	//2019/11/12 19:51:50 27578687
	//--- PASS: TestOne (10.03s)
}

// 原子操作 和 互斥锁测试(互斥锁)
func TestTwo(t *testing.T) {
	var i int64
	var mu sync.Mutex
	go func() {
		for {
			go func() {
				mu.Lock()
				i ++
				mu.Unlock()
			}()
		}
	}()

	time.Sleep(time.Second * 10)
	log.Println(i)

	//=== RUN   TestTwo
	//2019/11/12 19:52:13 28026506
	//--- PASS: TestTwo (10.03s)
}

// 为什么互斥锁比原子操作还快(⊙﹏⊙)