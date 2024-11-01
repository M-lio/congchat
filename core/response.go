package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// 定义状态码
const (
	ERROR   = 0
	SUCCESS = 1
)

type Api struct {
	Context *gin.Context
	orm     *gorm.DB
	Errors  error
}

// MakeContext 设置http上下文
func (e *Api) MakeContext(c *gin.Context) *Api {
	e.Context = c
	return e
}

func (e *Api) AddError(err error) {
	if e.Errors == nil {
		e.Errors = err
	} else if err != nil {
		e.Errors = fmt.Errorf("%v; %w", e.Errors, err)
	}
}

// Response 返回数据包装
type Rsp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Result 手动组装返回结果
func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Rsp{
		code,
		msg,
		data,
	})
}
