package potentiometer

import (
	"machine"
)

type PotLinearShape uint8

const (
	RotaryPotShape PotLinearShape = 4 << iota
	SlidePotShape
)

type PotLinearChange uint8

// PotLinear change interrupt constants for SetInterrupt.
const (
	PotLinearClockWise     PotLinearChange = 4 << iota
	PotLinearAntiClockWise PotLinearChange = 3 << iota
	PotLinearUp            PotLinearChange = 2 << iota
	PotLinearDown
	PotLinearRest = PotLinearClockWise | PotLinearAntiClockWise | PotLinearUp | PotLinearDown
)

var (
	isInitialized = false
)

type PotLinear interface {
	ID() int
	Change() PotLinearChange
	HasChange() bool
	Pin() machine.ADC
	SetInterrupt(change PotLinearChange, callback func(pot PotLinear))
	Shape() PotLinearShape
	Value() uint16
}

type Device struct {
	id       int
	callback func(pot PotLinear)
	change   PotLinearChange
	pin      machine.ADC
	shape    PotLinearShape
	value    uint16
}

func (ptl *Device) SetInterrupt(change PotLinearChange, callback func(pot PotLinear)) {
	ptl.callback = callback
}

func (ptl *Device) Change() PotLinearChange {
	return ptl.change
}

func (ptl *Device) Shape() PotLinearShape {
	return ptl.shape
}

func (ptl *Device) HasChange() bool {

	newValue := ptl.Pin().Get() >> 9

	if ptl.value != newValue {
		if ptl.value < newValue {
			if ptl.shape == RotaryPotShape { // switch
				ptl.change = PotLinearClockWise
			} else {
				ptl.change = PotLinearDown
			}
		} else {
			if ptl.shape == RotaryPotShape { // switch
				ptl.change = PotLinearAntiClockWise
			} else {
				ptl.change = PotLinearUp
			}
		}
		ptl.value = newValue
		return true
	}

	return false
}

func (ptl *Device) ID() int {
	return ptl.id
}

func (ptl *Device) Pin() machine.ADC {
	return ptl.pin
}

func (ptl *Device) Value() uint16 {
	return ptl.value
}

func New(pin machine.ADC, config machine.ADCConfig, shape PotLinearShape) PotLinear {

	if !isInitialized {
		machine.InitADC()
		isInitialized = true
	}

	pin.Configure(config)

	return &Device{
		id:     int(pin.Pin),
		change: PotLinearRest,
		pin:    pin,
		shape:  shape,
		value:  pin.Get() >> 9,
	}
}
