
package rsrc

type Resource interface {
  Type() string
  Units() string
  Qty() float64
  SetQty(float64)
  Clone() Resource
}

