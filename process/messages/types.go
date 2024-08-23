package messages

import (
	"github.com/safecility/go/lib"
)

type PowerProfile struct {
	PowerFactor float64 `datastore:"-" firestore:",omitempty"`
	Voltage     float64 `datastore:"-" firestore:",omitempty"`
}

type PowerDevice struct {
	lib.Device
	Profile *PowerProfile `datastore:"-" firestore:",omitempty"`
}
