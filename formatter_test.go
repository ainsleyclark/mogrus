// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"fmt"
	"github.com/ainsleyclark/errors"
	"github.com/sirupsen/logrus"
	"time"
)

func (t *LoggerTestSuite) TestFormatter() {
	now := time.Now()
	nowStr := now.Format(time.StampMilli)

	tt := map[string]struct {
		entry *logrus.Entry
		want  string
	}{
		"Debug": {
			&logrus.Entry{
				Level:   logrus.DebugLevel,
				Message: "message",
			},
			fmt.Sprintf(Prefix+" %s | KRA | [DEBUG] | [msg] message\n", nowStr),
		},
		"Info": {
			&logrus.Entry{
				Level:   logrus.InfoLevel,
				Message: "message",
			},
			fmt.Sprintf(Prefix+" %s | KRA | [INFO]  | [msg] message\n", nowStr),
		},
		"Warning": {
			&logrus.Entry{
				Level:   logrus.WarnLevel,
				Message: "message",
			},
			fmt.Sprintf(Prefix+" %s | KRA | [WARNING] | [msg] message\n", nowStr),
		},
		"Error": {
			&logrus.Entry{
				Level:   logrus.ErrorLevel,
				Message: "message",
			},
			fmt.Sprintf(Prefix+" %s | KRA | [ERROR] | [msg] message\n", nowStr),
		},
		"Fatal": {
			&logrus.Entry{
				Level:   logrus.FatalLevel,
				Message: "message",
			},
			fmt.Sprintf(Prefix+" %s | KRA | [FATAL] | [msg] message\n", nowStr),
		},
		"Panic": {
			&logrus.Entry{
				Level:   logrus.PanicLevel,
				Message: "message",
			},
			fmt.Sprintf(Prefix+" %s | KRA | [PANIC] | [msg] message\n", nowStr),
		},
		"Fields": {
			&logrus.Entry{
				Data: logrus.Fields{
					"fields": logrus.Fields{"key1": "test1"},
				},
				Level: logrus.InfoLevel,
			},
			fmt.Sprintf(Prefix+" %s | KRA | [INFO]  | key1: test1\n", nowStr),
		},
		"Print Error Pointer": {
			&logrus.Entry{
				Data: logrus.Fields{
					"error": &errors.Error{Code: "INTERNAL", Message: "message", Operation: "operation", Err: fmt.Errorf("error")},
				},
				Level: logrus.ErrorLevel,
			},
			fmt.Sprintf(Prefix+" %s | KRA | [ERROR] | [code] INTERNAL [msg] message [op] operation [error] error\n", nowStr),
		},
		"Print Error Non Pointer": {
			&logrus.Entry{
				Data: logrus.Fields{
					"error": errors.Error{Code: "INTERNAL", Message: "message", Operation: "operation", Err: fmt.Errorf("error")},
				},
				Level: logrus.ErrorLevel,
			},
			fmt.Sprintf(Prefix+" %s | KRA | [ERROR] | [code] INTERNAL [msg] message [op] operation [error] error\n", nowStr),
		},
		"Nil To Error": {
			&logrus.Entry{
				Data: logrus.Fields{
					"error": 1,
				},
				Level: logrus.ErrorLevel,
			},
			fmt.Sprintf(Prefix+" %s | KRA | [ERROR]\n", nowStr),
		},
		"Print Error": {
			&logrus.Entry{
				Data: logrus.Fields{
					"error": fmt.Errorf("error"),
				},
				Level: logrus.ErrorLevel,
			},
			fmt.Sprintf(Prefix+" %s | KRA | [ERROR] | [error] error\n", nowStr),
		},
		"Print Error String": {
			&logrus.Entry{
				Data: logrus.Fields{
					"error": "error",
				},
				Level: logrus.ErrorLevel,
			},
			fmt.Sprintf(Prefix+" %s | KRA | [ERROR] | [error] error\n", nowStr),
		},
		"Server Success": {
			&logrus.Entry{
				Data: logrus.Fields{
					"status_code":    200,
					"client_ip":      "127.0.0.1",
					"request_method": "GET",
					"request_url":    "/page",
					"data_length":    0,
				},
				Level: logrus.InfoLevel,
			},
			fmt.Sprintf(Prefix+" %s | 200 | [INFO]  | 127.0.0.1 |   GET    \"/page\"\n", nowStr),
		},
		"Server Not Found": {
			&logrus.Entry{
				Data: logrus.Fields{
					"status_code":    404,
					"client_ip":      "127.0.0.1",
					"request_method": "GET",
					"request_url":    "/page",
					"data_length":    0,
				},
				Level: logrus.InfoLevel,
			},
			fmt.Sprintf(Prefix+" %s | 404 | [INFO]  | 127.0.0.1 |   GET    \"/page\"\n", nowStr),
		},
		"Message": {
			&logrus.Entry{
				Data: logrus.Fields{
					"message": "message",
				},
				Level: logrus.InfoLevel,
			},
			fmt.Sprintf(Prefix+" %s | KRA | [INFO]  | [msg] message\n", nowStr),
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			test.entry.Time = now
			f := Formatter{
				Colours: false,
			}
			got, err := f.Format(test.entry)
			t.NoError(err)
			t.Equal(test.want, string(got))
		})
	}
}
