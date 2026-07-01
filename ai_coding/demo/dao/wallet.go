// Package dao 是数据访问层，负责钱包的存储。本 demo 用内存 map 实现，
// 是整个服务里唯一持有存储状态的地方。
package dao

import (
	"sync"

	"demo/coinwallet/model"
)

// Store 是基于内存 map 的钱包存储，并发安全。
//
// 共享 map 在多 goroutine（HTTP 并发请求）下读写，必须用锁保护——
// Go 的 map 并发读写会直接 fatal，不是可恢复的 panic。
type Store struct {
	mu      sync.RWMutex
	wallets map[string]model.Wallet
}

// NewStore 返回一个初始化好的空 Store。
func NewStore() *Store {
	return &Store{wallets: make(map[string]model.Wallet)}
}

// Get 返回指定用户的钱包；第二个返回值表示是否存在。
func (s *Store) Get(userID string) (model.Wallet, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w, ok := s.wallets[userID]
	return w, ok
}

// Save 写入（或覆盖）一个钱包。
func (s *Store) Save(w model.Wallet) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.wallets[w.UserID] = w
}
