package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const systemPrompt = "你是一个简洁、友好的中文助手。请直接回答用户的问题。"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Agent 保存对话历史，并直接调用 OpenAI 兼容接口。
type Agent struct {
	client   *http.Client
	endpoint string
	apiKey   string
	model    string
	messages []Message
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type chatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func NewAgent(apiKey, model, baseURL string) *Agent {
	return &Agent{
		client: &http.Client{
			Timeout: 2 * time.Minute,
		},
		endpoint: strings.TrimRight(baseURL, "/") + "/chat/completions",
		apiKey:   apiKey,
		model:    model,
		messages: []Message{{Role: "system", Content: systemPrompt}},
	}
}

func (a *Agent) Chat(ctx context.Context, input string) (string, error) {
	// 每轮都保存 user/assistant 消息，这就是最简单的短期记忆。
	a.messages = append(a.messages, Message{Role: "user", Content: input})

	body, err := json.Marshal(chatRequest{
		Model:    a.model,
		Messages: a.messages,
	})
	if err != nil {
		a.rollbackUserMessage()
		return "", fmt.Errorf("编码请求失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.endpoint, bytes.NewReader(body))
	if err != nil {
		a.rollbackUserMessage()
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		a.rollbackUserMessage()
		return "", fmt.Errorf("调用模型失败: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		a.rollbackUserMessage()
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result chatResponse
	if err := json.Unmarshal(responseBody, &result); err != nil {
		a.rollbackUserMessage()
		return "", fmt.Errorf("解析响应失败（HTTP %d）: %w", resp.StatusCode, err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		a.rollbackUserMessage()
		if result.Error != nil && result.Error.Message != "" {
			return "", fmt.Errorf("模型接口返回 HTTP %d: %s", resp.StatusCode, result.Error.Message)
		}
		return "", fmt.Errorf("模型接口返回 HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(responseBody)))
	}
	if len(result.Choices) == 0 || strings.TrimSpace(result.Choices[0].Message.Content) == "" {
		a.rollbackUserMessage()
		return "", errors.New("模型没有返回有效内容")
	}

	answer := result.Choices[0].Message.Content
	a.messages = append(a.messages, Message{Role: "assistant", Content: answer})
	return answer, nil
}

func (a *Agent) Clear() {
	a.messages = a.messages[:1]
}

func (a *Agent) rollbackUserMessage() {
	a.messages = a.messages[:len(a.messages)-1]
}

func main() {
	if err := loadEnv(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	model := os.Getenv("OPENAI_MODEL")
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if apiKey == "" || model == "" {
		fmt.Fprintln(os.Stderr, "请先设置 OPENAI_API_KEY 和 OPENAI_MODEL")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	agent := NewAgent(apiKey, model, baseURL)
	if len(os.Args) > 1 {
		answer, err := agent.Chat(ctx, strings.Join(os.Args[1:], " "))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Agent:", err)
			os.Exit(1)
		}
		fmt.Println("Agent:", answer)
		return
	}

	runConversation(ctx, agent, os.Stdin, os.Stdout)
}

func loadEnv() error {
	// 兼容在 agent 目录运行，以及 IDE 以仓库根目录作为 Working directory。
	for _, filename := range []string{".env", "agent/.env"} {
		if err := godotenv.Load(filename); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("加载 %s 失败: %w", filename, err)
		}
	}
	return nil
}

func runConversation(ctx context.Context, agent *Agent, input io.Reader, output io.Writer) {
	fmt.Fprintln(output, "开始对话。输入 /clear 清空上下文，输入 /exit 退出。")
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 1024), 1024*1024)

	for {
		fmt.Fprint(output, "You: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(output, "读取输入失败:", err)
			}
			return
		}

		question := strings.TrimSpace(scanner.Text())
		switch question {
		case "":
			continue
		case "/exit", "/quit":
			return
		case "/clear":
			agent.Clear()
			fmt.Fprintln(output, "Agent: 对话上下文已清空。")
			continue
		}

		answer, err := agent.Chat(ctx, question)
		if err != nil {
			fmt.Fprintln(output, "Agent:", err)
			if ctx.Err() != nil {
				return
			}
			continue
		}
		fmt.Fprintln(output, "Agent:", answer)
	}
}
