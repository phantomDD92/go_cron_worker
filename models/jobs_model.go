package models

import (
	"time"
  	"github.com/jinzhu/gorm/dialects/postgres"
	// "gorm.io/datatypes"
)

type Job struct {
	ID                uint `gorm:"primaryKey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	AccountId         uint           `json:"account_id"`
	SpiderName        string         `json:"spider_name"`
	SpiderStatId	  			 uint           `json:"spider_stat_id"`
	JobName           string         `json:"job_name"`
	JobGroupId        uint           `json:"job_group_id"`
	ServerId          uint           `json:"server_id"`
	JobType		      string         `json:"job_type"`
	JobStartEpoc      int64          `json:"job_start_epoc"`
	JobStartTime      time.Time      `json:"job_start_time"`
	JobFinishTime     time.Time      `json:"job_finish_time"`
	JobLastSdkUpdate  time.Time      `json:"job_last_sdk_update"`
	JobRunTime        int            `json:"job_run_time"`
	JobStatus         string         `json:"job_status"`
	JobFinishReason   string         `json:"job_finish_reason"`
	JobStatsFreq      int            `json:"job_stats_freq"` // if none defined, should use the accounts setting
	ScrapeopsVersion  string         `json:"scrapeops_version"`
	ScrapyVersion     string         `json:"scrapy_version"`
	PythonVersion     string         `json:"python_version"`
	SystemVersion     string         `json:"system_version"`
	MiddlewareEnabled bool           `json:"middleware_enabled"`
	OverallError      int            `json:"overall_error"`
	OverallWarning    int            `json:"overall_warning"`
	OverallCritical   int            `json:"overall_critical"`
	UniqueError          int         `json:"unique_error"`
	UniqueWarning        int         `json:"unique_warning"`
	JobSettings       postgres.Jsonb `json:"job_settings"`
	JobArgs           postgres.Jsonb `json:"job_args"`

	JobRequests             int            `json:"job_requests"`
	JobResponses             int            `json:"job_responses"`
	Job_200Responses        int            `json:"job_200_responses"`
	// JobNon_300Responses     int            `json:"job_non_300_responses"`
	JobBytes                float64        `json:"job_bytes"`
	JobSuccessRate          float64        `json:"job_success_rate"`
	JobAvgLatency           float64        `json:"job_avg_latency"`
	JobTotalLatency           float64        `json:"job_total_latency"`
	JobStatusCodes          postgres.Jsonb `json:"job_status_codes"`
	JobItems                uint64         `json:"job_items"`
	JobChecks               int            `json:"job_checks"`
	JobPassedChecks         int            `json:"job_passed_checks"`
	JobFailedUrls           int            `json:"job_failed_urls"`
	JobFailedUrlsEnabled    bool           `json:"job_failed_urls_enabled"`
	JobDataCoverage         postgres.Jsonb `json:"job_data_coverage"`
	JobInvalidItems         int            `json:"job_invalid_items"`
	JobFieldCoverage        int            `json:"job_field_coverage"`
	JobExcludeMovingAverage bool           `json:"job_exclude_moving_average"`
	JobScrapyStats          postgres.Jsonb `json:"job_scrapy_stats"`
	JobCustomGroups         postgres.Jsonb `json:"job_custom_groups"`
	JobMode                 string         `json:"job_mode"`
}



type JobRedisVerison struct {
	ID                uint `gorm:"primaryKey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	AccountId         uint           `json:"account_id"`
	SpiderName        string         `json:"spider_name"`
	SpiderStatId	  			 uint           `json:"spider_stat_id"`
	JobName           string         `json:"job_name"`
	JobGroupId        uint           `json:"job_group_id"`
	ServerId          uint           `json:"server_id"`
	JobType		      string         `json:"job_type"`
	JobStartEpoc      int64          `json:"job_start_epoc"`
	JobStartTime      time.Time      `json:"job_start_time"`
	JobFinishTime     time.Time      `json:"job_finish_time"`
	JobLastSdkUpdate  time.Time      `json:"job_last_sdk_update"`
	JobRunTime        int            `json:"job_run_time"`
	JobStatus         string         `json:"job_status"`
	JobFinishReason   string         `json:"job_finish_reason"`
	JobStatsFreq      int            `json:"job_stats_freq"` // if none defined, should use the accounts setting
	ScrapeopsVersion  string         `json:"scrapeops_version"`
	ScrapyVersion     string         `json:"scrapy_version"`
	PythonVersion     string         `json:"python_version"`
	SystemVersion     string         `json:"system_version"`
	MiddlewareEnabled bool           `json:"middleware_enabled"`
	OverallError      int            `json:"overall_error"`
	OverallWarning    int            `json:"overall_warning"`
	OverallCritical   int            `json:"overall_critical"`
	UniqueError          int         `json:"unique_error"`
	UniqueWarning        int         `json:"unique_warning"`
	JobSettings       postgres.Jsonb `json:"job_settings"`
	JobArgs           postgres.Jsonb `json:"job_args"`

	JobRequests             int            `json:"job_requests"`
	JobResponses             int            `json:"job_responses"`
	Job_200Responses        int            `json:"job_200_responses"`
	// JobNon_300Responses     int            `json:"job_non_300_responses"`
	JobBytes                float64        `json:"job_bytes"`
	JobSuccessRate          float64        `json:"job_success_rate"`
	JobAvgLatency           float64        `json:"job_avg_latency"`
	JobTotalLatency           float64        `json:"job_total_latency"`
	JobStatusCodes          map[string]int `json:"job_status_codes"`
	JobItems                uint64         `json:"job_items"`
	JobChecks               int            `json:"job_checks"`
	JobPassedChecks         int            `json:"job_passed_checks"`
	JobFailedUrls           int            `json:"job_failed_urls"`
	JobFailedUrlsEnabled    bool           `json:"job_failed_urls_enabled"`
	JobDataCoverage         map[string]interface{} `json:"job_data_coverage"`
	JobInvalidItems         int            `json:"job_invalid_items"`
	JobFieldCoverage        int            `json:"job_field_coverage"`
	JobExcludeMovingAverage bool           `json:"job_exclude_moving_average"`
	JobScrapyStats          postgres.Jsonb `json:"job_scrapy_stats"`
	JobCustomGroups         map[string]string `json:"job_custom_groups"`
	JobMode                 string         `json:"job_mode"`
}



