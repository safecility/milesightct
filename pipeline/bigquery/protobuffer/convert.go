package protobuffer

import (
	"github.com/safecility/iot/devices/milesightct/pipeline/bigquery/messages"
	"time"
)

func CreateProtobufMessage(r *messages.MilesightCTReading) *Milesight {
	bq := &Milesight{
		DeviceUID:            r.UID,
		Time:                 r.Time.Format(time.RFC3339),
		AccumulatedCurrent:   float64(r.Current.Total),
		InstantaneousCurrent: float64(r.Current.Value),
		MaximumCurrent:       float64(r.Current.Max),
		MinimumCurrent:       float64(r.Current.Min),
	}
	if r.Device != nil {
		bq.DeviceUID = r.DeviceUID
	}
	return bq
}
