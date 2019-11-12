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