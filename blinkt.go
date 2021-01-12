// Package blinkt provides control over a Blinkt! LED display
package blinkt

// Blinkt provides control over a Blinkt! LED display
type Blinkt struct {
	// Pixels is the number of pixels supported by the display
	Pixels       int
	remainOnExit bool
	cmdChan      chan func()
	closed       chan struct{}
	pp           []pixel
	apa          APA102
}

// New creates a Blinkt to control the display.
func New() Blinkt {

	bl := Blinkt{
		Pixels:  8,
		pp:      make([]pixel, 8),
		cmdChan: make(chan func()),
		closed:  make(chan struct{}),
	}
	bl.SetBrightness(50)
	bl.apa.Open()

	/* cmdLoop serialises writes to APA102 and closing */
	go func() {
		for {
			cmd := <-bl.cmdChan
			cmd()
			select {
			case <-bl.closed:
				return
			default:
			}
		}
	}()
	return bl
}

// Close the Blinkt interface.
func (bl *Blinkt) Close() {
	closer := func() {
		if !bl.remainOnExit {
			bl.apa.WritePixels(make([]pixel, 8))
		}
		bl.apa.Close()
		close(bl.closed)
	}
	select {
	case bl.cmdChan <- closer:
		<-bl.closed
	case <-bl.closed:
	}
}

// Clear sets all the pixels to off.
//
// Show must still be called to update the physical display.
func (bl *Blinkt) Clear() {
	for p := range bl.pp {
		bl.pp[p] &= 0xff000000
	}
}

// ClearPixel sets the pixel to off.
//
// Show must still be called to update the physical display.
func (bl *Blinkt) ClearPixel(p int) {
	bl.pp[p] &= 0xff000000
}

// SetClearOnExit controls the blanking of the display when the Blinkt is closed.
//
// When clearOneExit is true all pixels are turned off when Blinkt is closed.
// This is the default.
// When clearOneExit is false the display is left in its current state.
func (bl *Blinkt) SetClearOnExit(clearOnExit bool) {
	bl.remainOnExit = !clearOnExit
}

// Show updates the physical dispolay with the values from Set/Clear.
func (bl *Blinkt) Show() {
	pixels := append([]pixel(nil), bl.pp[:]...)
	show := func() {
		bl.apa.WritePixels(pixels)
	}
	select {
	case bl.cmdChan <- show:
	case <-bl.closed:
	}
}

// SetAll sets all pixels to specified r, g, b colour.
//
// Show must be called to update the LEDs.
func (bl *Blinkt) SetAll(r, g, b int) {

	for p := range bl.pp {
		bl.SetPixel(p, r, g, b)
	}
}

// SetPixel sets an individual pixel to specified r, g, b colour.
//
// Show must be called to update the LEDs.
func (bl *Blinkt) SetPixel(p, r, g, b int) {
	px := bl.pp[p] & 0xff000000
	px |= pixel((b << 16) | g<<8 | r)
	bl.pp[p] = px
}

// SetBrightness sets the brightness of all pixels.
//
// Brightness should be in percent, i.e. 0 to 100.0.
// Greater than or equal to 100 is assumed to mean full brightness.
// Less than or equal to 0 is assumed to mean off.
func (bl *Blinkt) SetBrightness(brightness float64) {

	brightnessInt := convertBrightnessToInt(brightness)

	for i, px := range bl.pp {
		px &^= 0xff000000
		px |= pixel(brightnessInt << 24)
		bl.pp[i] = px
	}
}

// SetPixelBrightness sets the brightness of pixel p.
//
// Brightness should be in percent, i.e. 0 to 100.0.
// Greater than or equal to 100 is assumed to mean full brightness.
// Less than or equal to 0 is assumed to mean off.
func (bl *Blinkt) SetPixelBrightness(p int, brightness float64) {
	brightnessInt := convertBrightnessToInt(brightness)
	bl.pp[p] = pixel(uint(bl.pp[p])&^uint(0xff000000) | brightnessInt<<24)
}

type pixel uint

func convertBrightnessToInt(brightness float64) uint {

	switch {
	case brightness <= 0:
		return 0
	case brightness >= 100:
		return 31
	default:
		return uint(brightness * 0.31)
	}
}
