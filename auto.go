
package main

import (
  "fmt"
  "github.com/rwcarlsen/goclus/books"
  "github.com/rwcarlsen/goclus/sim"
  "github.com/rwcarlsen/goclus/agents/fac"
  "github.com/rwcarlsen/goclus/agents/mkt"
)

func registerAgents(l *sim.Loader) {
  l.Register(fac.Fac{})
  l.Register(mkt.Mkt{})
}

func main() {
  // load simulation
  l := &sim.Loader{}
  registerAgents(l)
  err := l.LoadSim("input.json")
  if err != nil {
    fmt.Println(err)
    return
  }

  bks := &books.Books{Eng: l.Engine}
  bks.Collect()
  defer bks.Close()
  l.Engine.RegisterMsgNotify(bks)
  l.Engine.RegisterTransNotify(bks)

  l.Engine.Run()
  err = bks.Dump()
  if err != nil {
    fmt.Println(err)
  }
}

