package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go-gin-demo-service/database"
	"go-gin-demo-service/model"

	"github.com/gin-gonic/gin"
)

func CheckUpdate(c *gin.Context) {
	var info model.AppUpdateInfo
	err := database.DB.QueryRow(
		"SELECT has_update, version_code, version_name, download_url, changelog, force_update, file_size FROM app_updates ORDER BY id DESC LIMIT 1",
	).Scan(&info.HasUpdate, &info.VersionCode, &info.VersionName, &info.DownloadUrl, &info.Changelog, &info.ForceUpdate, &info.FileSize)
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "查询更新信息失败"))
		return
	}

	// 从 downloadUrl 中提取文件名，读取 static 目录中的实际文件大小
	fileName := info.DownloadUrl[strings.LastIndex(info.DownloadUrl, "/")+1:]
	staticPath := filepath.Join(staticDir, fileName)
	if f, err := os.Stat(staticPath); err == nil {
		info.FileSize = f.Size()
	}

	c.JSON(http.StatusOK, model.Ok(info))
}
