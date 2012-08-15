
package rsrc

const (
  EPS = 1e-6
  INFINITY = 1e25
)

type Resource interface {
  Type() string
  Units() string
  Qty() float64
  SetQty(float64)
  Clone() Resource
}

