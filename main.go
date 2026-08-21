package main

import (
	"fmt"
	"image/color"
	"machine"
	"strconv"
	"time"

	"github.com/felipemantoan/tinygo-midi-pad/controller"
	"tinygo.org/x/drivers/encoders"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/drivers/ws2812"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freesans"
)

var (
	neo      machine.Pin
	leds     [2]color.RGBA
	enc      = encoders.NewQuadratureViaInterrupt(machine.GPIO19, machine.GPIO20)
	noiseADC = 0.025
)

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
	sw.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		fmt.Println("SW Encoder: ", !sw.Get())

	})

	neo := machine.GPIO18
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.NewWS2812(neo)
	ws.SetBrightness(50)
	leds[0] = color.RGBA{R: 0, G: 223, B: 197}
	leds[1] = color.RGBA{R: 255, G: 100, B: 0}

	ws.WriteColors(leds[:])
	controller.Controller()
	go blinkBuiltInLed()

	for oldValue := 0; ; {

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
