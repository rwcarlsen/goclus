package econ

import (
	"github.com/rwcarlsen/goclus/sim"
	"github.com/rwcarlsen/goclus/trans"
	"time"
)

type Econ struct {
	id      string
	req     map[sim.Agent]map[time.Time]float64
	off     map[sim.Agent]map[time.Time]float64
	times   []time.Time
	Tracked []int
	eng     *sim.Engine
}

func (e *Econ) SetId(id string) {
	e.id = id
}

func (e *Econ) Id() string {
	return e.id
}

func (e *Econ) Start(eng *sim.Engine) {
	e.eng = eng
}

func (e *Econ) MsgNotify(m *sim.Message) {
	if !e.isTracked(m.Owner) || m.Trans == nil {
		return
	}

	qty := m.Trans.Resource().Qty()
	if m.Trans.Type() == trans.Offer {
		e.off[m.Owner][e.eng.Time()] += qty
	} else {
		e.req[m.Owner][e.eng.Time()] += qty
	}
	e.times = append(e.times, e.eng.Time())
}

func (e *Econ) isTracked(a sim.Agent) bool {
	for _, id := range e.Tracked {
		if a.Id() == id {
			return true
		}
	}
	return false
}

func (e *Econ) OfferQty(id string) (float64, error) {
	mkt, err := e.eng.GetService(id)
	if err != nil {
		return 0, err
	}

	if len(e.times) == 1 {
		return 0, nil
	}
	prev := e.times[len(e.times)-2]

	return e.off[mkt][prev], nil
}

func (e *Econ) RequestQty(id string) (float64, error) {
	mkt, err := e.eng.GetService(id)
	if err != nil {
		return 0, err
	}

	if len(e.times) == 1 {
		return 0, nil
	}
	prev := e.times[len(e.times)-2]

	return e.off[mkt][prev], nil
}

func (e *Econ) UnmetDemand(id string) (float64, error) {
	r, _ := e.RequestQty(id)
	o, err := e.OfferQty(id)
	if err != nil {
		return 0, err
	}
	return r - o, nil
}
