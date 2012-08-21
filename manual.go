
package main

import (
  "fmt"
  "time"
  "github.com/rwcarlsen/goclus/sim"
  "github.com/rwcarlsen/goclus/rsrc"
  "github.com/rwcarlsen/goclus/books"
  "github.com/rwcarlsen/goclus/agents/fac"
  "github.com/rwcarlsen/goclus/agents/mkt"
)

func main() {
  var month time.Duration = 43829 * time.Minute
  eng := &sim.Engine{
    Duration: 36 * month,
    Step: month,
  }
  config(eng)

  bks := &books.Books{Eng: eng}
  bks.Collect()
  defer bks.Close()
  eng.RegisterMsgNotify(bks)
  eng.RegisterTransNotify(bks)

  eng.Run()

  err := bks.Dump()
  if err != nil {
    fmt.Println(err)
  }
}

func config(eng *sim.Engine) {
  milk := "milk"
  cheese := "cheese"
  src := &fac.Fac{
    Name: "src",
    OutCommod: milk,
    OutUnits: milk,
    OutSize: 5,
    CreateRate: rsrc.INFINITY,
  }

  null := &fac.Fac{
    Name: "null",
    InCommod: milk,
    InUnits: milk,
    InSize: 5,
    OutCommod: cheese,
    OutUnits: cheese,
    OutSize: 5,
    ConvertAmt: 5,
    ConvertPeriod: 1,
    ConvertOffset: 0,
  }

  null2 := &fac.Fac{
    Name: "null2",
    InCommod: cheese,
    InUnits: cheese,
    InSize: 5,
    OutCommod: milk,
    OutUnits: milk,
    OutSize: 3,
    ConvertAmt: 5,
    ConvertPeriod: 1,
    ConvertOffset: 0,
  }

  snk := &fac.Fac{
    Name: "snk",
    InCommod: cheese,
    InUnits: cheese,
    InSize: rsrc.INFINITY,
  }

  milkMkt := &mkt.Mkt{Shuffle: true}
  cheeseMkt := &mkt.Mkt{Shuffle: true}

  eng.RegisterTickTock(src, snk, null, null2)
  eng.RegisterResolve(milkMkt, cheeseMkt)
  eng.RegisterComm(milk, milkMkt)
  eng.RegisterComm(cheese, cheeseMkt)
}
