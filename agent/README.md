# 用 Go 手写一个对话 Agent

这个示例不使用 Agent 框架，直接调用 OpenAI 兼容接口，并通过 `godotenv` 从本地 `.env` 读取配置。

它的核心执行链是：

```text
用户输入 -> 加入对话历史 -> 调用模型 HTTP API -> 保存回复 -> 等待下一轮输入
```

`Agent.messages` 保存当前进程中的消息历史。每次请求都会把历史一起发给模型，因此模型能够理解多轮对话上下文。

## 运行

复制示例配置：

```bash
cp .env.example .env
```

然后编辑 `.env`：

```dotenv
OPENAI_API_KEY="你的 API Key"
OPENAI_MODEL="你的模型名称"

# OpenAI 官方接口可省略；兼容接口需要设置
OPENAI_BASE_URL="https://api.openai.com/v1"
```

`.env` 已加入仓库根目录的 `.gitignore`，不会提交到 Git；`.env.example` 不包含密钥，可以正常提交。

启动多轮对话：

```bash
go run .
```

程序同时支持 IDE 将仓库根目录设为 Working directory，此时会自动读取 `agent/.env`。

对话命令：

- `/clear`：清空对话上下文
- `/exit` 或 `/quit`：退出程序

也可以直接完成一次问答：

```bash
go run . "请用一句话介绍 Go"
```

## 代码组成

1. `Message`：一条 system、user 或 assistant 消息。
2. `Agent`：保存配置和对话历史。
3. `Agent.Chat`：组装 JSON、发起 HTTP 请求并保存模型回复。
4. `runConversation`：实现命令行多轮输入循环。

当前记忆只保存在内存中，退出程序后会消失。后续可以继续加入工具调用、流式输出、历史裁剪和持久化。
