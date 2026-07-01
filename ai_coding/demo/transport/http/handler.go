// Package httpapi 是 HTTP 传输层：负责路由、请求解析、状态码映射，不写业务逻辑。
//
// 包名用 httpapi 而非 http，避免遮蔽标准库 net/http（目录仍叫 http/ 以对齐分层）。
package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"demo/coinwallet/service"
)

// NewHandler 装配路由并返回 http.Handler。使用 Go 1.22+ 的方法 + 通配路由。
func NewHandler(svc *service.Service) http.Handler {
	h := &handler{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.health)
	mux.HandleFunc("GET /wallets/{id}", h.balance)
	mux.HandleFunc("POST /wallets/{id}/grant", h.grant)
	mux.HandleFunc("POST /wallets/{id}/spend", h.spend)
	return mux
}

type handler struct {
	svc *service.Service
}

// amountRequest 是 grant/spend 的请求体。金额用 int64，与账本一致。
type amountRequest struct {
	Amount int64 `json:"amount"`
}

func (h *handler) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (h *handler) balance(w http.ResponseWriter, r *http.Request) {
	wallet, err := h.svc.Balance(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, wallet)
}

func (h *handler) grant(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeAmount(w, r)
	if !ok {
		return
	}
	wallet, err := h.svc.Grant(r.PathValue("id"), req.Amount)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, wallet)
}

func (h *handler) spend(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeAmount(w, r)
	if !ok {
		return
	}
	wallet, err := h.svc.Spend(r.PathValue("id"), req.Amount)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, wallet)
}

// decodeAmount 解析请求体；解析失败时直接写 400 并返回 false。
func decodeAmount(w http.ResponseWriter, r *http.Request) (amountRequest, bool) {
	var req amountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "invalid json body"})
		return amountRequest{}, false
	}
	return req, true
}

type errorBody struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError 把 service 层的 sentinel error 映射成 HTTP 状态码。
func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, service.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, service.ErrInvalidAmount):
		status = http.StatusBadRequest
	case errors.Is(err, service.ErrInsufficientBalance):
		status = http.StatusConflict
	}
	writeJSON(w, status, errorBody{Error: err.Error()})
}
