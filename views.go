package main

import (
	"time"

	"github.com/jinzhu/copier"
	"gopkg.in/macaron.v1"
)

func minersView(ctx *macaron.Context) {
	var miners []*Miner
	var mnrs []*MinersResponse

	db.Find(&miners)
	for _, m := range miners {
		mr := &MinersResponse{}
		copier.Copy(mr, m)
		mnrs = append(mnrs, mr)
	}

	ctx.JSON(200, mnrs)
}

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
	countRef := 0

	mr.ActiveMiners = count

	for _, m := range miners {
		db.Find(&referred, &Miner{ReferralID: m.ID})
		countRef += len(referred)
	}

	mr.MinRefCount = count + (countRef / 4)

	ctx.JSON(200, mr)
}

func ipView(ctx *macaron.Context) {
	var miners []*Miner

	ipcr := &IPCountResponse{}
	ip := ctx.Params("ip")
	height := getHeight()

	db.Where("ip = ? AND mining_height > ?", ip, height-2880).Find(&miners)
	ipcr.Count = len(miners)

	ctx.JSON(200, ipcr)
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

type MinersResponse struct {
	Address          string
	LastNotification time.Time
	TelegramId       int64
	MiningHeight     int64
}

type IPCountResponse struct {
	Count int `json:"count"`
}
