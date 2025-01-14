package main

import (
	"bufio"
	"os"
	"strings"
	"time"
)

const DefaultConfig = "/etc/miflora"

// Config holds the configration for our run. It's populated with some default values: adapter is "default".
// And duration is set to 1 * time.Hour.
type Config struct {
	Adapter  string
	Duration time.Duration
	Devices  map[string]string
}

func ParseConfig(file string) (Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	c := Config{Adapter: "default", Duration: time.Hour, Devices: map[string]string{}}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		left, right, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		left = strings.TrimSpace(left)
		right = strings.TrimSpace(right)

		switch left {
		case "adapter":
			c.Adapter = right
		case "duration":
			c.Duration, err = time.ParseDuration(right)
			if err != nil {
				return c, err
			}
		default:
			c.Devices[left] = right
		}
	}
	return c, scanner.Err()
}
