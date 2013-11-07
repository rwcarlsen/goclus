package sim

import (
	"errors"
	"fmt"
	"time"
)

type Engine struct {
	Duration  time.Duration
	Step      time.Duration
	Load      *Loader
	services  map[string]Agent
	tickers   []Ticker
	resolvers []Resolver
	tockers   []Tocker
	starters  []Starter
	enders    []Ender
	tm        time.Time // current time (in the simulation)
	nextId    int       // the next agent ID
}

// RegisterAll registers agent a to receive time-related notifications for
// all sim package interfaces implemented.
func (e *Engine) RegisterAll(a Agent) (ifaces []string) {
	e.nextId++
	a.SetId(e.nextId)

	if t, ok := a.(Ticker); ok {
		e.tickers = append(e.tickers, t)
		ifaces = append(ifaces, "Ticker")
	}
	if t, ok := a.(Tocker); ok {
		e.tockers = append(e.tockers, t)
		ifaces = append(ifaces, "Tocker")
	}
	if t, ok := a.(Resolver); ok {
		e.resolvers = append(e.resolvers, t)
		ifaces = append(ifaces, "Resolver")
	}
	if t, ok := a.(Starter); ok {
		t.Start(e)
	}
	if t, ok := a.(Ender); ok {
		e.enders = append(e.enders, t)
		ifaces = append(ifaces, "Ender")
	}
	return ifaces
}

// RegisterService registers an agent with a simulation-global list that
// can be accessed by all agents.  The agent's ID will be used as the
// retrival key.
func (e *Engine) RegisterService(a Agent) error {
	if e.services == nil {
		e.services = map[string]Agent{}
	}

	if _, ok := e.services[a.Name()]; ok {
		return errors.New("sim: duplicate service name '" + a.Name() + "'")
	}
	e.services[a.Name()] = a
	return nil
}

// GetService returns the agent registered under id or an error if no agent
// has registered under the given id.
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

func (e *Engine) Run() {
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
