package blinkt

// High level functions common to all APA102 implementations

// WritePixels writes a sequence of pixels to the display.
//
// The pixels are prefixed with a start frame and terminated with the end frame.
func (a *APA102) WritePixels(pixels []pixel) {
	// Start Frame
	a.WriteFrame(0)
	// LED frames
	for _, p := range pixels {
		a.WriteFrame(0xe0000000 | uint(p))
	}
	// End Frame
	a.WriteFrame(0xffffffff)
}

// WriteFrame writes a single LED frame to the APA102.
//
// The frame value is passed as a 32bit uint and it clocked out MSB first.
func (a *APA102) WriteFrame(val uint) {
	var mask = uint(0x80000000)
	var bit int
	for i := 0; i < 32; i++ {
		if val&mask == 0 {
			bit = 0
		} else {
			bit = 1
		}
		mask >>= 1
		a.WriteBit(bit)
	}
}
