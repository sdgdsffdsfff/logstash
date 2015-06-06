package utils

import (
	//"container/list"
	"errors"
)

var (
	Err_OverFlow = errors.New("cursor overfloat")
)

type IntCursor struct {
	index int
	ints  []int
}

func NewIntCursor(arr ...int) *IntCursor {
	return &IntCursor{
		index: 0,
		ints:  arr,
	}
}

func (c *IntCursor) Assert(ok bool, overflow func()) {
	if ok {
		c.Ahead(1)
	} else {
		c.Back(-1)
	}

	if c.Ended() {
		overflow()
		c.Reset()
	}
}

func (c *IntCursor) Walk(ahead bool, overflow func()) {
	if ahead {
		c.Ahead(1)
	} else {
		c.Back(-1)
	}
	if c.Ended() {
		overflow()
		c.Reset()
	}
}

func (c *IntCursor) Ahead(step int) {
	c.index += step
}

func (c *IntCursor) Reset() {
	c.index = 0
}

func (c *IntCursor) Back(step int) {
	c.index -= step
	if c.index < 0 {
		c.index = 0
	}
}

func (c *IntCursor) Ended() bool {
	return c.index > len(c.ints)-1
}

func (c *IntCursor) Val() (v int, err error) {
	if c.index > len(c.ints)-1 || c.index < 0 {
		return 0, Err_OverFlow
	}

	return c.ints[c.index], nil
}
