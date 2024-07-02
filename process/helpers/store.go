package helpers

import (
	"cloud.google.com/go/firestore"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/iot/devices/milesightct/process/store"
)

func GetStore(config *Config) store.DeviceStore {
	if config.Store.Firestore != nil {
		return getFirestore(config)
	}
	return getSql(config)
}

func getFirestore(config *Config) *store.DeviceFirestore {
	ctx := context.Background()

	id := config.Store.Firestore.Database

	fsClient, err := firestore.NewClientWithDatabase(ctx, config.ProjectName, *id)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create Firestore client")
	}

	fs := store.NewDeviceFirestore(fsClient, config.Store.Firestore.Deadline)

	return fs
}

func getSql(config *Config) store.DeviceStore {
	ctx := context.Background()
	secretsClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create secrets client")
	}
	defer func(secretsClient *secretmanager.Client) {
		err := secretsClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close secrets client")
		}
	}(secretsClient)
	sqlSecret := setup.GetNewSecrets(config.ProjectName, secretsClient)
	password, err := sqlSecret.GetSecret(config.Store.Sql.Secret)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get secret")
	}
	config.Store.Sql.Config.Password = string(password)

	s, err := setup.NewSafecilitySql(config.Store.Sql.Config)
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup safecility sql")
	}
	c, err := store.NewDeviceSql(s)
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup safecility device sql")
	}
	return c
}