type JobGroup struct {
	ID                           uint `gorm:"primaryKey"`
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
	JobGroupMultiServerUUID      string         `json:"job_group_multi_server_uuid"`
	AccountId                    uint           `json:"account_id"`
	SpiderId                     uint           `json:"spider_id"`
	SpiderName                   string         `json:"spider_name"`
	SpiderStatId	  			 uint           `json:"spider_stat_id"`
	ProjectId                    int            `json:"project_id"`
	ScheduledJobId               uint           `json:"scheduled_job_id"`
	JobGroupName                 string         `json:"job_group_name"`
	JobGroupType                 string         `json:"job_group_type"`
	JobGroupDayStartTime         time.Time      `json:"job_group_day_start_time"`
	JobGroupStartTime            time.Time      `json:"job_group_start_time"`
	JobGroupFinishTime           time.Time      `json:"job_group_finish_time"`
	JobGroupLastSdkUpdate        time.Time      `json:"job_group_last_sdk_update"`
	JobGroupRunTime              int            `json:"job_group_run_time"`
	JobGroupStatus               string         `json:"job_group_status"`
	JobGroupFinishReason         string         `json:"job_group_finish_reason"`
	JobGroupStatsFreq            int            `json:"job_group_stats_freq"` // if none defined, should use the accounts setting
	JobGroupSettings             postgres.Jsonb `json:"job_group_settings"`
	JobGroupArgs                 postgres.Jsonb `json:"job_group_args"`
	JobGroupOverallGroupError    int            `json:"job_group_overall_group_error"`
	JobGroupUniqueError          int            `json:"job_group_unique_error"`
	JobGroupOverallGroupWarning  int            `json:"job_group_overall_group_warning"`
	JobGroupUniqueWarning        int            `json:"job_group_unique_warning"`
	JobGroupOverallGroupCritical int            `json:"job_group_overall_group_critical"`
	JobGroupMultiServer          bool           `json:"job_group_multi_server"`
	JobGroupRequests             int            `json:"job_group_requests"`
	JobGroupResponses             int            `json:"job_group_responses"`
	JobGroup_200Responses        int            `json:"job_group_200_responses"`
	JobGroupNon_300Responses     int            `json:"job_group_non_300_responses"`
	JobGroupBytes                float64        `json:"job_group_bytes"`
	JobGroupSuccessRate          float64        `json:"job_group_success_rate"`
	JobGroupAvgLatency           float64        `json:"job_group_avg_latency"`
	JobGroupTotalLatency           float64        `json:"job_group_total_latency"`
	JobGroupStatusCodes          postgres.Jsonb `json:"job_group_status_codes"`
	JobGroupItems                uint64         `json:"job_group_items"`
	JobGroupChecks               int            `json:"job_group_checks"`
	JobGroupPassedChecks         int            `json:"job_group_passed_checks"`
	JobGroupFailedUrls           int            `json:"job_group_failed_urls"`
	JobGroupFailedUrlsEnabled    bool           `json:"job_group_failed_urls_enabled"`
	JobGroupDataCoverage         postgres.Jsonb `json:"job_group_data_coverage"`
	JobGroupInvalidItems         int            `json:"job_group_invalid_items"`
	JobGroupFieldCoverage        int            `json:"job_group_field_coverage"`
	JobGroupExcludeMovingAverage bool           `json:"job_group_exclude_moving_average"`
	JobGroupScrapyStats          postgres.Jsonb `json:"job_group_scrapy_stats"`
	JobGroupCustomGroups         postgres.Jsonb `json:"job_group_custom_groups"`

	JobGroupMode                 string         `json:"job_group_mode"`
	JobGroupNumJobs				 int            `json:"job_group_num_jobs"`
	JobGroupFinishedJobs		 int            `json:"job_group_finished_jobs"`
	JobGroupRunningJobs			 int            `json:"job_group_running_jobs"`
	JobGroupUnknownJobs			 int            `json:"job_group_unknown_jobs"`
	JobGroupShutdownJobs		 int            `json:"job_group_shutdown_jobs"`
}



