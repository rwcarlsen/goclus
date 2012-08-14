
package books

import "encoding/json"

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

// Dump needs more work - as it stands, it will likely cause infinite looping
func (b *Books) Dump(name string) error {
  data, err := json.Marshal(b)
  if err != nil {
    return err
  }

  f, err := os.Create(name)
  if err != nil {
    return err
  }
  defer f.Close()

  f.Write(data)
  return nil
}
