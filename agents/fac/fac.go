
package fac

import (
  "time"
  "github.com/rwcarlsen/goclus/rsrc"
  "github.com/rwcarlsen/goclus/rsrc/buff"
  "github.com/rwcarlsen/goclus/msg"
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

func (f *Fac) Tick(tm time.Duration) {
  // make offers
  qty := f.outBuff.Qty()
  if qty > rsrc.EPS {
    genMsg(f.OutCommod, qty, trans.Offer)
  }

  // make requests
  qty := f.inBuff.Space()
  if qty > rsrc.EPS {
    genMsg(f.InCommod, qty, trans.Request)
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
  approveOffers()
  convertRes()

  qty := math.Min(f.CreateRate, f.outBuff.Space())
  createRes(qty)
}

func (f *Fac) approveOffers() {
  if f.queuedOrders == nil {
    f.queuedOrders = []Resource{}
  }
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
  qty = math.Min(qty, f.inBuff.Quantity())

  now := int64(f.Sim.Eng.SinceStart())
  rem := (now + int64(f.ConvertOffset)) % int64(f.ConvertPeriod)
  if qty <= rsrc.EPS {
    return
  } else if rem > 0 {
    return
  }

  rs := f.inBuff.PopQty(qty)
  if f.InUnits == f.OutUnits {
    f.outBuff.PushAll(rs)
  } else {
    createRes(qty)
  }
}

func (f *Fac) Receive(m *msg.Message) {
  if f.queuedOrders == nil {
    f.queuedOrders = []Resource{}
  }
  if msg.Sender == f {
    f.queuedOrders = append(f.queuedOrders, m)
  }
}

func check(err error) {
  if err != nil {
    panic(err.Error())
  }
}

func (f *Fac) RemoveResource(tran *Transaction) {
  rs, err := f.outBuff.PopQty(tran.Resource().Qty())
  check(err)
  tran.Manifest = rs
}

func (f *Fac) AddResource(tran *Transaction) {
  err := f.inBuff.PushAll(tran.Manifest)
  check(err)
}
