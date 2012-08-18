
package main

import (
  "fmt"
  "github.com/rwcarlsen/goclus/books"
  "github.com/rwcarlsen/goclus/trans"
  "github.com/rwcarlsen/goclus/msg"
  "github.com/rwcarlsen/goclus/loader"
  "github.com/rwcarlsen/goclus/agents/fac"
  "github.com/rwcarlsen/goclus/agents/mkt"
)

func registerAgents(l *loader.Loader) {
  l.Register(fac.Fac{})
  l.Register(mkt.Mkt{})
}

func main() {
  // load simulation
  l := &loader.Loader{}
  registerAgents(l)
  err := l.LoadSim("input.json")
  if err != nil {
    fmt.Println(err)
  }

  // setup book-keeping
  transCh := make(chan *trans.Transaction)
  msgCh := make(chan *msg.Message)
  trans.ToOutput = transCh
  msg.ToOutput = msgCh
  b := books.Books{
    TransIn: transCh,
    MsgIn: msgCh,
    Eng: l.Engine,
  }
  b.Collect()

  l.Engine.Run()

  // finish book-keeping
  b.Close()
  err = b.Dump()
  if err != nil {
    fmt.Println(err)
  }
}

