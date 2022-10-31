package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/copier"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/proto"
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

	height := getHeight()

	// db.Find(&referred, &Miner{ReferralID: m.ID})
	db.Where("mining_height > ? AND referral_id = ?", height-1440, m.ID).Find(&referred)
	mr.ReferredCount = len(referred)

	db.Where("mining_height > ?", height-2880).Find(&miners)
	count = len(miners)
	countRef := 0

	mr.ActiveMiners = count

	for _, m := range miners {
		if m.ReferralID != 0 {
			countRef++
		}
	}

	mr.ActiveReferred = countRef
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

func checkConfirmationView(ctx *macaron.Context) {
	addr := ctx.Params("addr")
	m := &Miner{}
	db.First(m, &Miner{Address: addr})
	log.Println(prettyPrint(m))

	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a := proto.MustAddressFromString(addr)

	balance, _, err := cl.Addresses.Balance(c, a)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if balance.Balance >= Fee {
		m.Confirmed = true
		db.Save(m)
	}

	gr := &GeneralResponse{Success: true}
	ctx.JSON(200, gr)
}

func statsView(ctx *macaron.Context) {
	var miners []*Miner
	sr := &StatsResponse{}
	db.Find(&miners)
	height := getHeight()

	for _, m := range miners {
		if height-uint64(m.MiningHeight) <= 1440 {
			sr.ActiveMiners++
			if m.ReferralID != 0 {
				sr.ActiveReferred++
			}
		}

		if height-uint64(m.MiningHeight) <= 2880 {
			sr.PayoutMiners++
		}
	}

	sr.InactiveMiners = len(miners) - sr.PayoutMiners

	ctx.JSON(200, sr)
}

type MinerResponse struct {
	Address          string    `json:"address"`
	LastNotification time.Time `json:"last_notification"`
	TelegramId       int64     `json:"telegram_id"`
	MiningHeight     int64     `json:"mining_height"`
	ReferredCount    int       `json:"referred_count"`
	MinRefCount      int       `json:"min_ref_count"`
	ActiveMiners     int       `json:"active_miners"`
	ActiveReferred   int       `json:"active_referred"`
	Confirmed        bool      `json:"confirmed"`
}

type MinersResponse struct {
	Address          string    `json:"address"`
	LastNotification time.Time `json:"last_notification"`
	TelegramId       int64     `json:"telegram_id"`
	MiningHeight     int64     `json:"mining_height"`
}

type IPCountResponse struct {
	Count int `json:"count"`
}

type GeneralResponse struct {
	Success bool `json:"success"`
}

type StatsResponse struct {
	ActiveMiners   int `json:"active_miners"`
	ActiveReferred int `json:"active_referred"`
	PayoutMiners   int `json:"payout_miners"`
	InactiveMiners int `json:"inactive_miners"`
}
