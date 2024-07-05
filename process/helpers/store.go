package helpers

import (
	"cloud.google.com/go/firestore"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/iot/devices/milesightct/process/store"
)

type stores struct {
	config  *Config
	secrets *setup.Secrets
}

func GetStore(config *Config) (store.DeviceStore, error) {

	if config.Store.Rest != nil {
		return store.CreateDeviceClient(config)
	}

	s := stores{
		config: config,
	}
	defer s.close()

	db := s.getDb()

	if config.Store.Cache != nil {
		var key []byte
		var err error
		if config.Store.Cache.Secret != nil {
			key, err = s.getSecrets().GetSecret(*config.Store.Cache.Secret)
			return nil, err
		}
		rClient := redis.NewClient(&redis.Options{
			Addr:     config.Store.Cache.Address(),
			Password: string(key),
			DB:       0, // use default DB
		})
		return store.NewDeviceCache(rClient, db, config.ContextDeadline), nil
	}
	return db, nil
}

func (st stores) getSecrets() *setup.Secrets {
	if st.secrets != nil {
		return st.secrets
	}
	ctx := context.Background()
	secretsClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create secrets client")
	}
	defer func(secretsClient *secretmanager.Client) {
		err = secretsClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close secrets client")
		}
	}(secretsClient)
	secrets := setup.GetNewSecrets(st.config.ProjectName, secretsClient)
	st.secrets = secrets
	return secrets
}

func (st stores) close() error {
	if st.secrets != nil {
		return st.secrets.Close()
	}
	return nil
}

func (st stores) getDb() store.DeviceStore {
	if st.config.Store.Firestore != nil {
		return st.getFirestore()
	}
	return st.getSql()
}

func (st stores) getFirestore() *store.DeviceFirestore {
	ctx := context.Background()

	id := st.config.Store.Firestore.Database

	fsClient, err := firestore.NewClientWithDatabase(ctx, st.config.ProjectName, *id)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create Firestore client")
	}

	fs := store.NewDeviceFirestore(fsClient, st.config.ContextDeadline)

	return fs
}

func (st stores) getSql() store.DeviceStore {

	password, err := st.getSecrets().GetSecret(st.config.Store.Sql.Secret)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get secret")
	}
	st.config.Store.Sql.Config.Password = string(password)

	s, err := setup.NewSafecilitySql(st.config.Store.Sql.Config)
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup safecility sql")
	}
	c, err := store.NewDeviceSql(s)
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup safecility device sql")
	}
	return c
}
