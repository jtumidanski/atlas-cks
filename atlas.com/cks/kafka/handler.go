package kafka

import (
	"atlas-cks/kafka/consumers"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type HandlerFunc[E any] func(logrus.FieldLogger, opentracing.Span, E)

func Adapt[E any](eh HandlerFunc[E]) consumers.MessageHandler {
	return func(l logrus.FieldLogger, span opentracing.Span, msg kafka.Message) {
		var event E
		err := json.Unmarshal(msg.Value, &event)
		if err != nil {
			l.WithError(err).Errorf("Could not unmarshal event into %s.", msg.Value)
		} else {
			eh(l, span, event)
		}
	}
}
