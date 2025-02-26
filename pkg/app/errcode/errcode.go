package errcode

import (
	"fmt"
	"github.com/jinzhu/copier"
	"sync"
)

/*
常用的一些错误处理公共方法，标准化我们的错误输出
*/

type Err interface {
	Error() string
	ECode() int
	WithDetails(details ...string) Err
}

var globalMap map[int]Err
var once sync.Once

type myErr struct {
	Code    int      `json:"code,omitempty"`    //状态码，0：成功；其他代表失败
	Msg     string   `json:"msg,omitempty"`     //返回状态描述
	Details []string `json:"details,omitempty"` //详细信息
}

// ECode 返回错误码
func (m *myErr) ECode() int {
	return m.Code
}

// Error 返回错误
func (m *myErr) Error() string {
	return fmt.Sprintf("错误码：%v，错误信息：%v，详细信息：%v", m.Code, m.Msg, m.Details)
}

// WithDetails 返回带有详细信息的错误
func (m *myErr) WithDetails(details ...string) Err {
	var newErr = &myErr{}
	_ = copier.Copy(newErr, m) //深层次拷贝
	newErr.Details = append(newErr.Details, details...)
	return newErr
}

// NewErr 根据错误码和错误信息创建新的错误
func NewErr(code int, msg string) Err {
	once.Do(func() {
		globalMap = make(map[int]Err)
	})
	if _, ok := globalMap[code]; ok {
		panic("错误码已存在")
	}
	err := &myErr{Code: code, Msg: msg}
	globalMap[code] = err
	return err
}
