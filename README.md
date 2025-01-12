# NAME

miflora-prom - generate prometheus metrics from Xiaomi Mi flower care plant sensor.

# SYNOPSIS

**miflora-prom**

# DESCRIPTION

miflora-prom is a small application that connects to a Xiami Mi flower plant sensor over blue tooth
and converts the sensor data into prometheus metrics. If needs a config file in `/etc/miflora` to
operate.

The data retrieved and the prometheus metrics are:

* battery level, `miflora_meta_battery_percentage{name="<alias>"}`
* firmware version, `miflora_meta_firmware_version{name="<alias>", version="<version>"}`
* illumination level in Lux, `miflora_plant_illumination_lux{name="<alias>"}`
* moisture percentage, `miflora_plant_moisture_percentage{name="<alias>"}`
* ground conductivity, `miflora_plant_conductivity{name="<alias>"}`

When running as a non-root user the following capabilities are needed:
'cap_net_raw,cap_net_admin+eip' for `miflora-prom` to accessing bluetooth.

Under Linux, the application uses the mac address to connect to devices; under MacOs the UUID.

## Config file

The configuration file contains lines constisting of a `LHS <colon> RHS`. It defines the MAC
addreses of UUIDs of the sensors to be queried and two other variables:

    # adapter, defaults to 'default'
    adapter: default
    # how often to query the sensors
    duration: 1h
    myfirstsensor: 422b23155c369dfee0aea210d1a9bc37
    mysecondsensor: ...

# Acknowledgments

- Xiaomi for creating the sensor.
- [Creators of go-ble](https://github.com/go-ble/ble).
- [Creators of miflora wiki](https://github.com/ChrisScheffler/miflora/wiki/The-Basics).
