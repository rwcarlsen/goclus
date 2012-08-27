
package mat

import (
  "testing"
  "github.com/rwcarlsen/goclus/comp"
)

const qty1 float64 = 5

func mat1() *Material {
  cm := comp.Map{92235:1.1, 92238:3.9}
  cmp := comp.New(cm)
  return New(qty1, cmp)
}

func mat2() *Material {
  qty := 4.0
  cm := comp.Map{92235:1}
  cmp := comp.New(cm)
  return New(qty, cmp)
}

func mat3() *Material {
  qty := 2.31
  cm := comp.Map{92238:2.3, 94239:0.01}
  cmp := comp.New(cm)
  return New(qty, cmp)
}

func TestType(t *testing.T) {
  m := mat1()
  if m.Type() != Type {
    t.Fatalf("Type()=%v, want %v", m.Type(), Type)
  }
}

func TestUnits(t *testing.T) {
  m := mat1()
  if m.Units() != "kg" {
    t.Fatalf("Units()=%v, want %v", m.Units(), "kg")
  }
}

func TestQty(t *testing.T) {
  m := mat1()
  if m.Qty() != qty1 {
    t.Errorf("Qty()=%v, want %v", m.Qty(), qty1)
  }

  var qty float64 = 3
  m.SetQty(qty)
  if m.Qty() != qty {
    t.Errorf("Qty()=%v, want %v", m.Qty(), qty)
  }
}

func TestClone(t *testing.T) {
  m := mat1()
  clone := m.Clone().(*Material)
  if m.Comp != clone.Comp {
    t.Fatal("Clone should have same composition, but doesn't")
  }
}

func TestExtractMass_Good(t *testing.T) {
  m := mat1()
  cmp := m.Comp
  var qty float64 = 4.9
  e, err := m.ExtractMass(qty)
  if err != nil {
    t.Fatal("Extraction err: ", err)
  }

  if m.Qty() != qty1 - qty {
    t.Errorf("m.Qty()=%v, want %v", m.Qty(), qty1 - qty)
  }
  if e.Qty() != qty {
    t.Errorf("e.Qty()=%v, want %v", e.Qty(), qty)
  }
  if m.Comp != cmp {
    t.Errorf("m should have same composition, but doesn't")
  }
  if e.Comp != cmp {
    t.Errorf("e should have same composition as m, but doesn't")
  }
}

func TestExtractMass_Bad(t *testing.T) {
  m := mat1()
  cmp := m.Comp
  var qty float64 = 5.1
  _, err := m.ExtractMass(qty)
  if err == nil {
    t.Fatal("Expected extraction err but got nil")
  }

  if m.Qty() != qty1 {
    t.Errorf("m.Qty()=%v, want %v", m.Qty(), qty1)
  }
  if m.Comp != cmp {
    t.Errorf("m should have same composition, but doesn't")
  }
}

func TestExtractComp_Good(t *testing.T) {
  m := mat1()
  mcmp := m.Comp
  ecmp := mat2().Comp

  var qty float64 = 1.0
  e, err := m.ExtractComp(qty, ecmp)
  if err != nil {
    t.Fatal("Extraction err: ", err)
  }

  if m.Qty() != qty1 - qty {
    t.Errorf("m.Qty()=%v, want %v", m.Qty(), qty1 - qty)
  }
  if e.Qty() != qty {
    t.Errorf("e.Qty()=%v, want %v", e.Qty(), qty)
  }
  if m.Comp == mcmp {
    t.Errorf("m should have changed composition, but doesn't")
  }
  if e.Comp != ecmp {
    t.Errorf("e should have extracted composition, but doesn't")
  }
}

func TestExtractComp_Bad(t *testing.T) {
  m := mat1()
  mcmp := m.Comp
  ecmp := mat2().Comp

  var qty float64 = 1.2
  _, err := m.ExtractComp(qty, ecmp)
  if err == nil {
    t.Fatal("Expected extraction err but got nil")
  }

  if m.Qty() != qty1 {
    t.Errorf("m.Qty()=%v, want %v", m.Qty(), qty1)
  }
  if m.Comp != mcmp {
    t.Errorf("m should have same composition, but doesn't")
  }
}

func TestAbsorb(t *testing.T) {
}

