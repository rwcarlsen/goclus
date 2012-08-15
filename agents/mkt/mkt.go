
package mkt

import (
  "time"
  "math/rand"
  "github.com/rwcarlsen/goclus/msg"
  "github.com/rwcarlsen/goclus/rsrc"
  "github.com/rwcarlsen/goclus/trans"
)

type Group []*msg.Message

type Mkt struct {
  Shuffle bool
  Seed int64
  offers Group
  requests Group
}

func (m *Mkt) Receive(mg *msg.Message) {
  if mg.Trans.Type == trans.Offer {
    m.offers = append(m.offers, mg)
  } else {
    m.requests = append(m.requests, mg)
  }
}

func (m *Mkt) Resolve(tm *time.Duration) {
  if m.Shuffle {
    shuffle(m.offers)
    shuffle(m.requests)
  }

  for len(m.requests) + len(m.offers) > 0 {
    mg := m.requests[0]
    qty := mg.Trans.Resource().Qty()
    matched := m.extractQty(&m.offers, qty)
    m.matchAll(matched, mg)
    m.requests = m.requests[1:]
  }

  m.offers = Group{}
  m.requests = Group{}
}

func (m *Mkt) matchAll(group Group, mg *msg.Message) {
  for _, gpMem := range group {
    gpMem.Trans.MatchWith(mg.Trans)
    gpMem.Dir = msg.Down
    gpMem.SendOn()
  }
}

func (m *Mkt) extractQty(group *Group, qty float64) Group {
  var extracted Group
  unmet := qty
  for len(*group) > 0 && unmet >= rsrc.EPS {
    currMsg := (*group)[0]
    currQty := currMsg.Trans.Resource().Qty()
    if currQty <= unmet + rsrc.EPS {
      extracted = append(extracted, currMsg)
      group = &((*group)[1:])
      unmet -= currQty
    } else {
      split := m.extractFromMsg(currMsg, unmet)
      extracted = append(extracted, split)
      unmet = 0
    }
  }
  return extracted
}

func (m *Mkt) extractFromMsg(mg *msg.Message, qty float64) *msg.Message {
  extracted := mg.Clone()
  extracted.Trans.Resource().SetQty(qty)

  remainder := mg.Trans.Resource().Qty() - qty
  mg.Trans.Resource().SetQty(remainder)

  return extracted
}

func shuffle(gp Group) {
  inds := rand.Perm(len(gp))
  orig := make(Group, len(gp))
  copy(orig, gp)

  for i, ind := range inds {
    gp[ind] = orig[i]
  }
}

