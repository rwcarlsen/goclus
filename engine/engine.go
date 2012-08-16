
package engine

import "time"
import "fmt"

type Ticker interface {
  Tick(time.Duration)
}

type Tocker interface {
  Tock(time.Duration)
}

type TickTocker interface {
  Ticker
  Tocker
}

type Resolver interface {
  Resolve(time.Duration)
}

type Engine struct {
  Duration time.Duration
  Start time.Time
  Step time.Duration
  tickers []Ticker
  tockers []Tocker
  resolvers []Resolver
  tm time.Time // current time (in the simulation)
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
  fmt.Println("start-time: ", e.Start)
  end := e.Start.Add(e.Duration)
  for tm := e.Start; tm.Before(end); tm = tm.Add(e.Step) {
    fmt.Println("timestep: ", tm)
    fmt.Println("ticking...")
    for _, t := range e.tickers {
      t.Tick(e.Step)
    }
    fmt.Println("resolving...")
    for _, r := range e.resolvers {
      r.Resolve(e.Step)
    }
    fmt.Println("tocking...")
    for _, t := range e.tockers {
      t.Tock(e.Step)
    }
  }
}

func (e *Engine) Time() time.Time {
  return e.tm
}

func (e *Engine) SinceStart() time.Duration {
  return e.tm.Sub(e.Start)
}

