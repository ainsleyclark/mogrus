// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"bytes"
	"fmt"
	"github.com/ainsleyclark/errors"
	"github.com/sirupsen/logrus"
)

func (t *LoggerTestSuite) TestInit() {
	opts := Options{}
	New(opts)
	t.Equal(logrus.TraceLevel, logger.Level)
	logger = logrus.New()
}

func (t *LoggerTestSuite) TestLogger() {
	tt := map[string]struct {
		fn   func()
		want string
	}{
		"Trace": {
			func() {
				Trace("trace")
			},
			"trace",
		},
		"Debug": {
			func() {
				Debug("debug")
			},
			"debug",
		},
		"Info": {
			func() {
				Info("info")
			},
			"info",
		},
		"Warn": {
			func() {
				Warn("warning")
			},
			"warning",
		},
		"Error": {
			func() {
				Error("error")
			},
			"error",
		},
		"With Field": {
			func() {
				WithField("test", "with-field").Error()
			},
			"with-field",
		},
		"With Fields": {
			func() {
				WithFields(logrus.Fields{"test": "with-fields"}).Error()
			},
			"with-fields",
		},
		"With Error": {
			func() {
				WithError(&errors.Error{Code: "code", Message: "message", Operation: "op", Err: fmt.Errorf("err")}).Error()
			},
			"[code] code [msg] message [op] op [error] err",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			buf := t.Setup()
			test.fn()
			t.Contains(buf.String(), test.want)
		})
	}
}

func (t *LoggerTestSuite) TestLogger_Fatal() {
	buf := t.Setup() // nolint
	defer func() {
		logger = logrus.New()
	}()
	logger.ExitFunc = func(i int) {}
	Fatal("fatal")
	t.Contains(buf.String(), "fatal")
}

func (t *LoggerTestSuite) TestLogger_Panic() {
	buf := t.Setup()
	t.Panics(func() {
		Panic("panic")
	})
	t.Contains(buf.String(), "panic")
}

func (t *LoggerTestSuite) TestLogger_SetOutput() {
	buf := &bytes.Buffer{}
	SetOutput(buf)
	t.Equal(buf, logger.Out)
}

func (t *LoggerTestSuite) TestSetLevel() {
	defer func() {
		logger = logrus.New()
	}()
	SetLevel(logrus.WarnLevel)
	t.Equal(logrus.WarnLevel, logger.GetLevel())
}

func (t *LoggerTestSuite) TestSetLogger() {
	defer func() {
		logger = logrus.New()
	}()
	l := logger
	SetLogger(l)
	t.Equal(l, logger)
}
