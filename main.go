package main

import (
	"log"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
)

func main() {
	config, err := ParseConfig(DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	d, err := dev.NewDevice(config.Adapter)
	if err != nil {
		log.Fatal(err)
	}

	ble.SetDefaultDevice(d)

	r, err := Readings(config)
	if err != nil {
		log.Fatal(err)
	}

	// start the loopy loop

}
