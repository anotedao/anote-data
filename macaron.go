package main

import (
	"github.com/go-macaron/cache"
	macaron "gopkg.in/macaron.v1"
)

func initMacaron() *macaron.Macaron {
	m := macaron.Classic()

	m.Use(macaron.Renderer())
	m.Use(cache.Cacher())

	m.Get("/miners", minersView)
	m.Get("/miner/:addr", minerView)
	m.Get("/ping/:addr", pingView)
	m.Get("/ipcount/:ip", ipView)
	m.Get("/confirmation/:addr", checkConfirmationView)
	m.Get("/stats", statsView)

	return m
}
