// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"github.com/ainsleyclark/errors"
	"time"
)

// Entry defines a singular entry sent to Mongo
// when a Logrus event is fired.
type Entry struct {
	Level   string               `json:"level" bson:"level"`
	Time    time.Time            `json:"time" bson:"time"`
	Message string               `json:"string" bson:"string"`
	Data    map[string]any       `json:"data" bson:"data"`
	Error   *errors.Error        `json:"error" bson:"error"`
	Expiry  map[string]time.Time `json:"expiry" bson:"expiry"`
}

// HasError returns true if the Entry has an
// error attached to it.
func (e *Entry) HasError() bool {
	return e.Error != nil
}
