
package trans

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
  Res *Resource
  Sup Supplier
  Req Requester
  Manifest []*Resource
}

func NewOffer(sup Supplier) {
  return &Transaction{
    Type: Offer,
    Sup: sup,
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

func (t *Transaction) Clone() *Transaction {
  clone := *t
  clone.Res = t.Res.Clone()
  return &clone
}
