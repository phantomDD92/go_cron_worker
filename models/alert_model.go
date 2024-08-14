package models

import (
  "github.com/jinzhu/gorm"
  "time"
)

type AlertThreshold struct {
	gorm.Model
	Parameter 						string `json:"alert_param"`
	TriggerType 					string `json:"alert_trigger_type"`
	TriggerThreshold 				int	`json:"alert_threshold"`
	TriggerThresholdUnits 			string `json:"alert_threshold_units"`
	TriggerDataType					string `json:"trigger_data_type"`
	AlertLevel 						string `json:"alert_level"`
	AlertCommunicationMethodId		uint `json:"alert_communication_method_id"`
}

type AlertHistory struct {
	gorm.Model
	JobGroupId						uint `json:"job_group_id"`
	AccountId 					uint `json:"account_id"`
	PeriodicStatId				uint `json:"period_id"`
	Parameter 					string `json:"alert_param"`
	TriggerType 				string  `json:"alert_trigger_type"`
	TriggerThreshold 			float32	`json:"alert_threshold"`
	StatsValue		 			float32	`json:"stats_value"`
	AlertLevel 					string `json:"alert_level"`
	AlertMethod					string `json:"alert_method"`
	AlertStatus					string `json:"alert_status"`
	AlertCommunicationMethodId 	uint `json:"alert_communication_method_id"`
	AlertType					string `json:"alert_type"`
	AlertThresholdId			uint `json:"alert_threshold_id"`
}

type HealthCheckHistory struct {
	ID                    uint `gorm:"primaryKey"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
	JobGroupId            uint    `json:"job_group_id"`
	SpiderStatId		  uint    `json:"spider_stat_id"`
	AccountId             uint    `json:"account_id"`
	Parameter             string  `json:"alert_param"`
	TriggerType           string  `json:"alert_trigger_type"`
	TriggerThreshold      float32 `json:"alert_threshold"`
	TriggerThresholdUnits string  `json:"alert_threshold_units"`
	TriggerDataType       string  `json:"trigger_data_type"`
	StatsValue            float32 `json:"stats_value"`
	CheckPassed           bool    `json:"check_passed"`
	HealthCheckType       string  `json:"health_check_type"`
	JobOrSpider       	  string  `json:"job_or_spider"`
}


type AlertThresholdsAccount  struct {
	gorm.Model
	AccountId			uint `json:"account_id"`
	AlertThresholdId	uint `json:"alertthreshold_id"`
}

type AlertSpiderThresholds  struct {
	gorm.Model
	SpiderId			uint `json:"spider_id"`
	AlertThresholdId	uint `json:"alert_threshold_id"`
	AlertActive			bool `json:"alert_active"`				
}

type AlertJobThresholds struct {
	gorm.Model
	JobGroupId			uint	`json:"job_group_id"`
	AlertThresholdId	uint 	`json:"alert_threshold_id"`
	AlertSilenced		bool	`json:"alert_silenced"`
}

type AlertCommunicationMethod struct {
	gorm.Model
	AccountId		uint `json:"account_id"`
	CommMethod		string `json:"communication_method"`
}

type ChatTool  struct {
	gorm.Model
	AlertCommunicationMethodId		uint `json:"alert_communication_method_id"`
	ChatToolType					string `json:"chat_tool_type"`
	ChatToolToken					string `json:"chat_tool_token"`
	ChatToolChannel					string `json:"chat_tool_channel"`		
	ChatToolNotificationString		string `json:"chat_tool_notification_string"`						
}

type AlertEmails  struct {
	gorm.Model
	AlertCommunicationMethodId		uint `json:"alert_communication_method_id"`
	Email							string `json:"email"`
	EmailType						string `json:"email_type"`					
}


type AlertTexts  struct {
	gorm.Model
	AlertCommunicationMethodId		uint `json:"alert_communication_method_id"`
	Number							string `json:"number"`				
}







