package character

import (
	"atlas-cks/kafka/handler"
	"atlas-cks/keymap"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type createdEvent struct {
	WorldId     byte   `json:"worldId"`
	CharacterId uint32 `json:"characterId"`
	Name        string `json:"name"`
}

func CreatedEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &createdEvent{}
	}
}

func HandleCreatedEvent(db *gorm.DB) handler.EventHandler {
	return func(l logrus.FieldLogger, span opentracing.Span, e interface{}) {
		if event, ok := e.(*createdEvent); ok {
			err := keymap.CreateDefault(l, db)(event.CharacterId)
			if err != nil {
				l.WithError(err).Errorf("Unable to create default keymapping for character %d.", event.CharacterId)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
