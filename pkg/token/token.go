package token

import (
	"errors"
	"time"
)

/*
Paseto（Platform-Agnostic Security Tokens）令牌是一种用于安全通信的令牌格式，旨在提供简单、安全和可靠的令牌生成和验证机制。
相比于 JWT，Paseto 所做的改变：
- 不会向用户开放所有的加密算法
- header中不再含有 alg 字段，也不会有 none 算法
- payload 使用加密算法，而不是简单的编码
*/

var ErrSecretLen = errors.New("密钥长度不正确")

type MakerToken interface {
	// CreateToken 生成 MakerToken
	CreateToken(content []byte, expireDate time.Duration) (string, *Payload, error)
	// VerifyToken 解析 MakerToken
	VerifyToken(token string) (*Payload, error)
}
