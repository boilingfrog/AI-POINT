# Claude 编码 Agent 配置指南

本文档说明如何为项目配置 Claude Code，让 AI 编码助手更好地理解项目并提供精准帮助。

---

## 一、核心配置文件

### 1.1 CLAUDE.md - 项目说明文件

**位置**：项目根目录 `CLAUDE.md`

**作用**：每次对话自动加载，让 Claude 始终了解项目上下文

**必填内容**：

```markdown
# 项目名称

## 项目简介
一句话说明这个项目是做什么的

## 技术栈
- 语言：Go / Python / TypeScript 等
- 框架：具体框架和版本
- 数据库：PostgreSQL / MySQL 等
- 部署：Docker / Kubernetes 等

## 项目结构
```
src/          # 源代码
  api/        # API 接口
  service/    # 业务逻辑
  dao/        # 数据访问
tests/        # 测试
docs/         # 文档
```

## 开发命令
- 构建：`make build`
- 测试：`make test`
- 运行：`make run`
- Lint：`make lint`

## 代码规范
- 命名：驼峰 / 蛇形
- 注释：必须为公共 API 添加文档注释
- 测试：新功能必须有单元测试
- 提交：遵循 Conventional Commits

## 注意事项
- 敏感配置放在环境变量，不提交到 git
- 数据库迁移必须可回滚
- API 变更需要更新 OpenAPI 文档
```

### 1.2 .claude/settings.json - 项目配置

**位置**：`.claude/settings.json`

**作用**：项目级别的 Claude 行为配置

```json
{
  "gitCommitStyle": "conventional",
  "preferredLanguage": "zh-CN",
  "codeStyle": {
    "indentation": "spaces",
    "indentSize": 4,
    "lineLength": 120
  },
  "testing": {
    "framework": "pytest",
    "runBeforeCommit": true
  },
  "security": {
    "scanOnSave": true,
    "blockSecretsCommit": true
  }
}
```

### 1.3 .mcp.json - MCP 服务配置（项目共享）

**位置**：项目根目录 `.mcp.json`

**作用**：团队共享的 MCP 服务配置，提交到 git

**注意**：以下是配置结构示例，实际 URL 和认证方式需查阅各 MCP 服务的官方文档。

```json
{
  "mcpServers": {
    "github": {
      "transport": "http",
      "url": "https://api.github.com/mcp/",
      "headers": {
        "Authorization": "Bearer ${GITHUB_TOKEN}"
      }
    },
    "context7": {
      "transport": "http",
      "url": "https://api.context7.com/mcp/",
      "headers": {
        "Authorization": "Bearer ${CONTEXT7_TOKEN}"
      }
    }
  }
}
```

**配置说明**：
- 实际 MCP 服务 URL 请查阅官方文档或使用 `claude mcp add` 命令自动配置
- 推荐使用 `claude mcp add` 命令添加 MCP 服务，而不是手动编辑此文件
- 可以使用 `claude mcp list` 查看已配置的 MCP 服务

**环境变量说明**：
- `GITHUB_TOKEN`：GitHub Personal Access Token（从 https://github.com/settings/tokens 创建）
- `CONTEXT7_TOKEN`：Context7 API Token（从 https://context7.com 获取）

---

## 二、技能（Skills）配置

Skills 是可复用的工作流，封装高频任务。

### 2.1 创建 Skill

**目录结构**：
```
.claude/skills/
  <skill-name>/
    SKILL.md          # 技能定义
    examples/         # 使用示例（可选）
```

### 2.2 示例：代码审查技能

**文件**：`.claude/skills/security-review/SKILL.md`

```markdown
---
name: security-review
description: 对代码进行安全审查，检查常见漏洞。审查代码安全性时使用。
disable-model-invocation: true
allowed-tools: Read Grep Glob
arguments: [target]
argument-hint: [file-or-dir]
---

# Security Review Skill

## 任务
审查 `$target` 指定的文件或目录，检查以下安全问题：

1. SQL 注入风险
2. XSS 漏洞
3. 硬编码密钥
4. 不安全的反序列化
5. 路径遍历
6. 认证/授权缺陷

## 输出
生成安全审查报告，包含：
- 发现的问题清单（按严重程度排序）
- 每个问题的位置和修复建议
- 安全评分（1-10）

## 示例
\`/security-review src/api/user.go\` → `$target` = `src/api/user.go`
\`/security-review src/\` → `$target` = `src/`
```

