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
	m.Get("/ipcount/:ip", ipView)

	return m
}
