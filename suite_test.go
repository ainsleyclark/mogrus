// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"io"
	"testing"
)

// LoggerTestSuite defines the helper used for L
// testing.
type LoggerTestSuite struct {
	suite.Suite
}

// TestLogger asserts testing has begun.
func TestLogger(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

// TearDownTestSuite - Teardown logging after testing.
func (t *LoggerTestSuite) TearDownTestSuite() {
	New(Options{})
}

// Setup is a helper function for setting up the L
// suite.
func (t *LoggerTestSuite) Setup() *bytes.Buffer {
	buf := &bytes.Buffer{}
	L.SetLevel(logrus.TraceLevel)
	L.SetOutput(buf)
	L.SetFormatter(&Formatter{
		Options: Options{Colours: false},
	})
	return buf
}

// SetupHooks is a helper function function for setting up
// the hooks for testing.
func (t *LoggerTestSuite) SetupHooks(writer io.Writer) WriterHook {
	return WriterHook{
		Writer: writer,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	}
}
