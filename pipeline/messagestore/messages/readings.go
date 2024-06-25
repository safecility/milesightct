package messages

import (
	"github.com/safecility/go/lib"
	"time"
)

type PowerDevice struct {
	*lib.Device
	PowerFactor float64 `datastore:",omitempty"`
	Voltage     float64 `datastore:",omitempty"`
}

type MilesightCTReading struct {
	*PowerDevice
	UID     string
	Power   bool      `datastore:",omitempty"`
	Time    time.Time `datastore:",omitempty"`
	Version `datastore:",omitempty"`
	Current `datastore:",omitempty"`
}

type Alarms struct {
	t  bool
	tr bool
	r  bool
	rr bool
}

type Version struct {
	Ipso     string `datastore:",omitempty"`
	Hardware string `datastore:",omitempty"`
	Firmware string `datastore:",omitempty"`
}

type Current struct {
	Total  float32 `datastore:",omitempty"`
	Value  float32 `datastore:",omitempty"`
	Max    float32 `datastore:",omitempty"`
	Min    float32 `datastore:",omitempty"`
	Alarms `datastore:",omitempty"`
}
