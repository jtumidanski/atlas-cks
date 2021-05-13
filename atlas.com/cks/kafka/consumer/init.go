package consumer

import (
	"atlas-cks/character"
	"atlas-cks/kafka/handler"
	"atlas-cks/keymap"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateEventConsumers(l *logrus.Logger, db *gorm.DB) {
	cec := func(topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
		createEventConsumer(l, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_CHANGE_KEY_MAP", keymap.CommandCreator(), keymap.HandleChangeCommand(db))
	cec("TOPIC_CHARACTER_CREATED_EVENT", character.CreatedEventCreator(), character.HandleCreatedEvent(db))
}

func createEventConsumer(l *logrus.Logger, topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
	go NewConsumer(l, topicToken, "Character Keyboard Settings Service", emptyEventCreator, processor)
}
