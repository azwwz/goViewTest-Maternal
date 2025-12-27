# Channel concurrency safety (Go)

这是一个简化的面试题讲解项目，回答“channel 是否并发安全？”的问题。

## 结论
- **发送/接收是并发安全的**：多个 goroutine 可以同时对同一个 channel 发送/接收，不需要额外锁。
- **关闭不是并发安全的**：`close` 必须由单一 goroutine 负责，否则可能出现 `close of closed channel` 的 panic。
- **向已关闭的 channel 发送是不安全的**：会直接 panic。

## 为什么
Go 的 channel 实现内置了锁和调度机制，保证并发下发送/接收的原子性与正确阻塞/唤醒；但 `close` 是一次性状态变更，多个并发 close 会产生竞态。

## 运行
```bash
go run .
```

## 输出说明
- 第一个示例展示并发发送/接收没有问题。
- 第二个示例展示并发 close 的危险性（用 recover 捕获 panic）。
- 第三个示例展示“单一拥有者关闭”的安全模式。

## 面试回答模板（可背）
“Go 的 channel 在发送和接收上是并发安全的，多 goroutine 可以同时操作同一个 channel；但关闭 channel 不是并发安全的，必须确保只有一个关闭者，否则会 panic。典型做法是让发送方的 owner 负责 close，接收方只 range 读取。”
