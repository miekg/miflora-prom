package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-ble/ble"
)

func Readings(c Config) (readings []Reading, err error) {
	for alias, device := range c.Devices {
		log.Printf("Querying %s %s", alias, device)

		filter := func(a ble.Advertisement) bool {
			return strings.EqualFold(a.Addr().String(), device)
		}

		d := Device{Device: device, Identifier: alias}

		ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), 2*time.Second))
		cln, err := ble.Connect(ctx, filter)
		if err != nil {
			log.Printf("Failed to scan %s: %s", d.Identifier, err)
			continue
		}

		log.Printf("Connected to %s", d.Identifier)

		p, err := cln.DiscoverProfile(true)
		if err != nil {
			cln.CancelConnection()
			log.Printf("Failed to discover profile: %s", err)
		}

		err = d.enableSensorReadings(cln, p)
		if err != nil {
			log.Printf("Failed to enable sensor read: %s", err)
			continue
		}

		sys, err := systemReadings(cln, p)
		if err != nil {
			log.Printf("Failed to read system: %s", err)
		}

		sensor, err := sensorReadings(cln, p)
		if err != nil {
			log.Printf("Failed to read sensors: %s", err)
		}

		r := Reading{Alias: alias, System: sys, Sensor: sensor}
		readings = append(readings, r)
		if err = cln.CancelConnection(); err != nil {
			log.Printf("Failed to cancel connection: %s", err)
		}
		log.Printf("Disconnected from %s", d.Identifier)

	}
	return readings, nil
}

// Enable read of temperature, humidity, light and conductivity.
func (device *Device) enableSensorReadings(cln ble.Client, p *ble.Profile) (err error) {
	// UUID service and characteristic to enable read
	cu := ble.MustParse("00001a0000001000800000805f9b34fb")
	su := ble.MustParse("0000120400001000800000805f9b34fb")
	// bytes to enable read
	enableSensorReadingsBytes := []byte{0xa0, 0x1f}

	s := findService(p, su)
	if s == nil {
		return fmt.Errorf("service not found")
	}
	c := findCharacteristic(s, cu)
	if c == nil {
		return fmt.Errorf("characteristic not found")
	}
	if err = cln.WriteCharacteristic(c, enableSensorReadingsBytes, false); err != nil {
		return err
	}
	return nil
}

// Read the characteristics for battery level and firmware version.
func systemReadings(cln ble.Client, p *ble.Profile) (sys System, err error) {
	// UUID service and characteristic to read battery level and firmware version
	cu := ble.MustParse("00001a0200001000800000805f9b34fb")
	su := ble.MustParse("0000120400001000800000805f9b34fb")

	// find service and characteristic
	s := findService(p, su)
	if s == nil {
		return sys, fmt.Errorf("service not found")
	}

	c := findCharacteristic(s, cu)
	if c == nil {
		return sys, fmt.Errorf("characteristic not found")
	}

	b, err := cln.ReadCharacteristic(c)
	if err != nil {
		return sys, err
	}

	return System{Battery: uint16(b[0]), Firmware: string(b[2:])}, nil
}

// Read the characteristic to get conductivity, humidity, light and temperature sensor data.
func sensorReadings(cln ble.Client, p *ble.Profile) (sensor Sensor, err error) {
	// UUID of service and characteristic holding the sensor data
	cu := ble.MustParse("00001a0100001000800000805f9b34fb")
	su := ble.MustParse("0000120400001000800000805f9b34fb")

	s := findService(p, su)
	if s == nil {
		return sensor, fmt.Errorf("service not found")
	}

	c := findCharacteristic(s, cu)
	if c == nil {
		return sensor, fmt.Errorf("characteristic not found")
	}

	if (c.Property & ble.CharRead) == 0 {
		return sensor, nil
	}

	b, err := cln.ReadCharacteristic(c)
	if err != nil {
		return sensor, err
	}
	var subtrahend float32 = 10.0
	sensor = Sensor{
		Conductivity: binary.LittleEndian.Uint16(b[8:10]),
		Moisture:     uint16(b[7]),
		Illumination: binary.LittleEndian.Uint32(b[3:7]),
		Temperature:  float32(binary.LittleEndian.Uint16(b[0:2])) / subtrahend,
	}
	return sensor, nil
}

// Find the service based on the UUID.
func findService(p *ble.Profile, u ble.UUID) *ble.Service {
	for _, s := range p.Services {
		if s.UUID.Equal(u) {
			return s
		}
	}
	return nil
}

// Find the characteristic based on the UUID.
func findCharacteristic(s *ble.Service, u ble.UUID) *ble.Characteristic {
	for _, c := range s.Characteristics {
		if c.UUID.Equal(u) {
			return c
		}
	}
	return nil
}

// Reading is a struct that holds the readings from a single device
type Reading struct {
	Alias string
	System
	Sensor
}

// Device is a struct that holds the device information
type Device struct {
	Device     string
	Identifier string
}

// System is a struct that holds the battery and firmware information.
type System struct {
	Battery  uint16
	Firmware string
}

// Sensors is a struct that holds the conductivity, moisture, light and temperature information
type Sensor struct {
	Conductivity uint16
	Moisture     uint16
	Illumination uint32
	Temperature  float32
}
