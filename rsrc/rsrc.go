
package rsrc

type Interface interface {
  Type() string
  Units() string
  Name()
}

type resource interface {
  id int
  res Interface
}

func (r *resource) ID() int {
  return r.id
}

var id := 0

func New(r Interface) *resource {
  id++
  return &resource{
    id: id,
    res: r,
  }
}

