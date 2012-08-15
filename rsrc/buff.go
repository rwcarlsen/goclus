
package rsrc

const EPS = 1e-6

type Buffer struct {
  unlimited bool
  capacity float64
  res []Resource
}

func (b *Buffer) init() {
  if res == nil {
    res = []Resource{}
  }
}

func (b *Buffer) Capacity() float64 {
  if b.unlimited {
    return -1
  }
  return b.capacity
}

func (b *Buffer) SetCapacity(capacity float64) error {
  if b.Qty() - b.capacity > EPS {
    return errors.New("rsrc: New buffer capacity lower than existing quantity")
  }
  b.capacity = capacity
  return nil
}

func (b *Buffer) Count() int {
  return len(b.res)
}

func (b *Buffer) Qty() float64 {
  var tot float64
  for _, r := range b.res {
    tot += r.Qty()
  }
  return tot
}

func (b *Buffer) Space() float64 {
  if b.unlimited {
    return -1
  }
  return b.capacity - b.Qty()
}

func (b *Buffer) IsUnlimited() bool {
  return b.unlimited
}

func (b *Buffer) MakeUnlimited() {
  b.unlimited = true
}

func (b *Buffer) MakeLimited(capacity float64) error {
  err := b.SetCapacity(capacity)
  if err != nil {
    return err
  }
  b.unlimited = false
}

func (b *Buffer) PopQty(qty float64) ([]Resource, error) {
  
}

func (b *Buffer) PopN(num int) ([]Resource, error) {
  
}

func (b *Buffer) PopOne() (Resource, error) {
  
}

func (b *Buffer) PushOne(r Resource) error {
  
}

func (b *Buffer) PushAll(rs []Resource) error {
  
}

