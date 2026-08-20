package potentiometer

import (
	"fmt"
	"machine"
	"time"
)

type PotLinearShape uint8

const (
	RotaryPotShape PotLinearShape = 4 << iota
	SlidePotShape
)

type PotLinearChange uint8

// PotLinear change interrupt constants for SetInterrupt.
const (
	PotLinearRest PotLinearChange = 4 << iota
	PotLinearClockWise
	PotLinearAntiClockWise
	PotLinearUp
	PotLinearDown
	PotLinearToggle = PotLinearClockWise | PotLinearAntiClockWise | PotLinearUp | PotLinearDown
)

var (
	isInitialized = false
)

type PotLinear interface {
	ID() int
	Change() PotLinearChange
	HasChange() bool
	Scan()
	SetInterrupt(change PotLinearChange, callback func(pot PotLinear))
	Value() uint16
}

type PotLinearConfig struct {
	pin        machine.ADC
	shape      PotLinearShape
	shiftRight uint16
}

type Device struct {
	id       int
	callback func(pot PotLinear)
	change   PotLinearChange
	config   PotLinearConfig
	value    uint16
}

func (ptl *Device) Scan() {
	for {
		if ptl.HasChange() {
			fmt.Println("Pot:", ptl.id, ", Value:", ptl.Value(), ", Change: ", ptl.Change())
			ptl.callback(ptl) //interrupt
		}
		time.Sleep(25 * time.Millisecond)
	}
}

func (ptl *Device) SetInterrupt(change PotLinearChange, callback func(pot PotLinear)) {
	ptl.callback = callback
}

func (ptl *Device) Change() PotLinearChange {
	return ptl.change
}

func (ptl *Device) HasChange() bool {

	newValue := ptl.config.pin.Get() >> ptl.config.shiftRight

	if ptl.value != newValue {
		if ptl.value < newValue {
			if ptl.config.shape == RotaryPotShape { // switch
				ptl.change = PotLinearClockWise
			} else {
				ptl.change = PotLinearDown
			}
		} else {
			if ptl.config.shape == RotaryPotShape { // switch
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

func (ptl *Device) Value() uint16 {
	return ptl.value
}

func Configure(pin machine.ADC, shape PotLinearShape, shiftRight uint16) PotLinearConfig {
	return PotLinearConfig{
		pin:        pin,
		shape:      shape,
		shiftRight: shiftRight,
	}
}

func New(config PotLinearConfig) PotLinear {

	if !isInitialized {
		machine.InitADC()
		isInitialized = true
	}

	config.pin.Configure(machine.ADCConfig{})

	value := config.pin.Get()

	if config.shiftRight > 0 {
		value = config.pin.Get() >> config.shiftRight
	}

	return &Device{
		id:     int(config.pin.Pin),
		change: PotLinearRest,
		config: config,
		value:  value,
	}
}
