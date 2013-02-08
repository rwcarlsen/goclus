package sim

import (
	"errors"
	"fmt"
	"github.com/rwcarlsen/goclus/msg"
	"github.com/rwcarlsen/goclus/trans"
	"time"
)

type Engine struct {
	Duration    time.Duration
	Step        time.Duration
	Load        *Loader
	services    map[string]Agent
	msgListen   []msg.Listener
	transListen []trans.Listener
	tickers     []Ticker
	resolvers   []Resolver
	tockers     []Tocker
	starters    []Starter
	enders      []Ender
	tm          time.Time // current time (in the simulation)
}

func (e *Engine) RegisterAll(a Agent) (ifaces []string) {
	if t, ok := a.(Ticker) {
		e.tickers = append(e.tickers, t)
		ifaces = append(ifaces, "Ticker")
	}
	if t, ok := a.(Tocker) {
		e.tockers = append(e.tockers, t)
		ifaces = append(ifaces, "Tocker")
	}
	if t, ok := a.(Resolver) {
		e.resolvers = append(e.resolvers, t)
		ifaces = append(ifaces, "Resolver")
	}
	if t, ok := a.(Starter) {
		e.enders = append(e.enders, t)
		ifaces = append(ifaces, "Starter")
	}
	if t, ok := a.(Ender) {
		e.enders = append(e.enders, t)
		ifaces = append(ifaces, "Ender")
	}
	if t, ok := a.(msg.Listener) {
		e.msgListen = append(e.msgListen, t)
		ifaces = append(ifaces, "msg.Listener")
	}
	if t, ok := a.(trans.Listener) {
		e.transListen = append(e.transListen, t)
		ifaces = append(ifaces, "trans.Listener")
	}
	return ifaces
}

func (e *Engine) RegisterService(a Agent) error {
	if e.services == nil {
		e.services = map[string]Agent{}
	}

	if _, ok := e.services[a.Id()]; ok {
		return errors.New("sim: duplicate service id '" + a.Id() + "'")
	}
	e.services[a.Id()] = a
	return nil
}

func (e *Engine) GetService(id string) (Agent, error) {
	unreg := errors.New("sim: service id '" + id + "' not registered")
	if e.services == nil {
		return nil, unreg
	}
	a, ok := e.services[id]
	if !ok {
		return nil, unreg
	}
	return a, nil
}

func (e *Engine) GetComm(id string) (msg.Communicator, error) {
	v, err := e.GetService(id)
	if err == nil {
		if c, ok := v.(msg.Communicator); ok {
			return c, nil
		}
		return nil, errors.New("sim: cannot convert '" + id + "' to msg.Communicator")
	}
	return nil, err
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
	e.runTimeSteps()
	for _, en := range e.enders {
		en.End(e)
	}
}

func (e *Engine) runTimeSteps() {
	end := e.tm.Add(e.Duration)
	for ; e.tm.Before(end); e.tm = e.tm.Add(e.Step) {
		fmt.Println("timestep: ", e.tm)
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
}

func (e *Engine) Time() time.Time {
	return e.tm
}

func (e *Engine) SinceStart() time.Duration {
	return e.tm.Sub(time.Time{})
}
