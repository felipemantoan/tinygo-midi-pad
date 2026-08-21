package controller

import (
	"machine/usb/adc/midi"
)

func Controller() {
	midi.Port()
	go Pads()
	go CommandControls()
}
