
package comp

import (
  "testing"
  "math"
)

var tests = []struct{
  ratio float64
  m1, m2, m3 Map
  // m3 = expected[m1.Mix(ratio, m2)]
}{
  {
    ratio: 0,
    m1: Map{92235: 1},
    m2: Map{92235: 1, 92238: 2},
    m3: Map{92235: 1, 92238: 2},
  },{
    ratio: 1,
    m1: Map{92235: 1},
    m2: Map{92238: 2},
    m3: Map{92235: 1, 92238: 1},
  },{
    ratio: 1,
    m1: Map{92235: 1},
    m2: Map{92238: 2},
    m3: Map{92235: 0.5, 92238: 0.5},
  },{
    ratio: 1,
    m1: Map{92235: 1, 92238:3},
    m2: Map{92235: 2, 92238:3},
    m3: Map{92235: 0.65, 92238: 1.35},
  },{
    ratio: 2,
    m1: Map{92235: 1, 92238:3},
    m2: Map{92235: 2, 92238:3},
    m3: Map{92235: 0.9, 92238: 2.1},
  },{
    ratio: -2,
    m1: Map{92235: 1, 92238:3},
    m2: Map{92235: 2, 92238:3},
    m3: Map{92235: 0.1, 92238: 0.9},
  },{
    ratio: -1,
    m1: Map{92235: 1, 92238:3},
    m2: Map{92235: 2, 92238:3},
  },
}

func TestMapClone(t *testing.T) {
}

func TestCompClone(t *testing.T) {
}

func TestNew(t *testing.T) {
}

func TestNormalize(t *testing.T) {
}

func TestMix(t *testing.T) {
  for i, test := range tests {
    c1 := New(test.m1)
    c2 := New(test.m2)
    var c3 *Composition
    if test.m3 != nil {
      c3 = New(test.m3)
    }

    c4, err := c1.Mix(test.ratio, c2)
    if c3 != nil && err != nil {
      t.Fatalf("test %v threw error: %v, but expected nil.", i+1, err)
    } else if c3 == nil && err == nil {
      t.Fatalf("test %v should return non-nil error, but didn't", i+1)
    } else if c3 == nil && err != nil {
      return
    }

    c3.Normalize()
    c4.Normalize()
    want := c3.comp
    got := c4.comp

    for iso, v := range want {
      if v != math.Nextafter(got[iso], v) {
        t.Errorf("test %v failed for iso=%v: want %v, got %v", i+1, iso, v, got[iso])
      }
    }
  }
}

