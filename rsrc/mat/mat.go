
package mat

type Material struct {
  comp *Composition
  mass float64
}

func New(mass float64, comp *Composition) *Material {
  return &Material{
    comp: comp,
    mass: qty,
  }
}

func (m *Material) Type() string {
  return "Material"
}

func (m *Material) Units() string {
  return "kg"
}

func (m *Material) Qty() float64 {
  return mass
}

func (m *Material) SetQty(mass float64) {
  m.mass = mass
}

func (m *Material) Clone() Resource {
  clone := *m
  return &clone
}

func (m *Material) Extract(mass float64) *Material, error {
  if mass > m.mass {
    return nil, errors.New("rsrc: extraction amount too large")
  }

  cut := NewMaterial(mass, m.comp)
  m.mass -= mass
  return cut, nil
}

func (m *Material) Extract(mass float64) *Material, error {
  if mass > m.mass {
    return nil, errors.New("rsrc: extraction amount too large")
  }

  cut := NewMaterial(mass, m.comp)
  m.mass -= mass
  return cut, nil
}

func (m *Material) Absorb(other *Material) {
  if mass > m.mass {
    return nil, errors.New("rsrc: extraction amount too large")
  }
  m.comp := MixedComp(m.mass / other.mass, m.comp, other.comp)
  m.mass += other.mass
}

type Composition struct {

}

// Mix adjusts the composition by combining it with other where ratio is qty of
// the comp divided by the qty of other.
//
// Negative ratios imply subtracting/removal.
func (c *Composition) Mix(ratio float64, other *Composition) {
  
}