**frontmatter 字段说明**（对齐 [官方 Agent Skills 规范](https://code.claude.com/docs/en/skills)）：

| 字段 | 必填 | 说明 |
|------|------|------|
| `name` | 否 | 显示名，默认用目录名 |
| `description` | 是 | 描述用途，Claude 据此决定何时自动加载 |
| `disable-model-invocation` | 否 | `true` 表示禁止 Claude 自动触发，只能手动 `/name` 调用。默认 false |
| `user-invocable` | 否 | `false` 时从 `/` 菜单隐藏（仅作背景知识）。默认 true |
| `allowed-tools` | 否 | skill 激活时无需询问即可使用的工具（空格或逗号分隔） |
| `disallowed-tools` | 否 | skill 激活时禁用的工具 |
| `arguments` | 否 | 命名位置参数，用于内容里的 `$name` 替换 |
| `argument-hint` | 否 | 自动补全时提示期望的参数，如 `[filename] [format]` |
| `model` | 否 | skill 激活时使用的模型 |
| `context` | 否 | 设为 `fork` 时在子 agent 中执行 |

> ⚠️ 没有 `invocation` 和 `autoInvoke` 字段——要"禁止自动触发"请用 `disable-model-invocation: true`。

**参数替换**（在 SKILL.md 正文中可用）：

| 占位符 | 说明 |
|--------|------|
| `$ARGUMENTS` | 调用时传入的全部参数（原样字符串） |
| `$ARGUMENTS[N]` / `$N` | 按 0 起始的位置取参数，如 `$0` 是第一个 |
| `$name` | 需先在 frontmatter 的 `arguments` 列表声明，按顺序映射位置 |

> 没有 `$FILE` 这类内置变量。要接收文件路径，用 `$0`，或在 `arguments: [file]` 声明后用 `$file`。

### 2.3 常用技能推荐

| 技能名 | 作用 | 调用方式 |
|--------|------|----------|
| `code-review` | 代码质量审查 | `/code-review <file>` |
| `add-tests` | 为现有代码生成测试 | `/add-tests <file>` |
| `refactor` | 重构代码以提高可维护性 | `/refactor <file>` |
| `generate-docs` | 生成 API 文档 | `/generate-docs <file>` |
| `security-review` | 安全审查 | `/security-review <file>` |

---

## 三、Hooks 自动化

Hooks 在特定事件触发时自动执行任务。

### 3.1 配置 Hooks

**位置**：`.claude/hooks.json`（在项目的 `.claude` 目录下）

**作用**：定义在特定事件（保存、提交、推送等）时自动执行的任务

```json
{
  "beforeFileSave": [
    {
      "name": "format",
      "command": "gofmt -w ${FILE}",
      "description": "格式化 Go 代码"
    }
  ],
  "beforeCommit": [
    {
      "name": "lint",
      "command": "make lint",
      "description": "运行代码检查",
      "abortOnFailure": true
    },
    {
      "name": "test",
      "command": "make test-quick",
      "description": "运行快速测试",
      "abortOnFailure": true
    }
  ],
  "afterCommit": [
    {
      "name": "notify",
      "command": "echo '提交成功: ${COMMIT_HASH}'",
      "description": "提交成功通知"
    }
  ]
}
```

### 3.2 可用事件

| 事件 | 触发时机 | 可用变量 |
|------|----------|----------|
| `beforeFileSave` | 文件保存前 | `${FILE}` |
| `afterFileSave` | 文件保存后 | `${FILE}` |
| `beforeCommit` | Git 提交前 | `${FILES}` |
| `afterCommit` | Git 提交后 | `${COMMIT_HASH}` |
| `beforePush` | Git 推送前 | `${BRANCH}` |

---

## 四、自定义命令

### 4.1 创建命令

**目录结构**：
```
.claude/commands/
  deploy.md          # /deploy 命令
  db-migrate.md      # /db-migrate 命令
  code-stats.md      # /code-stats 命令
```

### 4.2 示例：部署命令

**文件**：`.claude/commands/deploy.md`

> 说明：`.claude/commands/` 与 `.claude/skills/` 现已合并——两者都能创建 `/deploy` 命令。旧的 commands 文件继续可用，但官方推荐用 skills（支持目录、附属文件、自动加载等）。部署这类手动触发的操作应加 `disable-model-invocation: true`，避免 Claude 自动执行。

```markdown
---
name: deploy
description: 部署应用到指定环境
disable-model-invocation: true
argument-hint: [dev|staging|prod]
---

# Deploy Command

部署应用到指定环境（dev / staging / prod）

## 使用方式
```
/deploy dev          # 部署到开发环境
/deploy staging      # 部署到预发布环境
/deploy prod         # 部署到生产环境（需要确认）
```

## 执行步骤

1. 检查当前分支和工作区状态
2. 运行测试套件
3. 构建 Docker 镜像
4. 推送到容器仓库
5. 更新 Kubernetes 部署
6. 验证服务健康状态

## 环境要求

`$ARGUMENTS` 必须是 `dev` / `staging` / `prod` 之一

生产环境部署需要：
- 主分支已合并
- 所有测试通过
- 用户明确确认

**注意**：`$ARGUMENTS` 变量自动包含命令后的所有参数，例如：
- `/deploy dev` → `$ARGUMENTS` = `dev`
- `/deploy staging --force` → `$ARGUMENTS` = `staging --force`

## 执行命令

```bash
#!/bin/bash
ENV=$ARGUMENTS

if [ "$ENV" != "dev" ] && [ "$ENV" != "staging" ] && [ "$ENV" != "prod" ]; then
  echo "错误：环境必须是 dev / staging / prod"
  exit 1
fi

if [ "$ENV" = "prod" ]; then
  echo "⚠️  即将部署到生产环境，请确认"
  # Claude 会要求用户确认
fi

echo "部署到 $ENV 环境..."
make deploy ENV=$ENV
```
```

---

## 五、Memory - 持久化记忆

Memory 让 Claude 记住项目特定的上下文和偏好。

### 5.1 手动创建记忆

**位置**：`.claude/memory/`

**文件格式**：Markdown with Frontmatter

**示例**：`.claude/memory/api-versioning.md`

```markdown
---
name: api-versioning
description: API 版本管理策略
type: project
---

# API 版本管理

本项目使用 URL 路径版本控制：

- `/api/v1/` - 稳定版本，向后兼容
- `/api/v2/` - 新版本，可能有破坏性变更

**版本升级规则：**
1. 不能删除现有字段，只能标记为 deprecated
2. 新增字段必须有默认值
3. 破坏性变更必须创建新版本
4. 旧版本至少支持 6 个月

**Why：** 确保客户端不会因 API 变更而中断

**How to apply：** 添加新 API 时，先检查是否可以向后兼容地扩展现有版本
```

### 5.2 通过对话创建记忆

在对话中说：
```
请记住：我们的日志格式使用结构化 JSON，字段包含 timestamp, level, message, context
```

Claude 会自动创建记忆文件。

---

## 六、完整项目配置检查清单

### 必备配置（优先级高）

- [ ] `CLAUDE.md` - 项目说明文件
- [ ] `.gitignore` - 排除 `.claude/settings.local.json`
- [ ] `.claude/settings.json` - 项目配置
- [ ] 配置基础 MCP 服务（如 GitHub、Context7）

### 推荐配置（优先级中）

- [ ] `.claude/hooks.json` - 自动化 lint 和 test
- [ ] `.claude/skills/` - 常用技能（code-review、add-tests）
- [ ] `.claude/commands/` - 项目特定命令（deploy、db-migrate）
- [ ] `.mcp.json` - 团队共享的 MCP 配置

### 进阶配置（优先级低）

- [ ] `.claude/memory/` - 项目特定知识
- [ ] `.claude/agents/` - 子 Agent 配置
- [ ] `.claude/output-styles/` - 自定义输出格式
- [ ] LSP 插件配置（如 `gopls`、`pyright`）

---

## 七、语言特定配置

### 7.1 Go 项目

**CLAUDE.md 增加**：
```markdown
## Go 特定
- Go 版本：1.21+
- 包管理：go mod
- 代码检查：golangci-lint
- 测试框架：标准库 testing + testify
- Mock：gomock

## 开发命令
- `go run ./cmd/server` - 启动服务
- `go test ./...` - 运行所有测试
- `golangci-lint run` - 代码检查
- `go mod tidy` - 整理依赖
```

**推荐 MCP**：
- `gopls` LSP 插件

### 7.2 Python 项目

**CLAUDE.md 增加**：
```markdown
## Python 特定
- Python 版本：3.11+
- 包管理：Poetry / pip
- 代码检查：ruff / pylint
- 测试框架：pytest
- 类型检查：mypy

## 开发命令
- `poetry run python -m app.main` - 启动服务
- `poetry run pytest` - 运行测试
- `poetry run ruff check .` - 代码检查
- `poetry run mypy .` - 类型检查
```

**推荐 MCP**：
- `pyright-lsp` 插件
- `poetry-mcp` 包管理

### 7.3 TypeScript/Node.js 项目

**CLAUDE.md 增加**：
```markdown
## Node.js 特定
- Node 版本：20+
- 包管理：pnpm / npm / yarn
- 代码检查：ESLint
- 测试框架：Vitest / Jest
- 类型检查：TypeScript

## 开发命令
- `pnpm dev` - 启动开发服务器
- `pnpm test` - 运行测试
- `pnpm lint` - 代码检查
- `pnpm build` - 构建生产版本
```

**推荐 MCP**：
- `typescript-lsp` 插件
- `playwright-mcp` E2E 测试

---

## 八、团队协作配置

### 8.1 提交到 Git 的文件

**应该提交**：
- `CLAUDE.md`
- `.claude/settings.json`
- `.claude/skills/`
- `.claude/commands/`
- `.claude/hooks.json`
- `.mcp.json`（移除敏感 token，使用环境变量）

**不应该提交**：
- `.claude/settings.local.json` - 个人配置
- `.claude/memory/` - 个人记忆（可选）
- `.claude/.cache/` - 缓存文件

### 8.2 .gitignore 规则

```gitignore
# Claude 本地配置和缓存
.claude/settings.local.json
.claude/.cache/
.claude/memory/  # 如果不想共享个人记忆
```

### 8.3 团队 README

在项目 `README.md` 中添加 Claude 配置说明：

```markdown
## 使用 Claude Code 开发

本项目配置了 Claude Code 集成。首次使用：

1. 安装 Claude Code：https://claude.ai/code
2. 配置 MCP 服务：
   ```bash
   export GITHUB_TOKEN="your_github_token"
   export CONTEXT7_TOKEN="your_context7_token"
   ```
3. 启动 Claude：`claude code`
4. 查看项目说明：`CLAUDE.md`
5. 可用技能：`/code-review`、`/add-tests`、`/security-review`
6. 可用命令：`/deploy`、`/db-migrate`

配置详见 `claude-coding-agent.md`
```

---

## 九、常见问题

### Q1：CLAUDE.md 和 README.md 有什么区别？

- **README.md**：给人看的，介绍项目、安装、使用
- **CLAUDE.md**：给 AI 看的，提供开发上下文、命令、规范

### Q2：Skills 和 Commands 有什么区别？

- **Skills**：复杂的工作流，可以调用多个工具，有自己的 prompt 和逻辑
- **Commands**：简单的命令封装，通常执行一个 bash 脚本或一组固定操作

### Q3：MCP 配置中的 Token 如何管理？

- 使用环境变量：`${GITHUB_TOKEN}`
- 团队共享 `.mcp.json` 结构，个人设置 token
- 敏感 token 不提交到 git

### Q4：如何测试 Claude 配置是否生效？

**检查 CLAUDE.md 是否加载**
```
启动 Claude 后，询问："这个项目的技术栈是什么？"
如果 Claude 能准确回答 CLAUDE.md 中记录的技术栈，说明加载成功
```

**检查 Skills 是否可用**
```
输入：/code-review
如果命令被识别（而不是报错"未知命令"），说明 Skill 加载成功
```

**检查 MCP 是否连接**
```
输入：/mcp list
会列出所有已配置的 MCP 服务及其状态
```

**检查 Hooks 是否生效**
```
修改一个文件并保存，观察是否：
- 自动格式化（如果配置了 beforeFileSave hook）
- 提交前自动运行测试（如果配置了 beforeCommit hook）
```

**检查 Commands 是否可用**
```
输入：/help
会列出所有可用的命令，包括自定义命令
```

---

## 十、参考资源

- [Claude Code 官方文档](https://code.claude.com/docs)
- [Skills 官方文档](https://code.claude.com/docs/en/skills)（frontmatter 字段、参数替换的权威来源）
- [MCP 协议文档](https://modelcontextprotocol.io)
- [Agent Skills 开放标准](https://agentskills.io)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)（Go 代码审查官方规范）
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [golangci-lint Linters](https://golangci-lint.run/docs/linters/)
