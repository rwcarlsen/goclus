
package sim

import (
  "time"
  "fmt"
  "github.com/rwcarlsen/goclus/msg"
  "github.com/rwcarlsen/goclus/trans"
  "errors"
)

type Agent interface {
  SetId(string)
  Id() string
  SetEngine(*Engine)
}

type Ticker interface {
  Tick()
}

type Tocker interface {
  Tock()
}

type Ender interface {
  End()
}

type TickTocker interface {
  Ticker
  Tocker
}

type Resolver interface {
  Resolve()
}

type Engine struct {
  Duration time.Duration
  Step time.Duration
  Load *Loader
  comms map[string]msg.Communicator
  msgListen []msg.Listener
  transListen []trans.Listener
  tickers []Ticker
  tockers []Tocker
  enders []Ender
  resolvers []Resolver
  tm time.Time // current time (in the simulation)
}

func (e *Engine) RegisterComm(name string, c msg.Communicator) error {
  if e.comms == nil {
    e.comms = map[string]msg.Communicator{}
  }

  if _, ok := e.comms[name]; ok {
    return errors.New("sim: duplicate name registration '" + name + "'")
  }
  e.comms[name] = c
  return nil
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
    e.RegisterTick(t)
    e.RegisterTock(t)
  }
}

func (e *Engine) RegisterResolve(rs ...Resolver) {
  e.resolvers = append(e.resolvers, rs...)
}

func (e *Engine) RegisterEnd(enders ...Ender) {
  e.enders = append(e.enders, enders...)
}

func (e *Engine) RegisterMsgNotify(l msg.Listener) {
  e.msgListen = append(e.msgListen, l)
}

func (e *Engine) RegisterTransNotify(l trans.Listener) {
  e.transListen = append(e.transListen, l)
}

func (e *Engine) MsgNotify(m *msg.Message) {
  for _, l := range e.msgListen {
    l.MsgNotify(m)
  }
}

func (e *Engine) TransNotify(t *trans.Transaction) {
  for _, l := range e.transListen {
    l.TransNotify(t)
  }
}

func (e *Engine) Run() {
  msg.ListenAll(e)
  trans.ListenAll(e)

  start := time.Time{}
  end := start.Add(e.Duration)
  for tm := start; tm.Before(end); tm = tm.Add(e.Step) {
    fmt.Println("timestep: ", tm)
    fmt.Println("ticking...")
    for _, t := range e.tickers {
      t.Tick()
    }
    fmt.Println("resolving...")
    for _, r := range e.resolvers {
      r.Resolve()
    }
    fmt.Println("tocking...")
    for _, t := range e.tockers {
      t.Tock()
    }
  }

  for _, en := range e.enders {
    en.End()
  }
}

func (e *Engine) Time() time.Time {
  return e.tm
}

func (e *Engine) SinceStart() time.Duration {
  return e.tm.Sub(time.Time{})
}

