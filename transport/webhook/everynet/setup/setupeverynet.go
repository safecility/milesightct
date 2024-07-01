package main

import (
	"cloud.google.com/go/pubsub"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/transports/everynet/helpers"
	"os"
	"time"
)

// main setup pubsub and output a jwt token to be used by everynet webhook
func main() {
	deployment, isSet := os.LookupEnv("Deployment")
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

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
	ourSecrets := setup.GetNewSecrets(config.ProjectName, secretsClient)
	sigSecret, err := ourSecrets.GetSecret(config.Secret)

	sig := hmac.New(sha256.New, sigSecret)

	hmacSecret := sig.Sum(nil)

	now := time.Now()
	expires := now.Add(config.ExpiresHours * time.Hour)
	log.Debug().Time("expires", expires).Msg("expires on")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"service": config.ApplicationName,
		"uplinks": config.Topics.Uplinks,
		"created": now.Format(time.RFC3339),
		"exp":     expires.Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		log.Error().Err(err).Msg("could not generate token")
		return
	}

	fo, err := os.Create(fmt.Sprintf("%s-%s", config.ApplicationName, "jwt.txt"))
	if err != nil {
		log.Error().Err(err).Msg("could not create file")
	}
	defer func() {
		if err := fo.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close output file")
		}
	}()
	write, err := fo.WriteString(tokenString)
	if err != nil {
		log.Error().Err(err).Msg("could not write to file")
	} else {
		log.Info().Msgf("wrote %d bytes", write)
	}

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if token == nil || !token.Valid {
		log.Error().Err(err).Msg("invalid token")
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Printf("token OK %+v \n", claims)
	} else {
		fmt.Println(err)
	}

	gpsClient, err := pubsub.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup pubsub")
	}

	uplinkTopic := gpsClient.Topic(config.Topics.Uplinks)
	exists, err := uplinkTopic.Exists(ctx)
	if !exists {
		uplinkTopic, err = gpsClient.CreateTopic(ctx, config.Topics.Uplinks)
		if err != nil {
			log.Fatal().Err(err).Msg("setup could not create topic")
		}

		log.Info().Str("topic", uplinkTopic.String()).Msg("created topic")
	}
	log.Info().Msg("setup complete")
}
