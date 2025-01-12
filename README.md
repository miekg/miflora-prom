# miflora-go

A Golang application for reading data from Xiaomi Mi Flora plant sensors and exporting the values
at prometheus metrics. Retrieves data from the sensor, such as battery, humidity, conductivity, soil
moisture and temperature.

Configuration is done via a small text file that has that lives in /etc/miflora.

    alias:mac-address

or

    alias:uuid

lines, the alias is used in the prometheus metrics. The following ones are exported:

* `miflora_meta_battery_percentage{name="<alias"}`
* `miflora_meta_firmware_version{name="<alias"}`
* `miflora_plant_illumination_lux{name="<alias"}`
* `miflora_plant_moisture_percentage{name="<alias"}`
* `miflora_plant_conductivity{name="<alias"}`

## Installation

1. Install Golang on your computer if you don't have it already installed
2. Clone the repository: `git clone https://github.com/darox/miflora-go`
3. Build the application: `cd cmd/miflora-go && go build`
4. Add capabilities to run as none-root user: `sudo setcap 'cap_net_raw,cap_net_admin+eip' miflora-go`

## Usage

Under Linux, the application uses the mac address to connect to devices; under MacOs the UUID.

1. Copy `config/config.yaml` to the same folder where miflora-go will run or use the param `--config-path` to specify the path of the config file
2. The application will now scan the device and printout the result

### Acknowledgments

- Xiaomi for creating the sensor
- [Creators of go-ble](https://github.com/go-ble/ble)
- [Creators of miflora wiki](https://github.com/ChrisScheffler/miflora/wiki/The-Basics)
