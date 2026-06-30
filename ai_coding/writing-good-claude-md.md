# 如何写好 CLAUDE.md

CLAUDE.md 是 Claude Code 每次会话自动加载的项目说明文件，写得好能让 AI 准确理解项目，写得差会浪费 token 还影响效果。本文总结写好 CLAUDE.md 的实用方法。

---

## 一、核心原则

### 1.1 把它当成「给同事的项目入职文档」

写 CLAUDE.md 就像给新同事写入职指南：
- ✅ 写项目是干什么的、怎么跑起来、有哪些约定
- ❌ 不写「请帮我写代码」「你是专家」这类对话式指令

**好例子**：
```markdown
## 项目简介
电商后台 API，用 Go + PostgreSQL，提供商品/订单/用户管理接口。

## 开发命令
- `make run` - 启动开发服务器（监听 8080 端口）
- `make test` - 运行单元测试
- `make lint` - 代码检查（golangci-lint）
```

**坏例子**：
```markdown
你是 Go 专家，请帮我写高质量代码。遇到问题先思考再回答。
```

### 1.2 控制长度：200 行以内

- **原因**：CLAUDE.md 全文加载进每次会话，越长消耗 token 越多，Claude 遵守的可靠性越低
- **超长怎么办**：拆到 `.claude/rules/`，按路径作用域加载（见第四节）

### 1.3 用事实，不要空泛建议

| 空泛 ❌ | 具体 ✅ |
|---------|---------|
| "代码要规范" | "函数名用驼峰，包名用小写单字" |
| "注意安全" | "SQL 查询必须用 `sqlx.In()` 防注入，敏感字段用 `json:\"-\"` 避免泄漏" |
| "测试很重要" | "新增 API handler 必须在 `handlers_test.go` 里加对应测试" |

---

## 二、必写内容

### 2.1 项目简介
**一句话说清楚项目是干什么的**，让 Claude 建立基本认知。

```markdown
# 电商后台 API

为前端 Web/App 提供商品、订单、用户管理接口，部署在 Kubernetes。
```

### 2.2 技术栈
列出语言、框架、数据库、部署方式，**带版本号**。

```markdown
## 技术栈
- 语言：Go 1.21+
- Web 框架：Gin 1.9
- 数据库：PostgreSQL 15（通过 sqlx 操作）
- 缓存：Redis 7
- 部署：Docker + Kubernetes
- CI/CD：GitHub Actions
```

### 2.3 项目结构
用树形图或注释说明主要目录的作用。

```markdown
## 项目结构
```
.
├── cmd/server/      # 程序入口
├── internal/
│   ├── api/         # HTTP handlers
│   ├── service/     # 业务逻辑
│   ├── model/       # 数据模型
│   └── db/          # 数据库操作
├── migrations/      # 数据库迁移（goose）
└── tests/           # 集成测试
```
```

### 2.4 开发命令
列出启动、测试、lint、构建等常用命令，**不要让 Claude 猜**。

```markdown
## 开发命令
- `make run` - 启动开发服务器（需先 `docker-compose up -d` 启动 DB）
- `make test` - 运行所有单元测试
- `make lint` - golangci-lint 检查
- `make migrate-up` - 应用数据库迁移
- `make build` - 构建生产二进制
```

### 2.5 代码规范
写具体的命名、格式、注释要求。

```markdown
## 代码规范
- 命名：函数/类型用驼峰，包名小写单字，常量全大写下划线
- 格式：`gofmt` + 行宽 120
- 注释：所有 exported 函数/类型必须有文档注释
- 错误处理：用 `errors.Wrap()` 包装，不吞错误
- 测试：handlers 必须有单元测试，service 层核心逻辑必须有测试
```

---

## 三、可选但实用的内容

### 3.1 架构决策
关键设计选择的理由，避免 Claude 改掉你精心设计的部分。

```markdown
## 架构约定
- 认证：JWT token 放在 `Authorization: Bearer <token>` header
- 错误响应：统一用 `{"error": {"code": "ERR_CODE", "message": "..."}}`
- ID 生成：用 UUID v4，不用自增 ID（分布式考虑）
- 事务：跨表操作必须在 `db.Transaction()` 里执行
```

### 3.2 外部依赖
第三方服务、API、认证方式。

```markdown
## 外部依赖
- 支付：对接 Stripe API（sandbox key 在 `.env.local`）
- 文件存储：AWS S3（`us-west-2` 区域）
- 短信：Twilio（测试环境用 mock，生产用真实 API）
```

### 3.3 常见问题
团队成员经常踩的坑。

```markdown
## 常见问题
- **数据库连接失败**：确认 `docker-compose up -d` 已启动 PostgreSQL
- **测试失败**：清理测试数据库 `make test-db-reset`
- **Lint 报错 `govet`**：本地跑 `go mod tidy` 更新依赖
```

