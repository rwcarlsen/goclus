
package comp

import (
  "testing"
)

const (
  zero float64 = 0
  qty1 float64 = 5
  qty2 float64 = 4
  qty3 float64 = 2.31
)

var tests = []struct{
  ratio float64
  m1 Map
  m2 Map
  m3 Map // m3 = expected[m1.Mix(ratio, m2)]
}{
  {
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
  },
}

func TestMapClone(t *testing.T) {
}

func TestCompClone(t *testing.T) {
}

func TestMix_Good(t *testing.T) {
  for i, test := range tests {
    c1 := New(test.m1)
    c2 := New(test.m2)
    c3 := New(test.m3)

    c4, err := c1.Mix(test.ratio, c2)
    if err != nil {
      t.Errorf("test %v threw error: %v, but expected nil.", i+1, err)
    }

    c3.Normalize()
    c4.Normalize()
    want := c3.comp
    got := c4.comp

    for iso, v := range want {
      if v != got[iso] {
        t.Errorf("test %v failed: want %v, got %v", i+1, want, got)
        break
      }
    }
  }
}

func TestMix_Bad(t *testing.T) {
}
