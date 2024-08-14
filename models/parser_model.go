package models

import (
	"time"
)

type AccountParser struct {
	ID                              uint      `gorm:"primaryKey"`
	AccountId                       uint      `json:"account_id"`
	ProxyPlanId                     uint      `json:"proxy_plan_id"`
	ProxyActivated                  bool      `json:"proxy_activated"`
	AccountParserTotalRequests      uint      `json:"account_parser_total_requests"`
	AccountParserSuccessfulRequests uint      `json:"account_parser_successful_requests"`
	AccountParserUsedCredits        int64     `json:"account_parser_used_credits"`
	AccountParserPlanRenewalDate    time.Time `json:"account_parser_renewal_date"`
	CreatedAt                       time.Time
	UpdatedAt                       time.Time
}

type AccountParserStat struct {
	ID                            uint      `gorm:"primaryKey"`
	AccountId                     uint      `json:"account_id"`
	AccountParserStatDayStartTime time.Time `json:"account_parser_stat_day_start_time"`
	AccountParserStatRequests     uint      `json:"account_parser_stat_requests"`
	AccountParserStatSuccessful   uint      `json:"account_parser_stat_successful"`
	AccountParserStatFailed       uint      `json:"account_parser_stat_failed"`
	AccountParserStatCredits      int64     `json:"account_parser_stat_credits"`
	AccountParserStatDataCoverage float64   `json:"account_parser_stat_data_coverage"`
	CreatedAt                     time.Time
	UpdatedAt                     time.Time
}

type AccountParserDomainStat struct {
	ID                                  uint      `gorm:"primaryKey"`
	AccountId                           uint      `json:"account_id"`
	AccountParserStatId                 uint      `json:"account_parser_stat_id"`
	AccountParserDomainStatDayStartTime time.Time `json:"account_parser_domain_stat_day_start_time"`
	AccountParserDomainStatDomain       string    `json:"account_parser_domain_stat_domain"`
	AccountParserDomainStatRequests     uint      `json:"account_parser_domain_stat_requests"`
	AccountParserDomainStatSuccessful   uint      `json:"account_parser_domain_stat_successful"`
	AccountParserDomainStatFailed       uint      `json:"account_parser_domain_stat_failed"`
	AccountParserDomainStatCredits      int64     `json:"account_parser_stat_credits"`
	AccountParserDomainStatDataCoverage float64   `json:"account_parser_stat_data_coverage"`
	CreatedAt                           time.Time
	UpdatedAt                           time.Time
}
