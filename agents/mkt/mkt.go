package mkt

import (
	"github.com/rwcarlsen/goclus/rsrc"
	"github.com/rwcarlsen/goclus/sim"
	"github.com/rwcarlsen/goclus/trans"
	"math/rand"
)

type Mkt struct {
	sim.Agenty
	Shuffle  bool
	Seed     int64
	offers   sim.MsgGroup
	requests sim.MsgGroup
}

func (m *Mkt) Receive(mg *sim.Message) {
	if mg.Trans.Type() == trans.Offer {
		m.offers = append(m.offers, mg)
	} else {
		m.requests = append(m.requests, mg)
	}
}

func (m *Mkt) Resolve() {
	if m.Shuffle {
		shuffle(m.offers)
		shuffle(m.requests)
	}

	var matched sim.MsgGroup
	for len(m.requests) > 0 && len(m.offers) > 0 {
		mg := m.requests[0]
		qty := mg.Trans.Resource().Qty()
		m.offers, matched = m.extractQty(m.offers, qty)
		m.matchAll(matched, mg)
		m.requests = m.requests[1:]
	}

	m.offers = sim.MsgGroup{}
	m.requests = sim.MsgGroup{}
}

func (m *Mkt) matchAll(group sim.MsgGroup, mg *sim.Message) {
	for _, gpMem := range group {
		err := gpMem.Trans.MatchWith(mg.Trans)
		if err != nil {
			panic(err.Error())
		}
		gpMem.Dir = sim.DownMsg
		gpMem.SendOn()
	}
}

func (m *Mkt) extractQty(group sim.MsgGroup, qty float64) (orig, extracted sim.MsgGroup) {
	unmet := qty
	for len(group) > 0 && unmet >= rsrc.EPS {
		currMsg := group[0]
		currQty := currMsg.Trans.Resource().Qty()
		if currQty <= unmet+rsrc.EPS {
			extracted = append(extracted, currMsg)
			group = group[1:]
			unmet -= currQty
		} else {
			split := m.extractFromMsg(currMsg, unmet)
			extracted = append(extracted, split)
			unmet = 0
		}
	}
	return group, extracted
}

func (m *Mkt) extractFromMsg(mg *sim.Message, qty float64) *sim.Message {
	extracted := mg.Clone()
	extracted.Trans.Resource().SetQty(qty)

	remainder := mg.Trans.Resource().Qty() - qty
	mg.Trans.Resource().SetQty(remainder)

	return extracted
}

func shuffle(gp sim.MsgGroup) {
	inds := rand.Perm(len(gp))
	orig := make(sim.MsgGroup, len(gp))
	copy(orig, gp)

	for i, ind := range inds {
		gp[ind] = orig[i]
	}
}
