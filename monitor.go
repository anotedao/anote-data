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

func (m *Monitor) loadMiners(minerType string) {
	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addr, err := proto.NewAddressFromString(minerType)
	if err != nil {
		log.Println(err)
	}

	entries, _, err := cl.Addresses.AddressesData(ctx, addr)
	if err != nil {
		log.Println(err)
	}

	for _, m := range entries {
		miner := &Miner{}
		db.FirstOrCreate(miner, &Miner{Address: m.GetKey()})

		if minerType == MobileAddress {
			miner.MiningHeight = m.ToProtobuf().GetIntValue()
		} else {
			encId := m.ToProtobuf().GetStringValue()
			telId := DecryptMessage(encId)
			telIdInt, err := strconv.Atoi(telId)
			if err != nil {
				log.Println(err)
			}
			miner.TelegramId = int64(telIdInt)
		}

		db.Save(miner)
	}
}

func (m *Monitor) loadReferrals() {
	var miners []*Miner
	db.Find(&miners)

	for _, m := range miners {
		referral, _ := getData("referral", m.Address)

		if referral != nil {
			ref := &Miner{}
			db.First(ref, &Miner{Address: referral.(string)})
			m.ReferralID = ref.ID
			db.Save(m)
		}
	}
}

func (m *Monitor) start() {
	for {
		m.loadMiners(MobileAddress)

		m.loadMiners(TelegramAddress)

		m.loadReferrals()

		log.Println("Done update.")

		time.Sleep(time.Second * MonitorTick)
	}
}

func initMonitor() {
	m := &Monitor{}
	go m.start()
}
