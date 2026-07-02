# coinwallet — 金币钱包 HTTP 服务（demo）

一个用 Go 标准库写的金币钱包微服务，**用来演示 reborn 那套 Claude Code 配置逻辑**
（分层结构 + 多 agent 工作流 + hooks + 合并前评审闸门）在一个能真实编译运行的小项目上
怎么落地。金币发放/消费属于"资金域"，正好对应 merger 的评审闸门 + 人审规则。

> 本目录的 `.claude/` 与本文件，只有在**以 `demo/` 为工作目录**打开 Claude Code 时才生效
> （hooks/设置按 cwd 向上发现）。在 ai_coding 根目录打开不会激活这里的配置。

## 技术栈

- Go 1.26 · **仅标准库**（`net/http` / `encoding/json` / `sync` / `testing`）
- **零外部依赖**：本项目离线开发，不能 `go mod download`，禁止引入任何需拉取的库

## 代码地图

```
cmd/server/     # 装配各层 + 启动 http.Server（PORT 环境变量，默认 8080）
transport/http/ # 传输层(package httpapi)：路由、JSON 编解码、状态码映射，不写业务
service/        # 业务逻辑 + 校验 + sentinel error + Repository 接口
dao/            # 数据访问：内存 map + sync.RWMutex（唯一持有状态处）
model/          # 数据结构：Wallet{UserID, Balance int64}
```

依赖严格单向：`cmd → transport/http → service → dao → model`，不允许反向。
`service` 只通过自己声明的 `Repository` 接口依赖 `dao`。

## 怎么验证改动

```bash
make lint    # gofmt 门禁：有未格式化文件即失败
make vet     # go vet ./...
make test    # go test ./...（service 表驱动 + transport httptest）
make run     # go run ./cmd/server，监听 :8080
# 冒烟：
curl -s -XPOST localhost:8080/wallets/u1/grant -d '{"amount":100}'   # {"user_id":"u1","balance":100}
curl -s localhost:8080/wallets/u1                                    # balance 100
curl -s -XPOST localhost:8080/wallets/u1/spend -d '{"amount":30}'    # balance 70
```

> 保存 .go 会自动 gofmt，`git commit` 前会跑 gofmt/vet 门禁（见 `.claude/settings.json` hooks）。

## 必守规则

- **分层单向依赖**：只在对应层改代码，不跨层反向依赖；handler 不写业务逻辑
- **金额用 int64**：账本金额一律 int64，**禁止 float**（精度问题，资金域绝不允许）
- **金额必校验**：发放/消费金额必须 > 0；扣减前必须查余额是否充足
- **错误用 sentinel**：service 层返回 sentinel error，上层用 `errors.Is` 映射状态码
- **共享 map 加锁**：并发访问的 map 必须用 `sync.RWMutex`（Go map 并发读写会 fatal）
- **先读后写、小步、改完自验证**：改前读懂相邻层，diff 最小，改完必跑 `make vet test`
- 🔴 **红线**：资金域（发放/消费）逻辑改动，必须过 review + 人审，见 `agent-merger.md`

## 按需查阅（相关任务才读）

- 通用编码纪律 → `.claude/docs/coding-rules.md`
- Go 代码规范（命名/注释/错误/并发/测试/提交）→ `.claude/docs/conventions.md`
- 开发工作流（persona）→ `.claude/docs/agent-developer.md`
- 审查工作流（persona）→ `.claude/docs/agent-code-reviewer.md`
- 合并工作流（persona，含合并前评审闸门）→ `.claude/docs/agent-merger.md`
- 并行 agent / 协作链 → `.claude/docs/parallel-agents.md`

## 三类扩展怎么区分

- **persona**（开发/审查/合并人格）：`.claude/docs/agent-*.md`，并行跑时由 plan 指定采用，
  **是约定，不是内置**，不能 `/` 调用
- **子代理 agent**（`.claude/agents/*.md`）：内置特性，主代理按需委派（如 `test-writer`）
- **skill / 命令**（`.claude/skills/`）：你手动敲，如 `/code-review-go`

## 并行协作链

`developer（提交）→ code-reviewer（merger 合并前评审闸门）→ merger（合并）`；
单人开发无此闸门（developer 自审即可），并行多 worker 才有；资金域评审闸门叠加人审。
详见 `.claude/docs/parallel-agents.md`。
