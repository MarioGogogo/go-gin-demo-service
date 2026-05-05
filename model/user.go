package model

// User 用户模型
type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required,gte=0,lte=150"`
}
