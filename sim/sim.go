package sim

import (
	"errors"
	"fmt"
	"github.com/rwcarlsen/goclus/msg"
	"github.com/rwcarlsen/goclus/trans"
	"time"
)

type Agent interface {
	SetId(string)
	Id() string
}

// Agenty is provided as a convenient way to automatically satisfy the Id and
// SetId methods of the Agent interface.  Simply embed Agenty in the sim 
// agent's struct:
//
//    type MyAgent struct {
//       sim.Agenty
//       ...
//    }
type Agenty string

// Id returns the value passed via SetId or the empty string if it hasn't
// been called.
func (a Agenty) Id() string {
	return string(a)
}

// SetId sets the value returned by Id.
func (a *Agenty) SetId(id string) {
	*a = Agenty(id)
}

type Ender interface {
	End(*Engine)
}

type Ticker interface {
	Tick()
}

type Tocker interface {
	Tock()
}

type Resolver interface {
	Resolve()
}

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
	enders      []Ender
	tm          time.Time // current time (in the simulation)
}

func (e *Engine) RegisterAll(a Agent) (ifaces []string) {
	switch t := a.(type) {
	case Ticker:
		e.tickers = append(e.tickers, t)
		ifaces = append(ifaces, "Ticker")
	}
	switch t := a.(type) {
	case Tocker:
		e.tockers = append(e.tockers, t)
		ifaces = append(ifaces, "Tocker")
	}
	switch t := a.(type) {
	case Resolver:
		e.resolvers = append(e.resolvers, t)
		ifaces = append(ifaces, "Resolver")
	}
	switch t := a.(type) {
	case Ender:
		e.enders = append(e.enders, t)
		ifaces = append(ifaces, "Ender")
	}
	switch t := a.(type) {
	case msg.Listener:
		e.msgListen = append(e.msgListen, t)
		ifaces = append(ifaces, "msg.Listener")
	}
	switch t := a.(type) {
	case trans.Listener:
		e.transListen = append(e.transListen, t)
		ifaces = append(ifaces, "trans.Listener")
	}
	return
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
