package main

import (
	"atlas-cks/database"
	"atlas-cks/kafka/consumers"
	"atlas-cks/keymap"
	"atlas-cks/logger"
	"atlas-cks/rest"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	l := logger.CreateLogger()
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	db := database.Connect(l, database.SetMigrations(keymap.Migration))

	consumers.CreateEventConsumers(l, db, ctx, wg)

	rest.CreateRestService(l, db, ctx, wg)

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
