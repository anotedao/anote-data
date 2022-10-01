package main

import (
	"log"

	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db = initDb()

	initMonitor()

	m := initMacaron()

	m.Run("127.0.0.1", Port)
}
