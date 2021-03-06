package character

import (
	"atlas-cks/kafka"
	"atlas-cks/keymap"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	consumerName = "character_created_event"
	topicToken   = "TOPIC_CHARACTER_CREATED_EVENT"
)

type createdEvent struct {
	WorldId     byte   `json:"worldId"`
	CharacterId uint32 `json:"characterId"`
	Name        string `json:"name"`
}

func NewConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[createdEvent](consumerName, topicToken, groupId, HandleCreatedEvent(db))
	}
}

func HandleCreatedEvent(db *gorm.DB) kafka.HandlerFunc[createdEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event createdEvent) {
		err := keymap.CreateDefault(l, db)(event.CharacterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to create default keymapping for character %d.", event.CharacterId)
		}
	}
}
