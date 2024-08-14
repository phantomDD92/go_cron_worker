
package models

import (
	"time"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type Spider struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	AccountId uint `json:"account_id"`
	ServerId      uint   `json:"server_id"`
	SpiderName    string `json:"spider_name"`
	SpiderType    string `json:"spider_type"`
	SpiderStatus  string `json:"spider_status"`
	SpiderVersion string `json:"spider_version"`
	SpiderSettings          postgres.Jsonb `json:"spider_settings"`
	SpiderGeneratedType     string         `json:"spider_generated_type"`
	BotName                 string         `json:"bot_name"`
	SpiderMiddlewareEnabled bool           `json:"spider_middleware_enabled"`
}

type SpiderStat struct {
	ID                           uint `gorm:"primaryKey"`
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
	AccountId                    uint           `json:"account_id"`
	SpiderName                   string         `json:"spider_name"`
	SpiderStatStartTime            time.Time      `json:"spider_stat_start_time"`
	SpiderStatLastSdkUpdate        time.Time      `json:"spider_stat_last_sdk_update"`
	SpiderStatRunTime              int            `json:"spider_stat_run_time"`
	SpiderStatRunning               int         	`json:"spider_stat_running"`
	SpiderStatFinished         int        	 `json:"spider_stat_finished"`
	SpiderStatOverallGroupError    int            `json:"spider_stat_overall_group_error"`
	SpiderStatOverallGroupWarning  int            `json:"spider_stat_overall_group_warning"`
	SpiderStatOverallGroupCritical int            `json:"spider_stat_overall_group_critical"`
	SpiderStatRequests             int            `json:"spider_stat_requests"`
	SpiderStatResponses             int            `json:"spider_stat_responses"`
	SpiderStat_200Responses        int            `json:"spider_stat_200_responses"`
	SpiderStatNon_300Responses     int            `json:"spider_stat_non_300_responses"`
	SpiderStatBytes                float64        `json:"spider_stat_bytes"`
	SpiderStatSuccessRate          float64        `json:"spider_stat_success_rate"`
	SpiderStatAvgLatency           float64        `json:"spider_stat_avg_latency"`
	SpiderStatTotalLatency         float64        `json:"spider_stat_total_latency"`
	SpiderStatStatusCodes          postgres.Jsonb `json:"spider_stat_status_codes"`
	SpiderStatItems                uint64            `json:"spider_stat_items"`
	SpiderStatChecks               int            `json:"spider_stat_checks"`
	SpiderStatPassedChecks         int            `json:"spider_stat_passed_checks"`
	SpiderStatFailedUrls           int            `json:"spider_stat_failed_urls"`
	SpiderStatDataCoverage         postgres.Jsonb `json:"spider_stat_data_coverage"`
	SpiderStatInvalidItems         int          `json:"spider_stat_invalid_items"`
	SpiderStatFieldCoverage		 	int	   		`json:"spider_stat_field_coverage"`
	SpiderStatNumJobs				int	   		`json:"spider_stat_num_jobs"`
	SpiderStatUniqueJobs			int	   		`json:"spider_stat_unique_jobs"`
	SpiderStatExcludeMovingAverage	bool	   		`json:"spider_stat_exclude_moving_average"`
	SpiderStatItemsArray			string	   	`json:"spider_stat_items_array"`	
	SpiderStatAggHour  				int            `json:"spider_stat_agg_hour"`
}


type SpiderStatRedisVersion struct {
	ID                           uint `gorm:"primaryKey"`
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
	AccountId                    uint           `json:"account_id"`
	SpiderName                   string         `json:"spider_name"`
	SpiderStatStartTime            time.Time      `json:"spider_stat_start_time"`
	SpiderStatLastSdkUpdate        time.Time      `json:"spider_stat_last_sdk_update"`
	SpiderStatRunTime              int            `json:"spider_stat_run_time"`
	SpiderStatRunning               int         	`json:"spider_stat_running"`
	SpiderStatFinished         int        	 `json:"spider_stat_finished"`
	SpiderStatOverallGroupError    int            `json:"spider_stat_overall_group_error"`
	SpiderStatOverallGroupWarning  int            `json:"spider_stat_overall_group_warning"`
	SpiderStatOverallGroupCritical int            `json:"spider_stat_overall_group_critical"`
	SpiderStatRequests             int            `json:"spider_stat_requests"`
	SpiderStatResponses             int            `json:"spider_stat_responses"`
	SpiderStat_200Responses        int            `json:"spider_stat_200_responses"`
	SpiderStatNon_300Responses     int            `json:"spider_stat_non_300_responses"`
	SpiderStatBytes                float64        `json:"spider_stat_bytes"`
	SpiderStatSuccessRate          float64        `json:"spider_stat_success_rate"`
	SpiderStatAvgLatency           float64        `json:"spider_stat_avg_latency"`
	SpiderStatTotalLatency         float64        `json:"spider_stat_total_latency"`
	SpiderStatStatusCodes          map[string]int `json:"spider_stat_status_codes"`
	SpiderStatItems                uint64            `json:"spider_stat_items"`
	SpiderStatChecks               int            `json:"spider_stat_checks"`
	SpiderStatPassedChecks         int            `json:"spider_stat_passed_checks"`
	SpiderStatFailedUrls           int            `json:"spider_stat_failed_urls"`
	SpiderStatDataCoverage         map[string]interface{} `json:"spider_stat_data_coverage"`
	SpiderStatInvalidItems         int          `json:"spider_stat_invalid_items"`
	SpiderStatFieldCoverage		 	int	   		`json:"spider_stat_field_coverage"`
	SpiderStatNumJobs				int	   		`json:"spider_stat_num_jobs"`
	SpiderStatUniqueJobs			int	   		`json:"spider_stat_unique_jobs"`
	SpiderStatExcludeMovingAverage	bool	   		`json:"spider_stat_exclude_moving_average"`
	SpiderStatItemsArray			string	   	`json:"spider_stat_items_array"`	
	SpiderStatAggHour  				int            `json:"spider_stat_agg_hour"`
}

