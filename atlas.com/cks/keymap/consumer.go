package keymap

import (
	"atlas-cks/kafka/handler"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

func CommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &changeCommand{}
	}
}

func HandleChangeCommand(db *gorm.DB) handler.EventHandler {
	return func(l logrus.FieldLogger, span opentracing.Span, e interface{}) {
		if event, ok := e.(*changeCommand); ok {
			err := ChangeKeyMap(l, db)(event.CharacterId, event.Changes)
			if err != nil {
				l.WithError(err).Errorf("Unable to update character %d keybinding.", event.CharacterId)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
