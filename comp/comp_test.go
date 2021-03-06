package comp

import (
	"math"
	"testing"
)

var tests = []struct {
	ratio      float64
	m1, m2, m3 Map
	// m3 = expected[m1.Mix(ratio, m2)]
}{
	{
		ratio: 0,
		m1:    Map{922350: 1},
		m2:    Map{922350: 1, 922380: 2},
		m3:    Map{922350: 1, 922380: 2},
	}, {
		ratio: 1,
		m1:    Map{922350: 1},
		m2:    Map{922380: 2},
		m3:    Map{922350: 1, 922380: 1},
	}, {
		ratio: 1,
		m1:    Map{922350: 1},
		m2:    Map{922380: 2},
		m3:    Map{922350: 0.5, 922380: 0.5},
	}, {
		ratio: 1,
		m1:    Map{922350: 1, 922380: 3},
		m2:    Map{922350: 2, 922380: 3},
		m3:    Map{922350: 0.65, 922380: 1.35},
	}, {
		ratio: 2,
		m1:    Map{922350: 1, 922380: 3},
		m2:    Map{922350: 2, 922380: 3},
		m3:    Map{922350: 0.9, 922380: 2.1},
	}, {
		ratio: -2,
		m1:    Map{922350: 1, 922380: 3},
		m2:    Map{922350: 2, 922380: 3},
		m3:    Map{922350: 0.1, 922380: 0.9},
	}, {
		ratio: -1,
		m1:    Map{922350: 1, 922380: 3},
		m2:    Map{922350: 2, 922380: 3},
	},
}

func TestMapClone(t *testing.T) {
}

func TestCompClone(t *testing.T) {
}

func TestNew(t *testing.T) {
}

func TestPartialMix(t *testing.T) {
	c := New(Map{922350: .1, 922380: .4, 942390: .5})
	part, frac := c.Partial(922350, 922380)
	rfrac := .5
	t.Log(part, ", ", frac)
	thinned, err := c.Mix(-1.0/(frac*rfrac), part)
	if err != nil {
		t.Errorf("error: %v", err)
	} else if thinned.comp[942390] != 2.0/3.0 {
		t.Errorf("942390 want %v got %v", 2.0/3.0, thinned.comp[942390])
	}
}

func TestMix(t *testing.T) {
	for i, test := range tests {
		c1 := New(test.m1)
		c2 := New(test.m2)
		c4, err := c1.Mix(test.ratio, c2)

		test.m3.normalize()
		want := test.m3

		if want != nil && err != nil {
			t.Fatalf("test %v threw error: %v, but expected nil.", i+1, err)
		} else if want == nil && err == nil {
			t.Fatalf("test %v should return non-nil error, but didn't", i+1)
		} else if want == nil && err != nil {
			return
		}

		got := c4.comp

		for iso, v := range want {
			if floatNe(v, got[iso]) {
				t.Errorf("test %v failed on iso=%v: want %v, got %v", i+1, iso, v, got[iso])
			}
		}
	}
}

func floatEq(a, b float64) bool {
	absTol := 0.0
	relTol := 1e-10
	return math.Abs(a-b) <= (absTol + relTol*math.Abs(b))
}

func floatNe(a, b float64) bool {
	absTol := 0.0
	relTol := 1e-10
	return math.Abs(a-b) > (absTol + relTol*math.Abs(b))
}
