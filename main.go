package main

import (
	"fmt"
	"image/color"
	"machine"
	"machine/usb/adc/midi"
	"strconv"
	"time"

	"github.com/felipemantoan/tinygo-midi-pad/matrix"
	"tinygo.org/x/drivers/encoders"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/drivers/ws2812"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freesans"
)

var (
	neo     machine.Pin
	leds    [2]color.RGBA
	voltage = 3.3 / 65535
	enc     = encoders.NewQuadratureViaInterrupt(machine.GPIO19, machine.GPIO20)
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
	midi.C3,
	midi.D3,
	midi.E3,
	midi.F3,
	midi.G3,
	midi.A3,
	midi.B3,
	midi.C4,
	midi.D4,
	midi.E4,
	midi.F4,
	midi.G4,
	midi.A4,
	midi.B4,
	midi.C5,
	midi.D5,
}

func blinkBuiltInLed() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for {
		time.Sleep(500 * time.Millisecond)
		led.High()
		time.Sleep(500 * time.Millisecond)
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

	machine.InitADC()
	sensorA := machine.ADC{Pin: machine.ADC0}
	sensorA.Configure(machine.ADCConfig{})

	sensorB := machine.ADC{Pin: machine.ADC1}
	sensorB.Configure(machine.ADCConfig{})

	sensorC := machine.ADC{Pin: machine.ADC2}
	sensorC.Configure(machine.ADCConfig{})

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
			err := m.NoteOff(0, 1, notes[c.ID()], 0x40)

			if err != nil {
				fmt.Println(err)
			}
		}
	})

	mesh.Scan()

	for oldValue := 0; ; {
		adc0Value := float64(sensorA.Get()) * voltage
		adc1Value := float64(sensorB.Get()) * voltage
		adc2Value := float64(sensorC.Get()) * voltage

		if adc0Value > noiseADC {
			fmt.Println("Valor ADC sensorA: ", adc0Value)

		}

		if adc1Value > noiseADC {
			fmt.Println("Valor ADC sensorB: ", adc1Value)

		}

		if adc2Value > noiseADC {
			fmt.Println("Valor ADC sensorC: ", adc2Value)

		}

		// fmt.Println("SW Encoder: ", !sw.Get())

		if newValue := enc.Position(); newValue != oldValue {
			device.ClearDisplay()
			oldValue = newValue
		}
		device.Display()

		tinyfont.WriteLine(device, &freesans.Bold18pt7b, 0, 31, strconv.Itoa(oldValue), color.RGBA{255, 255, 255, 1})
		// fmt.Println("value: ", oldValue)
		time.Sleep(10 * time.Millisecond)

	}
}
