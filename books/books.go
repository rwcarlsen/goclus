
package books

import (
  "os"
  "time"
  "encoding/json"
  "reflect"
  "github.com/rwcarlsen/goclus/msg"
  "github.com/rwcarlsen/goclus/trans"
  "github.com/rwcarlsen/goclus/sim"
)

type TransData struct {
  Id int
  TransId int
  SupId int
  ReqId int
  ResType string
  Qty float64
  Units string
}

type AgentData struct {
  Id int
  Type string
  Born time.Time
  ParentId int
}

type Books struct {
  id string
  eng *sim.Engine
  aId int
  done chan bool
  transIn chan *trans.Transaction
  msgIn chan *msg.Message
  miscIn chan interface{}
  tranDat []*TransData
  AgentDat map[interface{}]*AgentData
  MiscDat []interface{}
}

func (b *Books) SetId(id string) {
  b.id = id
}

func (b *Books) Id() string {
  return b.id
}

func (b *Books) SetEngine(e *sim.Engine) {
  b.init()
  b.msgIn<-m
}

func (b *Books) End() {
  b.done<-true
  b.saveData()
}

func (b *Books) MsgNotify(m *msg.Message) {
  b.init()
  b.msgIn<-m
}

func (b *Books) TransNotify(t *trans.Transaction) {
  b.init()
  b.transIn<-t
}

func (b *Books) init() {
  if b.done != nil {
    return
  }
  b.done = make(chan bool)
  b.transIn = make(chan *trans.Transaction)
  b.msgIn = make(chan *msg.Message)
  go func() {
    for {
      select {
        case t := <-b.transIn:
          b.regTrans(t)
        case m := <-b.msgIn:
          b.regComm(m.PrevOwner)
          b.regComm(m.Owner)
        case i := <-b.miscIn:
          b.MiscDat = append(b.MiscDat, i)
        case <-b.done:
          return
      }
    }
  }()
}

func (b *Books) regTrans(t *trans.Transaction) {
  id, tid := 0, 0
  if len(b.tranDat) > 0 {
    id = b.tranDat[len(b.tranDat)-1].Id + 1
    tid = b.tranDat[len(b.tranDat)-1].TransId + 1
  }

  b.regAgent(t.Sup)
  b.regAgent(t.Req)

  for _, r := range t.Manifest {
    tp := reflect.Indirect(reflect.ValueOf(r)).Type()
    tdat := &TransData{
      Id: id,
      TransId: tid,
      SupId: b.AgentDat[t.Sup].Id,
      ReqId: b.AgentDat[t.Req].Id,
      ResType: tp.PkgPath() + "." + tp.Name(),
      Qty: r.Qty(),
      Units: r.Units(),
    }
    b.tranDat = append(b.tranDat, tdat)
  }
}

func (b *Books) regComm(c msg.Communicator) {
  b.regAgent(c)

  // this comes last to prevent infinite looping
  if par := c.Parent(); par != nil {
    b.regComm(par)
    b.AgentDat[c].ParentId = b.AgentDat[par].Id
  }
}

func (b *Books) regAgent(a interface{}) {
  if b.AgentDat == nil {
    b.AgentDat = map[interface{}]*AgentData{}
  } else if _, ok := b.AgentDat[a]; ok {
    return
  }

  b.aId++
  tp := reflect.Indirect(reflect.ValueOf(a)).Type()

  b.AgentDat[a] = &AgentData{
    Id: b.aId,
    Type: tp.PkgPath() + "." + tp.Name(),
    Born: b.getTime(),
    ParentId: -1,
  }
}

func (b *Books) saveData() error {
  agents := []*AgentData{}
  for _, val := range b.AgentDat {
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
  if b.Eng != nil {
    return b.Eng.Time()
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

