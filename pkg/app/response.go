package app

import (
	"Goutils/pkg/app/errcode"
	//"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
构造 HTTP 响应，针对使用 Gin 框架的 HTTP 服务的一个封装。更加结构化和模块化的方式来处理 HTTP 响应，使得代码更易于理解、维护和发展
*/

type Response struct {
	c *gin.Context
}

// State 请求处理的状态码
type State struct {
	Code int         `json:"code,omitempty"` // 状态码，0：成功；否则失败
	Msg  string      `json:"msg,omitempty"`  // 状态的具体描述
	Data interface{} `json:"data,omitempty"` // 数据，失败时返回空
}

type List struct {
	List  interface{} `json:"list,omitempty"`  // 列表数据，可以是任意类型的集合
	Total int64       `json:"total,omitempty"` // 列表中数据项的总数
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{c: ctx}
}

// Reply 响应单个数据
func (r *Response) Reply(err errcode.Err, datas ...interface{}) {
	var data interface{}
	if len(datas) > 0 {
		data = datas[0]
	}
	if err == nil {
		err = errcode.StatusOk
	} else {
		data = nil
	}
	r.c.JSON(http.StatusOK, State{
		Code: err.ECode(),
		Msg:  err.Error(),
		Data: data,
	})
}

// ReplyList 响应列表数据
func (r *Response) ReplyList(err errcode.Err, total int64, data interface{}) {
	if err == nil {
		err = errcode.StatusOk
	} else {
		data = nil
	}
	r.c.JSON(http.StatusOK, State{
		Code: err.ECode(),
		Msg:  err.Error(),
		Data: List{
			List:  data,
			Total: total,
		},
	})
}
