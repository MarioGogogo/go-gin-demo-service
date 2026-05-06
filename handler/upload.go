package handler

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go-gin-demo-service/database"
	"go-gin-demo-service/model"

	"github.com/gin-gonic/gin"
)

const modulesDir = "./uploads/modules"

func init() {
	os.MkdirAll(modulesDir, 0755)
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// GetModules 获取模块列表，支持 ?type= 筛选
func GetModules(c *gin.Context) {
	filterType := c.Query("type")

	var rows *sql.Rows
	var err error
	if filterType != "" {
		rows, err = database.DB.Query("SELECT id, name, type, version, file_name, file_path, file_size, changelog, created_at, updated_at FROM modules WHERE type = $1 ORDER BY updated_at DESC", filterType)
	} else {
		rows, err = database.DB.Query("SELECT id, name, type, version, file_name, file_path, file_size, changelog, created_at, updated_at FROM modules ORDER BY updated_at DESC")
	}
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "查询失败"))
		return
	}
	defer rows.Close()

	type moduleRow struct {
		model.Module
	}
	var modules []model.Module
	for rows.Next() {
		var m model.Module
		if err := rows.Scan(&m.ID, &m.Name, &m.Type, &m.Version, &m.FileName, &m.FilePath, &m.FileSize, &m.Changelog, &m.CreatedAt, &m.UpdatedAt); err != nil {
			c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "读取数据失败"))
			return
		}
		// 查询版本历史
		hrows, err := database.DB.Query("SELECT version, file_name, file_size, changelog, created_at FROM version_histories WHERE module_id = $1 ORDER BY created_at DESC", m.ID)
		if err == nil {
			for hrows.Next() {
				var h model.VersionHistory
				hrows.Scan(&h.Version, &h.FileName, &h.FileSize, &h.Changelog, &h.CreatedAt)
				m.History = append(m.History, h)
			}
			hrows.Close()
		}
		modules = append(modules, m)
	}

	if modules == nil {
		modules = []model.Module{}
	}

	c.JSON(http.StatusOK, model.Ok(gin.H{
		"list":  modules,
		"total": len(modules),
	}))
}

// UploadModule 上传模块
func UploadModule(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeBadRequest, "请选择文件"))
		return
	}

	name := c.PostForm("name")
	moduleType := c.PostForm("type")
	version := c.PostForm("version")
	changelog := c.PostForm("changelog")
	code := c.PostForm("code")
	downloadUrl := c.PostForm("downloadUrl")

	if moduleType == "" || version == "" {
		c.JSON(http.StatusOK, model.Fail(model.CodeBadRequest, "类型和版本号不能为空"))
		return
	}
	if moduleType == "hotel-module" && name == "" {
		c.JSON(http.StatusOK, model.Fail(model.CodeBadRequest, "功能名称不能为空"))
		return
	}
	if moduleType == "app" {
		name = "App"
	}

	id := generateID()
	fileName := id + "_" + file.Filename
	savePath := filepath.Join(modulesDir, fileName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "文件保存失败"))
		return
	}

	now := time.Now()

	// App类型：查找已有app记录进行覆盖
	if moduleType == "app" {
		var existID string
		var existPath string
		err := database.DB.QueryRow("SELECT id, file_path FROM modules WHERE type = 'app' LIMIT 1").Scan(&existID, &existPath)
		if err == nil {
			// 追加版本历史
			database.DB.Exec("INSERT INTO version_histories (module_id, version, file_name, file_size, changelog, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
				existID, version, file.Filename, file.Size, changelog, now)
			// 更新模块信息
			database.DB.Exec("UPDATE modules SET version=$1, file_name=$2, file_path=$3, file_size=$4, changelog=$5, updated_at=$6 WHERE id=$7",
				version, file.Filename, savePath, file.Size, changelog, now, existID)
			if existPath != "" {
				os.Remove(existPath)
			}
			// 返回更新后的数据
			m := model.Module{ID: existID, Name: name, Type: moduleType, Version: version, FileName: file.Filename, FilePath: savePath, FileSize: file.Size, Changelog: changelog, UpdatedAt: now}
			c.JSON(http.StatusOK, model.OkWithMsg(m, "上传成功"))
			return
		}
	}

	// 分包类型：同名同类型覆盖
	if moduleType == "hotel-module" {
		var existID string
		var existPath string
		err := database.DB.QueryRow("SELECT id, file_path FROM modules WHERE name = $1 AND type = $2", name, moduleType).Scan(&existID, &existPath)
		if err == nil {
			database.DB.Exec("INSERT INTO version_histories (module_id, version, file_name, file_size, changelog, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
				existID, version, file.Filename, file.Size, changelog, now)
			database.DB.Exec("UPDATE modules SET version=$1, file_name=$2, file_path=$3, file_size=$4, changelog=$5, code=$6, download_url=$7, updated_at=$8 WHERE id=$9",
				version, file.Filename, savePath, file.Size, changelog, code, downloadUrl, now, existID)
			if existPath != "" {
				os.Remove(existPath)
			}
			m := model.Module{ID: existID, Name: name, Type: moduleType, Version: version, FileName: file.Filename, FilePath: savePath, FileSize: file.Size, Changelog: changelog, UpdatedAt: now}
			c.JSON(http.StatusOK, model.OkWithMsg(m, "上传成功"))
			return
		}
	}

	// 新记录
		_, err = database.DB.Exec("INSERT INTO modules (id, name, type, version, file_name, file_path, file_size, changelog, code, download_url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			id, name, moduleType, version, file.Filename, savePath, file.Size, changelog, code, downloadUrl, now, now)
	if err != nil {
		os.Remove(savePath)
		c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "保存失败"))
		return
	}
	database.DB.Exec("INSERT INTO version_histories (module_id, version, file_name, file_size, changelog, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		id, version, file.Filename, file.Size, changelog, now)

	m := model.Module{ID: id, Name: name, Type: moduleType, Version: version, FileName: file.Filename, FilePath: savePath, FileSize: file.Size, Changelog: changelog, CreatedAt: now, UpdatedAt: now}
	c.JSON(http.StatusOK, model.OkWithMsg(m, "上传成功"))
}

