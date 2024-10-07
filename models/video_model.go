package models

import (
	"time"
)

type MarketingWebsiteVideo struct {
	ID                   uint64    `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Website              string    `gorm:"type:varchar(255)" json:"website"`
	WebsitePageTypes     string    `gorm:"type:varchar(255)" json:"website_page_types"`
	VideoName            string    `gorm:"type:varchar(255)" json:"video_name"`
	VideoURL             string    `gorm:"type:varchar(255)" json:"video_url"`
	VideoChannel         string    `gorm:"type:varchar(255)" json:"video_channel"`
	VideoChannelURL      string    `gorm:"type:varchar(255)" json:"video_channel_url"`
	VideoDate            time.Time `gorm:"type:timestamp(6);not null" json:"video_date"`
	VideoSummary         string    `gorm:"type:varchar(255)" json:"video_summary"`
	VideoViews           int       `gorm:"type:integer" json:"video_views"`
	VideoLikes           int       `gorm:"type:integer" json:"video_likes"`
	VideoLength          int       `gorm:"type:integer" json:"video_length"`
	VideoPreviewImageUrl string    `gorm:"type:varchar(255)" json:"video_preview_image_url"`
	CodeLanguage         string    `gorm:"type:varchar(255)" json:"code_language"`
	CodeLibraries        string    `gorm:"type:varchar(255)" json:"code_libraries"`
	CodeLevel            string    `gorm:"type:varchar(255)" json:"code_level"`
	CodeWorks            string    `gorm:"type:varchar(255)" json:"code_works"`
	CodeGithubRepo       string    `gorm:"type:varchar(255)" json:"code_github_repo"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// TableName overrides the table name used by MarketingWebsiteVideo to `marketing_website_videos`
func (MarketingWebsiteVideo) TableName() string {
	return "marketing_website_videos"
}
