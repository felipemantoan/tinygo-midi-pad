package controller

import (
	"machine"
	"machine/usb/adc/midi"

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
	shiftRight = uint16(9)
)

func CommandControls() {
	m := midi.Port()
	pots := make([]potentiometer.PotLinear, len(adcPins))

	for i, adcPin := range adcPins {

		config := potentiometer.Configure(
			adcPin,
			potentiometer.RotaryPotShape,
			shiftRight)

		pots[i] = potentiometer.New(config)

		pots[i].SetInterrupt(potentiometer.PotLinearToggle, func(pot potentiometer.PotLinear) {
			if pot.HasChange() {
				m.ControlChange(0, 1, controls[i], uint8(pot.Value())) // Tornar isso configurável
			}
		})

		go pots[i].Scan()
	}
}
