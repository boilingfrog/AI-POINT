package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAgentKeepsConversationHistory(t *testing.T) {
	var requests []chatRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("请求路径 = %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization = %q", got)
		}

		var request chatRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("解析请求失败: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		requests = append(requests, request)

		answer := "第一轮回复"
		if len(requests) == 2 {
			answer = "第二轮回复"
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": Message{Role: "assistant", Content: answer}},
			},
		})
	}))
	defer server.Close()

	agent := NewAgent("test-key", "test-model", server.URL+"/v1")
	agent.client = server.Client()

	if _, err := agent.Chat(context.Background(), "第一轮问题"); err != nil {
		t.Fatal(err)
	}
	if _, err := agent.Chat(context.Background(), "第二轮问题"); err != nil {
		t.Fatal(err)
	}

	if len(requests) != 2 {
		t.Fatalf("请求次数 = %d，期望 2", len(requests))
	}
	second := requests[1]
	if second.Model != "test-model" {
		t.Errorf("model = %q", second.Model)
	}
	if len(second.Messages) != 4 {
		t.Fatalf("第二次请求消息数 = %d，期望 4", len(second.Messages))
	}
	if second.Messages[1].Content != "第一轮问题" ||
		second.Messages[2].Content != "第一轮回复" ||
		second.Messages[3].Content != "第二轮问题" {
		t.Errorf("第二次请求没有包含完整历史: %#v", second.Messages)
	}

	agent.Clear()
	if len(agent.messages) != 1 || agent.messages[0].Role != "system" {
		t.Errorf("Clear 后的消息 = %#v", agent.messages)
	}
}

func TestAgentRollsBackFailedMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"bad key"}}`))
	}))
	defer server.Close()

	agent := NewAgent("bad-key", "test-model", server.URL)
	agent.client = server.Client()

	if _, err := agent.Chat(context.Background(), "不会被保存"); err == nil {
		t.Fatal("期望接口错误，实际没有错误")
	}
	if len(agent.messages) != 1 {
		t.Errorf("失败消息没有回滚: %#v", agent.messages)
	}
}
