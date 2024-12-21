package controllers

import (
	"congchat-user/core"
	"congchat-user/service"
	"congchat-user/service/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SysTransferController struct {
	core.Api
}

// Transfer 转账接口
func (e SysTransferController) Transfer(c *gin.Context) {
	var req dto.TransferRequest
	var rsp core.Rsp
	s := new(service.SysTransf)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := s.CreateTransfer(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	rsp.Code = 0
	rsp.Msg = "转账成功"
	c.JSON(http.StatusOK, rsp)
}
