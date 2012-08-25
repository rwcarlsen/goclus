// Package msg declares types and interfaces used for inter-agent
// communication.
package msg

import (
	"github.com/rwcarlsen/goclus/trans"
)

var listeners []Listener

type msgDir int

const (
	// Up indicates a message passage to/through parent channels - toward the
	// message receiver.
	Up msgDir = iota
	// Down indicates a return-trip message retracing its "Up" path - toward
	// the message sender.
	Down
)

type Group []*Message

// Communicator is implemented by all agents that require the ability to
// communicate with other simulation agents.
// Note that this method should generally not be invoked directly;
// inter-agent message passing should be achieved via a message's SendOn
// method.
type Communicator interface {
	Receive(*Message)
	Parent() Communicator
	SetParent(Communicator)
}

// Listener is implemented by entities that desire to receive notifications
// every time a message is passed between any two simulation agents.
type Listener interface {
	MsgNotify(*Message)
}

// ListenAll adds l to a global list of agents that receive notifications
// for every message passed between any two agents (usually used by "special"
// agents e.g. book-keeper, etc.).
// These notifications are sent every time a message SendOn method is
// called - before the receiver actually receives the message. Simulation
// execution continues only after l's MsgNotify method returns.
func ListenAll(l Listener) {
	listeners = append(listeners, l)
}

func notifyListeners(m *Message) {
	for _, l := range listeners {
		l.MsgNotify(m)
	}
}

// Message is the canonical way to send information between simulation agents.
//
// Creating and sending a message:
//
//    recv := eng.GetComm("foo")
//    m := msg.New(a, recv)
//    m.SendOn()
// 
// Returning a message to its sender:
//
//    m.Dir = msg.Down
//    m.SendOn()
type Message struct {
	// Dir defaults to Up (sending a message toward its receiver).
	Dir msgDir
	// Trans is used to carry desired/matched transaction information between
	// agents.
	Trans *trans.Transaction
	// Payload can be used as desired to send arbitrary information.
	Payload   interface{}
	PrevOwner Communicator
	Owner     Communicator
	sender    Communicator
	receiver  Communicator
	pathStack []Communicator
	hasDest   bool
}

// New creates a new message with receiver as the intended final destination. The
// returned message is immediately sendable via the SendOn method.
func New(sender, receiver Communicator) *Message {
	if receiver == nil {
		panic("msg: cannot have nil message receiver")
	}
	return &Message{
		Dir:       Up,
		sender:    sender,
		receiver:  receiver,
		Owner:     sender,
		pathStack: []Communicator{sender},
	}
}

// Sender returns the communicator that originally sent this message.
func (m *Message) Sender() Communicator {
	return m.sender
}

// Receiver returns the original intended recipient of this message.
func (m *Message) Receiver() Communicator {
	return m.receiver
}

// Clone returns a shallow copy of this message except the copy has a clone
// of the message's transaction.
func (m *Message) Clone() *Message {
	clone := *m
	clone.Trans = m.Trans.Clone()
	return &clone
}

// SendOn sends the message toward its intended destination.
//
// If the message Dir is Up and SetNext is not called, the message is sent
// to the current communicator's parent. If the current communicator has no
// parent, the message is sent to its receiver as specified when the
// message was created.
// If the message Dir is Up and SetNext has been called, the message is sent to
// the receiver specified in the SetNext call.
//
// If the message Dir is down, the message retraces its upward path sending
// itself to each previous owner until it reaches its original sender.
func (m *Message) SendOn() {
	if !m.hasDest {
		m.autoSetNext()
	}

	m.validateForSend()

	if m.Dir == Down {
		m.pathStack = m.pathStack[:len(m.pathStack)-1]
	}

	next := m.pathStack[len(m.pathStack)-1]
	m.PrevOwner, m.Owner = m.Owner, next

	notifyListeners(m)
	m.hasDest = false
	next.Receive(m)
}

// SetNext allows manual specification of the next message receiver.
// If the message Dir is Down, calls to SetNext do nothing, and the message
// will continue retrace its previous path with each SendOn invocation as
// if SetNext had not been called.
func (m *Message) SetNext(dest Communicator) {
	if m.Dir == Down {
		return
	}
	m.pathStack = append(m.pathStack, dest)
	m.hasDest = true
}

func (m *Message) autoSetNext() {
	next := m.Owner.Parent()
	if next == nil {
		next = m.receiver
	}
	m.SetNext(next)
}

func (m *Message) validateForSend() {
	hasDest := false
	i := -1
	if m.Dir == Up {
		hasDest = len(m.pathStack) > 0
		i = len(m.pathStack) - 1
	} else if m.Dir == Down {
		hasDest = len(m.pathStack) > 1
		i = len(m.pathStack) - 2
	}

	if !hasDest {
		panic("msg: no message receiver")
	} else if next := m.pathStack[i]; next == m.Owner {
		panic("msg: Circular message send attempt")
	}
}
