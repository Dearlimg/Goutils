package token

import (
	"encoding/json"
	"github.com/Dearlimg/Goutils/pkg/utils"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strings"
	"testing"
	"time"
)

type M struct {
	UserID   int64  `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
}

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker([]byte(RandomString(32)))
	require.NoError(t, err)
	require.NotEmpty(t, maker)
	userID := RandomInt(1, 1000)
	userName := utils.RandomOwner()
	content := M{
		UserID:   userID,
		UserName: userName,
	}
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	data, err := json.Marshal(content)
	require.NoError(t, err)
	//生成 token
	token, _, err := maker.CreateToken(data, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	// 解析token
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	result := M{}
	err = json.Unmarshal(payload.Content, &result)
	require.NoError(t, err)
	require.Equal(t, content, result)
	//使用测试框架中的 require.WithinDuration 函数来验证 payload.IssuedAt 是否在 issuedAt 时间点的 ±1毫秒范围内
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Millisecond)

	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
}

func TestMaker(t *testing.T) {
	maker, err := NewPasetoMaker([]byte(RandomString(32)))
	require.NoError(t, err)
	require.NotEmpty(t, maker)
	userID := RandomInt(1, 1000)
	userName := utils.RandomOwner()
	content := M{
		UserID:   userID,
		UserName: userName,
	}
	data, err := json.Marshal(content)
	require.NoError(t, err)
	duration := time.Second
	//生成token
	token, _, err := maker.CreateToken(data, duration)
	require.NoError(t, err)
	result, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	time.Sleep(duration * 2)
	//此时 token 已经超时了
	result2, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, result2)
}

const aplphabetic = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt 返回 min 到 max 之间的一个随机数
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString 生成一个长度为 n 的随机字符串
func RandomString(n int) string {
	var sb strings.Builder
	k := len(aplphabetic)
	for i := 0; i < n; i++ {
		c := aplphabetic[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}
