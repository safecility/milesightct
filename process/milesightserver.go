package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/iot/devices/milesightct/process/helpers"
	"github.com/safecility/iot/devices/milesightct/process/server"
	"github.com/safecility/iot/devices/milesightct/process/store"
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
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
	}
	if gpsClient == nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
		return
	}
	defer func(gpsClient *pubsub.Client) {
		err = gpsClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close pubsub client")
		}
	}(gpsClient)

	uplinksSubscription := gpsClient.Subscription(config.Subscriptions.Uplinks)
	exists, err := uplinksSubscription.Exists(ctx)
	if !exists {
		log.Fatal().Str("subscription", config.Subscriptions.Uplinks).Msg("no uplinks subscription")
	}

	milesightTopic := gpsClient.Topic(config.Topics.Milesight)
	exists, err = milesightTopic.Exists(ctx)
	if !exists {
		log.Fatal().Str("topic", config.Topics.Milesight).Msg("no hotdrop topic")
	}
	if err != nil {
		log.Fatal().Err(err).Str("topic", config.Topics.Milesight).Msg("could not get topic")
	}
	defer milesightTopic.Stop()

	ds, err := helpers.GetStore(config)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get store")
	}
	defer func(ds store.DeviceStore) {
		err = ds.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close store")
		}
	}(ds)

	milesightServer := server.NewMilesightServer(ds, uplinksSubscription, milesightTopic, config.PipeAll)
	milesightServer.Start()

}
