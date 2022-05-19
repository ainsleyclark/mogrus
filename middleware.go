// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package logger

//func Log(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		ctx := context.WithValue(r.Context(), "request", r)
//		*r = *r.WithContext(ctx)
//
//		startTime := time.Now()
//		next.ServeHTTP(w, r)
//		endTime := time.Now()
//
//		fields := logrus.Fields{
//			//"status_code":
//			"latency_time":   endTime.Sub(startTime),
//			"client_ip":      r.RemoteAddr,
//			"request_method": r.Method,
//			"request_url":    r.RequestURI,
//			//"start_time":
//			//"message":        appMessage,
//			"error": r.Context().Value("error"),
//		}
//
//		logger.WithFields(fields).Info()
//	})
//}

//
//func Error(w http.ResponseWriter, r *http.Request, err error) {
//	if err == nil {
//		return
//	}
//
//	ctx := context.WithValue(r.Context(), "error", err)
//	*r = *r.WithContext(ctx)
//
//	errObj := struct {
//		Error string `json:"error"`
//	}{
//		Error: err.Error(),
//	}
//	if err := otohttp.Encode(w, r, http.StatusInternalServerError, errObj); err != nil {
//		log.Printf("failed to encode error: %s\n", err)
//	}
//}
