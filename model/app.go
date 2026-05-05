package model

type AppUpdateInfo struct {
	HasUpdate   bool   `json:"hasUpdate"`
	VersionCode int    `json:"versionCode"`
	VersionName string `json:"versionName"`
	DownloadUrl string `json:"downloadUrl"`
	Changelog   string `json:"changelog"`
	ForceUpdate bool   `json:"forceUpdate"`
	FileSize    int64  `json:"fileSize"`
}
