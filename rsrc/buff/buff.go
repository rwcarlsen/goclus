// Package buff provides a generic tool for managing resource
// inventories.
package buff

import (
	"errors"
	"github.com/rwcarlsen/goclus/rsrc"
)

var (
	OverCapErr  = errors.New("buff: cannot hold more than its capacity")
	TooSmallErr = errors.New("buff: operation results in negligible quantities")
)

// Buffer is a resource inventory that helps manage capacity, addition, and
// removal resources.
type Buffer struct {
	capacity float64
	res      []rsrc.Resource
}

// Capacity returns the maximum resource quantity this buffer can hold (units
// based on constituent resource objects' units).
func (b *Buffer) Capacity() float64 {
	return b.capacity
}

// SetCapacity sets the maximum quantity this store can hold.
//
// Returns an error if the new capacity is lower then the quantity currently
// residing in the buffer.
func (b *Buffer) SetCapacity(capacity float64) error {
	if b.Qty()-b.capacity > rsrc.EPS {
		return OverCapErr
	}
	b.capacity = capacity
	return nil
}

// Count returns the number of resource objects being held in the buffer.
func (b *Buffer) Count() int {
	return len(b.res)
}

// Qty returns the total resource quantity of constituent resource objects in the buffer.
func (b *Buffer) Qty() float64 {
	var tot float64
	for _, r := range b.res {
		tot += r.Qty()
	}
	return tot
}

// Space returns the quantity of space remaining in the buffer (Capacity - Qty).
func (b *Buffer) Space() float64 {
	return b.capacity - b.Qty()
}

// PopQty pops and returns the specified quantity of resources from the buffer.
//
// Resources are split if necessary in order to pop the exact quantity.
// Resources are retrieved in the order they were pushed (first in - first
// out).
func (b *Buffer) PopQty(qty float64) ([]rsrc.Resource, error) {
	if qty-b.Qty() > rsrc.EPS || qty < rsrc.EPS {
		return nil, TooSmallErr
	}

	left := qty
	popped := []rsrc.Resource{}
	for left > rsrc.EPS {
		r := b.res[0]
		b.res = b.res[1:]
		quan := r.Qty()
		if quan-left > rsrc.EPS {
			leftover := r.Clone()
			leftover.SetQty(quan - left)
			r.SetQty(left)
			b.res = append([]rsrc.Resource{leftover}, b.res...)
		}
		popped = append(popped, r)
		left -= quan
	}
	return popped, nil
}

// PopQty pops and returns the specified number of resources from the buffer.
//
// Resources are not split. Resources are retrieved in the order they were
// pushed (first in - first out).
func (b *Buffer) PopN(num int) ([]rsrc.Resource, error) {
	if len(b.res) < num {
		return nil, TooSmallErr
	}
	popped := b.res[:num]
	b.res = b.res[num:]
	return popped, nil
}

// PopOne pops and returns one resource object from the store.
//
// Resources are not split. Resources are retrieved in the order they were
// pushed (first in - first out).
func (b *Buffer) PopOne() (rsrc.Resource, error) {
	if len(b.res) < 1 {
		return nil, TooSmallErr
	}
	popped := b.res[0]
	b.res = b.res[1:]
	return popped, nil
}

// Push pushes one or more resource objects into the buffer.
//
// If the push would result in the buffer being over capacity, no resources are
// pushed, and an error is returned.
//
// Resource objects are never combined in the buffer.
func (b *Buffer) Push(rs ...rsrc.Resource) error {
	var tot float64
	for _, r := range rs {
		tot += r.Qty()
	}
	if tot-b.Space() > rsrc.EPS {
		return OverCapErr
	}
	b.res = append(b.res, rs...)
	return nil
}
