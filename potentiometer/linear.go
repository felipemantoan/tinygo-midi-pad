package potentiometer

import (
	"fmt"
	"machine"
	"time"
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

	hasChange := ptl.value != newValue

	if hasChange {
		ptl.change = determineMoviment(ptl.value, newValue, ptl.config.shape)
		ptl.value = newValue
		return true
	}

	return false

}

func determineMoviment(oldValue, newValue uint16, shape PotLinearShape) PotLinearChange {

	switch shape {
	case RotaryPotShape:

		if oldValue < newValue {
			return PotLinearClockWise
		} else {
			return PotLinearAntiClockWise
		}

	case SlidePotShape:
		if oldValue < newValue {
			return PotLinearDown
		} else {
			return PotLinearUp
		}

	default:
		return PotLinearRest

	}
}

func (ptl *Device) ID() int {
	return ptl.id
}

func (ptl *Device) Value() uint16 {
	return ptl.value
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
