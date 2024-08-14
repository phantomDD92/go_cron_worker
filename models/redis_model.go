package models

import (
  "time"
)


type RedisStat struct {
	ID         uint   `gorm:"primaryKey"`
	RedisStatType      string    `json:"redis_stat_type"` // cache, queue, etc.
	RedisStatPeriod    string    `json:"redis_stat_period"`
	RedisStatName      string    `json:"redis_stat_name"`
	RedisStatValue     int	     `json:"redis_stat_value"`
	RedisStatTimestamp time.Time	`json:"redis_stat_timestamp"`

	CreatedAt time.Time
	UpdatedAt time.Time
}