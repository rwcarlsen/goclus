
package msg

import (
  "github.com/rwcarlsen/goclus/trans"
)

var ToOutput chan *Message

type MsgDir int

const (
  Up MsgDir = iota
  Down
)

type Group []*Message

type Communicator interface {
  Receive(*Message)
  Parent() Communicator
}

type Message struct {
  Dir MsgDir
  Trans *trans.Transaction
  Sender Communicator
  Receiver Communicator
  Notes string
  PrevOwner Communicator
  Owner Communicator
  pathStack []Communicator
  hasDest bool
}

func New(sender, receiver Communicator) *Message {
  return &Message{
    Dir: Up,
    Sender: sender,
    Receiver: receiver,
    Owner: sender,
    pathStack: []Communicator{sender},
  }
}

func (m *Message) Clone() *Message {
  clone := *m
  clone.Trans = m.Trans.Clone()
  return &clone
}

func (m *Message) SendOn() {
  if !m.hasDest {
    m.autoSetDest()
  }

  m.validateForSend()

  if m.Dir == Down {
    m.pathStack = m.pathStack[:len(m.pathStack)-1]
  }

  next := m.pathStack[len(m.pathStack)-1]
  m.PrevOwner, m.Owner = m.Owner, next

  if ToOutput != nil {
    ToOutput<-m
  }

  m.hasDest = false
  next.Receive(m)
}

func (m *Message) SetDest(dest Communicator) {
  if m.Dir == Down {
    return
  }
  m.pathStack = append(m.pathStack, dest)
  m.hasDest = true
}

func (m *Message) autoSetDest() {
  next := m.Owner.Parent()
  if next == nil {
    next = m.Receiver
  }
  m.SetDest(next)
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
    panic("msg: No Message Receiver")
  } else if next := m.pathStack[i]; next == m.Owner {
    panic("msg: Circular message send attempt")
  }
}

