// Package books is an agent and transaction recording service for simulations.
package books

import (
	"encoding/json"
	"github.com/rwcarlsen/goclus/sim"
	"github.com/rwcarlsen/goclus/trans"
	"os"
	"reflect"
	"time"
)

// transData holds simulation transaction information in an
// output-write-ready format.
type transData struct {
	Id      int
	TransId int
	SupId   int
	ReqId   int
	ResType string
	Qty     float64
	Units   string
}

// transData holds simulation agent information in an
// output-write-ready format.
type agentData struct {
	Id       int
	Name     string
	Type     string
	Born     time.Time
	ParentId int
}

// Books is an agent that records participating-agent and transaction activity
// during a simulation.  Note that the Books agent is intended to be registered
// with a sim.Engine for services for which it implements appropriate
// interfaces. The public methods are NOT intended to be invoked by anything
// other than sim.Engine during the course of a simulation.
type Books struct {
	sim.Agenty
	eng      *sim.Engine
	tId      int // next trans id tracker
	eId      int // next trans entry id tracker
	done     chan bool
	transIn  chan *trans.Transaction
	msgIn    chan *sim.Message
	miscIn   chan interface{}
	tranDat  []*transData
	agentDat map[int]*agentData
	miscDat  []interface{}
}

// Start spins off a goroutine that book-keeps all transaction and agent
// information as provided via MsgNotify and TransNotify.
func (b *Books) Start(e *sim.Engine) {
	sim.ListenAllMsg(b)
	trans.ListenAll(b)

	b.done = make(chan bool)
	b.agentDat = map[int]*agentData{}
	b.transIn = make(chan *trans.Transaction)
	b.msgIn = make(chan *sim.Message)
	go func() {
		for {
			select {
			case t := <-b.transIn:
				b.regTrans(t)
			case m := <-b.msgIn:
				b.regAgent(m.PrevOwner)
				b.regAgent(m.Owner)
			case i := <-b.miscIn:
				b.miscDat = append(b.miscDat, i)
			case <-b.done:
				return
			}
		}
	}()
}

// End allows final recording operations to take place before the
// simulation closes; most notably, writing remaining collected information
// to an output file.
func (b *Books) End(e *sim.Engine) {
	b.done <- true
	b.saveData()
}

// MsgNotify is used to collect information about agents participating in a
// simulation from the sim.Engine.
func (b *Books) MsgNotify(m *sim.Message) {
	b.msgIn <- m
}

// TransNotify is used to collect information about matched, executed
// transactions as they occur through a simulation from the sim.Engine.
func (b *Books) TransNotify(t *trans.Transaction) {
	b.transIn <- t
}

func (b *Books) regTrans(t *trans.Transaction) {
	b.regAgent(t.Sup.(sim.Agent))
	b.regAgent(t.Req.(sim.Agent))
	for _, r := range t.Manifest {
		tp := reflect.Indirect(reflect.ValueOf(r)).Type()
		tdat := &transData{
			Id:      b.eId,
			TransId: b.tId,
			SupId:   t.Sup.(sim.Agent).Id(),
			ReqId:   t.Req.(sim.Agent).Id(),
			ResType: tp.PkgPath() + "." + tp.Name(),
			Qty:     r.Qty(),
			Units:   r.Units(),
		}
		b.eId++
		b.tranDat = append(b.tranDat, tdat)
	}
	b.tId++
}

func (b *Books) regAgent(a sim.Agent) {
	if _, ok := b.agentDat[a.Id()]; ok {
		return
	}

	tp := reflect.Indirect(reflect.ValueOf(a)).Type()
	b.agentDat[a.Id()] = &agentData{
		Id:   a.Id(),
		Name: a.Name(),
		Type: tp.PkgPath() + "." + tp.Name(),
		Born: b.getTime(),
	}

	// this comes last to prevent infinite looping
	if par := a.Parent(); par != nil {
		b.agentDat[a.Id()].ParentId = par.Id()
		b.regAgent(par)
	}
}

func (b *Books) saveData() error {
	agents := []*agentData{}
	for _, val := range b.agentDat {
		agents = append(agents, val)
	}

	err1 := dump("agents.out", agents)
	err2 := dump("trans.out", b.tranDat)
	if err1 != nil {
		return err1
	} else if err2 != nil {
		return err2
	}
	return nil
}

func (b *Books) getTime() time.Time {
	if b.eng != nil {
		return b.eng.Time()
	}
	return time.Time{}
}

func dump(name string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(data)
	return nil
}
