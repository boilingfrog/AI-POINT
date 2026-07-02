package service

import (
	"errors"
	"testing"

	"demo/coinwallet/dao"
)

// newService 用真实的内存 dao 构造 service，供各用例复用。
func newService() *Service {
	return New(dao.NewStore())
}

func TestGrant(t *testing.T) {
	tests := []struct {
		name       string
		grants     []int64 // 依次发放
		wantErr    error   // 最后一次发放期望的错误
		wantAmount int64   // 成功时期望余额
	}{
		{name: "新用户发放", grants: []int64{100}, wantAmount: 100},
		{name: "累加发放", grants: []int64{100, 50}, wantAmount: 150},
		{name: "金额为零非法", grants: []int64{0}, wantErr: ErrInvalidAmount},
		{name: "负数非法", grants: []int64{-5}, wantErr: ErrInvalidAmount},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService()
			var err error
			var balance int64
			for _, amt := range tt.grants {
				var w, e = s.Grant("u1", amt)
				err, balance = e, w.Balance
			}
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("期望错误 %v，得到 %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("未预期错误: %v", err)
			}
			if balance != tt.wantAmount {
				t.Fatalf("期望余额 %d，得到 %d", tt.wantAmount, balance)
			}
		})
	}
}

func TestSpend(t *testing.T) {
	tests := []struct {
		name       string
		grant      int64 // >0 时先发放
		spend      int64
		wantErr    error
		wantAmount int64
	}{
		{name: "正常扣减", grant: 100, spend: 30, wantAmount: 70},
		{name: "余额不足", grant: 100, spend: 999, wantErr: ErrInsufficientBalance},
		{name: "扣减非法金额", grant: 100, spend: 0, wantErr: ErrInvalidAmount},
		{name: "用户不存在", grant: 0, spend: 10, wantErr: ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService()
			if tt.grant > 0 {
				if _, err := s.Grant("u1", tt.grant); err != nil {
					t.Fatalf("准备阶段发放失败: %v", err)
				}
			}
			w, err := s.Spend("u1", tt.spend)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("期望错误 %v，得到 %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("未预期错误: %v", err)
			}
			if w.Balance != tt.wantAmount {
				t.Fatalf("期望余额 %d，得到 %d", tt.wantAmount, w.Balance)
			}
		})
	}
}

func TestBalanceNotFound(t *testing.T) {
	s := newService()
	if _, err := s.Balance("ghost"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("期望 ErrNotFound，得到 %v", err)
	}
}