// UpdateModule 编辑模块信息
func UpdateModule(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name      string `json:"name"`
		Version   string `json:"version"`
		Changelog string `json:"changelog"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeBadRequest, "参数错误"))
		return
	}

	var existID string
	err := database.DB.QueryRow("SELECT id FROM modules WHERE id = $1", id).Scan(&existID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, model.Fail(model.CodeNotFound, "模块不存在"))
		return
	}

	if req.Version != "" {
		database.DB.Exec("UPDATE modules SET version = $1, updated_at = $2 WHERE id = $3", req.Version, time.Now(), id)
	}
	if req.Name != "" {
		database.DB.Exec("UPDATE modules SET name = $1, updated_at = $2 WHERE id = $3", req.Name, time.Now(), id)
	}
	database.DB.Exec("UPDATE modules SET changelog = $1, updated_at = $2 WHERE id = $3", req.Changelog, time.Now(), id)

	// 返回更新后的数据
	var m model.Module
	database.DB.QueryRow("SELECT id, name, type, version, file_name, file_path, file_size, changelog, created_at, updated_at FROM modules WHERE id = $1", id).
		Scan(&m.ID, &m.Name, &m.Type, &m.Version, &m.FileName, &m.FilePath, &m.FileSize, &m.Changelog, &m.CreatedAt, &m.UpdatedAt)

	c.JSON(http.StatusOK, model.OkWithMsg(m, "更新成功"))
}

// DeleteModule 删除模块
func DeleteModule(c *gin.Context) {
	id := c.Param("id")

	var filePath string
	err := database.DB.QueryRow("SELECT file_path FROM modules WHERE id = $1", id).Scan(&filePath)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, model.Fail(model.CodeNotFound, "模块不存在"))
		return
	}

	database.DB.Exec("DELETE FROM modules WHERE id = $1", id)
	if filePath != "" {
		os.Remove(filePath)
	}

	c.JSON(http.StatusOK, model.OkWithMsg(nil, "删除成功"))
}

// GetModuleHistory 获取模块版本历史
func GetModuleHistory(c *gin.Context) {
	id := c.Param("id")

	rows, err := database.DB.Query("SELECT version, file_name, file_size, changelog, created_at FROM version_histories WHERE module_id = $1 ORDER BY created_at DESC", id)
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "查询失败"))
		return
	}
	defer rows.Close()

	var history []model.VersionHistory
	for rows.Next() {
		var h model.VersionHistory
		rows.Scan(&h.Version, &h.FileName, &h.FileSize, &h.Changelog, &h.CreatedAt)
		history = append(history, h)
	}

	if history == nil {
		history = []model.VersionHistory{}
	}

	c.JSON(http.StatusOK, model.Ok(gin.H{
		"list":  history,
		"total": len(history),
	}))
}

// DownloadModule 下载模块文件
func DownloadModule(c *gin.Context) {
	id := c.Param("id")

	var filePath, fileName string
	err := database.DB.QueryRow("SELECT file_path, file_name FROM modules WHERE id = $1", id).Scan(&filePath, &fileName)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, model.Fail(model.CodeNotFound, "模块不存在"))
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusOK, model.Fail(model.CodeNotFound, "文件不存在"))
		return
	}

	c.FileAttachment(filePath, fileName)
}

// AdminPage 管理页面
func AdminPage(c *gin.Context) {
	file, err := os.Open("./web/index.html")
	if err != nil {
		c.String(http.StatusNotFound, "页面不存在")
		return
	}
	defer file.Close()
	data, _ := io.ReadAll(file)
	c.Data(http.StatusOK, "text/html; charset=utf-8", data)
}
