package matrix

import (
	"fmt"
	"time"
)

type Matrix interface {
	Config() MatrixConfig
	Cell(column int, row int) Cell
	Scan()
}

type matrixImpl struct {
	cells  [][]Cell
	config MatrixConfig
}

func New(config MatrixConfig) Matrix {
	cells := make([][]Cell, len(config.Rows()))

	id := 0
	for i, row := range config.Rows() {

		cells[i] = make([]Cell, len(config.Columns()))

		for j, column := range config.Columns() {
			cells[i][j] = newCell(id, column, row)
			id++
		}
	}

	return &matrixImpl{
		cells:  cells,
		config: config,
	}
}

func (m *matrixImpl) Cell(row int, column int) Cell {
	return m.cells[row][column]
}

func (m *matrixImpl) Config() MatrixConfig {
	return m.config
}

func (m *matrixImpl) Scan() {
	for {
		for i, column := range m.config.Columns() {
			column.High()
			time.Sleep(20 * time.Millisecond)

			for j := range m.config.Rows() {

				if m.Cell(j, i).HasChange() {
					// debug
					fmt.Println(m.Cell(j, i).ID())
				}
			}

			column.Low()
		}
	}
}
