package fac

import (
	"fmt"
	"github.com/rwcarlsen/goclus/rsrc"
	"github.com/rwcarlsen/goclus/rsrc/buff"
	"github.com/rwcarlsen/goclus/sim"
	"github.com/rwcarlsen/goclus/trans"
	"math"
	"time"
)

type Fac struct {
	sim.Agenty
	queuedOrders []*sim.Message

	InCommod string
	InUnits  string
	InSize   float64
	inBuff   *buff.Buffer

	OutCommod string
	OutUnits  string
	OutSize   float64
	outBuff   *buff.Buffer

	CreateRate    float64
	ConvertAmt    float64
	ConvertPeriod time.Duration
	ConvertOffset time.Duration
	eng           *sim.Engine
}

func (f *Fac) Start(e *sim.Engine) {
	f.eng = e
	f.inBuff = &buff.Buffer{}
	f.inBuff.SetCapacity(f.InSize)
	f.outBuff = &buff.Buffer{}
	f.outBuff.SetCapacity(f.OutSize)
}

func (f *Fac) Tick() {
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

	mkt, _ := f.eng.GetService(commod)
	m := sim.NewMsg(f, mkt)
	m.Trans = tran
	m.SendOn()
}

func (f *Fac) Tock() {
	if f.ConvertPeriod == 0 {
		f.ConvertPeriod = f.eng.Step
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
	f.queuedOrders = sim.MsgGroup{}
}

func (f *Fac) createRes(qty float64) {
	if qty < rsrc.EPS {
		return
	}
	r := rsrc.NewGeneric(qty, f.OutUnits)
	f.outBuff.Push(r)
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
		f.outBuff.Push(rs...)
	} else {
		f.createRes(qty)
	}
}

func (f *Fac) Receive(m *sim.Message) {
	if m.Sender() == f {
		f.queuedOrders = append(f.queuedOrders, m)
	}
}

func (f *Fac) RemoveResource(tran *trans.Transaction) {
	fmt.Println(f.Id(), " sending qty=", tran.Resource().Qty(), "of", f.OutCommod)
	rs, err := f.outBuff.PopQty(tran.Resource().Qty())
	check(err)
	tran.Manifest = rs
}

func (f *Fac) AddResource(tran *trans.Transaction) {
	fmt.Println(f.Id(), " getting qty=", tran.Resource().Qty(), "of", f.InCommod)
	err := f.inBuff.Push(tran.Manifest...)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}
