// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"context"
	"github.com/ainsleyclark/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	// hooker is a hook that writes logs of specified
	// LogLevels to specified Writer.
	hooker struct {
		Options Options
	}
	// Options defines the configuration used for creating a
	// new Mogrus hook.
	Options struct {
		Collection *mongo.Collection
		UseAll     bool
		Forget     bool
		Expiry     time.Duration
		FireHook   FireHook
		// TODO Expiry for each level, with how many seconds
	}
	// Entry defines a singular entry sent to Mongo
	// when a Logrus event is fired.
	Entry struct {
		Level   string         `json:"level" bson:"level"`
		Time    time.Time      `json:"time" bson:"time"`
		Message string         `json:"string" bson:"string"`
		Data    map[string]any `json:"data" bson:"data"`
		Error   *errors.Error  `json:"error" bson:"error"`
		Expiry  time.Time      `json:"expiry" bson:"expiry"`
	}
	// FireHook is a hook function called just before a
	// entry is logged to Mongo.
	FireHook func(e Entry)
)

// Validate validates the options before creating a new Hook.
func (m Options) Validate() error {
	if m.Collection == nil {
		return errors.New("mongo collection nil")
	}
	return nil
}

// New creates a new Mogrus hooker.
// Returns errors.INVALID if the collection is nil.
// Returns errors.INTERNAL if the indexes could not be added.
func New(ctx context.Context, opts Options) (*hooker, error) {
	const op = "Mogrus.New"

	err := opts.Validate()
	if err != nil {
		return nil, errors.NewInvalid(err, "Error validating Options", op)
	}

	err = addIndexes(ctx, opts.Collection)
	if err != nil {
		return nil, errors.NewInternal(err, "Error creating indexes", op)
	}

	return &hooker{
		Options: opts,
	}, nil
}

// Fire sends the entry time, level and message and any
// other entry fields to the database.
// Returns errors.INTERNAL if the entry could not be written.
func (hook *hooker) Fire(entry *logrus.Entry) error {
	const op = "Mogrus.Fire"

	formatted := Entry{
		Level:   entry.Level.String(),
		Time:    entry.Time,
		Message: entry.Message,
		Data:    make(map[string]any),
		Error:   nil,
		Expiry:  time.Time{},
	}

	for k, v := range entry.Data {
		if logrus.ErrorKey == k && v != nil {
			formatted.Error = errors.ToError(v)
			continue
		}
		entry.Data[k] = v
	}

	//data["expiry"] = map[string]time.Time{}
	//
	//if data["level"] == "panic" {
	//	data["expiry"] = map[string]time.Time{
	//		"ttl60s": time.Now(),
	//	}
	//} else {
	//	data["expiry"] = map[string]time.Time{
	//		"ttl5s": time.Now(),
	//	}
	//}

	if hook.Options.FireHook != nil {
		hook.Options.FireHook(formatted)
	}

	_, err := hook.Options.Collection.InsertOne(context.Background(), formatted)
	if err != nil {
		return errors.NewInternal(err, "Error writing entry to Mongo Collection", op)
	}

	return nil
}

// Levels define on which log levels this hook would
// trigger.
func (hook *hooker) Levels() []logrus.Level {
	return logrus.AllLevels
}

func addIndexes(ctx context.Context, collection *mongo.Collection) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"expiry.ttl60s": 1},
			Options: options.Index().SetExpireAfterSeconds(60).SetSparse(true),
		},
		{
			Keys:    bson.M{"expiry.ttl5s": 1},
			Options: options.Index().SetExpireAfterSeconds(5).SetSparse(true),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}

	return nil
}
