package controller

import (
	"fmt"
	"machine"
	"machine/usb/adc/midi"

	"github.com/felipemantoan/tinygo-midi-pad/matrix"
	"github.com/felipemantoan/tinygo-midi-pad/potentiometer"
)

var (
	adcPins = []machine.ADC{
		machine.ADC{Pin: machine.ADC0},
		machine.ADC{Pin: machine.ADC1},
		machine.ADC{Pin: machine.ADC2},
	}
	controls = []uint8{
		midi.CCExpression,
		midi.CCTremeloLevel,
		midi.CCChorusLevel,
	}
	columns = []machine.Pin{
		machine.GPIO0,
		machine.GPIO1,
		machine.GPIO2,
		machine.GPIO3,
	}
	rows = []machine.Pin{
		machine.GPIO4,
		machine.GPIO5,
		machine.GPIO6,
		machine.GPIO7,
	}
	notes = []midi.Note{
		midi.C2,
		midi.D2,
		midi.E2,
		midi.F2,
		midi.G2,
		midi.A2,
		midi.B2,
		midi.C3,
		midi.D3,
		midi.E3,
		midi.F3,
		midi.G3,
		midi.A3,
		midi.B3,
		midi.CS3,
		midi.AS3,
	}
)

func Controller() {
	midi.Port()
	go Pads()
	go CCs()
}

func Pads() {
	mesh := matrix.New(matrix.Configure(rows[:], columns[:]))
	mesh.SetInterrupt(matrix.CellToggle, PadToggle)
	mesh.Scan()
}

func PadToggle(c matrix.Cell) {
	m := midi.Port()
	if c.Change() == matrix.CellFalling {
		fmt.Println("CallBack CellFalling") // call to action
		err := m.NoteOn(0, 1, notes[c.ID()], 0x40)

		if err != nil {
			fmt.Println(err)
		}
	}

	if c.Change() == matrix.CellRising {
		fmt.Println("CallBack CellRising")
		err := m.NoteOff(0, 1, notes[c.ID()], 0x0)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func CCs() {
	m := midi.Port()
	pots := make([]potentiometer.PotLinear, len(adcPins))

	for i, adcPin := range adcPins {
		pots[i] = potentiometer.New(adcPin, potentiometer.RotaryPotShape, 9)
		pots[i].SetInterrupt(potentiometer.PotLinearToggle, func(pot potentiometer.PotLinear) {
			if pot.HasChange() {
				m.ControlChange(0, 1, controls[i], uint8(pot.Value())) // Tornar isso configurável
			}
		})

		go pots[i].Scan()
	}
}
