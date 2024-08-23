package messages

import (
	"github.com/safecility/go/lib"
	"time"
)

type UUID string

type Structure struct {
	SystemUID string `datastore:",omitempty"`
	TenantUID string `datastore:",omitempty"`
}

type PowerProfile struct {
	PowerFactor float64 `datastore:",omitempty" json:"powerFactor"`
	Voltage     float64 `datastore:",omitempty" json:"voltage"`
}

type PowerDevice struct {
	lib.Device
	Structure *Structure    `datastore:",flatten"`
	Profile   *PowerProfile `datastore:",flatten"`
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
