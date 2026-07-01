// Package service 承载金币钱包的业务逻辑与校验。所有金额校验都在这一层，
// 上层 transport 只做协议转换，不写业务规则。
package service

import (
	"errors"

	"demo/coinwallet/model"
)

// 业务错误用 sentinel error 表达，上层用 errors.Is 判断并映射到 HTTP 状态码。
var (
	ErrNotFound            = errors.New("wallet not found")
	ErrInvalidAmount       = errors.New("amount must be positive")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

// Repository 是 service 依赖的存储抽象。service 只依赖这个接口，
// 不直接依赖 dao 的具体实现——便于替换和测试。dao.Store 实现它。
type Repository interface {
	Get(userID string) (model.Wallet, bool)
	Save(w model.Wallet)
}

// Service 是钱包业务逻辑的入口。
type Service struct {
	repo Repository
}

// New 用给定的存储构造 Service。
func New(repo Repository) *Service {
	return &Service{repo: repo}
}

// Balance 返回用户当前钱包；用户不存在返回 ErrNotFound。
func (s *Service) Balance(userID string) (model.Wallet, error) {
	w, ok := s.repo.Get(userID)
	if !ok {
		return model.Wallet{}, ErrNotFound
	}
	return w, nil
}

// Grant 给用户发放金币。amount 必须为正；用户不存在则新建账户。
func (s *Service) Grant(userID string, amount int64) (model.Wallet, error) {
	if amount <= 0 {
		return model.Wallet{}, ErrInvalidAmount
	}
	w, ok := s.repo.Get(userID)
	if !ok {
		w = model.Wallet{UserID: userID}
	}
	w.Balance += amount
	s.repo.Save(w)
	return w, nil
}

// Spend 扣减用户金币。amount 必须为正；用户不存在返回 ErrNotFound；
// 余额不足返回 ErrInsufficientBalance。
func (s *Service) Spend(userID string, amount int64) (model.Wallet, error) {
	if amount <= 0 {
		return model.Wallet{}, ErrInvalidAmount
	}
	w, ok := s.repo.Get(userID)
	if !ok {
		return model.Wallet{}, ErrNotFound
	}
	if w.Balance < amount {
		return model.Wallet{}, ErrInsufficientBalance
	}
	w.Balance -= amount
	s.repo.Save(w)
	return w, nil
}
