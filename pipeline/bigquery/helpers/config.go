package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	OSDeploymentKey = "MILESIGHT_DEPLOYMENT"
)

type Config struct {
	ProjectName string `json:"projectName"`
	Pubsub      struct {
		Topics struct {
			Milesight string `json:"milesight"`
			Bigquery  string `json:"bigquery"`
		} `json:"topics"`
		Subscriptions struct {
			BigQuery  string `json:"bigquery"`
			Milesight string `json:"milesight"`
		} `json:"subscriptions"`
	} `json:"pubsub"`
	BigQuery struct {
		Dataset string `json:"dataset"`
		Table   string `json:"table"`
		Schema  struct {
			Name     string `json:"name"`
			FilePath string `json:"filePath"`
		} `json:"schema"`
	} `json:"bigQuery"`
	StoreAll bool `json:"storeAll"`
}

// GetConfig creates a config for the specified deployment
func GetConfig(deployment string) *Config {
	fileName := fmt.Sprintf("%s-config.json", deployment)

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not find config file")
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Err(err).Msg("could not close config during defer")
		}
	}(file)
	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		log.Fatal().Err(err).Str("filename", fileName).Msg("could not decode pubsub config")
	}
	return config
}
