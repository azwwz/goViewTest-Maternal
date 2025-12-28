# Go GMP 理解与代码导览

这份文档配合 `main.go` 使用，通过一个“简化调度器模拟 + 真实运行时对比”的方式，帮助你从代码层面理解 Go 的 GMP 模型。

## 1. GMP 是什么

- G (Goroutine): 轻量级协程，保存栈、指令位置等执行现场。
- M (Machine): OS 线程，真正执行指令的实体。
- P (Processor): 调度器逻辑实体，维护可运行 G 的队列(runq)。一个 M 必须绑定一个 P 才能执行 G。

**核心关系**

- G 只是“待运行的任务”，不能直接跑在 OS 线程上。
- M 是线程，但必须拿到一个 P 才能执行 G。
- P 是“执行资格 + 运行队列”，决定哪些 G 先运行。

## 2. 代码结构概览

`main.go` 有两种模式：

1) `-mode=sim`: 自己实现的简化 GMP 调度模拟器。用于看清调度行为和队列变化。
2) `-mode=runtime`: 运行时真实调度的并行表现对比（GOMAXPROCS=1 vs NumCPU）。

## 3. 模拟调度器说明

模拟器是“概念模型”，重点帮助理解以下机制：

- P 拥有本地队列 `runq`
- M 绑定 P，从 P 的 `runq` 取 G 运行
- G 被时间片打断后回到 P 的队列尾部
- 当某个 P 队列空时，M 会从其他 P “偷”一个 G（work stealing）

### 对应代码位置

- G / P / M 结构体: `main.go`
- Scheduler.step(): 每个 tick 让每个 M 执行一个 G 的时间片
- Scheduler.steal(): 简化的 work stealing

### 运行方式

```bash
# 2 个 P, 2 个 M, 6 个 G
# 输出每个时间片的调度事件

go run . -mode=sim -p=2 -m=2 -g=6 -seed=1 -log=true
```

你会看到类似输出：

- `M` 拿到 `P` 的 G 执行一个时间片
- G 还有剩余步骤时被放回队列
- P 空了时触发偷取

这些日志帮助你直观观察：

- 任务在 P 之间分布
- M 只能执行其绑定 P 的队列
- work stealing 如何避免某个 P 空转

## 4. 真实运行时对比说明

`-mode=runtime` 使用真实 Go 调度器做一个“CPU 密集任务”的并行对比：

- `GOMAXPROCS=1`: 所有 G 共享 1 个 P，不能真正并行
- `GOMAXPROCS=NumCPU`: P 数量增多，能并行执行

```bash
# 运行时对比

go run . -mode=runtime -work=5000000
```

输出示例：

- `GOMAXPROCS=1` 比较慢
- `GOMAXPROCS=NumCPU` 更快（并行）

这对应 GMP 中 “P 决定并行度” 的关键结论。

## 5. 关键理解点总结

- G 是任务，M 是线程，P 是运行资格和调度队列
- M 必须绑定 P 才能跑 G
- P 的数量决定了并行度，而不是 G 的数量
- work stealing 让任务分布更均衡，减少空闲

## 6. 自己可做的实验

1) 调整 `-p` 和 `-m` 观察模拟器行为
2) 增大 `-g` 看队列与偷取效果
3) 调整 `-work` 看真实调度中的并行效果

