package models

import (
	"time"
)

type PeriodicJobReport struct {
	ID        				uint           `gorm:"primaryKey"`
  	CreatedAt 				time.Time
  	UpdatedAt 				time.Time
	AccountId                  uint   `json:"account_id"`
	CronString                 string `json:"cron_string"`
	ReportLength               uint   `json:"report_length"`
	AlertCommunicationMethodId uint   `json:"alert_communication_method_id"`
}
