
package fac

import (
  "fmt"
  "time"
  "math"
  "github.com/rwcarlsen/goclus/rsrc"
  "github.com/rwcarlsen/goclus/rsrc/buff"
  "github.com/rwcarlsen/goclus/trans"
  "github.com/rwcarlsen/goclus/msg"
  "github.com/rwcarlsen/goclus/sim"
)

type Fac struct {
  Name string
  queuedOrders []*msg.Message

  InCommod string
  InUnits string
  InSize float64
  inBuff *buff.Buffer

  OutCommod string
  OutUnits string
  OutSize float64
  outBuff *buff.Buffer

  CreateRate float64
  ConvertAmt float64
  ConvertPeriod time.Duration
  ConvertOffset time.Duration
  eng *sim.Engine
  Test []int
}

func (f *Fac) init() {
  if f.inBuff == nil {
    f.inBuff = &buff.Buffer{}
    f.inBuff.SetCapacity(f.InSize)
    f.outBuff = &buff.Buffer{}
    f.outBuff.SetCapacity(f.OutSize)
  }
}

func (f *Fac) Id() string {
  return f.Name
}

func (f *Fac) SetId(id string) {
  f.Name = id
}

func (f *Fac) Parent() msg.Communicator {
  return nil
}

func (f *Fac) SetParent(par msg.Communicator) {
}

func (f *Fac) Tick(eng *sim.Engine) {
  fmt.Print(f.Name + " ticking: ")
  fmt.Println(f.Test)
  if len(f.Test) > 0 {
    f.Test[0]++
  }

  f.init()
  f.eng = eng

  // make offers
  qty := f.outBuff.Qty()
  if qty > rsrc.EPS {
    f.genMsg(f.OutCommod, qty, trans.Offer)
  }

  // make requests
  qty = f.inBuff.Space()
  if qty > rsrc.EPS {
    f.genMsg(f.InCommod, qty, trans.Request)
  }
}

func (f *Fac) genMsg(commod string, qty float64, t trans.TransType) {
  units := f.InUnits
  tran := trans.NewRequest(f)
  if t == trans.Offer {
    units = f.OutUnits
    tran = trans.NewOffer(f)
  }
  r := rsrc.NewGeneric(qty, units)
  tran.SetResource(r)

  mkt, _ := f.eng.GetComm(commod)
  m := msg.New(f, mkt)
  m.Trans = tran
  m.SendOn()
}

func (f *Fac) Tock(eng *sim.Engine) {
  f.init()
  f.eng = eng
  if f.ConvertPeriod == 0 {
    f.ConvertPeriod = eng.Step
  }

  f.approveOffers()
  f.convertRes()

  qty := math.Min(f.CreateRate, f.outBuff.Space())
  f.createRes(qty)
}

func (f *Fac) approveOffers() {
  for _, m := range f.queuedOrders {
    m.Trans.Approve()
  }
  f.queuedOrders = msg.Group{}
}

func (f *Fac) createRes(qty float64) {
  if qty < rsrc.EPS {
    return
  }
  r := rsrc.NewGeneric(qty, f.OutUnits)
  f.outBuff.PushOne(r)
}

func (f *Fac) convertRes() {
  qty := math.Min(f.ConvertAmt, f.outBuff.Space())
  qty = math.Min(qty, f.inBuff.Qty())

  now := int64(f.eng.SinceStart())
  rem := (now + int64(f.ConvertOffset)) % int64(f.ConvertPeriod)
  if qty <= rsrc.EPS {
    return
  } else if rem > 0 {
    return
  }

  rs, err := f.inBuff.PopQty(qty)
  check(err)
  if f.InUnits == f.OutUnits {
    f.outBuff.PushAll(rs)
  } else {
    f.createRes(qty)
  }
}

func (f *Fac) Receive(m *msg.Message) {
  if m.Sender == f {
    f.queuedOrders = append(f.queuedOrders, m)
  }
}

func (f *Fac) RemoveResource(tran *trans.Transaction) {
  fmt.Println(f.Name, " sending qty=", tran.Resource().Qty(), "of", f.OutCommod)
  f.init()
  rs, err := f.outBuff.PopQty(tran.Resource().Qty())
  check(err)
  tran.Manifest = rs
}

func (f *Fac) AddResource(tran *trans.Transaction) {
  fmt.Println(f.Name, " getting qty=", tran.Resource().Qty(), "of", f.InCommod)
  f.init()
  err := f.inBuff.PushAll(tran.Manifest)
  check(err)
}

func check(err error) {
  if err != nil {
    panic(err.Error())
  }
}

