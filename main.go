package main

import (
	"log"

	"gorm.io/gorm"
)

var db *gorm.DB

var conf *Config

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conf = initConfig()

	db = initDb()

	initMonitor()

	// var dbminers []*Miner

	// db.Unscoped().Find(&dbminers)
	// for i, dbm := range dbminers {
	// 	if dbm.DeletedAt.Time.Day() == time.Now().Day() && strings.HasPrefix(dbm.Address, "3A") {
	// 		// key := dbm.Address
	// 		// tel := EncryptMessage(strconv.Itoa(int(dbm.TelegramId)))
	// 		// ip := EncryptMessage(dbm.IP)

	// 		// ref := ""

	// 		// if dbm.ReferralID != 0 {
	// 		// 	r := &Miner{}
	// 		// 	db.Unscoped().First(r, dbm.ReferralID)
	// 		// 	ref = r.Address
	// 		// }

	// 		// value := fmt.Sprintf("%%s%%d%%s%%s__%s__%d__%s__%s", tel, dbm.MiningHeight, ip, ref)

	// 		// dataTransaction(key, &value, nil, nil)

	// 		// log.Println(key + " " + value)
	// 		db.Unscoped().Model(&Miner{}).Where("id", dbm.ID).Update("deleted_at", nil)
	// 		log.Println(i)
	// 	}
	// }
}
