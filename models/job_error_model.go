package models

import (
	"time"
)


type JobError struct {
	ID         uint   `gorm:"primaryKey"`
	AccountId  uint   `json:"account_id"`
	JobGroupId uint   `json:"job_group_id"`
	JobId 	   uint   `json:"job_id"`
	JobName    string `json:"job_name"`
	SpiderName string `json:"spider_name"`

	JobErrorType      string    `json:"job_error_type"`
	JobErrorEngine    string    `json:"job_error_engine"`
	JobErrorName      string    `json:"job_error_name"`
	JobErrorCause     string    `json:"job_error_cause"`
	JobErrorTraceback string    `json:"job_error_traceback"`
	JobErrorFilePath  string    `json:"job_error_file_path"`
	JobErrorUrl       string    `json:"job_error_url"`
	JobErrorCount     uint      `json:"job_error_count"`
	JobErrorDatetime  time.Time `json:"job_error_datetime"` //first time recorded on the user server

	CreatedAt time.Time
	UpdatedAt time.Time
}


type JobGroupErrors struct {
	ID         uint   `gorm:"primaryKey"`
	AccountId  uint   `json:"account_id"`
	JobGroupId uint   `json:"job_group_id"`
	SpiderId 	uint   `json:"spider_id"` // delete
	SpiderName string `json:"spider_name"`
	JobGroupName string `json:"job_group_name"`

	JobGroupErrorType      string    `json:"job_group_error_type"`
	JobGroupErrorEngine    string    `json:"job_group_error_engine"`
	JobGroupErrorName      string    `json:"job_group_error_name"`
	JobGroupErrorCause     string    `json:"job_group_error_cause"`
	JobGroupErrorTraceback string    `json:"job_group_error_traceback"`
	JobGroupErrorFilePath  string    `json:"job_group_error_file_path"`
	JobGroupErrorUrl       string    `json:"job_group_error_url"`
	JobGroupErrorCount     uint      `json:"job_group_error_count"`
	JobGroupErrorDatetime  time.Time `json:"job_group_error_datetime"` //first time recorded on the user server

	CreatedAt time.Time
	UpdatedAt time.Time
}



type JobGroupErrorsMap struct {
	ID         uint   `gorm:"primaryKey"`
	AccountId  uint   `json:"account_id"`
	JobGroupId uint   `json:"job_group_id"`
	JobGroupErrorId 	uint   `json:"job_group_error_id"` // delete
	CreatedAt time.Time
	UpdatedAt time.Time
}


type JobGroupErrorStat struct {
	ID         uint   `gorm:"primaryKey"`
	AccountId  uint   `json:"account_id"`
	JobGroupId uint   `json:"job_group_id"`
	JobGroupUniqueError          int            `json:"job_group_unique_error"`
	JobGroupUniqueWarning        int            `json:"job_group_unique_warning"`
	CreatedAt time.Time
	UpdatedAt time.Time
}