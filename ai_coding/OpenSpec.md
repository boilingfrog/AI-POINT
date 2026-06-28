# OpenSpec 使用文档

OpenSpec 是**规格驱动开发（SDD）**框架：让 AI 先写清楚要做什么（规格），你确认后再照着实现，告别「凭感觉写代码」。

核心循环四步：`explore`（探索）→ `propose`（提案）→ `apply`（实现）→ `archive`（归档）。

---

## 1. 安装

**先装 Node.js 18+：**

- **macOS**：`brew install node`，或 [nodejs.org](https://nodejs.org/) 下载 `.pkg` 安装包
- **Windows**：[nodejs.org](https://nodejs.org/) 下载 `.msi`（勾选 Add to PATH），或 `winget install OpenJS.NodeJS.LTS`

验证：`node -v` 和 `npm -v` 都能打印版本号。

**再装 OpenSpec（两平台命令相同）：**

```bash
npm install -g @fission-ai/openspec@latest
openspec --version   # 能打印版本号即成功
```

> macOS 报 `EACCES` 用 `sudo`；Windows 用管理员 PowerShell 执行。

## 2. 初始化

```bash
cd your-project
openspec init --tools claude   # 指定 Claude Code，省去交互选择
```

生成 `.claude/commands/opsx/`（5 个斜杠命令）和 `openspec/`（`specs/` 现状规格、`changes/` 在途改动）。

> init 后**重启一次 IDE / AI 助手**，斜杠命令才生效。

## 3. 工作流（在 Claude Code 对话框里输入）

```
/opsx:explore  深色模式怎么做     # （可选）理清需求
/opsx:propose  add-dark-mode     # 生成提案 + 规格 + 任务
                                 # → 审阅：openspec show / validate <name>
/opsx:apply                      # 按任务实现代码
/opsx:archive                    # 验证通过后归档，合并进 specs/
```

**审阅是最关键的一步**：propose 后逐份检查 `proposal.md`（做什么）、`tasks.md`（拆解）、`design.md`（方案），不满意让 AI 改，确认无误再 apply。

## 4. 文件提交

`openspec init` 生成的文件中，**应该提交的**：

- `.claude/commands/opsx/` - 斜杠命令配置（团队成员也需要）
- `.claude/skills/openspec-*/` - 技能配置
- `openspec/specs/` - 已归档的规格文档（项目的规格说明）
- `openspec/config.yaml` - OpenSpec 配置文件

**不需要提交的**（添加到 `.gitignore`）：

- `openspec/changes/` - 工作中的提案和改动（类似工作区，未完成的内容）
- `.claude/settings.local.json` - 本地设置
- `openspec/.cache/` - 缓存文件（如果有的话）

**推荐的 .gitignore 规则**：

```gitignore
# OpenSpec 工作区（在途改动）
openspec/changes/

# Claude 本地设置
.claude/settings.local.json

# OpenSpec 缓存
openspec/.cache/
```

## 5. 排查

| 现象 | 处理 |
|------|------|
| `command not found` | 安装失败 / PATH 问题，重装或检查 npm 全局 bin |
| `/opsx:` 命令不识别 | 确认项目根目录跑过 `openspec init`，并**重启 IDE** |
| `validate` 报 Nothing to validate | 加范围：`openspec validate --all` 或带名字 |

---

参考：[GitHub](https://github.com/Fission-AI/OpenSpec) · [官网](https://openspec.pro/)（基于 v1.5.0 实测）
