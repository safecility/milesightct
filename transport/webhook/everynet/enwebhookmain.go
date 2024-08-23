package main

import (
	"cloud.google.com/go/pubsub"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/transports/everynet/helpers"
	"github.com/safecility/microservices/go/transports/everynet/server"
	"os"
)

func main() {
	deployment, isSet := os.LookupEnv(helpers.OSDeploymentKey)
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

	ctx := context.Background()
	secretsClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create secrets client")
	}
	secrets := setup.GetNewSecrets(config.ProjectName, secretsClient)
	jwtSecret, err := secrets.GetSecret(config.Secret)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get secret")
	}

	psClient, err := pubsub.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get pubsub client")
	}
	defer func(psClient *pubsub.Client) {
		err := psClient.Close()
		if err != nil {
			log.Err(err).Msg("Failed to close pubsub client")
		}
	}(psClient)

	uplinksTopic := psClient.Topic(config.Topics.Uplinks)

	exists, err := uplinksTopic.Exists(ctx)
	if !exists || err != nil {
		log.Fatal().Err(err).Str("name", config.Topics.Uplinks).Msg("topic does not exist or error")
	}

	jwtParser := lib.NewJWTParser(string(jwtSecret))

	s := server.NewEverynetServer(&jwtParser, uplinksTopic)

	s.Start()
}
