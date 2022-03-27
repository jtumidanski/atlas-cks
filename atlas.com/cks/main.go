package main

import (
	"atlas-cks/character"
	"atlas-cks/database"
	"atlas-cks/kafka"
	"atlas-cks/keymap"
	"atlas-cks/logger"
	"atlas-cks/rest"
	"atlas-cks/tracing"
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const serviceName = "atlas-cks"
const consumerGroupId = "Character Keyboard Settings Service"

func main() {
	l := logger.CreateLogger(serviceName)
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	tc, err := tracing.InitTracer(l)(serviceName)
	if err != nil {
		l.WithError(err).Fatal("Unable to initialize tracer.")
	}
	defer func(tc io.Closer) {
		err := tc.Close()
		if err != nil {
			l.WithError(err).Errorf("Unable to close tracer.")
		}
	}(tc)

	db := database.Connect(l, database.SetMigrations(keymap.Migration))

	kafka.CreateConsumers(l, ctx, wg,
		character.NewConsumer(db)(consumerGroupId),
		keymap.NewConsumer(db)(consumerGroupId))

	rest.CreateService(l, db, ctx, wg, "/ms/cks", keymap.InitResource)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Initiating shutdown with signal %s.", sig)
	cancel()
	wg.Wait()
	l.Infoln("Service shutdown.")
}
