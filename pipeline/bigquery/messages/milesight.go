package messages

import (
	"cloud.google.com/go/bigquery"
	"time"
)

type MilesightCTReading struct {
	*PowerDevice
	UID     string
	Power   bool      `datastore:",omitempty"`
	Time    time.Time `datastore:",omitempty"`
	Current `datastore:",omitempty"`
}

func GetBigqueryTableMetadata(name string) *bigquery.TableMetadata {
	sampleSchema := bigquery.Schema{
		{Name: "DeviceUID", Type: bigquery.StringFieldType},
		{Name: "Time", Type: bigquery.TimestampFieldType},
		{Name: "AccumulatedCurrent", Type: bigquery.FloatFieldType},
		{Name: "InstantaneousCurrent", Type: bigquery.FloatFieldType},
		{Name: "MaximumCurrent", Type: bigquery.FloatFieldType},
		{Name: "MinimumCurrent", Type: bigquery.FloatFieldType},
	}

	return &bigquery.TableMetadata{
		Name:   name,
		Schema: sampleSchema,
	}
}

type AmpHour float32

type Current struct {
	Total AmpHour `datastore:",omitempty"`
	Value float32 `datastore:",omitempty"`
	Max   float32 `datastore:",omitempty"`
	Min   float32 `datastore:",omitempty"`
}
