
// Package rsrc provides a generalized resource interface and a generic resource type.
package rsrc

const (
  // EPS is an effective quantity precision - quantities and deviations
  // smaller than EPS should be ignored.
  EPS = 1e-6 
  // INFINITY is a number that can be used when infinite
  // quantities need to be represented.
  INFINITY = 1e25
)

// Resource is an interface that must be implemented by all transactable
// resources.
type Resource interface {
  Type() string
  Units() string
  Qty() float64
  SetQty(float64)
  Clone() Resource
}

