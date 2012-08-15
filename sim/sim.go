
package sim

import (
  "github.com/rwcarlsen/goclus/msg"
  "github.com/rwcarlsen/goclus/engine"
)

type Sim struct {
  Eng *engine.Engine
  Mkts map[string]msg.Communicator
}

func New() *Sim {
  return &Sim{
    Eng: &engine.Engine{},
    Mkts: map[string]msg.Communicator{},
  }
}
