# PirageKit

Simple [HomeKit](https://www.apple.com/ios/home/) enabled [Go](https://golang.org)
program targeted for the [Raspberry Pi](https://www.raspberrypi.org) platform to
control garage doors connected via GPIO.

## Why?

Seemed ridiculous to replace my existing garage door openers or to buy an expensive
add-on just so I can have Siri open my garage door. Since I already had a Pi
laying around, I figured adding a few reed switches and relays would not hurt.

## Requirements

1. Raspberry Pi with at least 3 GPIO pins available per door you want to control
2. 2 reed switches (NO)
3. Some form of a relay (I use a [Sain Smart relay](https://www.sainsmart.com/products/8-channel-5v-relay-module))

## Building

    $ VERSION=sem.version.number make

## Setup

Use the wiring diagram below to assist you in wiring up the Pi to your relay and
switches. If you are not using a Sain Smart relay, you will most likely not need
the circuit documented below or will need a different circuit.

Once everything is wired up, configure piragekit and run it. Launch Home on iOS
and tap + to add a new device, tap Add Accessory, tap "Don't Have a Code or Can't
Scan?". You should see PirageKit show up as a square. Tap it and enter the code
you configured earlier. 

### Wiring Diagram

![Wiring Diagram](wiring_diagram.png)
