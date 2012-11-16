package mkt

import (
	"github.com/rwcarlsen/goclus/msg"
	"github.com/rwcarlsen/goclus/rsrc"
	"github.com/rwcarlsen/goclus/sim"
	"github.com/rwcarlsen/goclus/trans"
	"math/rand"
)

type Mkt struct {
	msg.Commy
	sim.Agenty
	Shuffle  bool
	Seed     int64
	offers   msg.Group
	requests msg.Group
}

func (m *Mkt) Receive(mg *msg.Message) {
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

	var matched msg.Group
	for len(m.requests) > 0 && len(m.offers) > 0 {
		mg := m.requests[0]
		qty := mg.Trans.Resource().Qty()
		m.offers, matched = m.extractQty(m.offers, qty)
		m.matchAll(matched, mg)
		m.requests = m.requests[1:]
	}

	m.offers = msg.Group{}
	m.requests = msg.Group{}
}

func (m *Mkt) matchAll(group msg.Group, mg *msg.Message) {
	for _, gpMem := range group {
		err := gpMem.Trans.MatchWith(mg.Trans)
		if err != nil {
			panic(err.Error())
		}
		gpMem.Dir = msg.Down
		gpMem.SendOn()
	}
}

func (m *Mkt) extractQty(group msg.Group, qty float64) (orig, extracted msg.Group) {
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

func (m *Mkt) extractFromMsg(mg *msg.Message, qty float64) *msg.Message {
	extracted := mg.Clone()
	extracted.Trans.Resource().SetQty(qty)

	remainder := mg.Trans.Resource().Qty() - qty
	mg.Trans.Resource().SetQty(remainder)

	return extracted
}

func shuffle(gp msg.Group) {
	inds := rand.Perm(len(gp))
	orig := make(msg.Group, len(gp))
	copy(orig, gp)

	for i, ind := range inds {
		gp[ind] = orig[i]
	}
}
