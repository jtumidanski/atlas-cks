package consumers

import (
	"atlas-cks/retry"
	"atlas-cks/topic"
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
	"time"
)

func Create(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup, configs ...Config) {
	for _, c := range configs {
		go create(l, ctx, wg, c)
	}
}

func NewConfiguration(name string, topicToken string, groupId string, handler MessageHandler) Config {
	return Config{
		name:       name,
		topicToken: topicToken,
		groupId:    groupId,
		maxWait:    500,
		handler:    handler,
	}
}

type Config struct {
	name       string
	topicToken string
	groupId    string
	maxWait    time.Duration
	handler    MessageHandler
}

type MessageHandler func(l logrus.FieldLogger, span opentracing.Span, msg kafka.Message)

func create(cl *logrus.Logger, ctx context.Context, wg *sync.WaitGroup, c Config) {
	initSpan := opentracing.StartSpan("consumer_init")
	t := topic.GetRegistry().Get(cl, initSpan, c.topicToken)
	initSpan.Finish()

	l := cl.WithFields(logrus.Fields{"originator": t, "type": "kafka_consumer"})

	l.Infof("Creating topic consumer.")

	wg.Add(1)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
		Topic:   t,
		GroupID: c.groupId,
		MaxWait: c.maxWait,
	})

	go func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		for {
			var msg kafka.Message
			readerFunc := func(attempt int) (bool, error) {
				var err error
				msg, err = r.ReadMessage(ctx)
				if err == io.EOF || err == context.Canceled {
					return false, err
				} else if err != nil {
					l.WithError(err).Warnf("Could not read message on topic %s, will retry.", r.Config().Topic)
					return true, err
				}
				return false, err
			}

			err := retry.Try(readerFunc, 10)
			if err == io.EOF || err == context.Canceled || len(msg.Value) == 0 {
				l.Infof("Reader closed, shutdown.")
				return
			} else if err != nil {
				l.WithError(err).Errorf("Could not successfully read message.")
			} else {
				l.Infof("Message received %s.", string(msg.Value))
				go func() {
					headers := make(map[string]string)
					for _, header := range msg.Headers {
						headers[header.Key] = string(header.Value)
					}

					spanContext, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(headers))
					span := opentracing.StartSpan(c.name, opentracing.FollowsFrom(spanContext))
					defer span.Finish()

					c.handler(l, span, msg)
				}()
			}
		}
	}()

	l.Infof("Start consuming topic.")
	<-ctx.Done()
	l.Infof("Shutting down topic consumer.")
	if err := r.Close(); err != nil {
		l.WithError(err).Errorf("Error closing reader.")
	}
	wg.Done()
	l.Infof("Topic consumer stopped.")
}
