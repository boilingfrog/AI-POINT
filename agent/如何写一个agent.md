# 如何写一个 Agent

> 本文从概念到实践，系统地介绍如何设计并实现一个 AI Agent。

---

## 一、什么是 Agent？

Agent（智能体）是一个能够**感知环境、做出决策、并执行行动**以达成目标的程序。

与普通的 LLM 调用不同，Agent 具备以下特征：

| 特征 | 说明 |
|------|------|
| **自主性** | 无需每步人工干预，自主规划执行路径 |
| **工具使用** | 可调用外部工具（搜索、代码执行、API 等）|
| **多步推理** | 能拆解复杂任务，分步完成 |
| **状态保持** | 在对话或任务过程中维护上下文和记忆 |

```
用户目标
   │
   ▼
┌─────────────────────────────────────┐
│              Agent Loop              │
│                                     │
│  感知(Perceive) → 思考(Think)        │
│       ↑               │             │
│       │          行动(Act)          │
│       └──────── 观察(Observe) ──────┘
└─────────────────────────────────────┘
```

---

## 二、Agent 的核心架构

一个完整的 Agent 由以下模块构成：

### 2.1 大脑（LLM）

LLM 是 Agent 的核心推理引擎，负责：
- 理解用户意图
- 规划任务步骤
- 选择合适的工具
- 生成最终回答

**选型建议：**
- 复杂推理任务 → `claude-opus-4-8` / `claude-sonnet-5`
- 高频低延迟任务 → `claude-haiku-4-5`
- 代码生成 → 优先选推理能力强的模型

### 2.2 工具（Tools）

工具是 Agent 与外界交互的接口。常见工具类型：

```python
# 工具定义示例（Python）
tools = [
    {
        "name": "web_search",
        "description": "搜索互联网上的最新信息",
        "input_schema": {
            "type": "object",
            "properties": {
                "query": {
                    "type": "string",
                    "description": "搜索关键词"
                }
            },
            "required": ["query"]
        }
    },
    {
        "name": "code_executor",
        "description": "执行 Python 代码并返回结果",
        "input_schema": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string",
                    "description": "要执行的 Python 代码"
                }
            },
            "required": ["code"]
        }
    }
]
```

### 2.3 记忆（Memory）

| 类型 | 说明 | 实现方式 |
|------|------|----------|
| **短期记忆** | 当前对话上下文 | messages 列表 |
| **长期记忆** | 跨会话持久化信息 | 文件 / 数据库 |
| **语义记忆** | 知识检索 | 向量数据库（RAG）|
| **情节记忆** | 历史行为记录 | 结构化日志 |

### 2.4 规划（Planning）

Agent 的规划能力决定了它处理复杂任务的上限：

- **ReAct**：Reason + Act，边推理边行动，最常用
- **Plan-and-Execute**：先制定完整计划，再逐步执行
- **Reflection**：执行后反思，迭代改进结果

---

## 三、Agent Loop 的实现

### 3.1 最简实现（ReAct 模式）

```python
import anthropic

client = anthropic.Anthropic()

def run_agent(user_message: str, tools: list, tool_executor: callable):
    """
    最简 Agent 循环
    
    Args:
        user_message: 用户输入
        tools: 工具定义列表
        tool_executor: 工具执行函数
    """
    messages = [{"role": "user", "content": user_message}]
    
    while True:
        # 1. 调用 LLM
        response = client.messages.create(
            model="claude-opus-4-8",
            max_tokens=4096,
            tools=tools,
            messages=messages
        )
        
        # 2. 将 assistant 回复加入历史
        messages.append({"role": "assistant", "content": response.content})
        
        # 3. 检查是否需要调用工具
        if response.stop_reason == "tool_use":
            tool_results = []
            
            for block in response.content:
                if block.type == "tool_use":
                    # 4. 执行工具
                    result = tool_executor(block.name, block.input)
                    tool_results.append({
                        "type": "tool_result",
                        "tool_use_id": block.id,
                        "content": str(result)
                    })
            
            # 5. 将工具结果返回给模型
            messages.append({"role": "user", "content": tool_results})
        
        else:
            # 6. 模型完成任务，返回最终结果
            final_text = next(
                (b.text for b in response.content if hasattr(b, "text")),
                ""
            )
            return final_text
```

