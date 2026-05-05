package model

import "time"

// Response 统一响应结构
type Response struct {
	Success   bool        `json:"success"`     // 是否成功
	Code      int         `json:"code"`        // 业务状态码
	Message   string      `json:"message"`     // 提示信息
	Data      interface{} `json:"data"`        // 数据
	Timestamp int64       `json:"timestamp"`   // 时间戳
}

// 分页响应
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// 业务状态码
const (
	CodeSuccess      = 0
	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeNotFound     = 404
	CodeServerError  = 500
)

// Ok 成功响应
func Ok(data interface{}) Response {
	return Response{
		Success:   true,
		Code:      CodeSuccess,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// OkWithMsg 成功响应（自定义消息）
func OkWithMsg(data interface{}, message string) Response {
	return Response{
		Success:   true,
		Code:      CodeSuccess,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// Fail 失败响应
func Fail(code int, message string) Response {
	return Response{
		Success:   false,
		Code:      code,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().Unix(),
	}
}