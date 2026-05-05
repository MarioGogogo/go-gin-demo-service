package model

import "time"

type Module struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Version     string    `json:"version"`
	FileName    string    `json:"fileName"`
	FilePath    string    `json:"filePath"`
	FileSize    int64     `json:"fileSize"`
	Changelog   string    `json:"changelog"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	History     []VersionHistory `json:"history,omitempty"`
}

type VersionHistory struct {
	Version   string    `json:"version"`
	FileName  string    `json:"fileName"`
	FileSize  int64     `json:"fileSize"`
	Changelog string    `json:"changelog"`
	CreatedAt time.Time `json:"createdAt"`
}
