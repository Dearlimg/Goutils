package token

import (
	"github.com/google/uuid"
	"time"
)

// Payload 负载
type Payload struct {
	ID        uuid.UUID // uuid.UUID 用于管理每个 token 的唯一标识符，保证了在分布式系统中生成的 ID 具有唯一性
	Content   []byte    `json:"content,omitempty"` // token 的内容信息可以是任何内容
	IssuedAt  time.Time `json:"issued_at"`         // token 的签发时间
	ExpiredAt time.Time `json:"expired_at"`        // token 的过期时间
}

// NewPayload 创建一个新的 Payload，传入 JWT 内容信息以及多长时间间隔之后过期
func NewPayload(content []byte, expireDate time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom() //随机生成一个 uuid，返回生成的 uuid 值以及一个可能的错误
	if err != nil {
		return nil, err
	}
	return &Payload{
		ID:        tokenID,
		Content:   content,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(expireDate),
	}, nil
}
