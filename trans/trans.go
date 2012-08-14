
package trans

type Supplier interface {
  RemoveResource(*Transaction) []*Resource
}

type Requester interface {
  AddResource(*Transaction, []*Resource)
}

type Transaction struct {
  Res *Resource
  Sup Supplier
  Req Requester
  manifest []Resource
}

func (t *Transaction) Manifest() {
  return t.manifest
}

func (t *Transaction) Approve() {
  manifest := t.Sup.RemoveResource(t)
  t.Req.AddResource(t, manifest)
}

func (t *Transaction) Clone() *Transaction {
  clone := *t
  clone.Res = t.Res.Clone()
  return &clone
}