### 3.2 工具执行器

```python
import subprocess
import json

def tool_executor(tool_name: str, tool_input: dict) -> str:
    """统一的工具分发执行器"""
    
    if tool_name == "web_search":
        return web_search(tool_input["query"])
    
    elif tool_name == "code_executor":
        return execute_code(tool_input["code"])
    
    elif tool_name == "read_file":
        return read_file(tool_input["path"])
    
    else:
        return f"错误：未知工具 {tool_name}"


def execute_code(code: str) -> str:
    """在沙箱中执行代码"""
    try:
        result = subprocess.run(
            ["python3", "-c", code],
            capture_output=True,
            text=True,
            timeout=10  # 超时保护
        )
        return result.stdout or result.stderr
    except subprocess.TimeoutExpired:
        return "错误：代码执行超时（>10s）"
    except Exception as e:
        return f"错误：{str(e)}"
```

---

## 四、System Prompt 的设计

System Prompt 是 Agent 行为的"宪法"，写好它至关重要。

### 核心原则

```markdown
# 角色定义
你是一个专业的数据分析 Agent，帮助用户分析数据、生成报告。

# 能力边界
你可以：
- 执行 Python 数据分析代码
- 搜索最新的统计数据
- 生成图表和可视化

你不可以：
- 访问用户未明确授权的系统
- 删除或修改用户数据（只读）

# 工作原则
1. 遇到不确定的情况，先澄清需求再执行
2. 每次工具调用前，说明你的意图
3. 任务完成后，总结执行结果

# 输出格式
分析结果请按以下结构返回：
- 执行摘要（1-2句）
- 关键发现（bullet points）
- 详细分析
- 建议与后续步骤
```

### System Prompt 的黄金法则

| 原则 | 说明 |
|------|------|
| **具体 > 抽象** | "分析前先确认数据格式"比"仔细分析"有效 |
| **边界清晰** | 明确能做什么、不能做什么 |
| **格式约定** | 规定输出结构，确保结果可解析 |
| **错误处理** | 告诉 Agent 遇到错误时如何响应 |
| **避免过长** | 超过 2000 token 的 System Prompt 效果递减 |

---

## 五、进阶设计模式

### 5.1 Plan-and-Execute（适合长任务）

```python
def plan_and_execute(goal: str):
    # 阶段一：制定计划
    plan = llm_call(
        system="你是一个任务规划专家，将复杂目标拆解为具体步骤",
        user=f"目标：{goal}\n请制定详细的执行计划，以 JSON 列表格式输出"
    )
    
    steps = json.loads(plan)
    results = []
    
    # 阶段二：逐步执行
    for i, step in enumerate(steps):
        print(f"执行步骤 {i+1}/{len(steps)}: {step['description']}")
        
        result = run_agent(
            user_message=step["instruction"],
            context={"previous_results": results}
        )
        results.append(result)
    
    # 阶段三：汇总结果
    return summarize(goal, results)
```

### 5.2 Reflection（提升输出质量）

```python
def agent_with_reflection(task: str, max_iterations: int = 3):
    result = run_agent(task)
    
    for i in range(max_iterations):
        # 让模型评估自己的输出
        critique = llm_call(
            system="你是一个严格的评审员",
            user=f"""
            任务：{task}
            
            当前结果：{result}
            
            请评估：
            1. 结果是否完整回答了任务要求？
            2. 有哪些不足或错误？
            3. 如果满意，回复 "DONE"；否则给出改进建议
            """
        )
        
        if "DONE" in critique:
            break
            
        # 基于批评改进
        result = run_agent(f"{task}\n\n之前的尝试有以下问题：{critique}\n请改进")
    
    return result
```

### 5.3 Multi-Agent（多智能体协作）

```
                  ┌─────────────┐
                  │  Orchestrator│  ← 分配任务、汇总结果
                  └──────┬──────┘
          ┌───────────────┼───────────────┐
          ▼               ▼               ▼
   ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
   │  Research   │ │   Coder     │ │   Writer    │
   │   Agent     │ │   Agent     │ │   Agent     │
   └─────────────┘ └─────────────┘ └─────────────┘
   搜索、收集信息     编写、执行代码     生成最终文档
```

