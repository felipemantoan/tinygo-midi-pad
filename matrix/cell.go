package matrix

import "machine"

type CellChange uint8

// Cell change interrupt constants for SetInterrupt.
const (
	// Edge falling
	CellFalling CellChange = 4 << iota
	// Edge rising
	CellRising

	CellToggle = CellFalling | CellRising
)

type Cell interface {
	ID() int
	IsActivated() bool
	SetState(state bool)
	State() bool
	HasChange() bool
}

type cellImpl struct {
	id     int
	column machine.Pin
	row    machine.Pin
	state  bool
}

func (c *cellImpl) HasChange() bool {
	// Daria pra fazer em 1 linha
	oldState := c.State()
	newState := c.IsActivated()

	if oldState != newState {
		// temos um novo estatdo mas qual?
		if oldState == false && newState == true {
			c.SetState(true)
			return true
		}

		c.SetState(false)
		return true
		// CellRising
	}

	return false
}

func (c *cellImpl) IsActivated() bool {
	return c.column.Get() && c.row.Get()
}

func (c *cellImpl) SetState(state bool) {
	c.state = state
}

func (c *cellImpl) State() bool {
	return c.state
}

func (c *cellImpl) ID() int {
	return c.id
}

func newCell(id int, column machine.Pin, row machine.Pin) Cell {
	return &cellImpl{id: id, column: column, row: row, state: false}
}

func (c *cellImpl) SetInterrupt(change CellChange, callback func(c Cell)) {
	if change == CellFalling && (c.State() == false && c.IsActivated()) {
		c.SetState(true)
		// save callback(c)
	}

	if change == CellRising && (c.State() && !c.IsActivated()) {
		c.SetState(false)
		// save callback(c)
	}

}
