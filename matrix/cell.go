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
	Change() CellChange
	IsActivated() bool
	HasChange() bool
	State() bool
}

type cellImpl struct {
	id     int
	change CellChange
	column machine.Pin
	row    machine.Pin
	state  bool
}

func (c *cellImpl) Change() CellChange {
	return c.change
}

func (c *cellImpl) HasChange() bool {
	oldState := c.State()
	newState := c.IsActivated()

	if oldState != newState {
		if oldState == false && newState == true {
			c.state = true
			c.change = CellFalling
		} else {
			c.state = false
			c.change = CellRising
		}
		return true
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
	return &cellImpl{id: id, change: CellToggle, column: column, row: row, state: false}
}