type JobGroupRedisVersion struct {
	ID                           uint `gorm:"primaryKey"`
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
	JobGroupMultiServerUUID      string         `json:"job_group_multi_server_uuid"`
	AccountId                    uint           `json:"account_id"`
	SpiderId                     uint           `json:"spider_id"`
	SpiderName                   string         `json:"spider_name"`
	SpiderStatId	  			 uint           `json:"spider_stat_id"`
	ProjectId                    int            `json:"project_id"`
	ScheduledJobId               uint           `json:"scheduled_job_id"`
	JobGroupName                 string         `json:"job_group_name"`
	JobGroupType                 string         `json:"job_group_type"`
	JobGroupDayStartTime         time.Time      `json:"job_group_day_start_time"`
	JobGroupStartTime            time.Time      `json:"job_group_start_time"`
	JobGroupFinishTime           time.Time      `json:"job_group_finish_time"`
	JobGroupLastSdkUpdate        time.Time      `json:"job_group_last_sdk_update"`
	JobGroupRunTime              int            `json:"job_group_run_time"`
	JobGroupStatus               string         `json:"job_group_status"`
	JobGroupFinishReason         string         `json:"job_group_finish_reason"`
	JobGroupStatsFreq            int            `json:"job_group_stats_freq"` // if none defined, should use the accounts setting
	JobGroupSettings             postgres.Jsonb `json:"job_group_settings"`
	JobGroupArgs                 postgres.Jsonb `json:"job_group_args"`
	JobGroupOverallGroupError    int            `json:"job_group_overall_group_error"`
	JobGroupUniqueError          int            `json:"job_group_unique_error"`
	JobGroupOverallGroupWarning  int            `json:"job_group_overall_group_warning"`
	JobGroupUniqueWarning        int            `json:"job_group_unique_warning"`
	JobGroupOverallGroupCritical int            `json:"job_group_overall_group_critical"`
	JobGroupMultiServer          bool           `json:"job_group_multi_server"`
	JobGroupRequests             int            `json:"job_group_requests"`
	JobGroupResponses             int            `json:"job_group_responses"`
	JobGroup_200Responses        int            `json:"job_group_200_responses"`
	JobGroupNon_300Responses     int            `json:"job_group_non_300_responses"`
	JobGroupBytes                float64        `json:"job_group_bytes"`
	JobGroupSuccessRate          float64        `json:"job_group_success_rate"`
	JobGroupAvgLatency           float64        `json:"job_group_avg_latency"`
	JobGroupTotalLatency         float64        `json:"job_group_total_latency"`
	JobGroupStatusCodes          map[string]int `json:"job_group_status_codes"`
	JobGroupItems                uint64         `json:"job_group_items"`
	JobGroupChecks               int            `json:"job_group_checks"`
	JobGroupPassedChecks         int            `json:"job_group_passed_checks"`
	JobGroupFailedUrls           int            `json:"job_group_failed_urls"`
	JobGroupFailedUrlsEnabled    bool           `json:"job_group_failed_urls_enabled"`
	JobGroupDataCoverage         map[string]interface{} `json:"job_group_data_coverage"`
	JobGroupInvalidItems         int            `json:"job_group_invalid_items"`
	JobGroupFieldCoverage        int            `json:"job_group_field_coverage"`
	JobGroupExcludeMovingAverage bool           `json:"job_group_exclude_moving_average"`
	JobGroupScrapyStats          postgres.Jsonb `json:"job_group_scrapy_stats"`
	JobGroupCustomGroups         postgres.Jsonb `json:"job_group_custom_groups"`

	JobGroupMode                 string         `json:"job_group_mode"`
	JobGroupNumJobs				 int            `json:"job_group_num_jobs"`
	JobGroupFinishedJobs		 int            `json:"job_group_finished_jobs"`
	JobGroupRunningJobs			 int            `json:"job_group_running_jobs"`
	JobGroupUnknownJobs			 int            `json:"job_group_unknown_jobs"`
	JobGroupShutdownJobs		 int            `json:"job_group_shutdown_jobs"`
}




