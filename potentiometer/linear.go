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

	if ptl.value == newValue {
		return false
	}

	ptl.value = newValue
	ptl.determineMoviment(ptl.value < newValue)

	return true
}

func (ptl *Device) determineMoviment(asc bool) {

	switch ptl.config.shape {
	case RotaryPotShape:
		if asc {
			ptl.change = PotLinearClockWise
		} else {
			ptl.change = PotLinearAntiClockWise
		}
		break
	case SlidePotShape:
		if asc {
			ptl.change = PotLinearDown
		} else {
			ptl.change = PotLinearUp
		}
		break
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
