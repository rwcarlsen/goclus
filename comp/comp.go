
package comp

import "errors"

type Map map[int]float64

func cloneMap(m Map) Map {
  clone := Map{}
  for key, val := range m {
    clone[key] = val
  }
  return clone
}

// Composition is an immutable representation of nuclear material composition
type Composition struct {
  comp Map
}

func New(m Map) *Composition {
  return &Composition{comp: m}
}

func (c *Composition) Clone() *Composition {
  return &Composition{comp: cloneMap(c.comp)}
}

// Mix creates a new composition by combining the composition and other where
// ratio is qty of the composition divided by the qty of other.
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

