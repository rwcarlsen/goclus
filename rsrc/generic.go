
package rsrc

type generic struct {
  units string
  qty float64
}

// NewGeneric returns a new generic resource initialized with qty of the given
// units.
//
// Note that the specified units will be immutable.
func NewGeneric(qty float64, units string) *Generic {
  return &Generic{
    units: units,
    qty: qty,
  }
}

// Type returns "Generic" for all generic resources.
func (g *Generic) Type() string {
  return "Generic"
}

// Units returns the units specified at creation.
func (g *Generic) Units() string {
  return g.units
}

// Qty returns the quantity the resource holds of [units].
func (g *Generic) Qty() float64 {
  return g.qty
}

// SetQty changes the resources quantity to qty.
func (g *Generic) SetQty(qty float64) {
  g.qty = qty
}

// Clone returns a deep-copy of the generic resource.
func (g *Generic) Clone() Resource {
  clone := *g
  return &clone
}
