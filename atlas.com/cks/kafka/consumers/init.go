package consumers

import (
	"atlas-cks/character"
	"atlas-cks/kafka/handler"
	"atlas-cks/keymap"
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
)

func CreateEventConsumers(l *logrus.Logger, db *gorm.DB, ctx context.Context, wg *sync.WaitGroup) {
	cec := func(topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
		createEventConsumer(l, ctx, wg, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_CHANGE_KEY_MAP", keymap.CommandCreator(), keymap.HandleChangeCommand(db))
	cec("TOPIC_CHARACTER_CREATED_EVENT", character.CreatedEventCreator(), character.HandleCreatedEvent(db))
}

func createEventConsumer(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup, topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
	wg.Add(1)
	go NewConsumer(l, ctx, wg, topicToken, "Character Keyboard Settings Service", emptyEventCreator, processor)
}
