
package fac

import (
  "time"
  "math"
  "github.com/rwcarlsen/goclus/rsrc"
  "github.com/rwcarlsen/goclus/rsrc/buff"
  "github.com/rwcarlsen/goclus/trans"
  "github.com/rwcarlsen/goclus/msg"
  "github.com/rwcarlsen/goclus/sim"
)

type Fac struct {
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

func (f *Fac) Parent() msg.Communicator {
  return nil
}

func (f *Fac) InSize(qty float64) error {
  return f.inBuff.SetCapacity(qty)
}

func (f *Fac) OutSize(qty float64) error {
  return f.outBuff.SetCapacity(qty)
}

func (f *Fac) Tick(tm time.Duration) {
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
  if t == trans.Offer {
    units = f.OutUnits
  }
  r := rsrc.NewGeneric(qty, units)
  tran := trans.NewRequest(f)
  tran.SetResource(r)

  m := msg.New(f, f.Sim.Mkts[commod])
  m.Trans = tran
  m.SendOn()
}

func (f *Fac) Tock(tm time.Duration) {
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

func check(err error) {
  if err != nil {
    panic(err.Error())
  }
}

func (f *Fac) RemoveResource(tran *trans.Transaction) {
  rs, err := f.outBuff.PopQty(tran.Resource().Qty())
  check(err)
  tran.Manifest = rs
}

func (f *Fac) AddResource(tran *trans.Transaction) {
  err := f.inBuff.PushAll(tran.Manifest)
  check(err)
}
