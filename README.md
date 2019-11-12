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
```go 
func AddInt32(addr *int32,delta int32)
func AddInt64
func AddUint32
func AddUint64
func AddUintptr

```

