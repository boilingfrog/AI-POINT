---
name: code-review-go
description: 对 coinwallet 的 Go 变更做分层代码审查，对齐 Go 官方 CodeReviewComments 与 Uber Go Style Guide，按严重程度分级输出。审查 diff 或合并前使用。
argument-hint: [--diff | <path>]
---

对给定的 Go 改动（`--diff` 审当前变更，或指定 `<path>`）逐维度审查，输出分级问题清单。

## 审查维度

1. **错误处理** — 是否用 sentinel + `errors.Is`；跨层是否 `%w` wrap；有无吞错；
   错误字符串小写无标点。
2. **正确性与资金安全** — panic/nil/切片越界风险；**金额是否 int64、禁 float、校验 > 0、
   扣减前查余额、有无溢出**（资金域重点）。
3. **并发安全** — 共享 map 是否加锁；有无数据竞争；goroutine 有无退出路径。
4. **分层依赖** — 是否单向 `transport → service → dao → model`；有无反向依赖；
   handler 里是否混入了业务逻辑。
5. **HTTP 语义** — 状态码映射是否正确（404/400/409/500）；输入是否校验；
   JSON 解码错误是否处理；是否设了请求超时。
6. **性能** — 无谓的拷贝/分配；循环内重复计算；能预分配的 slice/map。
7. **安全与可观测** — 输入校验；日志不泄敏感信息；金额操作可追溯。
8. **可维护与命名** — 命名（缩写大小写一致）；导出有文档注释；无死代码/注释代码/调试日志。

## 输出格式

按严重程度分级，每条精确到 `file:line`，附依据和修复建议：

```
🔴 严重  [file:line] 说明 + 依据 + 修复
🟠 高    [file:line] ...
🟡 中    [file:line] ...
🔵 建议  [file:line] ...
```

## 结论

- 金额用 float / 漏校验 / 余额可为负 / 并发 map 无锁 / 会 panic → **REJECT**
- 分层反向依赖 / 错误映射缺失 / 边界漏测 → **NEEDS_WORK**
- 仅中低问题 → **APPROVE**（附改进建议）
