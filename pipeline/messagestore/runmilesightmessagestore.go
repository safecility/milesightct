package main

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/device/milesightct/pipeline/messagestore/helpers"
	"github.com/safecility/microservices/go/device/milesightct/pipeline/messagestore/server"
	"github.com/safecility/microservices/go/device/milesightct/pipeline/messagestore/store"
	"os"
)

func main() {

	ctx := context.Background()

	deployment, isSet := os.LookupEnv(helpers.OSDeploymentKey)
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

	gpsClient, err := pubsub.NewClient(ctx, config.ProjectName)
	defer func(gpsClient *pubsub.Client) {
		err := gpsClient.Close()
		if err != nil {
			log.Err(err).Msg("Error closing pubsub client")
		}
	}(gpsClient)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
	}
	if gpsClient == nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
		return // this is here so golang doesn't complain about gpsClient being possibly nil
	}

	milesightSubscription := gpsClient.Subscription(config.Subscriptions.Milesight)
	exists, err := milesightSubscription.Exists(ctx)
	if !exists {
		log.Fatal().Str("subscription", config.Subscriptions.Milesight).Msg("no eastron subscription")
	}

	dsClient, err := datastore.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not start service")
	}
	d, err := store.NewDatastoreMilesite(dsClient)

	if err != nil {
		log.Fatal().Err(err).Msg("could not get datastore milesight")
	}

	hotDropServer := server.NewMilesightServer(d, milesightSubscription, config.StoreAll)
	hotDropServer.Start()
}
