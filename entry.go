// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"github.com/ainsleyclark/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type (
	// Entry defines a singular entry sent to Mongo
	// when a Logrus event is fired.
	Entry struct {
		ID      primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
		Level   string               `json:"level" bson:"level"`
		Time    time.Time            `json:"time" bson:"time"`
		Message string               `json:"message" bson:"message"`
		Data    map[string]any       `json:"data" bson:"data"`
		Error   *Error               `json:"error" bson:"error"`
		Expiry  map[string]time.Time `json:"expiry" bson:"expiry"`
	}
	// Error defines a custom Error for log entries, detailing
	// the file line, errors are returned as strings instead
	// of the stdlib error.
	Error struct {
		Code      string `json:"code" bson:"code"`
		Message   string `json:"message" bson:"message"`
		Operation string `json:"operation" bson:"op"`
		Err       string `json:"error" bson:"err"`
		FileLine  string `json:"file_line" bson:"file_line"`
	}
)

// HasError returns true if the Entry has an
// error attached to it.
func (e *Entry) HasError() bool {
	return e.Error != nil
}

// Error implements the stdlib error interface.
func (e *Error) Error() error {
	if e.Err == "" {
		return nil
	}
	return errors.New(e.Err)
}
