package main

import (
	"fmt"
	"image/color"
	"machine"
	"machine/usb/adc/midi"
	"time"

	"github.com/felipemantoan/tinygo-midi-pad/matrix"
	"github.com/felipemantoan/tinygo-midi-pad/potentiometer"
	"tinygo.org/x/drivers/encoders"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/drivers/ws2812"
)

var (
	neo     machine.Pin
	leds    [2]color.RGBA
	enc     = encoders.NewQuadratureViaInterrupt(machine.GPIO19, machine.GPIO20)
	ADCPins = []machine.ADC{
		machine.ADC{Pin: machine.ADC0},
		machine.ADC{Pin: machine.ADC1},
		machine.ADC{Pin: machine.ADC2},
	}
	columns = [4]machine.Pin{
		machine.GPIO0,
		machine.GPIO1,
		machine.GPIO2,
		machine.GPIO3,
	}
	rows = [4]machine.Pin{
		machine.GPIO4,
		machine.GPIO5,
		machine.GPIO6,
		machine.GPIO7,
	}
	noiseADC = 0.025
)

var notes = [16]midi.Note{
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

func blinkBuiltInLed() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for {
		time.Sleep(1000 * time.Millisecond)
		led.High()
		time.Sleep(1000 * time.Millisecond)
		led.Low()
	}
}

func newSSD1306Display() *ssd1306.Device {

	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
		SDA:       machine.GPIO16,
		SCL:       machine.GPIO17,
	})

	display := ssd1306.NewI2C(machine.I2C0)

	display.Configure(ssd1306.Config{
		Address: ssd1306.Address_128_32, // or ssd1306.Address
		Width:   128,
		Height:  32, // or 64
	})

	return display
}

func main() {

	device := newSSD1306Display()
	device.ClearDisplay()

	enc.Configure(encoders.QuadratureConfig{
		Precision: 4,
	})

	sw := machine.GPIO21
	sw.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	neo := machine.GPIO18
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.NewWS2812(neo)
	ws.SetBrightness(50)
	leds[0] = color.RGBA{R: 0, G: 223, B: 197}
	leds[1] = color.RGBA{R: 255, G: 100, B: 0}

	ws.WriteColors(leds[:])

	go blinkBuiltInLed()

	m := midi.New()

	mesh := matrix.New(matrix.Configure(rows[:], columns[:]))
	mesh.SetInterrupt(matrix.CellFalling, func(c matrix.Cell) {

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
	})

	go mesh.Scan()

	pots := make([]potentiometer.PotLinear, len(ADCPins))

	for i, adcPin := range ADCPins {
		pots[i] = potentiometer.New(adcPin, potentiometer.RotaryPotShape, 9)
	}

	for {

		for i, pot := range pots {
			if pot.HasChange() {
				fmt.Println("Pot:", i, ", Value:", pot.Value(), ", Change: ", pot.Change())
			}
		}

		time.Sleep(60 * time.Millisecond)
	}

	// for oldValue := 0; ; {
	// 	// go potentiometer.Scan()

	// 	// fmt.Println("SW Encoder: ", !sw.Get())

	// 	if newValue := enc.Position(); newValue != oldValue {
	// 		device.ClearDisplay()
	// 		oldValue = newValue
	// 	}
	// 	device.Display()

	// 	tinyfont.WriteLine(device, &freesans.Bold18pt7b, 0, 31, strconv.Itoa(oldValue), color.RGBA{255, 255, 255, 1})
	// 	// fmt.Println("value: ", oldValue)
	// 	time.Sleep(10 * time.Millisecond)

	// }
}
