
package main

import (
  "fmt"
  "github.com/rwcarlsen/goclus/sim"
  "github.com/rwcarlsen/goclus/agents/fac"
  "github.com/rwcarlsen/goclus/agents/mkt"
  "github.com/rwcarlsen/goclus/books"
)

func registerAgents(l *sim.Loader) {
  l.Register(fac.Fac{})
  l.Register(mkt.Mkt{})
  l.Register(books.Books{})
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
  l.Engine.Run()
}

