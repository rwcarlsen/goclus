// Package comp contains types for manipulation of nuclear material compositions.
package comp

import "errors"
import "math"

type Map map[int]float64

// Clone creates and returns a copy of the map.
func (m Map) Clone() Map {
	c := Map{}
	for key, val := range m {
		c[key] = val
	}
	return c
}

// normalize makes the sum of map's elements add to zero.
func (m Map) normalize() {
  var tot float64 = 0
  for _, val := range m {
    tot += val
  }
  for iso, _ := range m {
    m[iso] /= tot
  }
}

// Composition is an immutable representation of nuclear material
// composition.
// Assigning compositions to new variables is cheap, the internally
// maintained composition information is not duplicated.  If a copy is
// neaded, use the Clone method.
type Composition struct {
	comp Map
}

// New creates a new composition from m.
// Note that any modifications to m after it has been passed to a
// composition will be visible to the composition object.
func New(m Map) *Composition {
  var tot float64 = 0
  for _, val := range m {
    tot += val
  }

  comp := m.Clone()
  comp.normalize()
  return &Composition{comp: comp}
}

// Clone returns a copy of the composition.
func (c *Composition) Clone() *Composition {
  return &Composition{comp: c.comp.Clone()}
}

// Partial returns a comp map from the composition containing only the
// listed isotopes in ratios as they occur in the composition.  frac is the
// total fraction of the composition that is composed of the listed
// isotopes.
//
// A composition that results from removing rfrac of U235 and U238 could be
// obtained as follows:
//
//    part, frac := c1.Partial(92235, 92238)
//    thinned, err := c1.Mix(-1/(frac*rfrac), part)
func (c *Composition) Partial(isos ...int) (part *Composition, frac float64) {
  m := Map{}
  for _, iso := range isos {
    qty := c.comp[iso]
    frac += qty
    m[iso] = qty
  }
  return New(m), frac
}

// Mix creates a new composition by combining the composition and other where
// ratio is the quantity of the composition divided by the quantity of other.
// A negative ratio implies subtracting/removal of other from the composition.
func (c *Composition) Mix(ratio float64, other *Composition) (*Composition, error) {
	if ratio == 0 || c == other {
		return other, nil
	}

	mcomp := c.comp.Clone()
  for key, _ := range mcomp {
    mcomp[key] *= math.Abs(ratio)
  }

	if ratio > 0 {
		for key, qty := range other.comp {
			mcomp[key] += qty
		}
	} else {
		for key, qty := range other.comp {
			if mcomp[key] < qty {
				return nil, errors.New("comp: Mix ratio results in negative component")
			}
			mcomp[key] -= qty
		}
	}
	return New(mcomp), nil
}
