package models

import (
	"time"
)

type GithubRepo struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Website          string    `gorm:"size:255;column:website"`
	WebsitePageTypes string    `gorm:"size:255;column:website_page_types"`
	RepoName         string    `gorm:"size:255;column:repo_name"`
	RepoURL          string    `gorm:"size:255;column:repo_url"`
	RepoImageURL     string    `gorm:"size:255;column:repo_image_url"`
	RepoSummary      string    `gorm:"size:255;column:repo_summary"`
	RepoOwner        string    `gorm:"size:255;column:repo_owner"`
	RepoCreatedDate  time.Time `gorm:"not null;column:repo_created_date"`
	RepoLastUpdated  time.Time `gorm:"not null;column:repo_last_updated"`
	RepoMaintained   bool      `gorm:"not null;default:true;column:repo_maintained"`
	RepoStars        int       `gorm:"column:repo_stars"`
	RepoForks        int       `gorm:"column:repo_forks"`
	CodeLanguage     string    `gorm:"size:255;column:code_language"`
	CodeLibraries    string    `gorm:"size:255;column:code_libraries"`
	CodeLevel        string    `gorm:"size:255;column:code_level"`
	CodeWorks        string    `gorm:"size:255;column:code_works"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// TableName sets the default table name
func (GithubRepo) TableName() string {
	return "marketing_website_github_repos"
}
