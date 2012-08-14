
package trans

type Transaction struct {

}

func (t *Transaction) Clone() *Transaction {
  clone := *t

  // clone the resource

  return &clone
}
