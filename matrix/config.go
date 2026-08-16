package matrix

import "machine"

type MatrixConfig interface {
	Columns() []machine.Pin
	Rows() []machine.Pin
}

type matrixConfigImpl struct {
	columns []machine.Pin
	rows    []machine.Pin
}

func (msc *matrixConfigImpl) Columns() []machine.Pin {
	return msc.columns
}

func (msc *matrixConfigImpl) Rows() []machine.Pin {
	return msc.rows
}

func Configure(rows []machine.Pin, columns []machine.Pin) MatrixConfig {

	for i := range columns {
		columns[i].Configure(machine.PinConfig{Mode: machine.PinMode(machine.PinOutput)})
	}

	for i := range rows {
		rows[i].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	return &matrixConfigImpl{
		rows:    rows,
		columns: columns,
	}
}
