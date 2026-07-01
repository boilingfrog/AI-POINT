package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"demo/coinwallet/dao"
	"demo/coinwallet/model"
	"demo/coinwallet/service"
)

func newTestServer() http.Handler {
	return NewHandler(service.New(dao.NewStore()))
}

// do 发一个请求并返回 ResponseRecorder。
func do(h http.Handler, method, target, body string) *httptest.ResponseRecorder {
	var r *http.Request
	if body == "" {
		r = httptest.NewRequest(method, target, nil)
	} else {
		r = httptest.NewRequest(method, target, strings.NewReader(body))
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	return rec
}

func TestWalletFlow(t *testing.T) {
	h := newTestServer()

	// 未发放前查询 → 404
	if rec := do(h, http.MethodGet, "/wallets/u1", ""); rec.Code != http.StatusNotFound {
		t.Fatalf("查询不存在钱包期望 404，得到 %d", rec.Code)
	}

	// 发放 100 → 200，余额 100
	rec := do(h, http.MethodPost, "/wallets/u1/grant", `{"amount":100}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("发放期望 200，得到 %d", rec.Code)
	}
	if got := decodeWallet(t, rec).Balance; got != 100 {
		t.Fatalf("发放后期望余额 100，得到 %d", got)
	}

	// 扣减 30 → 200，余额 70
	rec = do(h, http.MethodPost, "/wallets/u1/spend", `{"amount":30}`)
	if got := decodeWallet(t, rec).Balance; rec.Code != http.StatusOK || got != 70 {
		t.Fatalf("扣减后期望 200/70，得到 %d/%d", rec.Code, got)
	}
}

func TestErrorMapping(t *testing.T) {
	h := newTestServer()
	// 准备：给 u1 发放 100，供后续用例扣减。
	if rec := do(h, http.MethodPost, "/wallets/u1/grant", `{"amount":100}`); rec.Code != http.StatusOK {
		t.Fatalf("准备阶段发放失败，状态 %d", rec.Code)
	}

	cases := []struct {
		name, method, target, body string
		want                       int
	}{
		{"余额不足 409", http.MethodPost, "/wallets/u1/spend", `{"amount":9999}`, http.StatusConflict},
		{"不存在 404", http.MethodGet, "/wallets/ghost", "", http.StatusNotFound},
		{"坏 JSON 400", http.MethodPost, "/wallets/u1/grant", `{bad`, http.StatusBadRequest},
		{"非法金额 400", http.MethodPost, "/wallets/u1/grant", `{"amount":0}`, http.StatusBadRequest},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if rec := do(h, c.method, c.target, c.body); rec.Code != c.want {
				t.Fatalf("期望 %d，得到 %d", c.want, rec.Code)
			}
		})
	}
}

func decodeWallet(t *testing.T, rec *httptest.ResponseRecorder) model.Wallet {
	t.Helper()
	var w model.Wallet
	if err := json.NewDecoder(rec.Body).Decode(&w); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	return w
}
