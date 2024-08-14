package models

import (
  "time"
)


type Account struct {
	ID        				uint           `gorm:"primaryKey"`
  	CreatedAt 				time.Time
  	UpdatedAt 				time.Time
	Name					string `json:"name"`
	CreatedUserId			uint `json:"created_user_id"`
	PlanId					uint `json:"plan_id"`
	APIKey 					string `json:"api_key"`
	JobStartNotification 	int `json:"job_start_notification"`
	JobFinishNotification 	int `json:"job_finish_notification"`
	JobPeriodicReport		uint `json:"job_periodic_report"`
	Timezone 				string `json:"timezone"`
	DayEndUtcOffset			int `json:"day_end_utc_offset"`
	AggregationFrequency	uint `json:"aggregation_frequency"`
	NumJobsRun 				uint `json:"num_jobs_run"`
	NumJobsRunLast_24h		uint `json:"num_jobs_run_last_24h"`
	
}