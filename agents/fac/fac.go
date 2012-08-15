
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
  inBuff *buff.Buffer

  OutCommod string
  OutUnits string
  outBuff *buff.Buffer

  CreateRate float64
  ConvertAmt float64
  ConvertPeriod time.Duration
  ConvertOffset time.Duration

  Sim *sim.Sim
}

func (f *Fac) init() {
  if f.inBuff == nil {
    f.inBuff = &buff.Buffer{}
    f.outBuff = &buff.Buffer{}
  }
  if f.ConvertPeriod == 0 {
    f.ConvertPeriod = f.Sim.Eng.Step
  }
}

func (f *Fac) Parent() msg.Communicator {
  return nil
}

func (f *Fac) InSize(qty float64) error {
  f.init()
  return f.inBuff.SetCapacity(qty)
}

func (f *Fac) OutSize(qty float64) error {
  f.init()
  return f.outBuff.SetCapacity(qty)
}

func (f *Fac) Tick(tm time.Duration) {
  fmt.Println("Ticking")
  f.init()
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
    tran = trans.NewOffer(f)
    units = f.OutUnits
  }

  r := rsrc.NewGeneric(qty, units)
  tran.SetResource(r)

  m := msg.New(f, f.Sim.Mkts[commod])
  m.Trans = tran
  m.SendOn()
}

func (f *Fac) Tock(tm time.Duration) {
  f.init()
  f.approveOffers()
  f.convertRes()

  qty := math.Min(f.CreateRate, f.outBuff.Space())
  f.createRes(qty)
}

func (f *Fac) approveOffers() {
  for _, m := range f.queuedOrders {
    m.Trans.Approve()
  }
  f.queuedOrders = []*msg.Message{}
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

  now := int64(f.Sim.Eng.SinceStart())
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
  fmt.Println(f.Name, " sending stuff")
  f.init()
  rs, err := f.outBuff.PopQty(tran.Resource().Qty())
  check(err)
  tran.Manifest = rs
}

func (f *Fac) AddResource(tran *trans.Transaction) {
  fmt.Println(f.Name, " getting stuff: ", tran.Manifest)
  f.init()
  err := f.inBuff.PushAll(tran.Manifest)
  check(err)
}

func check(err error) {
  if err != nil {
    panic(err.Error())
  }
}

