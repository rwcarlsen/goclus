// Package mat provides a nuclear-material resource.
package mat

import (
	"errors"
	"github.com/rwcarlsen/goclus/comp"
	"github.com/rwcarlsen/goclus/rsrc"
)

const Type = "Material"

// Material is a resource for tracking, and manipulating nuclear materials.
type Material struct {
  // Comp represents the nuclear composition of the material.
	Comp *comp.Composition
	qty float64
}

// New creates and returns a new material of the given qty with
// composition Comp.
// Note that the specified units will be immutable.
func New(qty float64, Comp *comp.Composition) *Material {
	return &Material{
		Comp: Comp,
		qty: qty,
	}
}

// Type returns "Material" for all material resources.
func (m *Material) Type() string {
	return Type
}

// Units returns the units specified at creation.
func (m *Material) Units() string {
	return "kg"
}

// Qty returns the quantity of the material in [units].
func (m *Material) Qty() float64 {
	return m.qty
}

// SetQty changes the material's quantity to qty.
func (m *Material) SetQty(qty float64) {
	m.qty = qty
}

// Clone returns a shallow-copy of the material.
func (m *Material) Clone() rsrc.Resource {
	clone := *m
	return &clone
}

// ExtractMass creates and returns a new, compositionally identical
// material of the given qty.
// An error is returned if the extraction would result in a negative qty
// remaining in the material.
func (m *Material) ExtractMass(qty float64) (*Material, error) {
	if qty > m.qty {
		return nil, errors.New("rsrc: extraction amount too large")
	}

	cut := New(qty, m.Comp)
	m.qty -= qty
	return cut, nil
}

// ExtractComp creates and returns a new material of the given qty with composition
// Comp by extracting the corresponding amounts from this material.
// An error is returned if the extraction would result in a negative qty
// remaining in the material.
//
// A material that results from removing rfrac of U235 and U238 could be
// obtained as follows:
//
//     c, frac := m1.Comp.Partial(92235)
//     extracted, err := m1.ExtractComp(m1.Qty() * frac * rfrac, c)
func (m *Material) ExtractComp(qty float64, comp *comp.Composition) (*Material, error) {
	newComp, err := m.Comp.Mix(-m.qty/qty, comp)
	if err != nil {
		return nil, err
	}

	m.Comp = newComp
	m.qty -= qty

	return New(qty, comp), nil
}

// Absorb adds/combines other into the material.
func (m *Material) Absorb(other *Material) {
  if other.Comp != m.Comp {
    m.Comp, _ = m.Comp.Mix(m.qty/other.qty, other.Comp)
  }
	m.qty += other.qty
  other.qty = 0
}
