package main

import (
	"time"

	"github.com/jinzhu/copier"
	"gopkg.in/macaron.v1"
)

func minerView(ctx *macaron.Context) {
	addr := ctx.Params("addr")
	var referred []*Miner
	var miners []*Miner
	var count int

	m := &Miner{}
	db.First(m, &Miner{Address: addr})

	mr := &MinerResponse{}
	copier.Copy(mr, m)

	db.Find(&referred, &Miner{ReferralID: m.ID})
	mr.ReferredCount = len(referred)

	height := getHeight()
	db.Where("mining_height > ?", height-2880).Find(&miners)
	count = len(miners)

	mr.ActiveMiners = count

	for _, m := range miners {
		db.Find(&referred, &Miner{ReferralID: m.ID})
		count += len(referred)
	}

	mr.MinRefCount = count

	ctx.JSON(200, mr)
}

type MinerResponse struct {
	Address          string
	LastNotification time.Time
	TelegramId       int64
	MiningHeight     int64
	ReferredCount    int
	MinRefCount      int
	ActiveMiners     int
}
