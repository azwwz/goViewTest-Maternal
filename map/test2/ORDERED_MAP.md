## 有序 map 的现实场景示例

Go 的 map 本身是无序的。如果你需要“有序输出”，常见做法是：
1) 先收集所有 key 到切片
2) 对切片排序
3) 按排序后的 key 依次访问 map

本目录的 `ordered_map_example.go` 模拟了一个库存报表：输出时要求按产品编号（key）排序。

### 关键步骤

- 收集 key
- 排序 key
- 依序读取 map

### 运行示例

```bash
go run ordered_map_example.go
```

预期输出类似：

```
Inventory report (sorted by product code):
- apple (Apple): 18
- banana (Banana): 12
- orange (Orange): 7
```

### 说明与建议

- 如果你需要按 value 排序，可以先构造一个切片（例如 []Item 或 []struct{K string; V Item}），再按自定义规则排序。
- 如果你需要稳定且频繁的有序插入/删除，考虑用切片维护顺序，并用 map 做索引（或引入第三方有序 map 实现）。
