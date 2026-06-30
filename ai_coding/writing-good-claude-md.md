<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [如何编写优秀的 CLAUDE.md](#%E5%A6%82%E4%BD%95%E7%BC%96%E5%86%99%E4%BC%98%E7%A7%80%E7%9A%84-claudemd)
  - [原则：LLM （基本上）是无状态的](#%E5%8E%9F%E5%88%99llm-%E5%9F%BA%E6%9C%AC%E4%B8%8A%E6%98%AF%E6%97%A0%E7%8A%B6%E6%80%81%E7%9A%84)
  - [CLAUDE.md 让 Claude 了解你的代码库](#claudemd-%E8%AE%A9-claude-%E4%BA%86%E8%A7%A3%E4%BD%A0%E7%9A%84%E4%BB%A3%E7%A0%81%E5%BA%93)
  - [Claude 经常忽略 CLAUDE.md](#claude-%E7%BB%8F%E5%B8%B8%E5%BF%BD%E7%95%A5-claudemd)
  - [创建优秀的 CLAUDE.md 文件](#%E5%88%9B%E5%BB%BA%E4%BC%98%E7%A7%80%E7%9A%84-claudemd-%E6%96%87%E4%BB%B6)
    - [少即是多（指令）](#%E5%B0%91%E5%8D%B3%E6%98%AF%E5%A4%9A%E6%8C%87%E4%BB%A4)
    - [CLAUDE.md 文件长度和适用性](#claudemd-%E6%96%87%E4%BB%B6%E9%95%BF%E5%BA%A6%E5%92%8C%E9%80%82%E7%94%A8%E6%80%A7)
    - [渐进式披露](#%E6%B8%90%E8%BF%9B%E5%BC%8F%E6%8A%AB%E9%9C%B2)
    - [Claude（不）是昂贵的 linter](#claude%E4%B8%8D%E6%98%AF%E6%98%82%E8%B4%B5%E7%9A%84-linter)
    - [不要使用 /init 或自动生成你的 CLAUDE.md](#%E4%B8%8D%E8%A6%81%E4%BD%BF%E7%94%A8-init-%E6%88%96%E8%87%AA%E5%8A%A8%E7%94%9F%E6%88%90%E4%BD%A0%E7%9A%84-claudemd)
  - [总结](#%E6%80%BB%E7%BB%93)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 如何编写优秀的 CLAUDE.md

注：本文也适用于 AGENTS.md，它是 CLAUDE.md 的开源等价物，适用于 OpenCode、Zed、Cursor 和 Codex 等 agent 和框架。

---

## 原则：LLM （基本上）是无状态的

LLM 是无状态函数。它们的权重在用于推理时已经冻结，所以不会随时间学习。模型对你代码库的唯一了解，就是你输入给它的 token。

类似地，像 Claude Code 这样的编程 agent 框架通常要求你显式管理 agent 的记忆。CLAUDE.md（或 AGENTS.md）是默认情况下会被包含进你与 agent 的每一次对话的唯一文件。

这有三个重要含义：

- 编程 agent 在每次会话开始时对你的代码库**一无所知**
- 必须在每次开始会话时告诉 agent 关于代码库的任何重要信息
- CLAUDE.md 是实现这一点的首选方式

---

## CLAUDE.md 让 Claude 了解你的代码库

由于 Claude 在每次会话开始时对你的代码库一无所知，你应该使用 CLAUDE.md 来让 Claude 熟悉你的代码库。从高层次看，这意味着它应该涵盖：

- **是什么（WHAT）**：告诉 Claude 技术栈、项目结构。给 Claude 一张代码库地图。这在 monorepo 中尤其重要！告诉 Claude 有哪些应用、有哪些共享包，以及每样东西的用途，这样它就知道去哪里找东西。

- **为什么（WHY）**：告诉 Claude 项目的目的，以及仓库中各部分在做什么。项目不同部分的目的和功能是什么？

- **怎么做（HOW）**：告诉 Claude 应该如何在项目上工作。例如，你用 bun 而不是 node？你需要包含它在项目上实际做有意义工作所需的所有信息。Claude 如何验证自己的改动？如何运行测试、类型检查和编译步骤？

但实现方式很重要！不要试图把 Claude 可能需要运行的每一条命令都塞进 CLAUDE.md 文件——你会得到次优结果。

---

## Claude 经常忽略 CLAUDE.md

无论你使用哪个模型，你可能会注意到 Claude 经常忽略 CLAUDE.md 文件的内容。

你可以自己调查这一点，通过使用 ANTHROPIC_BASE_URL 在 claude code CLI 和 Anthropic API 之间放置一个日志代理。Claude code 会在用户消息中注入以下系统提醒和你的 CLAUDE.md 文件：

```
<system-reminder>
  重要：此上下文可能与你的任务相关，也可能无关。
  除非它与你的任务高度相关，否则你不应该回应此上下文。
</system-reminder>
```

因此，如果 Claude 判定 CLAUDE.md 的内容与当前任务不相关，它就会忽略你的指令。文件中不普遍适用于你让它工作的任务的信息越多，Claude 就越有可能忽略文件中的指令。

Anthropic 为什么要加这个？很难确定，但我们可以推测一下。我们遇到的大多数 CLAUDE.md 文件都包含一堆不广泛适用的指令。许多用户把这个文件当作添加"热修复"的方式，通过追加大量不一定广泛适用的指令来修正他们不喜欢的行为。

我们只能假设 Claude Code 团队发现，通过告诉 Claude 忽略不好的指令，框架实际上产生了更好的结果。

---

## 创建优秀的 CLAUDE.md 文件

以下部分提供了一些关于如何遵循上下文工程最佳实践编写优秀 CLAUDE.md 文件的建议。

效果因人而异。并非所有这些规则都必然适用于每种设置。像其他任何事情一样，一旦你理解何时以及为什么可以打破规则，并且有充分理由这样做，就可以随意打破规则。

### 少即是多（指令）

把 Claude 可能需要运行的每一条命令，以及你的代码标准和风格指南都塞进 CLAUDE.md 可能很诱人。我们建议不要这样做。

虽然这个话题还没有以非常严格的方式进行研究，但已经有一些研究表明：

- 前沿的思考型 LLM 可以合理一致地遵循约 150-200 条指令。较小的模型能关注的指令比大型模型少，非思考型模型能关注的指令比思考型模型少。

- **较小的模型变差得更快、更严重**。具体来说，随着指令数量增加，较小模型的指令遵循性能往往呈现指数衰减，而大型前沿思考模型呈现线性衰减。因此，我们建议不要对多步骤任务或复杂实现计划使用较小模型。

- LLM 偏向于提示外围的指令：在最开始（Claude Code 系统消息和 CLAUDE.md）和最末尾（最近的用户消息）

- 随着指令数量增加，指令遵循质量**均匀下降**。这意味着当你给 LLM 更多指令时，它不是简单地忽略新的（"文件中更靠后的"）指令——它开始均匀地忽略所有指令

我们对 Claude Code 框架的分析表明，Claude Code 的系统提示包含约 50 条独立指令。根据你使用的模型，这已经是你的 agent 能可靠遵循的指令的近三分之一——而这还是在规则、插件、技能或用户消息之前。

这意味着你的 CLAUDE.md 文件应该包含尽可能少的指令——理想情况下只包含普遍适用于你任务的指令。

### CLAUDE.md 文件长度和适用性

在其他条件相同的情况下，当 LLM 的上下文窗口充满了集中、相关的上下文（包括示例、相关文件、工具调用和工具结果）时，它在任务上的表现会比上下文窗口有大量无关上下文时更好。

由于 CLAUDE.md 会进入每一次会话，你应该确保其内容尽可能普遍适用。

例如，避免包含关于（比如）如何构建新数据库 schema 的指令——当你在做其他不相关的事情时，这不会有用，还会分散模型的注意力！

长度方面，"少即是多"原则同样适用。虽然 Anthropic 没有关于 CLAUDE.md 文件应该多长的官方建议，但普遍共识是 **< 300 行最好，更短更好**。

在 HumanLayer，我们根目录的 CLAUDE.md 文件不到六十行。

### 渐进式披露

编写一个简洁的、涵盖你想让 Claude 知道的所有内容的 CLAUDE.md 文件可能很有挑战性，特别是在大型项目中。

为了解决这个问题，我们可以利用**渐进式披露**原则，确保 Claude 只在需要时才看到特定于任务或项目的指令。

与其在 CLAUDE.md 文件中包含所有关于构建项目、运行测试、代码规范或其他重要上下文的不同指令，我们建议将特定于任务的指令保存在项目中某处的单独 markdown 文件中，并使用自描述的名称。

例如：

```
agent_docs/
  |- building_the_project.md
  |- running_tests.md
  |- code_conventions.md
  |- service_architecture.md
  |- database_schema.md
  |- service_communication_patterns.md
```

然后，在 CLAUDE.md 文件中，你可以包含这些文件的列表和每个文件的简要描述，并指示 Claude 决定哪些（如果有的话）相关，并在开始工作前阅读它们。或者，要求 Claude 先向你展示它想要阅读的文件以获得批准，然后再阅读它们。

**优先使用指针而不是副本**。如果可能，不要在这些文件中包含代码片段——它们会很快过时。相反，包含文件:行引用来将 Claude 指向权威上下文。

从概念上讲，这与 Claude Skills 的预期工作方式非常相似，尽管 skills 更专注于工具使用而非指令。

### Claude（不）是昂贵的 linter

我们看到人们在 CLAUDE.md 文件中放置的最常见内容之一是代码风格指南。**永远不要让 LLM 去做 linter 的工作**。与传统的 linter 和 formatter 相比，LLM 相对昂贵且速度极慢。我们认为你应该始终尽可能使用确定性工具。

代码风格指南不可避免地会在你的上下文窗口中添加一堆指令和大部分不相关的代码片段，降低 LLM 的性能和指令遵循能力，并占用你的上下文窗口。

LLM 是上下文学习者！如果你的代码遵循某一套风格指南或模式，你应该会发现，凭借对代码库的几次搜索（或一份好的研究文档！），你的 agent 应该倾向于遵循现有的代码模式和规范，而无需被告知。

如果你对此非常坚持，你甚至可以考虑设置一个 Claude Code Stop hook，运行你的 formatter 和 linter，并将错误呈现给 Claude 让它修复。不要让 Claude 自己去找格式问题。

额外加分：使用可以自动修复问题的 linter（我们喜欢 Biome），并仔细调整关于什么可以安全自动修复的规则，以实现最大（安全）覆盖率。

你也可以创建一个 Slash Command，包含你的代码指南，并将 Claude 指向版本控制中的更改，或你的 git status 等。这样，你可以分别处理实现和格式化。结果是两者都会有更好的表现。

### 不要使用 /init 或自动生成你的 CLAUDE.md

Claude Code 和其他带有 OpenCode 的框架都提供了自动生成 CLAUDE.md 文件（或 AGENTS.md）的方法。

因为 CLAUDE.md 会进入与 Claude code 的每一次会话，它是框架的最高杠杆点之一——无论好坏，取决于你如何使用它。

一行糟糕的代码就是一行糟糕的代码。实现计划中的一行糟糕内容有可能创建大量糟糕的代码行。误解系统工作方式的研究中的一行糟糕内容有可能导致计划中的大量糟糕行，因此也会产生更多糟糕的代码行。

但 CLAUDE.md 文件会影响你工作流程的每一个阶段以及它产生的每一个产物。因此，我们认为你应该花一些时间非常仔细地思考进入其中的每一行。

---

## 总结

- **CLAUDE.md 用于让 Claude 熟悉你的代码库**。它应该定义项目的为什么（WHY）、是什么（WHAT）和怎么做（HOW）。

- **少（指令）即是多**。虽然你不应该省略必要的指令，但应该在文件中包含尽可能少的合理指令。

- **保持 CLAUDE.md 的内容简洁且普遍适用**。

- **使用渐进式披露** - 不要告诉 Claude 你可能希望它知道的所有信息。相反，告诉它如何找到重要信息，这样它可以找到并使用它，但仅在需要时使用，以避免膨胀上下文窗口或指令数量。

- **Claude 不是 linter**。使用 linter 和代码格式化工具，并根据需要使用 Hooks 和 Slash Commands 等其他功能。

- **CLAUDE.md 是框架的最高杠杆点**，所以避免自动生成它。你应该精心制作其内容以获得最佳结果。

## 参考

【writing-a-good-claude-md】https://www.humanlayer.dev/blog/writing-a-good-claude-md   