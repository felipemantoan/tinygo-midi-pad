package potentiometer

import "machine"

type PotLinearShape uint8

const (
	RotaryPotShape PotLinearShape = 4 << iota
	SlidePotShape
)

type PotLinearChange uint8

// PotLinear change interrupt constants for SetInterrupt.
const (
	PotLinearRest PotLinearChange = iota
	PotLinearClockWise
	PotLinearAntiClockWise
	PotLinearUp
	PotLinearDown
	PotLinearToggle = PotLinearClockWise | PotLinearAntiClockWise | PotLinearUp | PotLinearDown
)

type PotLinearConfig struct {
	pin        machine.ADC
	shape      PotLinearShape
	shiftRight uint16
}

func Configure(pin machine.ADC, shape PotLinearShape, shiftRight uint16) PotLinearConfig {
	return PotLinearConfig{
		pin:        pin,
		shape:      shape,
		shiftRight: shiftRight,
	}
}
