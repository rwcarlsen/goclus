
package rsrc

type Generic struct {
  units string
  qty float64
}

func NewGeneric(qty float64, units string) *Generic {
  return &Generic{
    units: units,
    qty: qty,
  }
}

func (g *Generic) Type() string {
  return "Generic"
}

func (g *Generic) Units() string {
  return g.units
}

func (g *Generic) Qty() float64 {
  return g.qty
}

func (g *Generic) SetQty(qty float64) {
  g.qty = qty
}

func (g *Generic) Clone() Resource {
  clone := *g
  return &clone
}
