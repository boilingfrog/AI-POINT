# AI 编程

## 前言

这里来聊聊如何使用 AI 来提高编码的效率。

我主要使用 **Claude Code** 进行编码，本文梳理一份「能配合 Claude Code 使用、增强其编码能力」的插件 / 扩展清单。Claude Code 的扩展体系大致分五块：**官方插件市场**、**MCP 服务**、**内置扩展机制**、**IDE 集成**、**内置技能与命令**。

> 安全提示：插件 / MCP 会以你的权限执行代码或访问数据，**只安装可信来源**的扩展。

## 本项目文档

| 文档 | 内容 |
|------|------|
| [README.md](README.md) | 本文。Claude Code 扩展体系清单（插件 / MCP / 内置扩展 / IDE / 技能命令）+ 推荐配置 |
| [claude-coding-agent.md](claude-coding-agent.md) | Claude 项目配置完整方案：核心配置文件、Skills、Hooks、自定义命令、Memory、语言特定配置、团队协作（含本项目说明） |
| [OpenSpec.md](OpenSpec.md) | OpenSpec 规格驱动开发指南：为什么用、最大作用、完整工作流、文件提交策略 |

---

## 一、插件市场（Plugins）

插件是把 技能 / 子 Agent / Hooks / MCP 打包在一起的扩展，用 `/plugin` 命令安装管理。

**官方市场**（启动 Claude Code 时自动可用）

```bash
/plugin install <name>@claude-plugins-official
```

**社区市场**（先添加再安装）

```bash
/plugin marketplace add anthropics/claude-plugins-community
/plugin install <name>@claude-community
```

常见可用插件分类：

| 分类 | 示例插件 | 作用 |
|------|---------|------|
| **代码智能（LSP）** | `typescript-lsp`、`pyright-lsp`、`rust-analyzer-lsp`、`gopls`… | 接入 Language Server，提供内联诊断、跳转、补全（需本地装对应语言服务器） |
| **外部集成** | `github`、`gitlab`、`figma`、`slack`、`sentry`、`vercel`、`supabase`、`firebase` | 预配好的 MCP，直接对接平台 |
| **安全** | `security-guidance` | 每次改文件自动做安全审查，可加项目自定义规则 |
| **开发工作流** | `commit-commands`、`pr-review-toolkit`、`plugin-dev`、`agent-sdk-dev` | Git 提交流程、PR 审查 Agent、插件 / SDK 开发 |
| **输出风格** | `explanatory-output-style`、`learning-output-style` | 改变回答风格（讲解式 / 学习式） |

安装范围：`user`（所有项目）、`project`（随 git 共享）、`local`（仅本机本项目）。

---

## 二、MCP 服务（最实用的能力扩展）

MCP（Model Context Protocol）是连接外部工具与数据源的开放标准，是给 Claude Code 「长手脚」最直接的方式。用 `claude mcp add` 接入。

| 用途 | 服务 | 接入示例 |
|------|------|---------|
| **代码库 / Git** | GitHub | `claude mcp add --transport http github https://api.githubcopilot.com/mcp/ --header "Authorization: Bearer <PAT>"` |
| **数据库** | PostgreSQL | `claude mcp add --transport stdio db -- npx -y @bytebase/dbhub --dsn "postgresql://..."` |
| **浏览器自动化** | Playwright / Puppeteer | 跑 E2E、抓页面、调前端 |
| **文档检索** | Context7 | 让模型查到最新的库 / 框架官方文档，减少幻觉 |
| **设计稿** | Figma | 设计稿转代码、读取组件 |
| **错误监控** | Sentry | `claude mcp add --transport http sentry https://mcp.sentry.dev/mcp` |
| **项目管理** | Jira / Linear / Notion | 读写任务、关联 Issue |

传输方式：`http`（远程，推荐）、`stdio`（本地进程）、`sse`（已弃用）。
配置存放：本地 `~/.claude.json`；项目级 `.mcp.json`（放仓库根目录，随 git 共享）。

---

## 三、内置扩展机制（不装插件也能定制）

这几样写在项目 `.claude/` 目录里就能用，是日常提效的核心：

| 机制 | 位置 | 作用 |
|------|------|------|
| **CLAUDE.md** | 仓库根 / 子目录 | 项目说明、约定、构建/测试命令，启动自动加载，让模型始终带上下文 |
| **Skills** | `.claude/skills/<name>/SKILL.md` | 把高频工作流封装成可复用技能，`/name` 触发或自动调用 |
| **Subagents** | `.claude/agents/` | 子 Agent 隔离上下文、并行处理大任务，主上下文保持干净 |
| **Hooks** | `.claude/settings.json` / `hooks.json` | 事件触发自动化：保存后跑 lint、提交前跑测试等 |
| **Slash Commands** | `.claude/commands/*.md` | 自定义斜杠命令，支持 `$ARGUMENTS`、`$FILE`、执行 bash |
| **Output Styles** | `.claude/output-styles/` | 持久改变回答格式 |

---

## 四、IDE 集成

把 Claude Code 接进编辑器，获得可视化 diff、@ 引用、会话历史等：

| IDE | 安装 | 亮点 |
|-----|------|------|
| **VS Code 扩展**（官方） | 扩展市场搜 "Claude Code" | 图形聊天面板、并排 diff、`@` 引用行号、权限模式（normal/plan/auto）、checkpoint 回滚 |
| **JetBrains 插件**（官方） | Settings → Plugins 搜 "Claude Code" | 适配 IDEA / GoLand / PyCharm / WebStorm 等，终端式聊天、远程开发 / WSL |
| **Cursor 及 VS Code 分支** | 同 VS Code 扩展 | 在 Cursor / Kiro 等分支里直接装 VS Code 版扩展 |

> 你在用 GoLand 系（看路径是 Go 项目），可以直接装 JetBrains 插件。

---

## 五、内置技能与命令（开箱即用）

不用安装，直接 `/` 调用的高频项：

- `/code-review` — 代码质量与安全审查
- `/simplify` — 重构 / 精简代码
- `/verify` — 跑起应用并观察行为
- `/run` — 启动并驱动应用
- `/deep-research` — 多源联网研究 + 带引用综述
- `/batch` — 在多个 worktree 里并行改动
- `/loop` — 按间隔重复执行某命令
- 管理类：`/mcp`、`/plugin`、`/agents`、`/hooks`、`/memory`、`/compact`、`/rewind`

---

## 六、给主用 Claude Code 的推荐配置

按「投入小、收益大」排序：

1. **写好 CLAUDE.md** — 零成本，立刻让每次对话都带上项目上下文。
2. **接 Context7 MCP** — 查最新官方文档，显著减少 API 幻觉。
3. **接 GitHub MCP + `pr-review-toolkit`** — 打通 Issue / PR 工作流。
4. **装对应语言的 LSP 插件**（Go 项目用 `gopls`） — 拿到内联诊断，改代码更准。
5. **配 Hooks** — 提交前自动 lint / 测试，挡住低级错误。
6. **装 IDE 集成**（GoLand 装 JetBrains 插件） — 可视化 diff 比纯终端更顺手。
7. **需要时接 Playwright MCP / 数据库 MCP** — 调前端、查库直接在对话里完成。

---

## 参考来源

- [Claude Code 官方文档](https://code.claude.com/docs)
- [官方插件市场](https://claude.com/plugins)
- [社区插件市场（GitHub）](https://github.com/anthropics/claude-plugins-community)
- [MCP 目录](https://claude.ai/directory)
- [Agent Skills 标准](https://agentskills.io)
