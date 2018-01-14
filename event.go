package main

import (
	"os"
	"time"

	"github.com/brutella/hc/log"
	ga "github.com/ipstatic/piragekit/accessory"
	"github.com/stianeikeland/go-rpio"
)

const (
	open    = 0
	closed  = 1
	opening = 2
	closing = 3
	stopped = 4
)

type Door struct {
	Name    string
	Current int
	Target  int
	Top     rpio.Pin
	Bottom  rpio.Pin
	Relay   rpio.Pin
}

// currentState will read the current value of the reed switches and associate
// it with a state
func currentState(d *Door) {
	previous := d.Current
	top := d.Top.Read()
	bottom := d.Bottom.Read()

	if top == 1 && bottom == 0 {
		d.Current = closed
	} else if top == 0 && bottom == 1 {
		d.Current = open
	} else if top == 1 && bottom == 1 && (d.Current == closed || d.Current == opening) {
		d.Current = opening
	} else if top == 1 && bottom == 1 && (d.Current == open || d.Current == closing) {
		d.Current = closing
	} else {
		d.Current = stopped
	}

	if previous != d.Current {
		log.Info.Printf("current door state of %s: %v", d.Name, d.Current)
	}
}

// rectifyState changes the current state to the target state
func rectifyState(d *Door, acc *ga.GarageDoor, r chan bool) {
	for d.Current != d.Target {
		// refresh current state
		r <- true
		currentState(d)
		acc.GarageDoorOpener.CurrentDoorState.SetValue(d.Current)

		// if current state is stopped, we need to just exit as something is wrong
		if d.Current == stopped {
			log.Info.Printf("%s is in a stopped state, not performing asked action", d.Name)
			break
		}

		// if current state does not equal the in progress action of the target state
		// we need to trigger it to do so
		action := d.Target + 2
		if d.Current != action {
			log.Info.Printf("triggering relay on %s", d.Name)
			d.Relay.Toggle()
			time.Sleep(time.Second)
			d.Relay.Toggle()
		}

		time.Sleep(time.Second * 2)
	}

	// unblock local polling for changes
	r <- false
}

// subscribeEvents both listens for remote state as well as tracks local state
func subscribeEvents(dc DoorConfig, acc *ga.GarageDoor, c chan os.Signal) {
	// assign GPIO pins to the Door struct
	d := Door{Name: dc.Name}
	d.Top = rpio.Pin(dc.TopGPIO)
	d.Bottom = rpio.Pin(dc.BottomGPIO)
	d.Relay = rpio.Pin(dc.RelayGPIO)

	// set reed switches as inputs and the relay as an output
	d.Top.Input()
	d.Bottom.Input()
	d.Relay.Output()

	// setup "r" channel to block local polling when remote events occur
	r := make(chan bool)

	// set initial state and ensure Home knows it
	d.Current = closed
	d.Target = closed
	acc.GarageDoorOpener.CurrentDoorState.SetValue(d.Current)
	acc.GarageDoorOpener.TargetDoorState.SetValue(d.Target)

	// subscribe to remote updates from Home
	acc.GarageDoorOpener.TargetDoorState.OnValueRemoteUpdate(func(state int) {
		if state != d.Current {
			d.Target = state
			rectifyState(&d, acc, r)
		}
	})

	// subscribe to local updates and push them to Home
	// listening on the c channel in case we need to terminate
	// listening on the r channel in case we need to suspend until remote state
	// has been met
	for {
		select {
		case <-c:
			return
		case <-r:
			log.Info.Printf("remote state change occuring, suspending local monitoring of %s", d.Name)
			continue
		default:
			currentState(&d)
			if d.Current != acc.GarageDoorOpener.CurrentDoorState.GetValue() {
				log.Info.Printf("local state different from remote state for %s, updating remote to %v", d.Name, d.Current)
				acc.GarageDoorOpener.CurrentDoorState.SetValue(d.Current)
			}

			if d.Current == opening || d.Current == closing {
				d.Target = d.Current - 2

				if d.Target != acc.GarageDoorOpener.TargetDoorState.GetValue() {
					log.Info.Printf("local target different from remote state for %s, updating remote to %v", d.Name, d.Target)
					acc.GarageDoorOpener.TargetDoorState.SetValue(d.Target)
				}
			}
		}
	}
}
