# Go map 并发安全说明（库存场景）

## 结论
Go 语言内置的 `map` 不是并发安全的：多 goroutine 同时读写/写写会触发数据竞争，严重时直接 panic（`fatal error: concurrent map read and map write`）。在真实业务里必须加锁或使用并发安全结构。

## 真实场景
电商库存缓存（内存加速）：  
- 下单服务会并发扣减库存  
- 商品详情页会并发查询库存  
同一个 SKU 会在读写交错下被访问

## 代码说明
`main.go` 提供两种模式：
- `-mode=safe`：使用 `map + RWMutex`，读可并发，写互斥，安全且性能合理  
- `-mode=unsafe`：直接使用原生 `map`，在并发读写场景下可能 panic  

## 为什么需要锁
`map` 的内部结构不是并发可重入的：写操作可能触发扩容或重排，读操作并发进入会读取到不一致的数据结构，从而触发崩溃或数据错乱。  
`RWMutex` 的做法是：
- 读操作加读锁，可多读并发  
- 写操作加写锁，确保写期间没有读/写在进行

## 运行方式
```bash
go run . -mode=safe
```

危险示例（可能崩溃，仅用于演示）：
```bash
go run . -mode=unsafe
```

## 取舍建议
- 读多写少：优先 `map + RWMutex`  
- 并发高且访问模式特殊：可考虑 `sync.Map`，但一般业务缓存更推荐 `map + RWMutex`，可控性更强
