package main

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/gbigquery"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/go/setup"
	"github.com/safecility/iot/devices/milesightct/pipeline/bigquery/helpers"
	"os"
	"time"
)

func main() {
	deployment, isSet := os.LookupEnv(helpers.OSDeploymentKey)
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to BigQuery")
	}
	defer func(client *bigquery.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close bigquery.Client")
		}
	}(client)

	bqc := gbigquery.NewBQTable(client)

	t, err := bqc.CheckOrCreateBigqueryTable(&config.BigQuery)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create BigQuery table")
	}

	log.Info().Msg("finished BigQuery setup")

	sClient, err := pubsub.NewSchemaClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create schema client")
	}
	defer func(sClient *pubsub.SchemaClient) {
		err := sClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("could not close schema client")
		}
	}(sClient)

	schema, err := sClient.Schema(ctx, config.BigQuery.Schema.Name, pubsub.SchemaViewFull)
	if err != nil || schema == nil {
		schema, err = gbigquery.CreateProtoSchema(sClient, config.BigQuery.Schema.Name, config.BigQuery.Schema.FilePath)
		if err != nil {
			log.Fatal().Err(err).Msg("could not create schema")
		}
	}

	gpsClient, err := pubsub.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup pubsub")
	}

	bigqueryTopic := gpsClient.Topic(config.Pubsub.Topics.Bigquery)
	exists, err := bigqueryTopic.Exists(ctx)
	if !exists {
		bigqueryTopic, err = gbigquery.CreateBigqueryTopic(gpsClient, config.Pubsub.Topics.Bigquery, schema)
		if err != nil {
			log.Fatal().Str("sub", config.Pubsub.Subscriptions.BigQuery).Err(err).Msg("could not create bigquery topic")
		}
		log.Info().Msg("bigquery topic created")
	}
	bigQuerySubscription := gpsClient.Subscription(config.Pubsub.Subscriptions.BigQuery)
	exists, err = bigQuerySubscription.Exists(ctx)
	if !exists {
		err = gbigquery.CreateBigQuerySubscription(gpsClient, config.Pubsub.Subscriptions.BigQuery, t.FullID, bigqueryTopic)
		if err != nil {
			log.Fatal().Err(err).Msg("could not create bigquery subscription")
		}
		log.Info().Msg("created bigquery subscription")
	}

	milesightSubscription := gpsClient.Subscription(config.Pubsub.Subscriptions.Milesight)
	exists, err = milesightSubscription.Exists(ctx)
	if !exists {
		milesightTopic := gpsClient.Topic(config.Pubsub.Topics.Milesight)
		if exists, err = milesightTopic.Exists(ctx); err != nil {
			log.Fatal().Err(err).Msg("could not check if milesight topic exists")
		}
		if !exists {
			milesightTopic, err = gbigquery.CreateBigqueryTopic(gpsClient, config.Pubsub.Topics.Milesight, schema)
			if err != nil {
				log.Fatal().Err(err).Msg("could not create milesight topic")
			}
			log.Info().Msg("created milesight topic")
		}

		r, err := time.ParseDuration("1h")
		if err != nil {
			log.Fatal().Err(err).Msg("could not parse duration")
		}
		subConfig := stream.GetDefaultSubscriptionConfig(milesightTopic, r)
		milesightSubscription, err = gpsClient.CreateSubscription(ctx, config.Pubsub.Subscriptions.Milesight, subConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("setup could not create subscription")
		}
		log.Info().Str("sub", config.Pubsub.Subscriptions.Milesight).Msg("created subscription")
	}
	log.Info().Msg("finished pubsub setup")

}
