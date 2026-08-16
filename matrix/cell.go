package matrix

import "machine"

type Cell interface {
	ID() int
	IsActivated() bool
	SetState(state bool)
	State() bool
}

type cellImpl struct {
	id     int
	column machine.Pin
	row    machine.Pin
	state  bool
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
