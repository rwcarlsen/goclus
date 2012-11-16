package main

import (
	"github.com/rwcarlsen/goclus/agents/fac"
	"github.com/rwcarlsen/goclus/agents/mkt"
	"github.com/rwcarlsen/goclus/books"
	"github.com/rwcarlsen/goclus/rsrc"
	"github.com/rwcarlsen/goclus/sim"
	"time"
)

func main() {
	var month time.Duration = 43829 * time.Minute
	eng := &sim.Engine{
		Duration: 36 * month,
		Step:     month,
	}
	config(eng)

	eng.Run()
}

func config(eng *sim.Engine) {
	milk := "milk"
	cheese := "cheese"
	src := &fac.Fac{
		OutCommod:  milk,
		OutUnits:   milk,
		OutSize:    5,
		CreateRate: rsrc.INFINITY,
	}
	src.SetId("src")

	null := &fac.Fac{
		InCommod:      milk,
		InUnits:       milk,
		InSize:        5,
		OutCommod:     cheese,
		OutUnits:      cheese,
		OutSize:       5,
		ConvertAmt:    5,
		ConvertPeriod: 1,
		ConvertOffset: 0,
	}
	null.SetId("null")

	null2 := &fac.Fac{
		InCommod:      cheese,
		InUnits:       cheese,
		InSize:        5,
		OutCommod:     milk,
		OutUnits:      milk,
		OutSize:       3,
		ConvertAmt:    5,
		ConvertPeriod: 1,
		ConvertOffset: 0,
	}
	null2.SetId("null2")

	snk := &fac.Fac{
		InCommod: cheese,
		InUnits:  cheese,
		InSize:   rsrc.INFINITY,
	}
	snk.SetId("snk")

	milkMkt := &mkt.Mkt{Shuffle: true}
	milkMkt.SetId(milk)
	cheeseMkt := &mkt.Mkt{Shuffle: true}
	cheeseMkt.SetId(cheese)

	bks := &books.Books{}

	eng.RegisterAll(bks)
	eng.RegisterAll(src)
	eng.RegisterAll(snk)
	eng.RegisterAll(null)
	eng.RegisterAll(null2)
	eng.RegisterAll(milkMkt)
	eng.RegisterAll(cheeseMkt)

	eng.RegisterService(milkMkt)
	eng.RegisterService(cheeseMkt)
}
