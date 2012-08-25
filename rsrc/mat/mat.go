// Package mat provides a nuclear-material resource.
package mat

import (
	"errors"
	"github.com/rwcarlsen/goclus/comp"
	"github.com/rwcarlsen/goclus/rsrc"
)

// Material is a resource for tracking, and manipulating nuclear materials.
type Material struct {
	cmp *comp.Composition
	qty float64
}

// New creates and returns a new material of the given qty with
// composition cmp.
// Note that the specified units will be immutable.
func New(qty float64, cmp *comp.Composition) *Material {
	return &Material{
		cmp: cmp,
		qty: qty,
	}
}

// Type returns "Material" for all material resources.
func (m *Material) Type() string {
	return "Material"
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

	cut := New(qty, m.cmp)
	m.qty -= qty
	return cut, nil
}

// ExtractComp creates and returns a new material of the given qty with composition
// cmp by extracting the corresponding amounts from this material.
// An error is returned if the extraction would result in a negative qty
// remaining in the material.
func (m *Material) ExtractComp(qty float64, cmp *comp.Composition) (*Material, error) {
	if qty > m.qty {
		return nil, errors.New("rsrc: extraction amount too large")
	}

	newcmp, err := m.cmp.Mix(-m.qty/qty, cmp)
	if err != nil {
		return nil, err
	}

	m.cmp = newcmp
	m.qty -= qty

	return New(qty, cmp), nil
}

// Absorb adds/combines other into the material.
func (m *Material) Absorb(other *Material) {
	m.cmp, _ = m.cmp.Mix(m.qty/other.qty, other.cmp)
	m.qty += other.qty
}
