package messages

import (
	"github.com/safecility/go/lib"
)

type PowerProfile struct {
	PowerFactor float64 `datastore:"-" firestore:"powerFactor,omitempty" json:"powerFactor,omitempty"`
	Voltage     float64 `datastore:"-" firestore:"voltage,omitempty" json:"voltage,omitempty"`
}

type PowerDevice struct {
	lib.Device
	Profile *PowerProfile `datastore:"-" firestore:"profile,omitempty" json:"profile,omitempty"`
}