type JobGroupAvg struct {
	ID                           uint `gorm:"primaryKey"`
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
	
	AccountId                    uint           `json:"account_id"`
	JobGroupName                 string         `json:"job_group_name"`
	SpiderName                   string         `json:"spider_name"`

	JobGroupAvgWindowType           string            `json:"job_group_avg_window_type"`
	JobGroupAvgWindowLength          int            `json:"job_group_avg_window_length"`

	JobGroupAvgRunTime              int            `json:"job_group_avg_run_time"`
	JobGroupAvgOverallGroupError    int            `json:"job_group_avg_overall_group_error"`
	JobGroupAvgOverallGroupWarning  int            `json:"job_group_avg_overall_group_warning"`
	JobGroupAvgOverallGroupCritical int            `json:"job_group_avg_overall_group_critical"`
	JobGroupAvgRequests             int            `json:"job_group_avg_requests"`
	JobGroupAvg_200Responses        int            `json:"job_group_avg_200_responses"`
	JobGroupAvgNon_300Responses     int            `json:"job_group_avg_non_300_responses"`
	JobGroupAvgBytes                float64        `json:"job_group_avg_bytes"`
	JobGroupAvgSuccessRate          float64        `json:"job_group_avg_success_rate"`
	JobGroupAvgAvgLatency           float64        `json:"job_group_avg_avg_latency"`
	JobGroupAvgStatusCodes          postgres.Jsonb `json:"job_group_avg_status_codes"`
	JobGroupAvgItems                uint64            `json:"job_group_avg_items"`
	JobGroupAvgChecks               int            `json:"job_group_avg_checks"`
	JobGroupAvgPassedChecks         int            `json:"job_group_avg_passed_checks"`
	JobGroupAvgFailedUrls           int            `json:"job_group_avg_failed_urls"`
	JobGroupAvgInvalidItems         int          `json:"job_group_avg_invalid_items"`
	JobGroupAvgFieldCoverage		 int	   		`json:"job_group_avg_field_coverage"`
}


