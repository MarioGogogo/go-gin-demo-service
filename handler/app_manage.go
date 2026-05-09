package handler

import (
	"fmt"
	"net/http"
	"os"

	"go-gin-demo-service/database"
	"go-gin-demo-service/model"

	"github.com/gin-gonic/gin"
)

// GetAppList 获取 App 更新记录列表
func GetAppList(c *gin.Context) {
	rows, err := database.DB.Query(
		"SELECT id, has_update, version_code, version_name, download_url, changelog, force_update, file_size FROM app_updates ORDER BY id DESC")
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "查询失败"))
		return
	}
	defer rows.Close()

	type AppRow struct {
		ID          int    `json:"id"`
		HasUpdate   bool   `json:"hasUpdate"`
		VersionCode int    `json:"versionCode"`
		VersionName string `json:"versionName"`
		DownloadUrl string `json:"downloadUrl"`
		Changelog   string `json:"changelog"`
		ForceUpdate bool   `json:"forceUpdate"`
		FileSize    int64  `json:"fileSize"`
	}

	var list []AppRow
	for rows.Next() {
		var r AppRow
		rows.Scan(&r.ID, &r.HasUpdate, &r.VersionCode, &r.VersionName, &r.DownloadUrl, &r.Changelog, &r.ForceUpdate, &r.FileSize)
		list = append(list, r)
	}
	if list == nil {
		list = []AppRow{}
	}

	c.JSON(http.StatusOK, model.Ok(gin.H{
		"list":  list,
		"total": len(list),
	}))
}

// UploadApp 上传 App 文件，写入 app_updates 表
func UploadApp(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeBadRequest, "请选择文件"))
		return
	}

	version := c.PostForm("version")
	changelog := c.PostForm("changelog")
	if version == "" {
		c.JSON(http.StatusOK, model.Fail(model.CodeBadRequest, "版本号不能为空"))
		return
	}

	savePath := appDir + "/app.apk"
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "文件保存失败"))
		return
	}

	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	downloadUrl := fmt.Sprintf("%s://%s/download/app.apk", scheme, c.Request.Host)

	_, err = database.DB.Exec(
		`INSERT INTO app_updates (has_update, version_code, version_name, download_url, changelog, force_update, file_size)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		true, 0, version, downloadUrl, changelog, false, file.Size)
	if err != nil {
		os.Remove(savePath)
		c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "保存失败"))
		return
	}

	c.JSON(http.StatusOK, model.OkWithMsg(nil, "上传成功"))
}

// DeleteApp 删除 App 更新记录
func DeleteApp(c *gin.Context) {
	id := c.Param("id")

	var filePath string
	err := database.DB.QueryRow("SELECT download_url FROM app_updates WHERE id = $1", id).Scan(&filePath)
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeNotFound, "记录不存在"))
		return
	}

	database.DB.Exec("DELETE FROM app_updates WHERE id = $1", id)

	c.JSON(http.StatusOK, model.OkWithMsg(nil, "删除成功"))
}