适合场景：任务可以并行拆分，且各子任务相对独立。

---

## 六、错误处理与鲁棒性

### 6.1 必须处理的错误类型

```python
import time
from typing import Optional

def robust_tool_call(
    tool_fn: callable,
    tool_input: dict,
    max_retries: int = 3,
    timeout: int = 30
) -> Optional[str]:
    """带重试和超时的工具调用"""
    
    for attempt in range(max_retries):
        try:
            result = tool_fn(**tool_input)
            return result
            
        except TimeoutError:
            print(f"工具超时，第 {attempt+1} 次重试...")
            time.sleep(2 ** attempt)  # 指数退避
            
        except RateLimitError:
            wait_time = 60  # API 限速，等待 1 分钟
            print(f"触发限速，等待 {wait_time}s...")
            time.sleep(wait_time)
            
        except Exception as e:
            print(f"工具执行失败: {e}")
            if attempt == max_retries - 1:
                return f"工具调用失败（已重试 {max_retries} 次）：{str(e)}"
    
    return None
```

### 6.2 防止无限循环

```python
def run_agent_safe(user_message: str, max_steps: int = 20):
    """加入步数限制，防止 Agent 陷入死循环"""
    
    messages = [{"role": "user", "content": user_message}]
    step_count = 0
    
    while step_count < max_steps:
        step_count += 1
        response = call_llm(messages)
        
        if response.stop_reason != "tool_use":
            return response  # 正常结束
            
        # 执行工具...
        messages = update_messages(messages, response)
    
    # 超过最大步数，强制结束
    return "任务超出最大步数限制，可能陷入循环。请检查任务描述或增加 max_steps。"
```

### 6.3 人工介入（Human-in-the-Loop）

```python
SENSITIVE_TOOLS = {"delete_file", "send_email", "make_payment"}

def tool_executor_with_confirmation(tool_name: str, tool_input: dict) -> str:
    # 高风险操作需要人工确认
    if tool_name in SENSITIVE_TOOLS:
        print(f"\n⚠️  Agent 请求执行高风险操作：")
        print(f"   工具: {tool_name}")
        print(f"   参数: {json.dumps(tool_input, ensure_ascii=False, indent=2)}")
        
        confirm = input("是否允许？(y/N): ").strip().lower()
        if confirm != "y":
            return "用户取消了此操作"
    
    return execute_tool(tool_name, tool_input)
```

---

## 七、最佳实践清单

### 设计阶段

- [ ] **明确任务边界**：Agent 解决什么问题？不解决什么问题？
- [ ] **最小工具集**：只给 Agent 完成任务必需的工具，避免权限过大
- [ ] **工具描述要准确**：LLM 依靠描述选择工具，模糊描述导致误调用
- [ ] **设计失败路径**：每个工具都要有错误返回格式

### 开发阶段

- [ ] **先测单步**：在 Agent Loop 之前，验证每个工具单独工作正常
- [ ] **记录所有调用**：保存完整的 messages 历史，便于调试
- [ ] **设置超时和重试**：所有外部调用加保护
- [ ] **限制最大步数**：防止无限循环

### 部署阶段

- [ ] **费用监控**：设置 token 用量告警
- [ ] **输入过滤**：防止 Prompt Injection 攻击
- [ ] **输出审核**：对生成内容做安全检查（尤其是代码执行）
- [ ] **审计日志**：记录 Agent 的每次工具调用，便于事后审查

---

## 八、常见问题与解决方案

**Q: Agent 经常选错工具怎么办？**

A: 优化工具的 `description` 字段，使其描述更精确，并给出使用示例。如果工具数量超过 10 个，考虑用"工具路由"先分类再调用。

---

**Q: Agent 在简单任务上也调用太多步骤？**

A: 在 System Prompt 中加入"能一步完成的任务不要拆步骤"的指令。同时检查工具粒度是否太细。

---

**Q: 如何评估 Agent 的质量？**

A: 建立评测集（evals）：覆盖典型任务、边界情况、错误恢复场景。关注以下指标：
- 任务完成率
- 平均步数（越少越好，说明规划效率高）
- 工具调用准确率
- 成本（token 用量）

