package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

type Monitor struct {
}

func (m *Monitor) loadMiners() {
	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	addr, err := proto.NewAddressFromString(MobileAddress)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	entries, _, err := cl.Addresses.AddressesData(ctx, addr)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	for _, m := range entries {
		miner := &Miner{}
		db.FirstOrCreate(miner, &Miner{Address: m.GetKey()})

		minerData := m.ToProtobuf().GetStringValue()

		tel := parseItem(minerData, 0)
		mh := parseItem(minerData, 1)
		encIp := parseItem(minerData, 2)
		ref := parseItem(minerData, 3)

		if encIp != nil {
			ip := DecryptMessage(encIp.(string))
			if len(ip) > 0 {
				miner.IP = ip
			}
		}

		telId := DecryptMessage(tel.(string))
		telIdInt, err := strconv.Atoi(telId)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}
		miner.TelegramId = int64(telIdInt)
		if mh != nil {
			miner.MiningHeight = int64(mh.(int))
		}
		if ref != nil {
			refdb := &Miner{}
			db.First(refdb, &Miner{Address: ref.(string)})
			if refdb.ID != 0 {
				miner.ReferralID = refdb.ID
			}
		}

		if !miner.Confirmed {
			cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
			if err != nil {
				log.Println(err)
				logTelegram(err.Error())
			}

			c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			a := proto.MustAddressFromString(miner.Address)

			balance, _, err := cl.Addresses.Balance(c, a)
			if err != nil {
				log.Println(err)
				logTelegram(err.Error())
			}

			if balance.Balance >= Fee {
				miner.Confirmed = true
				miner.Balance = balance.Balance
			}
		}

		db.Save(miner)
	}

	var dbminers []*Miner

	db.Find(&dbminers)
	for _, dbm := range dbminers {
		found := false

		for _, m := range entries {
			if m.GetKey() == dbm.Address {
				found = true
			}
		}

		if !found {
			db.Delete(&dbm)
		}
	}
}

func (m *Monitor) start() {
	for {
		m.loadMiners()

		log.Println("Done update.")

		time.Sleep(time.Second * MonitorTick)
	}
}

func initMonitor() {
	m := &Monitor{}
	go m.start()
}
