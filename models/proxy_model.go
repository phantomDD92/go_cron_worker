package models

import (
	"encoding/json"
	"time"
)

type AccountProxy struct {
	ID                             uint      `gorm:"primaryKey"`
	AccountId                      uint      `json:"account_id"`
	ProxyPlanId                    uint      `json:"proxy_plan_id"`
	ProxyActivated                 bool      `json:"proxy_activated"`
	AccountProxyTotalRequests      uint      `json:"account_proxy_total_requests"`
	AccountProxySuccessfulRequests uint      `json:"account_proxy_successful_requests"`
	AccountProxyUsedCredits        int64     `json:"account_proxy_used_credits"`
	AccountProxyPlanRenewalDate    time.Time `json:"account_proxy_renewal_date"`
	CreatedAt                      time.Time
	UpdatedAt                      time.Time
}

type AccountProxyStat struct {
	ID                           uint      `gorm:"primaryKey"`
	AccountId                    uint      `json:"account_id"`
	AccountProxyStatDayStartTime time.Time `json:"account_proxy_stat_day_start_time"`
	AccountProxyStatRequests     uint      `json:"account_proxy_stat_requests"`
	AccountProxyStatSuccessful   uint      `json:"account_proxy_stat_successful"`
	AccountProxyStatFailed       uint      `json:"account_proxy_stat_failed"`
	AccountProxyStatCredits      int64     `json:"account_proxy_stat_credits"`
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
}

type AccountProxyDomainStat struct {
	ID                                 uint      `gorm:"primaryKey"`
	AccountId                          uint      `json:"account_id"`
	AccountProxyStatId                 uint      `json:"account_proxy_stat_id"`
	AccountProxyDomainStatDayStartTime time.Time `json:"account_proxy_domain_stat_day_start_time"`
	AccountProxyDomainStatDomain       string    `json:"account_proxy_domain_stat_domain"`
	AccountProxyDomainStatRequests     uint      `json:"account_proxy_domain_stat_requests"`
	AccountProxyDomainStatSuccessful   uint      `json:"account_proxy_domain_stat_successful"`
	AccountProxyDomainStatFailed       uint      `json:"account_proxy_domain_stat_failed"`
	AccountProxyDomainStatCredits      int64     `json:"account_proxy_stat_credits"`
	CreatedAt                          time.Time
	UpdatedAt                          time.Time
}

type ProxyFailedResponse struct {
	TimeWindow string `json:"time_window"`
	Method     string `json:"method"`

	Proxy           string `json:"proxy"`
	ProxyUrl        string `json:"proxy_url"`
	ProxyStatusCode int    `json:"proxy_status_code"`
	ProxyApiCredits int64  `json:"proxy_api_credits"`
	ProxyUniqueId   string `json:"proxy_unique_id"`

	BodyString  string `json:"body_string"`
	ContentType string `json:"content_type"`

	Url    string `json:"url"`
	Domain string `json:"domain"`

	FailedReason string `json:"failed_reason"`
	BlockType    string `json:"block_type"`
	Block        string `json:"block"`
}

type SopsProxyProvider struct {
	ID                                uint      `gorm:"primaryKey"`
	SopsProxyName                     string    `json:"sops_proxy_name"`
	SopsProxyProviderType             string    `json:"sops_proxy_provider_type"`
	SopsProxyProviderRenewalDate      time.Time `json:"sops_proxy_provider_renewal_date"`
	SopsProxyProviderApiCreditLimit   uint      `json:"sops_proxy_provider_api_credit_limit"`
	SopsProxyProviderConcurrencyLimit uint      `json:"sops_proxy_provider_concurrency_limit"`
	SopsProxyProviderRequests         uint      `json:"sops_proxy_provider_requests"`
	SopsProxyProviderSuccessful       uint      `json:"sops_proxy_provider_successful"`
	SopsProxyProviderFailed           uint      `json:"sops_proxy_provider_failed"`
	SopsProxyProviderFailedValidation uint      `json:"sops_proxy_provider_failed_validation"`
	SopsProxyProviderCredits          uint      `json:"sops_proxy_provider_credits"`
	CreatedAt                         time.Time
	UpdatedAt                         time.Time
}

