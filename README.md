# Go-Competitive-relationship
深入Golang竞争关系 Competitive relationship

### 处理共享内存时的竞争关系
- 互斥锁
- 原子操作
- 读写锁

这边只要讲解Atomic下的原子操作 和  读取写锁 

###  读写锁
``` go
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
```
读写锁   可以多读,但是只能单写  (如果有写操作  读会加锁)
适合 读多写少的情况

### 原子操作
写操作
```go 
func AddInt32(addr *int32,delta int32)
func AddInt64
func AddUint32
func AddUint64
func AddUintptr

```
读操作
```go 
atomic.LoadInt32
atomic.LoadInt64
atomic.LoadUint64
atomic.LoadUint64
```
比较并交换
```go 
func CompareAndSwapInt32(addr *int32,old, new int32) bool
CompareAndSwapInt64(addr *int64,old,new int64)

var ic int32
bol := atomic.CompareAndSwapInt32(&ic,0,100)
// 如果 ic == 0 就把ic赋值成100
```
交换
``` go
atomic.SwapInt32(addr *int32,new int32) (old int32)
```
存储
``` go 
StoreInt32(addr *int32,new int32)
```

### 深入 context
实现超时控制
```go  
func (d *DnsCoon) DnsParse(ctx context.Context, domain string) (string, error) {
	chaEr := make(chan error, 1)
	go func(chaEr chan error) {
		msg := &dns.Msg{}
		// 创建一个查询消息体
		msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)
		// 与当前dns建立连接
		e := d.dns.SetWriteDeadline(time.Now().Add(time.Millisecond * 200))
		if e != nil {
			chaEr <- e
			return
		}

		// 发送查询数据
		e = d.dns.WriteMsg(msg)
		if e != nil {
			chaEr <- e
			return
		}

		// 以下是接受数据阶段
		readMsg, e := d.dns.ReadMsg()
		if e != nil || len(readMsg.Question) == 0 {
			// 如果没有数据
			if e != nil {
				chaEr <- e
				return
			} else {
				chaEr <- errors.New("not data")
				return
			}
		}
		if len(readMsg.Answer) == 0 {
			chaEr <- NoDomain
			return
		} else {
			chaEr <- nil
			return
		}
	}(chaEr)

	for {
		select {
		case <-ctx.Done():
			return d.host, TimeOut
		case err := <-chaEr:
			return d.host, err
		}
	}
}
```

### 关于包倒入
倒入一个包会执行同包下所有.go的init 然后 执行 init.go 下的init

# 基本同步原语
### 1.Mutex
互斥锁
可以一个gorutine Lock 另一个gorutine Unlock但这种方式慎用

### 2.RWMutex
读写锁对于读做了优化，适合`写少读多`的场景，对于并发的读很适合
- 如果已经有一组reader持有读写锁，这个时候如果writer调用Lock，它会被阻塞。接着如果有reader调用RLock，等前面那组readerUnlock后，writer优先获取锁
``` golang
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
```
如果当前在写 读取操作就会被锁住，如果是读取就不会被锁住
### 3.Cond
类似Monitor,提供通知机制，可以Broadcast通知所有Waot的gorotine，也可以Signal通知某一个Wait的goroutine
Cond初始的时候要传入一个Locker接口的实现，这个Locker可以用来保护条件的变量。
Broadcast和Signal不需要加锁调用，但是调用Wait的时候需要加锁
Wait执行中有解锁重加锁的过程，在这个期间对临界区是没有保护的
一定要使用for循环来检查条件是否满足，因为随时都可以触发通知

# 扩展同步原语
对基本同步原语的补充，适用于额外的场景，由Go扩展包(试验包)和第三方提供。
### 1.ReentrantLock
标准库sync下的Mutex是不能重入的，如果想实现重入的话，可以利用:

*   goid:用来标记谁持有了当前锁，重入了几次
*   全局id:或者你自己维护一个全局id,但是实现的结构不再满足Locker接口

可重入锁也叫做递归锁，但是叫可重入锁更准确些，因为可重入可不只递归这么一种情况。

### 2.Semaphore
Dijkstra提出并发访问通用资源的并发原语，使用PV原语提供对临界区的保护。

二进制(取值0,1)的semaphore提供了锁的功能。  
计数器semaphore提供了对一组资源的保护。

包`golang.org/x/sync/semaphore`。

标准库内部的semaphore提供了休眠/唤醒的功能，用来实现基本同步原语的阻塞。
### 3.SingleFlight
并发的访问同一组资源的时候，只允许一个请求进行，这个请求把结果告诉其它等待者，避免雪崩的现象。

比如cache 失效的时候，只允许一个goroutine从数据库中捞数据回种，避免雪崩对数据库的影响。

扩展库中提供。
### ErrGroup
应用于 half sync/half async的场景(这个设计模式以后有机会再介绍)。

有一组异步的任务需要处理，利用这个原语可以等待所有的异步任务完成，并获取第一个错误。

如果想得到所有的错误，利用额外的map变量进行记录。

使用Context可以实现遇到第一个错误就返回。

扩展包中提供。

bilibili扩展了这个原语，提供了限定并发数量的功能。

### SpinLock
*   自旋锁
*   有时候很高效,因为当前CPU中运行的goroutine更有机会获取到锁
*   不公平
*   需要处理器忙等待
*   应用于竞争不是很激烈的状态

### fslock
文件锁， 可以控制多个进程之间的同步。

### concurrent-map
类似Java中的ConcurrentMap的设计思想，将key划分成一定数量的shard,每个shard一个锁，减少锁的竞争。

相对于sync.Map，可以应用写/删除/添加更频繁的场景。

# 原子操作
### 操作的数据
*   int32
*   int64
*   uint32
*   uint64
*   uintptr
*   unsafe.Pointer

### 操作方法
*   AddXXX (整数类型)
*   CompareAndSwapXXX：cas
*   LoadXXX：读取
*   StoreXXX：存储
*   SwapXXX：交换

### Subtract
*   有符号的类型，可以使用Add负数
*   无符号的类型，可以使用AddUint32(&x, ^uint32(c-1)),AddUint64(&x, ^uint64(c-1))
*   无符号类型减一， AddUint32(&x, ^uint32(0))， AddUint64(&x, ^uint64(0))

### Value

一个通用的对象，可以很方便的对struct等类型进行原子存储和加载。

由于不同的架构下对原子操作的支持是不一样的，有些架构师是不支持的。

