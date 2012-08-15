
package buff

import (
  "errors"
  "github.com/rwcarlsen/goclus/rsrc"
)

var (
  OverCapErr = errors.New("buff: cannot hold more than its capacity")
  TooSmallErr = errors.New("buff: operation results in negligible quantities")
)

type Buffer struct {
  capacity float64
  res []rsrc.Resource
}

func (b *Buffer) Capacity() float64 {
  return b.capacity
}

func (b *Buffer) SetCapacity(capacity float64) error {
  if b.Qty() - b.capacity > rsrc.EPS {
    return OverCapErr
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
  return b.capacity - b.Qty()
}

func (b *Buffer) PopQty(qty float64) ([]rsrc.Resource, error) {
  if qty - b.Qty() > rsrc.EPS || qty < rsrc.EPS {
    return nil, TooSmallErr
  }

  left := qty
  popped := []rsrc.Resource{}
  for left > rsrc.EPS {
    r := b.res[0]
    b.res = b.res[1:]
    quan := r.Qty()
    if quan - left > rsrc.EPS {
      leftover := r.Clone()
      leftover.SetQty(quan - left)
      r.SetQty(left)
      b.res = append([]rsrc.Resource{leftover}, b.res...)
    }
    popped = append(popped, r)
    left -= quan
  }
  return popped, nil
}

func (b *Buffer) PopN(num int) ([]rsrc.Resource, error) {
  if len(b.res) < num {
    return nil, TooSmallErr
  }
  popped := b.res[:num]
  b.res = b.res[num:]
  return popped, nil
}

func (b *Buffer) PopOne() (rsrc.Resource, error) {
  if len(b.res) < 1 {
    return nil, TooSmallErr
  }
  popped := b.res[0]
  b.res = b.res[1:]
  return popped, nil
}

func (b *Buffer) PushOne(r rsrc.Resource) error {
  if r.Qty() - b.Space() > rsrc.EPS {
    return OverCapErr
  }
  b.res = append(b.res, r)
  return nil
}

func (b *Buffer) PushAll(rs []rsrc.Resource) error {
  var tot float64
  for _, r := range rs {
    tot += r.Qty()
  }
  if tot - b.Space() > rsrc.EPS {
    return OverCapErr
  }
  b.res = append(b.res, rs...)
  return nil
}

