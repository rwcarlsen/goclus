package rsrc

type generic struct {
	units string
	qty   float64
}

// NewGeneric returns a new generic resource initialized with qty of the given
// units.
// Note that the specified units will be immutable.
func NewGeneric(qty float64, units string) *generic {
	return &generic{
		units: units,
		qty:   qty,
	}
}

// Type returns "Generic" for all generic resources.
func (g *generic) Type() string {
	return "Generic"
}

// Units returns the units specified at creation.
func (g *generic) Units() string {
	return g.units
}

// Qty returns the quantity the resource holds of [units].
func (g *generic) Qty() float64 {
	return g.qty
}

// SetQty changes the resources quantity to qty.
func (g *generic) SetQty(qty float64) {
	g.qty = qty
}

// Clone returns a deep-copy of the generic resource.
func (g *generic) Clone() Resource {
	clone := *g
	return &clone
}
