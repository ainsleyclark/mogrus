// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"time"
)

// Entry defines a singular entry sent to Mongo
// when a Logrus event is fired.
type (
	Entry struct {
		Level   string               `json:"level" bson:"level"`
		Time    time.Time            `json:"time" bson:"time"`
		Message string               `json:"message" bson:"message"`
		Data    map[string]any       `json:"data" bson:"data"`
		Error   *Error               `json:"error" bson:"error"`
		Expiry  map[string]time.Time `json:"expiry" bson:"expiry"`
	}
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
