
package main

import (
  "time"
  "github.com/rwcarlsen/goclus/sim"
  "github.com/rwcarlsen/goclus/rsrc"
  "github.com/rwcarlsen/goclus/agents/fac"
  "github.com/rwcarlsen/goclus/agents/mkt"
)

func main() {
  simul := sim.New()
  config(simul)

  var month time.Duration = 43829 * time.Minute
  simul.Eng.Duration = 36 * month
  simul.Eng.Start = time.Now()
  simul.Eng.Step = month

  simul.Eng.Run()
}

func config(simul *sim.Sim) {
  src := &fac.Fac{
    OutCommod: "milk",
    OutUnits: "milk",
    CreateRate: rsrc.INFINITY,
    Sim: simul,
  }
  src.OutSize(5)

  null := &fac.Fac{
    InCommod: "milk",
    InUnits: "milk",
    OutCommod: "cheese",
    OutUnits: "cheese",
    ConvertAmt: 5,
    ConvertPeriod: 1,
    ConvertOffset: 0,
    Sim: simul,
  }
  null.InSize(5)
  null.OutSize(5)

  null2 := &fac.Fac{
    InCommod: "cheese",
    InUnits: "cheese",
    OutCommod: "milk",
    OutUnits: "milk",
    ConvertAmt: 5,
    ConvertPeriod: 1,
    ConvertOffset: 0,
    Sim: simul,
  }
  null2.InSize(5)
  null2.OutSize(3)

  snk := &fac.Fac{
    InCommod: "cheese",
    InUnits: "cheese",
    Sim: simul,
  }
  snk.OutSize(rsrc.INFINITY)

  milkMkt := &mkt.Mkt{
    Commod: "milk",
    Shuffle: true,
  }

  cheeseMkt := &mkt.Mkt{
    Commod: "cheese",
    Shuffle: true,
  }

  simul.Eng.RegisterTickTock(src, snk, null, null2)
  simul.Eng.RegisterResolve(milkMkt, cheeseMkt)
}
