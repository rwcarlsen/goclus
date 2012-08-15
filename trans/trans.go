
package trans

type TransType int

import "github.com/rwcarlsen/goclus/rsrc"

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
  res *Resource
  Sup Supplier
  Req Requester
  Manifest []rsrc.Resource
}

func NewOffer(sup Supplier) {
  return &Transaction{
    Type: Offer,
    Sup: sup,
  }
}

func NewRequest(req Requester) {
  return &Transaction{
    Type: Offer,
    Req: req,
  }
}

func (t *Transaction) MatchWith(other *Transaction) error {
  if t.Type != other.Type {
    return errors.New("trans: Incompatible transaction types")
  }

  if t.Type == Offer {
    t.Req = other.Req
    other.Sup = t.Sup
  } else {
    t.Sup = other.Sup
    other.Req = t.Req
  }
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
  clone.Res = t.Res.Clone()
  return &clone
}
