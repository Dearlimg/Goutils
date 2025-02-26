package token

import (
	"errors"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	"time"
)

type PasetoMaker struct {
	paseto *paseto.V2 //Paseto 实例，用于生成和验证 Paseto 令牌
	key    []byte     //用于加密和解密的密钥
}

// NewPasetoMaker 创建 PasetoMaker 实例
func NewPasetoMaker(key []byte) (MakerToken, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, ErrSecretLen
	}
	return &PasetoMaker{
		paseto: paseto.NewV2(),
		key:    key,
	}, nil
}

// CreateToken 生成 Token
func (p *PasetoMaker) CreateToken(content []byte, expireDate time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(content, expireDate)
	if err != nil {
		return "", nil, err
	}
	//使用 Paseto 实例的 Encrypt(加密) 方法，使用密钥 p.Key 对 payload 进行加密，并返回生成的令牌
	token, err := p.paseto.Encrypt(p.key, payload, nil)
	if err != nil {
		return "", nil, err
	}
	return token, payload, nil
}

// VerifyToken 解析 Token
func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	// 使用 Paseto 实例的 Decrypt 方法解密令牌 token
	err := p.paseto.Decrypt(token, p.key, payload, nil)
	if err != nil {
		return nil, err
	}
	// 验证 token 是否已经超过过期时间
	if payload.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("超时错误")
	}
	return payload, nil
}
