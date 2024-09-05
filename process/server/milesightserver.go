package server

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/iot/devices/milesightct/process/messages"
	"github.com/safecility/iot/devices/milesightct/process/store"
	"net/http"
	"os"
	"strings"
)

type MilesightServer struct {
	cache          store.DeviceStore
	sub            *pubsub.Subscription
	milesightTopic *pubsub.Topic
	pipeAll        bool
}

func NewMilesightServer(cache store.DeviceStore, sub *pubsub.Subscription, eagle *pubsub.Topic, pipeAll bool) MilesightServer {
	return MilesightServer{sub: sub, cache: cache, milesightTopic: eagle, pipeAll: pipeAll}
}

func (es *MilesightServer) Start() {
	go es.receive()
	es.serverHttp()
}

func (es *MilesightServer) receive() {
	log.Debug().Str("sub", es.sub.String()).Msg("listening for messages")
	err := es.sub.Receive(context.Background(), func(ctx context.Context, message *pubsub.Message) {
		sm := &stream.SimpleMessage{}
		log.Debug().Str("data", fmt.Sprintf("%s", message.Data)).Msg("raw data")
		err := json.Unmarshal(message.Data, sm)
		message.Ack()
		if err != nil {
			log.Err(err).Msg("could not unmarshall data")
			return
		}

		mr, err := messages.ReadMilesightCT(sm.Payload)
		if err != nil {
			log.Err(err).Msg("could not read milesight CT")
			return
		}

		deviceUID := strings.Replace(sm.DeviceUID, "/", ":", 1)

		log.Debug().Str("messageID", sm.DeviceUID).Msg("milesight ct message")
		var pd *messages.PowerDevice
		if es.cache != nil {
			pd, err = es.cache.GetDevice(deviceUID)
			if err != nil {
				log.Warn().Err(err).Str("uid", sm.DeviceUID).Msg("could not get device")
			}
			if pd == nil {
				log.Debug().Str("uid", sm.DeviceUID).Msg("device not found")
			}
			mr.PowerDevice = pd
		}

		if mr.PowerDevice == nil && !es.pipeAll {
			log.Debug().Str("device", sm.DeviceUID).Msg("no device in cache and pipeAll == false")
			return
		}
		mr.Time = sm.Time
		// we thought we were getting this in the message but it's only in *some* messages
		mr.UID = sm.DeviceUID

		topic, err := stream.PublishToTopic(mr, es.milesightTopic)
		if err != nil {
			log.Err(err).Msg("could not publish usage to topic")
			return
		}
		log.Debug().Str("topic", *topic).Msg("published milesight ct to topic")
	})
	if err != nil {
		log.Err(err).Msg("could not receive from sub")
		return
	}
}

func (es *MilesightServer) serverHttp() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "started")
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("could write to http.ResponseWriter"))
		}
	})

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "running")
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("could write to http.ResponseWriter"))
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8089"
	}
	log.Debug().Msg(fmt.Sprintf("starting http server port %s", port))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start http")
	}
}
