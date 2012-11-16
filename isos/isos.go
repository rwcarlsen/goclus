package isos

import "fmt"
import "errors"

var info map[Iso]*Info
var groups map[string][]Iso

func init() {
	info = make(map[Iso]*Info)
	groups = make(map[string][]Iso)
}

func Load(path string) {
	// ...
}

func AddGroup(name string, isos ...Iso) {
	groups[name] = isos
}

func Group(name string) []Iso {
	return groups[name]
}

type Iso int

// Z returns the isotope's atomic number
func (i Iso) Z() int {
	return int(i) / 10000
}

// A returns the isotope's atomic weight
func (i Iso) A() int {
	return (int(i) - i.Z()) / 10
}

// Is returns the isotope's isomeric state. IS=0 for ground state
func (i Iso) Is() int {
	return int(i) - i.Z() - i.A()
}

func (i Iso) Info() (*Info, error) {
	if in, ok := info[i]; ok {
		return in, nil
	}
	return nil, errors.New("isos: no info for isotope " + fmt.Sprint(i))
}

type Info struct {
	EltName   string
	EltSymbol string
	HalfLife  float64
	A         float64
	Z         int
	IS        int // isomeric state
}

func (info *Info) Iso() Iso {
	return Iso(10000*info.Z + 10*int(info.A) + info.IS)
}
