package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/zuchi/go-qa/pkg/infra/api"
)

func main() {
	logCtx := log.WithFields(log.Fields{"component": "main"})
	ctx := context.Background()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	mongoUrl := os.Getenv("MONGO_URL")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("%s", mongoUrl)))

	if err != nil {
		logCtx.Panicf("cannot connect to mongodb: %v", err)
	}

	server := api.NewServer(ctx, client)

	go server.Initialize()

	<-quit

	ctxc, timeoutCtx := context.WithTimeout(ctx, 5*time.Second)
	err = server.Shutdown(ctxc)

	if err != nil {
		logCtx.WithError(err).Error("shutdown is not possible")
	}
	logCtx.Info("waiting for shutdown services")
	defer func() {
		err = client.Disconnect(ctxc)
		if err != nil {
			logCtx.WithError(err).Warnf("cannot disconnect from mongodb")
		}
	}()
	defer timeoutCtx()

	logCtx.Info("Bye")
}
