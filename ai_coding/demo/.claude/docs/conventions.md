# Go 代码规范

## 命名

- 包名小写、单数、无下划线（`service`、`dao`、`httpapi`）
- 文件名全小写；多词直接连写或用下划线均可，测试文件以 `_test.go` 结尾
- 导出标识符驼峰（`NewHandler`）；缩写全大写或全小写（`ID` / `id`、`HTTP`、`JSON`、`URL`）
- 常量按语义命名，避免魔法数字

## 注释

- 导出的类型/函数必须有**以名字开头**的文档注释（`// Wallet 表示...`）
- 业务约束（如"金额用 int64 禁 float"）在注释里写清原因，不只写"是什么"
- 复杂逻辑加行内注释解释"为什么"，简单代码不啰嗦

## 错误处理

- 业务错误用 **sentinel error**（`var ErrNotFound = errors.New(...)`），上层用 `errors.Is` 判断
- 跨层传递需要加上下文时用 `fmt.Errorf("...: %w", err)` wrap，保留可判定性
- 不吞错误（不写 `_ = doSomething()` 除非确实要忽略且注明）
- 错误字符串小写、结尾不带标点
- transport 层用 `errors.Is` 把 sentinel 映射成 HTTP 状态码，不在 handler 里判断字符串

## 并发

- 并发访问的共享 `map` 必须用 `sync.RWMutex` 保护（Go map 并发读写是 fatal，不可恢复）
- 状态只在 dao 层持有，service/transport 无状态
- goroutine 要有明确退出路径，不 fire-and-forget

## 测试

- 表驱动测试；文件 `_test.go`，与被测包同包或 `_test` 包
- HTTP 层用 `net/http/httptest`，不起真实端口
- **仅用标准库 `testing`**，不引入 testify 等断言库（离线约束）
- 边界必覆盖：金额非法（0/负）、余额不足、用户不存在、坏 JSON

## 提交规范

- Conventional Commits：`feat` / `fix` / `docs` / `test` / `refactor` / `chore`
- 一次提交一件事，信息具体到改了什么、为什么
