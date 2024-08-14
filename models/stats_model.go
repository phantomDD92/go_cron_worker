package models

import (
	"time"
  	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type StatsPostBody struct {
	ID                uint `gorm:"primaryKey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	JobId             uint                   `json:"job_id"`
	JobGroupId        uint                   `json:"job_group_id"`
	JobGroupName      string                 `json:"job_group_name"`
	Type              string                 `json:"type"`
	PeriodConcurrency int                    `json:"period_concurrency"`
	PeriodStartTime   int64                  `json:"period_start_time"`
	PeriodFinishTime  int64                  `json:"period_finish_time"`
	PeriodRunTime     int                    `json:"period_run_time"`
	PeriodCount       uint                   `json:"period_count"`
	JobFinishTime     int64                  `json:"job_finish_time"`
	JobRunTime        int            		 `json:"job_run_time"`
	JobStatus         string                 `json:"job_status"`
	JobFinishReason   string                 `json:"job_finish_reason"`
	SDKRunTime        int                    `json:"sdk_run_time"`
	Periodic          map[string]interface{} `json:"periodic"`
	Overall           map[string]interface{} `json:"overall"`
	OverallError      int                    `json:"overall_errors"`
	OverallWarning    int                    `json:"overall_warnings"`
	OverallCritical   int                    `json:"overall_criticals"`
	PeriodicError     int                    `json:"periodic_errors"`
	PeriodicWarning   int                    `json:"periodic_warnings"`
	PeriodicCritical  int                    `json:"periodic_criticals"`
	MultiServer       bool                   `json:"multi_server"`
	//LoggingData       LoggingInfo            `json:"logging_data"`
	FailedUrlsCount   int      `json:"failed_urls_count"`
	FailedUrlsEnabled bool     `json:"failed_urls_enabled"`
	FailedUrlsList    []string `json:"failed_urls_list"`

	DataCoverage         postgres.Jsonb         `json:"data_coverage"`
	InvalidItemsCount    int                    `json:"invalid_items_count"`
	InvalidItemsUrlsList map[string]interface{} `json:"invalid_items_urls_list"`
	FieldCoverage        int                    `json:"field_coverage"`

	ScrapyStats     map[string]string `json:"scrapy_stats"`
	JobCustomGroups map[string]string `json:"job_custom_groups"`

	ErrorDetails string `json:"error_details"`
}

type OverallResponseStats struct {
	gorm.Model
	JobGroupId 				uint 
	JobId 					uint //`mapstructure:"periodic_id"`
	Method       string         `mapstructure:"method"`
	Proxy        string         `mapstructure:"proxy"`
	ProxySetup   string         `mapstructure:"proxy_setup"`
	Domain       string         `mapstructure:"domain"`
	PageType     string         `mapstructure:"page_type"`
	StatusCode   string         `mapstructure:"status"`
	Validation   string         `mapstructure:"validation"`
	Geo          string         `mapstructure:"geo"`
	CustomTag    string         `mapstructure:"custom_tag"`
	CustomSignal string         `mapstructure:"custom_signal"`
	Count        uint64            `mapstructure:"count"`
	Bytes        uint64            `mapstructure:"bytes"`
	TotalLatency float32        `mapstructure:"total_latency"`
	MaxLatency   float32        `mapstructure:"max_latency"`
	MinLatency   float32        `mapstructure:"min_latency"`
	AvgLatency   float32        `mapstructure:"avg_latency"`
	Redirects    int            `mapstructure:"redirects"`
	Retries      int            `mapstructure:"retries"`
	Items        uint64            `mapstructure:"items"`
	//Coverage     postgres.Jsonb `mapstructure:"coverage"`
}

type OverallGroupResponseStats struct {
	gorm.Model
	JobGroupId 				uint `mapstructure:"job_group_id"`
	Method       string  `mapstructure:"method"`
	Proxy        string  `mapstructure:"proxy"`
	ProxySetup   string  `mapstructure:"proxy_setup"`
	Domain       string  `mapstructure:"domain"`
	PageType     string  `mapstructure:"page_type"`
	StatusCode   string  `mapstructure:"status"`
	Validation   string  `mapstructure:"validation"`
	Geo          string  `mapstructure:"geo"`
	CustomTag    string  `mapstructure:"custom_tag"`
	CustomSignal string  `mapstructure:"custom_signal"`
	Count        uint64     `mapstructure:"count"`
	Bytes        uint64     `mapstructure:"bytes"`
	TotalLatency float32 `mapstructure:"total_latency"`
	MaxLatency   float32 `mapstructure:"max_latency"`
	MinLatency   float32 `mapstructure:"min_latency"`
	AvgLatency   float32 `mapstructure:"avg_latency"`
	Redirects    int     `mapstructure:"redirects"`
	Retries      int     `mapstructure:"retries"`
	Items        uint64     `mapstructure:"items"`
}