
package books

import (
  "os"
  "time"
  "encoding/json"
  "reflect"
  "github.com/rwcarlsen/goclus/trans"
  "github.com/rwcarlsen/goclus/engine"
  "github.com/rwcarlsen/goclus/msg"
)

var Eng *engine.Engine

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
  aId int
  TranDat []*TransData
  AgentDat map[interface{}]*AgentData
  done chan bool
}

func (b *Books) Close() {
  b.done<-true
}

// Collect dispatches a goroutine that records data fed into transIn and commIn
// terminating when the Close method is called.
func (b *Books) Collect(transIn chan *trans.Transaction, commIn chan msg.Communicator) {
  b.done = make(chan bool)
  go func() {
    for {
      select {
        case t := <-transIn:
          b.regTrans(t)
        case c := <-commIn:
          b.regComm(c)
        case <-b.done:
          return
      }
    }
  }()
}

func (b *Books) init() {
  if b.AgentDat == nil {
    b.AgentDat = map[interface{}]*AgentData{}
  }
}

func (b *Books) regTrans(t *trans.Transaction) {
  id, tid := 0, 0
  if len(b.TranDat) > 0 {
    id = b.TranDat[len(b.TranDat)-1].Id + 1
    tid = b.TranDat[len(b.TranDat)-1].TransId + 1
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
    b.TranDat = append(b.TranDat, tdat)
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
  b.init()
  if _, ok := b.AgentDat[a]; ok {
    return
  }

  b.aId++
  tp := reflect.Indirect(reflect.ValueOf(a)).Type()

  b.AgentDat[a] = &AgentData{
    Id: b.aId,
    Type: tp.PkgPath() + "." + tp.Name(),
    Born: getTime(),
    ParentId: -1,
  }
}

func (b *Books) Dump() error {
  agents := []*AgentData{}
  for _, val := range b.AgentDat {
    agents = append(agents, val)
  }

  err1 := dump("agents.out", agents)
  err2 := dump("trans.out", b.TranDat)
  if err1 != nil {
    return err1
  } else if err2 != nil {
    return err2
  }
  return nil
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

func getTime() time.Time {
  if Eng != nil {
    return Eng.Time()
  }
  return time.Time{}
}

