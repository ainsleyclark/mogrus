// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"context"
	"fmt"
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
		Options
	}
	// ExpirationLevels defines the map of log levels mapped to
	// a duration a LevelIndex.
	ExpirationLevels map[logrus.Level]time.Duration
	// FireHook defines the function used for firing entry
	// to a call back function.
	FireHook func(e Entry)
)

const (
	// DefaultExpiryKey is the key and index stored within
	// Mongo for log levels that have a duration.
	DefaultExpiryKey = "ttl-%s"
)

// New creates a new Mogrus hooker.
// Returns errors.INVALID if the collection is nil.
// Returns errors.INTERNAL if the indexes could not be added.
func New(ctx context.Context, opts Options) (logrus.Hook, error) {
	const op = "Mogrus.New"

	err := opts.Validate()
	if err != nil {
		return nil, errors.NewInvalid(err, "Error validating Options", op)
	}

	err = addIndexes(ctx, opts)
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

	formatted := ToEntry(entry)

	// Add expiry to levels.
	for level := range hook.ExpirationLevels {
		if level == entry.Level {
			key := fmt.Sprintf(DefaultExpiryKey, level.String())
			formatted.Expiry[key] = time.Now()
		}
	}

	// Fire callback.
	if hook.FireHook != nil {
		hook.FireHook(formatted)
	}

	// Insert into Mongo.
	_, err := hook.Collection.InsertOne(context.Background(), formatted)
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

// indexes returns a collection of mongo.IndexModels with the
// appropriate expiry durations assigned.
func (e ExpirationLevels) indexes() []mongo.IndexModel {
	var indexes []mongo.IndexModel
	for level, duration := range e {
		key := fmt.Sprintf(DefaultExpiryKey, level.String())
		indexes = append(indexes, mongo.IndexModel{
			Keys:    bson.M{"expiry." + key: 1},
			Options: options.Index().SetExpireAfterSeconds(int32(duration.Seconds())).SetSparse(true),
		})
	}
	return indexes
}

// addIndexes is responsible for injecting the indexes to
// the Mongo collection.
func addIndexes(ctx context.Context, opts Options) error {
	// Bail if there are no expiration levels set.
	if len(opts.ExpirationLevels) == 0 {
		return nil
	}

	// Create the indexes within the Mongo collection.
	_, err := opts.Collection.Indexes().CreateMany(ctx, opts.ExpirationLevels.indexes())
	if err != nil {
		return err
	}

	return nil
}

// ToEntry transforms a logrus.Entry to a mogrus.Entry
func ToEntry(entry *logrus.Entry) Entry {
	// Construct a formatted Mongo Entry.
	formatted := Entry{
		Level:   entry.Level.String(),
		Time:    entry.Time,
		Message: entry.Message,
		Expiry:  make(map[string]time.Time),
	}

	// Range over the entries data and assign an error if
	// it exists, otherwise construct a map with the field data.
	for k, v := range entry.Data {
		if logrus.ErrorKey == k {
			e := errors.ToError(v)
			if e == nil {
				continue
			}
			formatted.Error = &Error{
				Code:      e.Code,
				Message:   e.Message,
				Operation: e.Operation,
				Err:       e.Err.Error(),
				FileLine:  e.FileLine(),
			}
			continue
		}
		if formatted.Data == nil {
			formatted.Data = make(map[string]any)
		}
		formatted.Data[k] = v
	}

	return formatted
}
