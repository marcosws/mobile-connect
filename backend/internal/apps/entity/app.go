package entity

import "time"

type App struct {
	ID           string `json:"id"`
	FileName     string `json:"fileName"`
	OriginalName string `json:"originalName"`

	PackageName string `json:"packageName"`
	VersionName string `json:"versionName"`
	VersionCode string `json:"versionCode"`

	MinSDK    int `json:"minSdk"`
	TargetSDK int `json:"targetSdk"`

	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`

	CreatedAt time.Time `json:"createdAt"`
}
