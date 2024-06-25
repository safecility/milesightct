package store

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/microservices/go/device/milesightct/pipeline/messagestore/messages"
)

type DatastoreMilesight struct {
	client *datastore.Client
}

func NewDatastoreMilesite(client *datastore.Client) (*DatastoreMilesight, error) {
	rd := &DatastoreMilesight{client: client}
	return rd, nil
}

func (d *DatastoreMilesight) AddMilesightMessage(m *messages.MilesightCTReading) error {
	ctx := context.Background()
	k := datastore.IncompleteKey("MilesightCT", nil)
	k, err := d.client.Put(ctx, k, m)
	if err != nil {
		return err
	}
	log.Debug().Str("uid", m.UID).Msg("putting new milesight ct message")
	return nil
}
