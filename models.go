package main

import (
	"time"

	"gorm.io/gorm"
)

type Miner struct {
	gorm.Model
	Address          string `gorm:"size:255;uniqueIndex"`
	LastNotification time.Time
	TelegramId       int64
	MiningHeight     int64
	ReferralID       uint   `gorm:"index"`
	IP               string `gorm:"index;default:127.0.0.1"`
	Confirmed        bool   `gorm:"default:false"`
}
