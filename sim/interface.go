
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

