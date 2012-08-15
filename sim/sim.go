
package sim

type Sim struct {
  Eng *engine.Engine
  Mkts map[string]Communicator
}

func New(e *engine.Engine) *Sim {
  return &Sim{
    Eng: e,
    Mkts: map[string]Communicator{},
  }
}
