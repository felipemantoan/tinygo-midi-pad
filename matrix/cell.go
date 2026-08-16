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
			c.state = true
			return true
		}

		c.state = false
		return true
		// CellRising
	}

	return false
}

func (c *cellImpl) IsActivated() bool {
	return c.column.Get() && c.row.Get()
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
	// switch && save callback
}