type CustomGroup struct {
	ID                           uint `gorm:"primaryKey"`
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
	AccountId                    uint           `json:"account_id"`
	CustomGroupKey          	string         `json:"custom_group_key"`
	CustomGroupNumValues          uint         `json:"custom_group_num_values"`
	
}


type CustomGroupStat struct {
	ID                           uint `gorm:"primaryKey"`
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
	AccountId                    uint           `json:"account_id"`
	CustomGroupStatType                   string         `json:"custom_group_stat_type"`
	CustomGroupStatName                   string         `json:"custom_group_stat_name"`
	CustomGroupStatStartTime            time.Time      `json:"custom_group_stat_start_time"`
	CustomGroupStatLastSdkUpdate        time.Time      `json:"custom_group_stat_last_sdk_update"`
	CustomGroupStatRunTime              int            `json:"custom_group_stat_run_time"`
	CustomGroupStatRunning               int         	`json:"custom_group_stat_running"`
	CustomGroupStatFinished         int        	 `json:"custom_group_stat_finished"`
	CustomGroupStatOverallGroupError    int            `json:"custom_group_stat_overall_group_error"`
	CustomGroupStatOverallGroupWarning  int            `json:"custom_group_stat_overall_group_warning"`
	CustomGroupStatOverallGroupCritical int            `json:"custom_group_stat_overall_group_critical"`
	CustomGroupStatRequests             int            `json:"custom_group_stat_requests"`
	CustomGroupStat_200Responses        int            `json:"custom_group_stat_200_responses"`
	CustomGroupStatNon_300Responses     int            `json:"custom_group_stat_non_300_responses"`
	CustomGroupStatBytes                float64        `json:"custom_group_stat_bytes"`
	CustomGroupStatSuccessRate          float64        `json:"custom_group_stat_success_rate"`
	CustomGroupStatAvgLatency           float64        `json:"custom_group_stat_avg_latency"`
	CustomGroupStatStatusCodes          postgres.Jsonb `json:"custom_group_stat_status_codes"`
	CustomGroupStatItems                uint64            `json:"custom_group_stat_items"`
	CustomGroupStatChecks               int            `json:"custom_group_stat_checks"`
	CustomGroupStatPassedChecks         int            `json:"custom_group_stat_passed_checks"`
	CustomGroupStatFailedUrls           int            `json:"custom_group_stat_failed_urls"`
	CustomGroupStatDataCoverage         postgres.Jsonb `json:"custom_group_stat_data_coverage"`
	CustomGroupStatInvalidItems         int          `json:"custom_group_stat_invalid_items"`
	CustomGroupStatFieldCoverage		 	int	   		`json:"custom_group_stat_field_coverage"`
	CustomGroupStatNumJobs				int	   		`json:"custom_group_stat_num_jobs"`
	CustomGroupStatUniqueJobs			int	   		`json:"custom_group_stat_unique_jobs"`
	CustomGroupStatExcludeMovingAverage	bool	   		`json:"custom_group_stat_exclude_moving_average"`
	CustomGroupStatItemsArray			string	   	`json:"custom_group_stat_items_array"`	
	CustomGroupStatAggHour  				int            `json:"custom_group_stat_agg_hour"`

}





