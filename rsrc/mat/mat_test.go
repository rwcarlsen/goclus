
package mat

import (
  "testing"
  "github.com/rwcarlsen/goclus/util/assert"
  "github.com/rwcarlsen/goclus/comp"
)

const (
  zero float64 = 0
  qty1 float64 = 5
  qty2 float64 = 4
  qty3 float64 = 2.31
)

func mat1() *Material {
  cm := comp.Map{922350:1.1, 922380:3.9}
  cmp := comp.New(cm)
  return New(qty1, cmp)
}

func mat2() *Material {
  cm := comp.Map{922350:1}
  cmp := comp.New(cm)
  return New(qty2, cmp)
}

func mat3() *Material {
  cm := comp.Map{922380:2.3, 942390:0.01}
  cmp := comp.New(cm)
  return New(qty3, cmp)
}

func TestType(t *testing.T) {
  assert.Eq(t, mat1().Type(), Type)
}

func TestUnits(t *testing.T) {
  assert.Eq(t, mat1().Units(), "kg")
}

func TestQty(t *testing.T) {
  m := mat1()
  assert.Eq(t, m.Qty(), qty1)

  var qty float64 = 3
  m.SetQty(qty)
  assert.Eq(t, m.Qty(), qty)
}

func TestClone(t *testing.T) {
  m := mat1()
  clone := m.Clone().(*Material)
  assert.Eq(t, m.Comp, clone.Comp)
}

func TestExtractMass_Good(t *testing.T) {
  m := mat1()
  cmp := m.Comp
  var qty float64 = 4.9
  e, err := m.ExtractMass(qty)

  assert.NoErr(t, err).Fatal()
  assert.Eq(t, m.Qty(), qty1 - qty)
  assert.Eq(t, e.Qty(), qty)
  assert.Eq(t, m.Comp, cmp)
  assert.Eq(t, e.Comp, cmp)
}

func TestExtractMass_Bad(t *testing.T) {
  m := mat1()
  cmp := m.Comp
  var qty float64 = 5.1
  _, err := m.ExtractMass(qty)

  assert.Err(t, err).Fatal()
  assert.Eq(t, m.Qty(), qty1)
  assert.Eq(t, m.Comp, cmp)
}

func TestExtractComp_Good(t *testing.T) {
  m := mat1()
  mcmp := m.Comp
  ecmp := mat2().Comp

  var qty float64 = 1.0
  e, err := m.ExtractComp(qty, ecmp)

  assert.NoErr(t, err).Fatal()
  assert.Eq(t, m.Qty(), qty1 - qty)
  assert.Eq(t, e.Qty(), qty)
  assert.Ne(t, m.Comp, mcmp)
  assert.Eq(t, e.Comp, ecmp)
}

func TestExtractComp_Bad(t *testing.T) {
  m := mat1()
  mcmp := m.Comp
  ecmp := mat2().Comp

  var qty float64 = 1.2
  _, err := m.ExtractComp(qty, ecmp)

  assert.Err(t, err).Fatal()
  assert.Eq(t, m.Qty(), qty1)
  assert.Eq(t, m.Comp, mcmp)
}

func TestAbsorb(t *testing.T) {
  m1 := mat1()
  cmp := m1.Comp
  m3 := mat3()
  m1.Absorb(m3)

  assert.Eq(t, m1.Qty(), qty1 + qty3)
  assert.Eq(t, m3.Qty(), zero)
  assert.Ne(t, m1.Comp, cmp)
}

