package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go-gin-demo-service/model"

	"github.com/gin-gonic/gin"
)

// 资源目录路径
const staticDir = "./static"

// Download 下载资源
func Download(c *gin.Context) {
	filename := c.Param("filename")

	// 安全检查：防止路径穿越
	if strings.Contains(filename, "..") {
		c.JSON(http.StatusOK, model.Fail(model.CodeBadRequest, "非法文件名"))
		return
	}

	filePath := filepath.Join(staticDir, filename)

	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		c.JSON(http.StatusOK, model.Fail(model.CodeNotFound, "文件不存在"))
		return
	}

	c.FileAttachment(filePath, filename)
}

// ListFiles 列出可下载的资源
func ListFiles(c *gin.Context) {
	entries, err := os.ReadDir(staticDir)
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "读取目录失败"))
		return
	}

	var files []gin.H
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, _ := entry.Info()
		files = append(files, gin.H{
			"name": entry.Name(),
			"size": info.Size(),
			"url":  "/download/" + entry.Name(),
		})
	}

	c.JSON(http.StatusOK, model.Ok(gin.H{
		"list":  files,
		"total": len(files),
	}))
}