type SpiderStatAvg struct {
	ID                           uint `gorm:"primaryKey"`
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
	
	AccountId                    uint           `json:"account_id"`
	SpiderName                 string         `json:"spider_name"`

	SpiderStatAvgWindowType           string            `json:"spider_stat_avg_window_type"`
	SpiderStatAvgWindowLength          int            `json:"spider_stat_avg_window_length"`
	SpiderStatAvgNumJobs          int            `json:"spider_stat_avg_num_jobs"`
	SpiderStatAvgRunning          int            `json:"spider_stat_avg_running"`
	SpiderStatAvgFinished          int            `json:"spider_stat_avg_finished"`

	SpiderStatAvgRunTime              int            `json:"spider_stat_avg_run_time"`
	SpiderStatAvgOverallGroupError    int            `json:"spider_stat_avg_overall_group_error"`
	SpiderStatAvgOverallGroupWarning  int            `json:"spider_stat_avg_overall_group_warning"`
	SpiderStatAvgOverallGroupCritical int            `json:"spider_stat_avg_overall_group_critical"`
	SpiderStatAvgRequests             int            `json:"spider_stat_avg_requests"`
	SpiderStatAvg_200Responses        int            `json:"spider_stat_avg_200_responses"`
	SpiderStatAvgNon_300Responses     int            `json:"spider_stat_avg_non_300_responses"`
	SpiderStatAvgBytes                float64        `json:"spider_stat_avg_bytes"`
	SpiderStatAvgSuccessRate          float64        `json:"spider_stat_avg_success_rate"`
	SpiderStatAvgAvgLatency           float64        `json:"spider_stat_avg_avg_latency"`
	SpiderStatAvgItems                uint64            `json:"spider_stat_avg_items"`
	SpiderStatAvgChecks               int            `json:"spider_stat_avg_checks"`
	SpiderStatAvgPassedChecks         int            `json:"spider_stat_avg_passed_checks"`
	SpiderStatAvgFailedUrls           int            `json:"spider_stat_avg_failed_urls"`
	SpiderStatAvgInvalidItems         int          `json:"spider_stat_avg_invalid_items"`
	SpiderStatAvgFieldCoverage		 int	   		`json:"spider_stat_avg_field_coverage"`
}



type SpiderError struct {
	ID         uint   `gorm:"primaryKey"`
	AccountId  uint   `json:"account_id"`
	SpiderStatId uint   `json:"spider_stat_id"`
	SpiderName string `json:"spider_name"`
	SpiderErrorStartTime            time.Time      `json:"spider_error_start_time"`

	SpiderErrorType      string    `json:"spider_error_type"`
	SpiderErrorEngine    string    `json:"spider_error_engine"`
	SpiderErrorName      string    `json:"spider_error_name"`
	SpiderErrorCause     string    `json:"spider_error_cause"`
	SpiderErrorTraceback string    `json:"spider_error_traceback"`
	SpiderErrorFilePath  string    `json:"spider_error_file_path"`
	SpiderErrorUrl       string    `json:"spider_error_url"`
	SpiderErrorCount     uint      `json:"spider_error_count"`
	SpiderErrorDatetime  time.Time `json:"spider_error_datetime"` //first time recorded on the user server

	CreatedAt time.Time
	UpdatedAt time.Time
}


type SpiderStatErrorsMap struct {
	ID         uint   `gorm:"primaryKey"`
	AccountId  uint   `json:"account_id"`
	SpiderStatId uint   `json:"spider_stat_id"`
	SpiderErrorId 	uint   `json:"spider_error_id"` 
	CreatedAt time.Time
	UpdatedAt time.Time
}



type SpiderStatErrorStat struct {
	ID         uint   `gorm:"primaryKey"`
	AccountId  uint   `json:"account_id"`
	SpiderStatId uint   `json:"spider_stat_id"`
	SpiderStatUniqueError          int            `json:"spider_stat_unique_error"`
	SpiderStatUniqueWarning        int            `json:"spider_stat_unique_warning"`
	CreatedAt time.Time
	UpdatedAt time.Time
}