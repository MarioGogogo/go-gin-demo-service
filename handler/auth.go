package handler

import (
	"net/http"

	"go-gin-demo-service/model"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginData struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeBadRequest, "用户名和密码不能为空"))
		return
	}

	// 硬编码账户密码
	if req.Username == "admin" && req.Password == "123456" {
		data := LoginData{
			Token:    "fake-token-admin",
			Username: req.Username,
		}
		c.JSON(http.StatusOK, model.OkWithMsg(data, "登录成功"))
		return
	}

	c.JSON(http.StatusOK, model.Fail(model.CodeUnauthorized, "用户名或密码错误"))
}
