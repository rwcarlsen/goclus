
package sim

import (
  "time"
  "fmt"
  "github.com/rwcarlsen/goclus/msg"
  "errors"
)

type Ticker interface {
  Tick(*Engine)
}

type Tocker interface {
  Tock(*Engine)
}

type TickTocker interface {
  Ticker
  Tocker
}

type Resolver interface {
  Resolve(*Engine)
}

type Engine struct {
  Duration time.Duration
  Step time.Duration
  comms map[string]msg.Communicator
  tickers []Ticker
  tockers []Tocker
  resolvers []Resolver
  tm time.Time // current time (in the simulation)
}

func (e *Engine) RegisterComm(name string, c msg.Communicator) {
  if e.comms == nil {
    e.comms = map[string]msg.Communicator{}
  }
  e.comms[name] = c
}

func (e *Engine) GetComm(name string) (msg.Communicator, error) {
  unreg := errors.New("sim: name not registered")
  if e.comms == nil {
    return nil, unreg
  } else if _, ok := e.comms[name]; !ok {
    return nil, unreg
  }
  return e.comms[name], nil
}

func (e *Engine) RegisterTick(ts ...Ticker) {
  e.tickers = append(e.tickers, ts...)
}

func (e *Engine) RegisterTock(ts ...Tocker) {
  e.tockers = append(e.tockers, ts...)
}

func (e *Engine) RegisterTickTock(ts ...TickTocker) {
  for _, t := range ts {
    e.tickers = append(e.tickers, t.(Ticker))
    e.tockers = append(e.tockers, t.(Tocker))
  }
}

func (e *Engine) RegisterResolve(rs ...Resolver) {
  e.resolvers = append(e.resolvers, rs...)
}

func (e *Engine) Run() {
  start := time.Time{}
  end := start.Add(e.Duration)
  for tm := start; tm.Before(end); tm = tm.Add(e.Step) {
    fmt.Println("timestep: ", tm)
    fmt.Println("ticking...")
    for _, t := range e.tickers {
      t.Tick(e)
    }
    fmt.Println("resolving...")
    for _, r := range e.resolvers {
      r.Resolve(e)
    }
    fmt.Println("tocking...")
    for _, t := range e.tockers {
      t.Tock(e)
    }
  }
}

func (e *Engine) Time() time.Time {
  return e.tm
}

func (e *Engine) SinceStart() time.Duration {
  return e.tm.Sub(time.Time{})
}

