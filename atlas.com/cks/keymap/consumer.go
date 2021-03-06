package keymap

import (
	"atlas-cks/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	consumerName = "change_key_map_command"
	topicToken   = "TOPIC_CHANGE_KEY_MAP"
)

type changeCommand struct {
	CharacterId uint32   `json:"characterId"`
	Changes     []Change `json:"changes"`
}

type Change struct {
	Key        int32 `json:"key"`
	ChangeType int8  `json:"changeType"`
	Action     int32 `json:"action"`
}

func NewConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[changeCommand](consumerName, topicToken, groupId, HandleChangeCommand(db))
	}
}

func HandleChangeCommand(db *gorm.DB) kafka.HandlerFunc[changeCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command changeCommand) {
		err := ChangeKeyMap(l, db)(command.CharacterId, command.Changes)
		if err != nil {
			l.WithError(err).Errorf("Unable to update character %d keybinding.", command.CharacterId)
		}
	}
}