---

**Q: Context 窗口不够用怎么办？**

A: 
1. 对历史 tool_result 做摘要压缩
2. 只保留最近 N 轮对话
3. 用长期记忆（数据库/文件）替代上下文传递历史

---

## 九、完整示例：一个文件分析 Agent

```python
#!/usr/bin/env python3
"""
文件分析 Agent：能读取文件并回答相关问题
"""

import os
import anthropic

client = anthropic.Anthropic()

# --- 工具定义 ---
TOOLS = [
    {
        "name": "read_file",
        "description": "读取本地文件内容。适用于分析代码、文档、数据文件。",
        "input_schema": {
            "type": "object",
            "properties": {
                "path": {"type": "string", "description": "文件的绝对或相对路径"},
                "max_lines": {"type": "integer", "description": "最多读取行数，默认100"}
            },
            "required": ["path"]
        }
    },
    {
        "name": "list_directory",
        "description": "列出目录下的文件和子目录",
        "input_schema": {
            "type": "object",
            "properties": {
                "path": {"type": "string", "description": "目录路径"}
            },
            "required": ["path"]
        }
    }
]

# --- 工具实现 ---
def read_file(path: str, max_lines: int = 100) -> str:
    try:
        with open(path, "r", encoding="utf-8") as f:
            lines = f.readlines()
        if len(lines) > max_lines:
            content = "".join(lines[:max_lines])
            return f"{content}\n... （已截断，共 {len(lines)} 行，只显示前 {max_lines} 行）"
        return "".join(lines)
    except FileNotFoundError:
        return f"错误：文件 {path} 不存在"
    except Exception as e:
        return f"错误：{str(e)}"

def list_directory(path: str) -> str:
    try:
        items = os.listdir(path)
        result = []
        for item in sorted(items):
            full_path = os.path.join(path, item)
            item_type = "📁" if os.path.isdir(full_path) else "📄"
            result.append(f"{item_type} {item}")
        return "\n".join(result) if result else "（空目录）"
    except Exception as e:
        return f"错误：{str(e)}"

def execute_tool(name: str, input_data: dict) -> str:
    if name == "read_file":
        return read_file(**input_data)
    elif name == "list_directory":
        return list_directory(**input_data)
    return f"未知工具：{name}"

# --- Agent 主循环 ---
def run_file_agent(question: str):
    print(f"\n🤖 任务：{question}\n")
    
    messages = [{"role": "user", "content": question}]
    step = 0
    
    while step < 15:
        step += 1
        
        response = client.messages.create(
            model="claude-opus-4-8",
            max_tokens=4096,
            system="""你是一个文件分析专家。
你可以读取文件和列出目录来回答用户问题。
分析时先了解目录结构，再读取相关文件。
回答要具体，引用文件中的实际内容。""",
            tools=TOOLS,
            messages=messages
        )
        
        messages.append({"role": "assistant", "content": response.content})
        
        if response.stop_reason != "tool_use":
            # 提取文字回答
            answer = next(
                (b.text for b in response.content if hasattr(b, "text")), ""
            )
            print(f"✅ 结果：\n{answer}")
            return answer
        
        # 执行工具调用
        tool_results = []
        for block in response.content:
            if block.type == "tool_use":
                print(f"  🔧 调用工具: {block.name}({block.input})")
                result = execute_tool(block.name, block.input)
                tool_results.append({
                    "type": "tool_result",
                    "tool_use_id": block.id,
                    "content": result
                })
        
        messages.append({"role": "user", "content": tool_results})
    
    return "超出最大步数限制"


if __name__ == "__main__":
    run_file_agent("分析当前目录的结构，并总结这个项目是做什么的")
```

---

## 十、延伸阅读

- [Anthropic Agent 官方文档](https://docs.anthropic.com/en/docs/agents)
- [Building Effective Agents（Anthropic 博客）](https://www.anthropic.com/research/building-effective-agents)
- [ReAct 论文](https://arxiv.org/abs/2210.03629)：Reason + Act 原始论文
- [LangGraph 文档](https://langchain-ai.github.io/langgraph/)：复杂 Agent 编排框架

---