
package comp

import "errors"

type Map map[int]float64

// Clone creates and returns a copy of the map.
func (m *Map) Clone() Map {
  c := Map{}
  for key, val := range m {
    c[key] = val
  }
  return c
}

// Composition is an immutable representation of nuclear material
// composition.
//
// Assigning compositions to new variables is cheap, the internally
// maintained composition information is not duplicated.  If a copy is
// neaded, use the Clone method.
type Composition struct {
  comp Map
}

// New creates a new composition from m.
//
// Note that any modifications to m after it has been passed to a
// composition will be visible to the composition object.
func New(m Map) *Composition {
  return &Composition{comp: m}
}

// Clone returns a copy of the composition.
func (c *Composition) Clone() *Composition {
  return &Composition{comp: c.comp.Clone()}
}

// Mix creates a new composition by combining the composition and other where
// ratio is the quantity of the composition divided by the quantity of other.
//
// A negative ratio implies subtracting/removal of other from the composition.
func (c *Composition) Mix(ratio float64, other *Composition) (*Composition, error) {
  if ratio == 0 {
    return other, nil
  }

  mixed := c.Clone()
  if ratio > 0 {
    for key, qty := range other.comp {
      mixed.comp[key] *= ratio
      mixed.comp[key] += qty
    }
  } else {
    for key, qty := range other.comp {
      mixed.comp[key] *= -1 * ratio
      if mixed.comp[key] < qty {
        return nil, errors.New("comp: Mix ratio results in negative component")
      }
      mixed.comp[key] -= qty
    }
  }
  return mixed, nil
}

