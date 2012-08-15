
package trans

import (
  "errors"
  "github.com/rwcarlsen/goclus/rsrc"
)

type TransType int

const (
  Offer TransType = iota
  Request
)

type Supplier interface {
  RemoveResource(*Transaction)
}

type Requester interface {
  AddResource(*Transaction)
}

type Transaction struct {
  Type TransType
  res rsrc.Resource
  Sup Supplier
  Req Requester
  Manifest []rsrc.Resource
}

func NewOffer(sup Supplier) *Transaction {
  return &Transaction{
    Type: Offer,
    Sup: sup,
  }
}

func NewRequest(req Requester) *Transaction {
  return &Transaction{
    Type: Request,
    Req: req,
  }
}

func (t *Transaction) MatchWith(other *Transaction) error {
  if t.Type == other.Type {
    return errors.New("trans: Non-complementary transaction types")
  }

  if t.Type == Offer {
    t.Req = other.Req
    other.Sup = t.Sup
  } else {
    t.Sup = other.Sup
    other.Req = t.Req
  }
  return nil
}

func (t *Transaction) Approve() {
  t.Sup.RemoveResource(t)
  t.Req.AddResource(t)
}

func (t *Transaction) Resource() rsrc.Resource {
  return t.res
}

func (t *Transaction) SetResource(r rsrc.Resource) {
  t.res = r.Clone()
}

func (t *Transaction) Clone() *Transaction {
  clone := *t
  clone.res = t.res.Clone()
  return &clone
}
