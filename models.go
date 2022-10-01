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
}

type KeyValue struct {
	gorm.Model
	Key      string `gorm:"size:255;uniqueIndex"`
	ValueInt uint64 `gorm:"type:int"`
	ValueStr string `gorm:"type:string"`
}