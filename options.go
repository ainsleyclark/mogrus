// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"github.com/ainsleyclark/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// Options defines the configuration used for creating a
// new Mogrus hook.
type Options struct {
	// Collection is the Mongo collection to write to when
	// a log is fired.
	Collection *mongo.Collection
	// FireHook is a hook function called just before an
	// entry is logged to Mongo.
	FireHook FireHook
	// ExpirationLevels allows for the customisation of expiry
	// time for each Logrus level by default entries do not expire.
	// There may be instances where you want to keep Panics in
	// the Mongo collection for longer than trace levels.
	// For example:
	/*
		var levels = ExpirationLevels{
			// Expire trace levels after 10 hours.
			logrus.TraceLevel: time.Hour * 10
			// Expire info levels after 24 hours.
			logrus.InfoLevel:t ime.Hour * 24,
			// Expire panic levels after 1 week.
			logrus.InfoLevel:t ime.Hour * 24 * 6,
		}
	*/
	ExpirationLevels ExpirationLevels
}

// Validate validates the options before creating a new Hook.
func (o Options) Validate() error {
	if o.Collection == nil {
		return errors.New("nil mongo collection")
	}
	return nil
}
