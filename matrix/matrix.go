package matrix

import (
	"errors"
	"fmt"
	"time"
)

const SCAN_TIME = 5_000_000

type Matrix interface {
	Config() MatrixConfig
	Cell(column int, row int) Cell
	Scan()
	SetInterrupt(change CellChange, callback func(c Cell))
}

type matrixImpl struct {
	cells        [][]Cell
	config       MatrixConfig
	callback     func(c Cell)
	callBackType CellChange
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

func (m *matrixImpl) scanMainLoop() {
	for {
		for i, column := range m.config.Columns() {
			column.High()
			time.Sleep(SCAN_TIME)

			for j := range m.config.Rows() {
				err := m.interrupt(j, i)
				if err != nil {
					fmt.Println(err)
				}
			}

			column.Low()
		}
	}
}

func (m *matrixImpl) SetInterrupt(change CellChange, callback func(c Cell)) {
	m.callBackType = change
	m.callback = callback
}

func (m *matrixImpl) interrupt(row int, column int) error {
	cell := m.Cell(row, column)

	if m.callback == nil {
		return errors.New("Deu Problema")
	}

	if cell.HasChange() {
		m.callback(cell)
	}

	return nil
}

func (m *matrixImpl) Scan() {
	m.scanMainLoop()
}
