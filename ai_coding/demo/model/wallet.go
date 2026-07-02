// Package model 定义金币钱包的核心数据结构。纯数据，不含逻辑、不依赖其它内部包。
package model

// Wallet 表示一个用户的金币账户。
//
// Balance 用 int64 记账（最小单位，如"分"或"个"），禁止用 float——
// 浮点会引入精度误差，资金域绝不允许。
type Wallet struct {
	UserID  string `json:"user_id"`
	Balance int64  `json:"balance"`
}
