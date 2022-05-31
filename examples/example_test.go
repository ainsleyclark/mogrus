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

	clientOptions := options.Client().
		ApplyURI(os.Getenv("MONGO_CONNECTION")).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalln(err)
	}

	opts := mogrus.Options{
		Collection: client.Database("logs").Collection("col"),
		// FireHook is a hook function called just before an
		// entry is logged to Mongo.
		FireHook: func(e mogrus.Entry) {
			fmt.Printf("%+v\n", e)
		},
		// ExpirationLevels allows for the customisation of expiry
		// time for each Logrus level by default entries do not expire.
		ExpirationLevels: mogrus.ExpirationLevels{
			logrus.DebugLevel: time.Second * 5,
			logrus.InfoLevel:  time.Second * 15,
			logrus.ErrorLevel: time.Second * 30,
		},
	}

	// Create the new Mogrus hook, returns an error if the
	// collection is nil.
	hook, err := mogrus.New(context.Background(), opts)
	if err != nil {
		log.Fatalln(err)
	}

	// Add the hook to the Logrus instance.
	l.AddHook(hook)

	l.Debug("Debug level")
	l.WithField("key", "value").Info("Info level")
	l.WithError(errors.NewInternal(errors.New("error"), "message", "op")).Error("Error level")
}
