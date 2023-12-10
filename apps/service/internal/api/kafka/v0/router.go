package v0

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kitanoyoru/kitaDriveBot/apps/service/internal/config"
	"github.com/kitanoyoru/kitaDriveBot/apps/service/internal/service"
)

type KafkaAPI struct {
	config  *config.KafkaConfig
	service *service.Service
	logger  watermill.LoggerAdapter // TODO: change to the zap logger
}

func NewKafkaAPI(config *config.KafkaConfig, service *service.Service) *KafkaAPI {
	logger := watermill.NewStdLogger(false, false)
	return &KafkaAPI{
		config,
		service,
		logger,
	}
}

func (api *KafkaAPI) GetRouter() (*message.Router, error) {
	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:     api.config.BrokerList,
			Unmarshaler: kafka.DefaultMarshaler{},
		},
		api.logger,
	)
	if err != nil {
		return nil, err
	}

	router, err := message.NewRouter(message.RouterConfig{}, api.logger)
	if err != nil {
		return nil, err
	}

	// Add handlers...

	return router, nil
}