type SopsDayProxyStat struct {
	ID                                      uint      `gorm:"primaryKey"`
	SopsProxyProviderId                     uint      `json:"sops_proxy_provider_id"`
	SopsProxyName                           string    `json:"sops_proxy_name"`
	SopsDayProxyStatDayStartTime            time.Time `json:"sops_day_proxy_stat_day_start_time"`
	SopsDayProxyStatRequests                uint      `json:"sops_day_proxy_stat_requests"`
	SopsDayProxyStatSuccessful              uint      `json:"sops_day_proxy_stat_successful"`
	SopsDayProxyStatFailed                  uint      `json:"sops_day_proxy_stat_failed"`
	SopsDayProxyStatFailedValidation        uint      `json:"sops_day_proxy_stat_failed_validation"`
	SopsDayProxyStatLatency                 float64   `json:"sops_day_proxy_stat_latency"`
	SopsDayProxyStatCredits                 uint      `json:"sops_day_proxy_stat_credits"`
	SopsDayProxyStatCreditsFailedValidation uint      `json:"sops_day_proxy_stat_credits_failed_validation"`
	CreatedAt                               time.Time
	UpdatedAt                               time.Time
}

type SopsDayProxyDomainStat struct {
	ID                                            uint      `gorm:"primaryKey"`
	SopsProxyProviderId                           uint      `json:"sops_proxy_provider_id"`
	SopsProxyName                                 string    `json:"sops_proxy_name"`
	SopsDayProxyStatId                            uint      `json:"sops_day_proxy_stat_id"`
	SopsDayProxyDomainStatDayStartTime            time.Time `json:"sops_day_proxy_domain_stat_day_start_time"`
	SopsDayProxyDomainStatDomain                  string    `json:"sops_day_proxy_domain_stat_domain"`
	SopsDayProxyDomainStatRequests                uint      `json:"sops_day_proxy_domain_stat_requests"`
	SopsDayProxyDomainStatSuccessful              uint      `json:"sops_day_proxy_domain_stat_successful"`
	SopsDayProxyDomainStatFailed                  uint      `json:"sops_day_proxy_domain_stat_failed"`
	SopsDayProxyDomainStatFailedValidation        uint      `json:"sops_day_proxy_domain_stat_failed_validation"`
	SopsDayProxyDomainStatLatency                 float64   `json:"sops_day_proxy_domain_stat_latency"`
	SopsDayProxyDomainStatCredits                 uint      `json:"sops_day_proxy_stat_credits"`
	SopsDayProxyDomainStatCreditsFailedValidation uint      `json:"sops_day_proxy_domain_stat_credits_failed_validation"`
	CreatedAt                                     time.Time
	UpdatedAt                                     time.Time
}

type AccountPPGBOverallStats struct {
	ID                 uint      `gorm:"primaryKey"`
	AccountID          uint      `json:"account_id"`
	Date               time.Time `json:"date"`
	SuccessfulRequests uint      `json:"successful_requests"`
	FailedRequests     uint      `json:"failed_requests"`
	TotalRequests      uint      `json:"total_requests"`
	BytesUsed          uint      `json:"bytes_used"`
	CreditsUsed        uint      `json:"credits_used"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type AccountPPGBDailyStats struct {
	ID                 uint      `gorm:"primaryKey"`
	AccountID          uint      `json:"account_id"`
	Date               time.Time `json:"date"`
	SuccessfulRequests uint      `json:"successful_requests"`
	FailedRequests     uint      `json:"failed_requests"`
	TotalRequests      uint      `json:"total_requests"`
	BytesUsed          uint      `json:"bytes_used"`
	CreditsUsed        uint      `json:"credits_used"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type AccountPPGBDailyDomainStats struct {
	ID                     uint      `gorm:"primaryKey"`
	AccountID              uint      `json:"account_id"`
	Date                   time.Time `json:"date"`
	SuccessfulRequests     uint      `json:"successful_requests"`
	FailedRequests         uint      `json:"failed_requests"`
	TotalRequests          uint      `json:"total_requests"`
	BytesUsed              uint      `json:"bytes_used"`
	CreditsUsed            uint      `json:"credits_used"`
	Domain                 string    `json:"domain"`
	AccountPPGBDailyStatID uint      `json:"account_ppgb_daily_stat_id"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type SopsProxyTestResult struct {
	ID                    uint            `json:"id"`
	URL                   string          `json:"url"`
	Domain                string          `json:"domain"`
	Method                string          `json:"method"`
	QueryParams           json.RawMessage `json:"query_params"`
	PostBody              json.RawMessage `json:"post_body"`
	CustomHeadersBool     bool            `json:"custom_headers_bool"`
	CustomHeaders         json.RawMessage `json:"custom_headers"`
	NumRequests           int             `json:"num_requests"`
	NumThreads            int             `json:"num_threads"`
	TestRunBy             string          `json:"test_run_by"`
	TestMode              string          `json:"test_mode"`
	TestString            string          `json:"test_string"`
	TestStatus            string          `json:"test_status"`
	TestResults           json.RawMessage `json:"test_results"`
	TestName              string          `json:"test_name"`
	TestNotes             string          `json:"test_notes"`
	TestType              string          `json:"test_type"`
	TestSuggestedSequence string          `json:"test_suggested_sequence"`
	FinalSequence         string          `json:"final_sequence"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}
