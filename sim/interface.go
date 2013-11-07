package sim

import (
	"fmt"
)

type Agent interface {
	SetId(int)
	Id() int
	Name() string
	SetName(string)
	Parent() Agent
	SetParent(Agent)
	// Receive should should generally not be invoked directly; inter-agent
	// message passing should be achieved via a message's SendOn method.
	Receive(*Message)
}

// Agenty is provided as a convenient way to automatically satisfy the Id and
// SetId methods of the Agent interface.  Simply embed Agenty in the sim
// agent's struct:
//
//    type MyAgent struct {
//       sim.Agenty
//       ...
//    }
type Agenty struct {
	id     int
	name   string
	parent Agent
}

// Id returns the value passed via SetId or the empty string if it hasn't
// been called.
func (a *Agenty) Id() int {
	return a.id
}

// SetId sets the value returned by Id.  Panics if called after already being
// set.
func (a *Agenty) SetId(id int) {
	if a.id != 0 {
		panic(fmt.Sprintf("duplicate id set on agent '%v', id=%v", a.name, a.id))
	}
	a.id = id
}

// Name returns the agent's name.  Used in simulation service registration and
// pretty printing.
func (a *Agenty) Name() string {
	return a.name
}

// SetName sets the agent's potentially non-unique name.
func (a *Agenty) SetName(name string) {
	a.name = name
}

// Parent returns the value passed via SetParent or nil if SetParent hasn't
// been called.
func (a *Agenty) Parent() Agent {
	return a.parent
}

// SetParent sets the value returned by Parent.
func (a *Agenty) SetParent(p Agent) {
	a.parent = p
}

// Receive does nothing
func (a *Agenty) Receive(*Message) {}

type Starter interface {
	Start(*Engine)
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
