package handler

import (
	"net/http"

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

	c.JSON(http.StatusOK, model.Ok(info))
}
