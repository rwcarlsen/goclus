
package engine

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
  tm time.Time // current time.(in the simulaeon)
}

func (e *Engine) init() {
  if e.tickers == nil || e.tockers == nil || e.resolvers == nil {
    e.tickers = []Ticker{}
    e.tockers = []Tocker{}
    e.resolvers = []Resolver{}
  }
}

func (e *Engine) RegisterTick(t Ticker) {
  e.init()
  e.tickers = append(e.tickers, t)
}

func (e *Engine) RegisterTock(t Tocker) {
  e.init()
  e.tockers = append(e.tockers, t)
}

func (e *Engine) RegisterResolve(r Resolver) {
  e.init()
  e.resolvers = append(e.resolvers, r)
}

func (e *Engine) RunSim() {
  for tm := e.Start; !tm.After(e.Start.Add(e.Duration)); tm = tm.Add(e.Step) {
    for _, t := range e.tickers {
      t.Tick(e.Step)
    }
    for _, t := range e.tockers {
      t.Tock(e.Step)
    }
    for _, r := range e.resolvers {
      r.Resolve(e.Step)
    }
  }
}

func (e *Engine) Time() time.Time {
  return e.tm
}

func (e *Engine) SinceStart() time.Duration {
  return e.Sub(e.Start)
}

