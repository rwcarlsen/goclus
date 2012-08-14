
package timer

var TI := timer{}

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

func RegisterTick(t Ticker) {
  TI.RegisterTick(t)
}

func RegisterTock(t Tocker) {
  TI.RegisterTock(t)
}

func RegisterResolve(r Resolver) {
  TI.RegisterResolve(r)
}

type Timer struct {
  Duration time.Duration
  Start time.Time
  Step time.Duration
  tickers []Ticker
  tockers []Tocker
  resolvers []Resolver
  tm time.Time // current time (in the simulation)
}

func (ti *Timer) init() {
  if ti.tickers == nil || ti.tockers == nil || ti.resolvers == nil {
    ti.tickers = []Ticker{}
    ti.tockers = []Tocker{}
    ti.resolvers = []Resolver{}
  }
}

func (ti *Timer) RegisterTick(t Ticker) {
  ti.init()
  ti.tickers = append(ti.tickers, t)
}

func (ti *Timer) RegisterTock(t Tocker) {
  ti.init()
  ti.tockers = append(ti.tockers, t)
}

func (ti *Timer) RegisterResolve(r Resolver) {
  ti.init()
  ti.resolvers = append(ti.resolvers, r)
}

func (ti *Timer) RunSim() {
  for tm := ti.Start; !tm.After(ti.Start.Add(ti.Duration)); tm = tm.Add(ti.Step) {
    for _, t := range ti.tickers {
      t.Tick(ti.Step)
    }
    for _, t := range ti.tockers {
      t.Tock(ti.Step)
    }
    for _, r := range ti.resolvers {
      r.Resolve(ti.Step)
    }
  }
}

func (ti *Timer) Now() time.Time {
  return ti.tm
}

func (ti *Timer) SinceStart() time.Duration {
  return ti.Sub(ti.Start)
}

