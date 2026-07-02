# Claude 编码 Agent 配置指南

本文档说明如何为项目配置 Claude Code，让 AI 编码助手更好地理解项目并提供精准帮助。

---

## 关于本项目（AI 编程学习项目）

本项目用于学习和整理 AI 辅助编程的最佳实践，包括 Claude Code、OpenSpec 等工具的使用方法和配置指南。

### 项目结构

```
.
├── README.md                    # 项目总览和扩展清单
├── OpenSpec.md                  # OpenSpec 使用文档
└── claude-coding-agent.md       # Claude 编码配置指南（本文件，含项目说明）
```

### 项目目标

- 整理 Claude Code 的插件、MCP、技能等扩展体系
- 记录 OpenSpec 规格驱动开发的工作流
- 提供可复用的配置模板和最佳实践
- 帮助团队快速上手 AI 辅助编程

### 内容组织

- **README.md**：介绍 Claude Code 的五大扩展机制（插件市场、MCP 服务、内置扩展、IDE 集成、内置技能与命令）
- **OpenSpec.md**：规格驱动开发（SDD）的核心流程——为什么用、最大作用、完整工作流、Git 提交策略
- **claude-coding-agent.md**：Claude 项目配置完整方案（本文件）——核心配置文件、Skills、Hooks、命令、Memory、语言特定配置、团队协作

### 本项目规范

- 文档使用中文，Markdown 格式
- 代码示例用 ``` 代码块
- 提交信息遵循 Conventional Commits
- 本项目是文档项目，不含可执行代码；配置示例中的 Token 用环境变量占位符，所有模板按实际项目调整

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

> ⚠️ settings.json 是**扁平结构**，字段名以[官方文档](https://code.claude.com/docs/en/settings)为准。不存在 `gitCommitStyle`、`preferredLanguage`、`codeStyle`、`testing`、`security` 这类字段——别凭直觉编。

```json
{
  "model": "claude-sonnet-4-5",
  "language": "zh-CN",
  "outputStyle": "default",
  "env": {
    "SOME_VAR": "value"
  },
  "permissions": {
    "allow": ["Bash(go test:*)", "Bash(golangci-lint run:*)"],
    "ask": ["Bash(git push:*)"],
    "deny": ["Read(./.env)", "Read(./secrets/**)"]
  },
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [
          { "type": "command", "command": "gofmt -w ." }
        ]
      }
    ]
  }
}
```

**常用字段**（节选，完整见官方文档）：

| 字段 | 作用 |
|------|------|
| `model` | 覆盖默认模型 |
| `language` | 回答语言（不是 `preferredLanguage`） |
| `outputStyle` | 输出风格 |
| `env` | 注入所有会话的环境变量 |
| `permissions` | 工具权限的 `allow` / `ask` / `deny` 规则 |
| `hooks` | 生命周期事件处理（见第三节） |

> 代码风格（缩进、行宽、命名）这类约定写进 `CLAUDE.md`，不是 settings.json 的字段。

### 1.3 .mcp.json - MCP 服务配置（项目共享）

**位置**：项目根目录 `.mcp.json`

**作用**：团队共享的 MCP 服务配置，提交到 git

**注意**：以下是配置结构示例，实际 URL 和认证方式需查阅各 MCP 服务的官方文档。

```json
{
  "mcpServers": {
    "github": {
      "type": "http",
      "url": "https://api.githubcopilot.com/mcp/",
      "headers": {
        "Authorization": "Bearer ${GITHUB_TOKEN}"
      }
    },
    "context7": {
      "type": "http",
      "url": "https://mcp.context7.com/mcp",
      "headers": {
        "CONTEXT7_API_KEY": "${CONTEXT7_TOKEN}"
      }
    }
  }
}
```

**配置说明**：
- 字段名是 `type`（`http` / `stdio` / `sse`），不是 `transport`
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

Hooks 在特定生命周期事件触发时自动执行 shell 命令。与靠 Claude「自觉遵守」的 CLAUDE.md 不同，hooks 由客户端强制执行。

### 3.1 配置 Hooks

**位置**：`.claude/settings.json` 的 `hooks` 键下（**没有**独立的 `hooks.json` 文件）

**作用**：在特定事件触发时自动执行命令，比如改完文件跑格式化、跑 lint

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [
          {
            "type": "command",
            "command": "gofmt -w ."
          }
        ]
      }
    ],
    "Stop": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "make lint"
          }
        ]
      }
    ]
  }
}
```

**结构说明**：
- 顶层 key 是**事件名**（如 `PostToolUse`），值是数组
- 每个元素有 `matcher` 和 `hooks` 列表
- `matcher` 只匹配**工具名**（如 `Bash`、`Edit|Write`）。要按命令模式过滤（如只在 `git commit`
  时跑），用 hook 上的 `if` 字段（权限规则语法，如 `"if": "Bash(git commit*)"`），**别写进 `matcher`**
