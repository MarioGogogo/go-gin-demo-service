package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"go-gin-demo-service/database"
	"go-gin-demo-service/model"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, age FROM users ORDER BY id")
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "查询失败"))
		return
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "读取数据失败"))
			return
		}
		users = append(users, u)
	}

	c.JSON(http.StatusOK, model.Ok(gin.H{
		"list":  users,
		"total": len(users),
	}))
}

func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeBadRequest, "无效的ID"))
		return
	}

	var u model.User
	err = database.DB.QueryRow("SELECT id, name, age FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Age)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, model.Fail(model.CodeNotFound, "用户不存在"))
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "查询失败"))
		return
	}

	c.JSON(http.StatusOK, model.Ok(u))
}

func CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeBadRequest, err.Error()))
		return
	}

	var id uint
	err := database.DB.QueryRow("INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id", user.Name, user.Age).Scan(&id)
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "创建失败"))
		return
	}
	user.ID = id

	c.JSON(http.StatusOK, model.OkWithMsg(user, "创建成功"))
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeBadRequest, "无效的ID"))
		return
	}

	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeBadRequest, err.Error()))
		return
	}

	result, err := database.DB.Exec("UPDATE users SET name = $1, age = $2 WHERE id = $3", input.Name, input.Age, id)
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "更新失败"))
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusOK, model.Fail(model.CodeNotFound, "用户不存在"))
		return
	}

	input.ID = uint(id)
	c.JSON(http.StatusOK, model.OkWithMsg(input, "更新成功"))
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeBadRequest, "无效的ID"))
		return
	}

	result, err := database.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "删除失败"))
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusOK, model.Fail(model.CodeNotFound, "用户不存在"))
		return
	}

	c.JSON(http.StatusOK, model.OkWithMsg(nil, "删除成功"))
}