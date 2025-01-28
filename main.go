package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config, err := ParseConfig(DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}
	if len(config.Devices) == 0 {
		log.Fatalf("No devices found in config: %s", DefaultConfig)
	}

	log.Printf("Found %d device(s) from the config", len(config.Devices))

	d, err := dev.NewDevice(config.Adapter)
	if err != nil {
		log.Fatal(err)
	}
	ble.SetDefaultDevice(d)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":9095", nil)
	}()

	go func() {
		i := 0
		for {
			if i > 0 {
				time.Sleep(config.Duration)
			}
			i++

			readings, err := Readings(config)
			if err != nil {
				log.Printf("Failed to get readings: %s", err)
				// metrics here too
				continue
			}

			for _, r := range readings {
				mifloraBattery.WithLabelValues(r.Alias).Set(float64(r.System.Battery))
				mifloraFirmware.WithLabelValues(r.Alias, r.System.Firmware).Set(1)

				mifloraIllumination.WithLabelValues(r.Alias).Set(float64(r.Sensor.Illumination))
				mifloraMoisture.WithLabelValues(r.Alias).Set(float64(r.Sensor.Moisture))
				mifloraConductivity.WithLabelValues(r.Alias).Set(float64(r.Sensor.Conductivity))
				mifloraTemperature.WithLabelValues(r.Alias).Set(float64(r.Sensor.Temperature))
			}
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}

var (
	mifloraBattery = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "miflora_meta_battery_percentage",
		Help: "The battery level in percent",
	}, []string{"name"})
	mifloraFirmware = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "miflora_meta_firmware_version",
		Help: "The version of the firmware",
	}, []string{"name", "version"})
	mifloraIllumination = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "miflora_plant_illumination_lux",
		Help: "The current illumination in lux",
	}, []string{"name"})
	mifloraMoisture = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "miflora_plant_moisture_percentage",
		Help: "The current moisture level in percent",
	}, []string{"name"})
	mifloraConductivity = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "miflora_plant_conductivity",
		Help: "The current conductivity level in µS/cm",
	}, []string{"name"})
	mifloraTemperature = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "miflora_plant_temperature_celsius",
		Help: "The current temperature in °C",
	}, []string{"name"})
)
