// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"context"
	"github.com/ainsleyclark/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	tt := map[string]struct {
		input func(t *mtest.T) Options
		want  any
	}{
		"Success": {
			func(t *mtest.T) Options {
				return Options{Collection: t.Coll}
			},
			nil,
		},
		"Validation Failed": {
			func(t *mtest.T) Options {
				return Options{}
			},
			"Error validating Options",
		},
		"Index Error": {
			func(t *mtest.T) Options {
				return Options{
					Collection: t.Coll,
					ExpirationLevels: ExpirationLevels{
						logrus.PanicLevel: time.Second * 10,
					},
				}
			},
			"Error creating indexes",
		},
		"Index Success": {
			func(t *mtest.T) Options {
				t.AddMockResponses(mtest.CreateSuccessResponse(
					bson.D{{"key", "value"}}...)) //nolint
				return Options{
					Collection: t.Coll,
					ExpirationLevels: ExpirationLevels{
						logrus.PanicLevel: time.Second * 10,
					},
				}
			},
			nil,
		},
	}

	for name, test := range tt {
		mt.Run(name, func(t *mtest.T) {
			_, err := New(context.TODO(), test.input(t))
			if err != nil {
				msg := errors.Message(err)
				if !reflect.DeepEqual(test.want, msg) {
					t.Fatalf("expecting %s, got %s", test.want, msg)
				}
			}
		})
	}
}

func TestHooker_Fire(t *testing.T) {
	var (
		mt    = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		now   = time.Now()
		entry = logrus.Entry{
			Level:   logrus.PanicLevel,
			Time:    now,
			Message: "message",
		}
	)

	defer mt.Close()

	tt := map[string]struct {
		input logrus.Entry
		mock  func(t *mtest.T) Options
		want  any
	}{
		"Simple": {
			entry,
			func(t *mtest.T) Options {
				t.AddMockResponses(mtest.CreateSuccessResponse(
					bson.D{{"key", "value"}}...)) //nolint
				return Options{Collection: t.Coll}
			},
			nil,
		},
		"With Error": {
			logrus.Entry{
				Level:   logrus.PanicLevel,
				Time:    now,
				Message: "message",
				Data: map[string]any{
					logrus.ErrorKey: errors.NewInvalid(errors.New("error"), "message", "op"),
				},
			},
			func(t *mtest.T) Options {
				t.AddMockResponses(mtest.CreateSuccessResponse(
					bson.D{{"key", "value"}}...)) //nolint
				return Options{Collection: t.Coll}
			},
			nil,
		},
		"With Data": {
			logrus.Entry{
				Level:   logrus.PanicLevel,
				Time:    now,
				Message: "message",
				Data: map[string]any{
					"key": "value",
				},
			},
			func(t *mtest.T) Options {
				t.AddMockResponses(mtest.CreateSuccessResponse(
					bson.D{{"key", "value"}}...)) //nolint
				return Options{Collection: t.Coll}
			},
			nil,
		},
		"With Expiry": {
			logrus.Entry{
				Level:   logrus.PanicLevel,
				Time:    now,
				Message: "message",
			},
			func(t *mtest.T) Options {
				t.AddMockResponses(mtest.CreateSuccessResponse(
					bson.D{{"key", "value"}}...)) //nolint
				t.AddMockResponses(mtest.CreateSuccessResponse(
					bson.D{{"key", "value"}}...)) //nolint
				return Options{
					Collection:       t.Coll,
					ExpirationLevels: ExpirationLevels{logrus.PanicLevel: time.Second * 1},
				}
			},
			nil,
		},
		"FireHook": {
			entry,
			func(t *mtest.T) Options {
				t.AddMockResponses(mtest.CreateSuccessResponse(
					bson.D{{"key", "value"}}...)) //nolint
				return Options{
					Collection: t.Coll,
					FireHook: func(e Entry) {
						if !reflect.DeepEqual(e.Message, entry.Message) {
							t.Fatalf("expecting %+v, got %+v", e.Message, entry.Message)
						}
						if !reflect.DeepEqual(e.Level, entry.Level.String()) {
							t.Fatalf("expecting %+v, got %+v", e.Level, entry.Level)
						}
					},
				}
			},
			nil,
		},
		"Mongo Error": {
			logrus.Entry{Level: logrus.PanicLevel, Time: now, Message: "message"},
			func(t *mtest.T) Options {
				t.AddMockResponses(bson.D{{"ok", 0}}) //nolint
				return Options{Collection: t.Coll}
			},
			"Error writing entry to Mongo Collection",
		},
	}

	for name, test := range tt {
		mt.Run(name, func(t *mtest.T) {
			h, err := New(context.TODO(), test.mock(t))
			if err != nil {
				t.Fatalf("error creating hooker")
			}

			err = h.Fire(&test.input)
			if err != nil {
				msg := errors.Message(err)
				if !reflect.DeepEqual(test.want, msg) {
					t.Fatalf("expecting %s, got %s", test.want, msg)
				}
			}
		})
	}
}

func TestHooker_Levels(t *testing.T) {
	h := &hooker{}
	got := h.Levels()
	want := logrus.AllLevels
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expecting %+v, got %+v", want, got)
	}
}
