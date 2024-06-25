package sections

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/iot/devices/milesightct/pipeline/bigquery/helpers"
	"os"
	"time"
)

func SetupPubsub(config *helpers.Config, t *bigquery.TableMetadata) {

	ctx := context.Background()
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
		schema, err = createProtoSchema(sClient, config.BigQuery.Schema.Name, config.BigQuery.Schema.FilePath)
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
		bigqueryTopic, err = createBigqueryTopic(gpsClient, config.Pubsub.Topics.Bigquery, schema)
		if err != nil {
			log.Fatal().Str("sub", config.Pubsub.Subscriptions.BigQuery).Err(err).Msg("could not create bigquery topic")
		}
		log.Info().Msg("bigquery topic created")
	}
	bigQuerySubscription := gpsClient.Subscription(config.Pubsub.Subscriptions.BigQuery)
	exists, err = bigQuerySubscription.Exists(ctx)
	if !exists {
		err = createBigQuerySubscription(gpsClient, config.Pubsub.Subscriptions.BigQuery, t.FullID, bigqueryTopic)
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
			milesightTopic, err = createBigqueryTopic(gpsClient, config.Pubsub.Topics.Milesight, schema)
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
	}
	log.Info().Msg("finished pubsub setup")
}

// createProtoSchema creates a schema resource from a schema proto file.
func createProtoSchema(client *pubsub.SchemaClient, schemaID, protoFile string) (*pubsub.SchemaConfig, error) {
	protoSource, err := os.ReadFile(protoFile)
	if err != nil {
		return nil, fmt.Errorf("error reading from file: %s", protoFile)
	}

	config := pubsub.SchemaConfig{
		Type:       pubsub.SchemaProtocolBuffer,
		Definition: string(protoSource),
	}

	ctx := context.Background()
	s, err := client.CreateSchema(ctx, schemaID, config)
	if err != nil {
		return nil, fmt.Errorf("CreateSchema: %w", err)
	}
	log.Debug().Str("schema", s.Name).Msg("Schema created")
	return s, nil
}

func createBigqueryTopic(client *pubsub.Client, topicName string, schema *pubsub.SchemaConfig) (*pubsub.Topic, error) {
	ctx := context.Background()
	bigqueryTopic, err := client.CreateTopicWithConfig(ctx, topicName, &pubsub.TopicConfig{
		SchemaSettings: &pubsub.SchemaSettings{
			Schema:   schema.Name,
			Encoding: pubsub.EncodingBinary,
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("setup could not create topic")
	}
	log.Info().Str("topic", bigqueryTopic.String()).Msg("created topic")

	return bigqueryTopic, nil
}

// createBigQuerySubscription creates a Pub/Sub subscription that exports messages to BigQuery.
func createBigQuerySubscription(client *pubsub.Client, subscriptionName, table string, topic *pubsub.Topic) error {
	ctx := context.Background()

	sub, err := client.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{
		Topic: topic,
		BigQueryConfig: pubsub.BigQueryConfig{
			Table:             table,
			WriteMetadata:     false,
			UseTopicSchema:    true,
			DropUnknownFields: true,
		},
	})
	if err != nil {
		return err
	}
	log.Debug().Str("subscription", sub.ID()).Msg("Created BigQuery subscription")

	return nil
}
