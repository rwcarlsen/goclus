
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
//
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
//
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

// Message is canonical way to send information between simulation agents.
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
  sender Communicator
  receiver Communicator
  // Payload can be used as desired to send arbitrary information.
  Payload interface{}
  PrevOwner Communicator
  Owner Communicator
  pathStack []Communicator
  hasDest bool
}

// New creates a new message with receiver as the intended final destination. The
// returned message is immediately sendable via the SendOn method.
func New(sender, receiver Communicator) *Message {
  if receiver == nil {
    panic("msg: cannot have nil message receiver")
  }
  return &Message{
    Dir: Up,
    sender: sender,
    receiver: receiver,
    Owner: sender,
    pathStack: []Communicator{sender},
  }
}

// Sender returns the communicator that originally sent this message.
func (m *Message) Sender() Communicator {
  return m.sender
}

// Receiver returns the original intended recipient of this message.
func (m *Message) Receiver() *Message {
  return m.receiver
}

// Clone
func (m *Message) Clone() *Message {
  clone := *m
  clone.Trans = m.Trans.Clone()
  return &clone
}

// SendOn
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

// SetNext 
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

