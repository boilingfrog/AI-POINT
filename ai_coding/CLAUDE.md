# AI 编程（ai_coding）

## 项目简介

整理 AI 辅助编程的最佳实践，重点是 Claude Code 和 OpenSpec 的使用方法与配置指南。纯文档项目，不含可执行代码。

## 文档结构

- `README.md` — Claude Code 扩展体系清单（插件 / MCP / 内置扩展 / IDE / 技能命令）
- `claude-coding-agent.md` — Claude 项目配置完整方案
- `OpenSpec.md` — OpenSpec 规格驱动开发指南

## 写作规范

- 文档用中文、Markdown 格式
- 缩进用 2 空格，正文行宽控制在 100 字符以内
- 代码示例用 ``` 代码块
- 提交信息遵循 Conventional Commits（如 `docs:`、`fix:`、`feat:`）

## 内容准确性要求

- 涉及 Claude Code 的字段名、配置位置、事件名、命令时，以[官方文档](https://code.claude.com/docs)为准，不要凭直觉编造
- settings.json 是扁平结构，没有 `gitCommitStyle` / `preferredLanguage` / `codeStyle` 这类字段
- 配置示例中的 token 用环境变量占位符，不写真实凭据
