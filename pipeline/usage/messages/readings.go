package messages

import (
	"fmt"
	"github.com/rs/zerolog/log"
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
	UID   string
	Power bool
	Time  time.Time
	Version
	Current
}

type MeterReading struct {
	*lib.Device
	ReadingKWH float64
	Time       time.Time
}

func (mc MilesightCTReading) Usage() (*MeterReading, error) {
	if mc.PowerDevice == nil {
		return nil, fmt.Errorf("device does not have its PowerDevice definitions")
	}
	if mc.Current.Total == 0 {
		log.Info().Str("reading", fmt.Sprintf("%+v", mc)).Msg("zero usage - check device is new")
	}
	kWh := float64(mc.Current.Total) * mc.Voltage * mc.PowerFactor / 1000.0
	mr := &MeterReading{
		ReadingKWH: kWh,
		Time:       mc.Time,
	}
	if mc.PowerDevice.Device == nil {
		log.Warn().Str("UID", mc.UID).Msg("device does not have device definitions")
		mr.Device = &lib.Device{
			DeviceUID: mc.UID,
		}
	} else {
		mr.Device = mc.Device
	}

	return mr, nil
}

type Alarms struct {
	t  bool
	tr bool
	r  bool
	rr bool
}

type Version struct {
	Ipso     string
	Hardware string
	Firmware string
}

type Current struct {
	Total float32
	Value float32
	Max   float32
	Min   float32
	Alarms
}
