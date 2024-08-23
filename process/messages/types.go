package messages

import (
	"github.com/safecility/go/lib"
)

type PowerProfile struct {
	PowerFactor float64 `datastore:"-" firestore:"powerFactor,omitempty"`
	Voltage     float64 `datastore:"-" firestore:"voltage,omitempty"`
}

type PowerDevice struct {
	lib.Device
	Profile *PowerProfile `datastore:"-" firestore:"profile,omitempty"`
}
