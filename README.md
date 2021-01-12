blinkt
======

This library is a fork of Alex Ellis'
[blinkt_go](https://github.com/alexellis/blinkt_go) that I was intending to port
to my [gpiod](https://github.com/warthog618/gpiod) library, but ended up going
off on a bit of a tangent and also using it as a testbed for comparing various
GPIO implementations, and incorporating a merge/port of Alex's
[blinkt_go_examples](https://github.com/alexellis/blinkt_go_examples).

## Quick Start

```go
import "github.com/warthog619/blinkt"

// create a blinkt to control the display
bl := blinkt.New()

// set the display brightness
bl.SetBrightness(23)

// write to pixels
bl.SetPixel(0, r, g, b)
...
bl.SetPixel(7, r, g, b)

// show the pixels on the display
bl.Show()

// clear the display
bl.Clear()
bl.Show()

// finish with the display
bl.Close()
```

## Usage

The following is a short summary of the Blinkt API.
#### New

Create a new Blinkt object using *New":

```go
bl := blinkt.New()
```

#### Close

If you no longer need to use the display, or your application is exiting, you
should close the display:

```go
bl.Close()
```

#### SetClearOnExit

By default the display is cleared when closed.  This can be changed with the
*SetClearOnExit* method:

```go
bl.SetClearOnExit(false)
```

### Whole Display

The following methods act on all pixels in the display.

Note that nothing is actually written to the physical display until *Show* is
called.

#### SetAll

Set the colour of all pixels at once using *SetAll*:

```go
bl.SetAll(r, g, b)
```

Setting the colour does not alter the brightness, which is independent.
#### SetBrightness

Set the brightness of all pixels at once using *SetBrightness*:

```go
bl.SetBrightness(42)
```

The brightness value is a percent, so 0-100. By default the display brightness
is set to 50.

#### Clear

Turn off all pixels using the *Clear* method:

```go
bl.Clear()
```

#### Show

Update the display with the current pixel states using *Show*:

```go
bl.Show()
```

The other methoss only update the pixels in memory - the *Show* command commits
that state to the display.

The *Show* method performs the update in a separate goroutine using a snapshot
of the pixel state, so the main goroutine can immediately begin updating the
pixel state while the snapshot is being written to the display.

The *Show* method will block if the background goroutine is still updating the
display with the previous *Show*.

### Individual Pixels

The following methods apply to an individual pixel.

The pixel parameter which identifies the pixel, and is denoted by **p** in these
snippets, is in the range 0-7.

#### ClearPixel

Turn off the pixel using the *ClearPixel* method:

```go
bl.ClearPixel(p)
```

#### SetPixel

Set the colour of the pixel using *SetPixel*:

```go
bl.SetPixel(p, r, g, b)
```

Setting the colour does not alter the brightness, which is independent.

#### SetPixelBrightness

Set the brightness of the pixel using *SetPixelBrightness*:

```go
bl.SetPixelBrightness(68)
```

### Goroutine Safety

The Blinkt is not safe to call from multiple goroutines, with the exception of
the *Close* method, which is made goroutine safe so it can be called from signal
handlers, e.g. when the application is closing down.


## Backends

The library provides a selection of backends to drive the GPIO, with the
selection being performed by a build tag. e.g.

```go build -tags blinkt_wiringpi```

will build an application using the WiringPi backend.

The default backend is my [gpiod](https://github.com/warthog618/gpiod) library,
which makes use of the official Linux GPIO interface.

|Backend|Build tag|Interface|
|---|---|---|
|[gpiod](https://github.com/warthog618/gpiod)|none (default)|Linux GPIO character device (*/dev/gpiochip0*) uAPI |
|[gpio](https://github.com/warthog618/gpio)|blinkt_gpio|Raspberry Pi */dev/gpiomem* direct access to the BCM hardware|
|[wiringpi](https://github.com/alexellis/rpi/)| blinkt_wiringpi|WiringPi cgo wrapper|
|[sysfs](https://github.com/alexellis/blinkt_go/)| blinkt_sysfs|deprecated GPIO SYSFS interface|

## Examples

The examples directory contains a collection of examples that demonstrate the
library and display.

These were originally ported from Alex's
[blinkt_go_examples](https://github.com/) which are themselves largely ports of
some of the [Pimoroni examples](https://github.com/pimoroni/blinkt/blob/master/examples).

I've ported a few more of those, such as the binary clocks and rainbow.

## Benchmarks

These are the results from a Raspberry Pi Zero W running Linux v5.10 and built
with go1.15.6:

```
$ go test -test.bench=.*
goos: linux
goarch: arm
pkg: github.com/warthog618/blinkt
BenchmarkShow 	     106	  10547923 ns/op
PASS
ok  	github.com/warthog618/blinkt	1.683s
```

A summary of the results for each of the backends:

|Backend|ns/op|Shows/sec|
|---|---|---|
|gpiod|10547923|94|
|gpio|102205|9784|
|wiringpi|3435945|291|
|sysfs|66638981|15|


For comparison the Python/WiringPi implementation can perform around 9 shows/sec
on the same platform.

My recommendation is to use the **gpiod** backend unless you have a serious need
for speed or very low CPU utilization (and aren't concerned about conflicting
with other BCM drivers), in which case the **gpio** library is a clear winner.

I could only recommend the Python interface for very low bandwidth applications.

