package potentiometer

import (
	"fmt"
	"machine"
)

var (
	voltage = 3.3 / 65535

	ADCConfig = machine.ADCConfig{
		Reference: 33000,
	}

	ADCPins = []machine.ADC{
		machine.ADC{Pin: machine.ADC0},
		machine.ADC{Pin: machine.ADC1},
		machine.ADC{Pin: machine.ADC2},
	}
)

func initPins() []machine.ADC {
	machine.InitADC()

	for _, pin := range ADCPins {
		pin.Configure(ADCConfig)
	}

	return ADCPins
}

func Scan() {
	pins := initPins()

	for i, pin := range pins {
		pinValue := pin.Get() >> 9
		if pinValue > 5 {
			fmt.Println("Pin: ", i, " Value: ", pinValue)
		}
	}
}
