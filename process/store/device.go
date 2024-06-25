package store

import (
	"github.com/safecility/iot/devices/milesightct/process/messages"
)

type DeviceStore interface {
	GetDevice(uid string) (*messages.PowerDevice, error)
}
