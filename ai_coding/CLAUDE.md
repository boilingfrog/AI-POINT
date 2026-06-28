# AI 编程学习项目

## 项目简介

本项目用于学习和整理 AI 辅助编程的最佳实践，包括 Claude Code、OpenSpec 等工具的使用方法和配置指南。

## 项目结构

```
.
├── README.md                    # 项目总览和扩展清单
├── OpenSpec.md                  # OpenSpec 使用文档
├── claude-coding-agent.md       # Claude 编码配置指南
└── CLAUDE.md                    # 本文件，项目说明
```

## 项目目标

- 整理 Claude Code 的插件、MCP、技能等扩展体系
- 记录 OpenSpec 规格驱动开发的工作流
- 提供可复用的配置模板和最佳实践
- 帮助团队快速上手 AI 辅助编程

## 内容组织

### README.md
介绍 Claude Code 的五大扩展机制：
1. 插件市场（官方和社区）
2. MCP 服务（外部工具集成）
3. 内置扩展（Skills、Commands、Hooks）
4. IDE 集成（VS Code、JetBrains）
5. 内置技能与命令

### OpenSpec.md
记录规格驱动开发（SDD）的核心流程：
- 为什么使用 OpenSpec（解决 AI 直接写代码的三大痛点）
- OpenSpec 的最大作用（让 AI 工作可预测、可审查、可追溯）
- 完整的工作流程和配置说明
- Git 提交策略

### claude-coding-agent.md
提供完整的 Claude 项目配置方案：
- 核心配置文件（CLAUDE.md、settings.json、.mcp.json）
- 技能（Skills）配置和示例
- Hooks 自动化
- 自定义命令
- Memory 持久化记忆
- 语言特定配置（Go、Python、TypeScript）
- 团队协作配置

## 代码规范

- 文档使用中文
- Markdown 格式，遵循通用最佳实践
- 代码示例使用 ``` 代码块
- 提交信息遵循 Conventional Commits

## 注意事项

- 本项目是文档项目，不包含可执行代码
- 配置示例中的 Token 使用环境变量占位符
- 所有配置模板仅供参考，需根据实际项目调整
