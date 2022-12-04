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

	m := &Miner{}
	db.First(m, &Miner{Address: addr})

	mr := &MinerResponse{}
	copier.Copy(mr, m)

	height := getHeight()

	if m.ID != 0 {
		db.Where("mining_height > ? AND referral_id = ? AND confirmed = 1", height-1440, m.ID).Find(&referred)
		mr.ReferredCount = len(referred)
		mr.Exists = true
	} else {
		mr.Exists = false
	}

	ctx.JSON(200, mr)
}

func pingView(ctx *macaron.Context) {
	addr := ctx.Params("addr")

	m := &Miner{}
	db.First(m, &Miner{Address: addr})

	pr := &PingResponse{}

	m.PingCount++
	m.LastPing = time.Now()
	db.Save(m)

	pr.LastPing = m.LastPing
	pr.PingCount = m.PingCount

	ctx.JSON(200, pr)
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
		m.Balance = balance.Balance
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
			if m.ReferralID != 0 && m.Confirmed {
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
	Confirmed        bool      `json:"confirmed"`
	Exists           bool      `json:"exists"`
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

type PingResponse struct {
	PingCount int64     `json:"ping_count"`
	LastPing  time.Time `json:"last_ping"`
}
