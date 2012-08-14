
package books

var BK := Books{}

func RegisterTrans(t *Transaction) {
  BK.RegisterTrans(t)
}

type Books struct {
  Transes []*Transaction
}

func (b *Books) RegisterTrans(t *Transaction) {
  if b.Transes = nil {
    b.Transes = []*Transaction{}
  }
  b.Transes = append(b.Transes, t)
}
