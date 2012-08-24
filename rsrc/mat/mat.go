
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
  mass float64
}

// New creates and returns a new material of the given mass with
// composition cmp.
func New(mass float64, cmp *comp.Composition) *Material {
  return &Material{
    cmp: cmp,
    mass: mass,
  }
}

func (m *Material) Type() string {
  return "Material"
}

func (m *Material) Units() string {
  return "kg"
}

func (m *Material) Qty() float64 {
  return m.mass
}

func (m *Material) SetQty(mass float64) {
  m.mass = mass
}

func (m *Material) Clone() rsrc.Resource {
  clone := *m
  return &clone
}

// ExtractMass creates and returns a new, compositionally identical
// material of the given mass.
//
// An error is returned if the extraction would result in a negative mass
// remaining in the material.
func (m *Material) ExtractMass(mass float64) (*Material, error) {
  if mass > m.mass {
    return nil, errors.New("rsrc: extraction amount too large")
  }

  cut := New(mass, m.cmp)
  m.mass -= mass
  return cut, nil
}

// ExtractComp creates and returns a new material of the given mass with composition
// cmp by extracting the corresponding amounts from this material.
//
// An error is returned if the extraction would result in a negative mass
// remaining in the material.
func (m *Material) ExtractComp(mass float64, cmp *comp.Composition) (*Material, error) {
  m.cmp, err = m.cmp.Mix(-m.mass / mass, cmp)
  if err != nil {
    return nil, err
  }

  return New(mass, cmp), nil
}

// Absorb adds/combines other into the material.
func (m *Material) Absorb(other *Material) {
  m.cmp, _ = m.cmp.Mix(m.mass / other.mass, other.cmp)
  m.mass += other.mass
}