---

## 四、进阶：拆分大型 CLAUDE.md

当 CLAUDE.md 超过 200 行，用 `.claude/rules/` 按作用域拆分。

### 4.1 无条件加载的规则
放在 `.claude/rules/` 根目录，无需 frontmatter，每次会话都加载。

**示例**：`.claude/rules/api-design.md`
```markdown
# API 设计规范

所有 API 遵循 RESTful 风格：
- GET /resources - 列表
- GET /resources/:id - 详情
- POST /resources - 创建
- PUT /resources/:id - 更新
- DELETE /resources/:id - 删除
```

### 4.2 按路径条件加载的规则
带 `paths` frontmatter，只在操作匹配文件时加载。

**示例**：`.claude/rules/database.md`
```markdown
---
paths:
  - "internal/db/**/*.go"
  - "migrations/*.sql"
---

# 数据库操作规范

1. 所有查询用参数化，禁止字符串拼接 SQL
2. 事务里不能调用外部 API（会阻塞连接池）
3. 迁移文件命名：`YYYYMMDDHHMMSS_description.sql`
```

**好处**：只在改 `internal/db/` 或 `migrations/` 下的文件时，这条规则才进上下文，节省 token。

---

## 五、避坑指南

### ❌ 别把操作手册都塞进去
**错误做法**：
```markdown
## 部署流程
1. SSH 到服务器 `ssh user@prod-server`
2. 拉取代码 `git pull origin main`
3. 构建 `docker build -t app:latest .`
4. 重启容器 `docker-compose up -d`
5. 检查日志 `docker logs -f app`
```

**正确做法**（部署操作封装成 Skill）：
- CLAUDE.md 只写：`部署：执行 /deploy [env] 命令`
- 详细步骤写在 `.claude/skills/deploy/SKILL.md`，手动触发

### ❌ 别写「Claude 应该怎么做」
**错误做法**：
```markdown
- 写代码前先理解需求
- 多问用户确认细节
- 遇到不确定的地方要查文档
```

**为什么错**：这些是 Claude 自身行为，不是项目特定的知识。CLAUDE.md 只写「项目的事实」。

### ❌ 别和 settings.json 重复
**错误做法**（在 CLAUDE.md 里写）：
```markdown
## 配置
- 回答用中文
- 使用 sonnet-4 模型
```

**正确做法**：这些配在 `.claude/settings.json`：
```json
{
  "language": "zh-CN",
  "model": "claude-sonnet-4-5"
}
```

---

## 六、检查清单

写完 CLAUDE.md 后，用这份清单自查：

**必备项**：
- [ ] 一句话项目简介
- [ ] 技术栈（带版本号）
- [ ] 开发命令（启动/测试/lint）
- [ ] 代码规范（命名/格式/注释）

**质量检查**：
- [ ] 长度 ≤ 200 行（超长考虑拆 rules）
- [ ] 没有空泛建议，都是可验证的事实
- [ ] 没有重复 settings.json 的配置
- [ ] 没有写操作流程（应该封装成 Skill）

**验证方式**：
在 Claude Code 里问：「这个项目的技术栈是什么？」「怎么跑测试？」——如果 Claude 能准确回答，说明 CLAUDE.md 生效了。

---

## 七、示例模板

### 7.1 最小模板（适合小项目）

```markdown
# 项目名

简短描述项目用途。

## 技术栈
- 语言 + 版本
- 主要框架

## 开发命令
- 启动
- 测试
- lint

## 代码规范
- 命名
- 格式
```

### 7.2 完整模板（适合团队项目）

```markdown
# 项目名

## 项目简介
一句话说明用途和部署位置。

## 技术栈
- 语言：版本
- 框架：版本
- 数据库：版本
- 部署：方式

## 项目结构
```
目录树 + 注释
```

## 开发命令
- 启动
- 测试
- lint
- 构建
- 迁移

## 代码规范
- 命名
- 格式
- 注释
- 错误处理
- 测试覆盖

## 架构约定
- 认证方式
- 错误响应格式
- 关键设计决策

## 外部依赖
- 第三方服务
- API 端点
- 认证方式

## 注意事项
- 常见坑
- 环境要求
```

---

## 八、相关资源

- [Claude Code 官方文档 - Memory](https://code.claude.com/docs/en/memory)
- [本项目：claude-coding-agent.md](./claude-coding-agent.md) - 完整配置方案
- [本项目：README.md](./README.md) - Claude Code 扩展清单

---

**记住**：CLAUDE.md 是给 AI 看的项目说明书，不是对话指令。写好它，Claude 就能像懂项目的同事一样工作。
