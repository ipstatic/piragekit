package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/log"
	ga "github.com/ipstatic/piragekit/accessory"
	"github.com/stianeikeland/go-rpio"
)

var (
	configFile = flag.String("config.file", "piragekit.yml", "configuration file")
	verbose    = flag.Bool("verbose", false, "verbose output")
	version    string
)

func main() {
	flag.Parse()
	if *verbose {
		log.Debug.Enable()
	}
	log.Info.Printf("PirageKit %s starting up...", version)

	config, err := loadConfig(*configFile)
	if err != nil {
		log.Info.Fatalf("Config error, aborting: %s", err)
	}

	err = rpio.Open()
	if err != nil {
		log.Info.Fatalf("Could not get GPIO state, aborting: %s", err)
	}
	defer rpio.Close()

	accessories := []*accessory.Accessory{}
	c := make(chan os.Signal, 1)

	for _, door := range config.Doors {
		log.Debug.Printf("Subscribing to %s door events", door.Name)
		info := accessory.Info{
			Name:         door.Name,
			Manufacturer: door.Manufacturer,
			Model:        door.Model,
		}
		acc := ga.NewGarageDoor(info)
		accessories = append(accessories, acc.Accessory)

		go subscribeEvents(door, acc, c)
	}

	transportConfig := hc.Config{
		Pin:         config.HomeKit.PIN,
		StoragePath: config.HomeKit.StoragePath,
	}

	transport, err := hc.NewIPTransport(transportConfig, ga.NewBridge().Accessory, accessories...)
	if err != nil {
		log.Info.Fatalf("HomeKit transport error, aborting: %s", err)
	}

	// trap Ctrl+C and call cancel on the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
			log.Debug.Printf("shutting down...")
			<-transport.Stop()
		case <-ctx.Done():
		}
	}()

	transport.Start()
}