- 每个 hook 是 `{ "type": "command", "command": "..." }`（可选 `if`）
- 命令通过 **stdin 收到 JSON 事件数据**（含 `tool_name`、`tool_input`——要改哪个文件从
  `tool_input.file_path` 取，bash 命令从 `tool_input.command` 取）；环境变量只有
  `$CLAUDE_PROJECT_DIR` 等少数几个，**没有 `$CLAUDE_FILE_PATHS`**（官方文档未记载，别用）

### 3.2 常用事件

| 事件 | 触发时机 | 说明 |
|------|----------|------|
| `PreToolUse` | 工具调用前 | `matcher` 匹配工具名 + 可选 `if` 过滤命令；`exit 2` 阻断（stderr 反馈给 Claude） |
| `PostToolUse` | 工具调用成功后 | 改文件后跑格式化/lint 用这个（`matcher: "Edit\|Write"`） |
| `UserPromptSubmit` | 用户提交提示词时 | 可注入额外上下文 |
| `SessionStart` | 会话开始/恢复 | 加载初始上下文 |
| `Stop` | Claude 完成一轮回复时 | 收尾检查（如跑测试） |
| `PreCompact` | 上下文压缩前 | — |

> ⚠️ **不存在** `beforeFileSave`、`afterFileSave`、`beforeCommit`、`afterCommit`、`beforePush` 这些事件，也没有 `${FILE}`、`${COMMIT_HASH}`、`${BRANCH}` 这类变量。要在「改文件后」做事，用 `PostToolUse` + `matcher`；要在「提交前」拦截，用 `PreToolUse` + `matcher: "Bash"` + `"if": "Bash(git commit*)"`（命令里 `exit 2` 阻断提交），或用 git 自身的 `pre-commit` 钩子。完整事件列表见[官方文档](https://code.claude.com/docs/en/hooks)。

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

Memory 让 Claude 跨会话记住项目特定的上下文和偏好。Claude Code 有两套记忆机制：

- **CLAUDE.md**：你手写的指令（见第一节），每次会话全量加载
- **自动记忆（Auto memory）**：Claude 根据你的纠正和偏好自己记下来的笔记

### 5.1 自动记忆的存储位置

**位置**：`~/.claude/projects/<project>/memory/`（**不在**项目的 `.claude/memory/`）

`<project>` 由 git 仓库推导，同一仓库的所有 worktree / 子目录共享一份。目录结构：

```
~/.claude/projects/<project>/memory/
├── MEMORY.md          # 简洁索引，每次会话加载前 200 行 / 25KB
├── debugging.md       # 调试笔记等话题文件（按需加载）
└── api-conventions.md
```

- 自动记忆默认开启，可在 `/memory` 里切换，或在 settings.json 设 `"autoMemoryEnabled": false`
- 要换存储位置，设 `autoMemoryDirectory`（绝对路径或 `~/` 开头）
- 自动记忆是**机器本地**的，不随 git 共享，也不跨机器同步

### 5.2 记忆文件的格式

话题文件是普通 Markdown，frontmatter 形如：

```markdown
---
name: api-versioning
description: API 版本管理策略
metadata:
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

写完话题文件后，在 `MEMORY.md` 里加一行指针即可。

### 5.3 通过对话创建记忆

在对话中说：
```
请记住：我们的日志格式使用结构化 JSON，字段包含 timestamp, level, message, context
```

Claude 会把它存进自动记忆。也可以说「加到 CLAUDE.md」让它写进项目说明文件。用 `/memory` 可浏览和编辑所有记忆文件。

---

## 六、完整项目配置检查清单

### 必备配置（优先级高）

- [ ] `CLAUDE.md` - 项目说明文件
- [ ] `.gitignore` - 排除 `.claude/settings.local.json`
- [ ] `.claude/settings.json` - 项目配置
- [ ] 配置基础 MCP 服务（如 GitHub、Context7）

### 推荐配置（优先级中）

- [ ] `.claude/settings.json` 配置 `hooks` 键 - 自动化格式化 / lint
- [ ] `.claude/skills/` - 常用技能（code-review、add-tests）
- [ ] `.claude/commands/` - 项目特定命令（deploy、db-migrate）
- [ ] `.mcp.json` - 团队共享的 MCP 配置

### 进阶配置（优先级低）

- [ ] `.claude/rules/` - 按路径作用域拆分的项目规则
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
- `.claude/settings.json`（含 `hooks` 配置）
- `.claude/skills/`
- `.claude/commands/`
- `.claude/rules/`
- `.mcp.json`（移除敏感 token，使用环境变量）

**不应该提交**：
- `.claude/settings.local.json` - 个人配置
- `.claude/.cache/` - 缓存文件

> 自动记忆存在 `~/.claude/projects/<project>/memory/`（机器本地，不在项目目录内），本来就不进 git，无需在 `.gitignore` 里处理。

### 8.2 .gitignore 规则

```gitignore
# Claude 本地配置和缓存
.claude/settings.local.json
.claude/.cache/
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

## 九、多 agent 工作流与子代理

单个 agent 够用就别上多 agent。只有当改动要一次铺开很多处（跨服务统一替换、
批量补测试/埋点），或想让「开发—审查—合并」各司其职时，才考虑拆给多个 agent。

先分清三种「角色载体」——它们在 Claude Code 里是不同的东西，别混：

| 载体 | 位置 | 是什么 | 怎么触发 |
|------|------|--------|----------|
| **子代理 agent** | `.claude/agents/*.md` | 内置特性：带独立 system prompt 的任务型代理 | 主代理按需委派，或你点名 |
| **persona（约定，非内置）** | 任意 `.md`（如 `.claude/docs/agent-*.md`） | 一份「工作流人格」文档，规定某类活怎么干 | 在 plan/提示里指明「采用某 persona」 |
| **skill / 命令** | `.claude/skills/`、`.claude/commands/` | 可复用的技能 / 命令 | 你手动敲 `/xxx` |

> ⚠️ **persona 不是 Claude Code 的内置概念**，只是「把一类工作的流程写成文档、让
> agent 采用」的约定用法。子代理才是内置特性。别把 persona 塞进 `.claude/agents/`
> 当子代理用——除非它真要作为独立委派单元。

### 9.1 子代理（内置）

放在 `.claude/agents/<name>.md`，frontmatter 三个字段：

```markdown
---
name: api-doc-updater
description: 接口有改动后同步 api/ 文档源并重生成产物。改完接口忘了同步时主动使用。
tools: Read, Grep, Glob, Edit, Bash   # 可选；省略则继承主代理全部工具
---

（正文即该子代理的 system prompt：职责、步骤、红线）
```

要点：
- **description 写清「何时用」**：主代理靠它判断要不要委派，写成触发条件，别只写功能
- **tools 给最小集**：只读型代理别给 `Edit`/`Bash`
- 子代理有**独立上下文**，适合「脏活隔离」（大范围搜索、生成物重建），产出汇报回主代理

### 9.2 persona 工作流：把一条流水线拆成角色

大改动想让职责分明时，为每类活写一份 persona 文档，各自规定流程，再由并行的
worker 分别采用。典型是「开发—审查—合并」三段：

```
developer（N worker，各自 worktree，持续小步提交）
      │
      ▼
code-reviewer  ── merger 合并前评审闸门
      │            REJECT/NEEDS_WORK → 退回 developer；APPROVE → 才合
      ▼
merger（评审 → git merge → lint+build+test → 循环）
```

三个要点，是这套能跑顺的关键：

1. **评审要「串进链路」，别悬空**：光写一份 reviewer persona 没用——必须在 merger
   的「合并前」显式插一步调它（`/code-review <diff>` 或采用 reviewer persona），
   否则 review 永远不会发生。
2. **规则要强制前置**：persona 开头用 `> 见 coding-rules.md` 这种软引用，agent 常
   常不读。把「先读规则」提成每个 persona 的**显式第一步**，比软引用可靠得多。
3. **高风险改动叠人审**：评审闸门是自动兜底，不替代人审；资金/权限等高风险域，
   闸门过了也要人确认再合。

### 9.3 hook 能把 agent 串起来吗？

能，但有边界，别用错：

- `SubagentStop` 事件在**子代理结束时**触发，可用来在**同一会话内**接力（如某子代理
  跑完自动触发一次检查）。
- 但如果 worker 是**各自独立进程**（tmux + worktree 多开 `claude`），它们不在同一会话
  里，而 `Stop`/`SubagentStop` 按进程/会话触发，**串不起跨会话的多个 worker**。这种
  编排该由**外部脚本**（启动/盯梢/合并）来做，不是 settings.json 的 hook。

一句话：**同一会话内的子代理接力用 `SubagentStop`；跨会话多进程编排用外部脚本 +
persona 约定。**

### 9.4 什么时候别上多 agent

- 各并行单元**文件不相交**才安全，否则合并冲突；大仓按「一个服务/一个包」切最稳
- 涉及资金/权限/数据一致性等**需逐行把关**的改动，串行 + 人审，别并行
- 单元改完**必须自验证**（build + test）再交给 merger 合

---

## 十、常见问题

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

## 十一、参考资源

- [Claude Code 官方文档](https://code.claude.com/docs)
- [Skills 官方文档](https://code.claude.com/docs/en/skills)（frontmatter 字段、参数替换的权威来源）
- [MCP 协议文档](https://modelcontextprotocol.io)
- [Agent Skills 开放标准](https://agentskills.io)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)（Go 代码审查官方规范）
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [golangci-lint Linters](https://golangci-lint.run/docs/linters/)
