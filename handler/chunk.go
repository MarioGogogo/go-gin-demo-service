package handler

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"

	"go-gin-demo-service/database"
	"go-gin-demo-service/model"

	"github.com/gin-gonic/gin"
)

type ChunkModuleItem struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	DownloadUrl string `json:"downloadUrl"`
	Version     string `json:"version"`
	FileSize    int64  `json:"fileSize"`
}

func GetChunkModuleList(c *gin.Context) {
	rows, err := database.DB.Query(
		"SELECT name, code, version, download_url FROM modules WHERE type = 'hotel-module' ORDER BY updated_at DESC")
	if err != nil {
		c.JSON(http.StatusOK, model.Fail(model.CodeServerError, "查询失败"))
		return
	}
	defer rows.Close()

	var results []ChunkModuleItem
	for rows.Next() {
		var name, code, version string
		var downloadUrl sql.NullString
		if err := rows.Scan(&name, &code, &version, &downloadUrl); err != nil {
			continue
		}

		var fileSize int64
		staticPath := filepath.Join(staticDir, code+".chunk.bundle")
		if info, err := os.Stat(staticPath); err == nil {
			fileSize = info.Size()
		}

		item := ChunkModuleItem{
			Name:    name,
			Code:    code,
			Version: version,
			FileSize: fileSize,
		}
		if downloadUrl.Valid {
			item.DownloadUrl = downloadUrl.String
		}
		results = append(results, item)
	}

	if results == nil {
		results = []ChunkModuleItem{}
	}

	c.JSON(http.StatusOK, model.OkWithMsg(results, "请求成功"))
}
