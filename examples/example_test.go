// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package examples

import (
	"context"
	"fmt"
	"github.com/ainsleyclark/errors"
	"github.com/ainsleyclark/mogrus"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)

	collection := connectMongo()

	hook, err := mogrus.New(context.Background(), mogrus.Options{
		Collection: collection,
		FireHook: func(e mogrus.Entry) {
			fmt.Printf("%+v\n", e)
		},
		ExpirationLevels: mogrus.ExpirationLevels{
			logrus.DebugLevel: time.Second * 5,
			logrus.InfoLevel:  time.Second * 15,
			logrus.ErrorLevel: time.Second * 30,
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	l.AddHook(hook)

	l.Debug("Debug level")
	l.WithField("key", "value").Info("Info level")
	l.WithError(errors.NewInternal(errors.New("error"), "message", "op")).Error("Error level")
}

func connectMongo() *mongo.Collection {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)

	uri := os.Getenv("MONGO_CONNECTION")
	fmt.Println(uri)

	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalln(err)
	}

	return client.Database("logs").Collection("info")
}
