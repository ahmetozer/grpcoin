// Copyright 2021 Ahmet Alp Balkan
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	stackdriver "github.com/tommy351/zap-stackdriver"
	"go.uber.org/zap"
)

type reqZapCtx struct{}

var reqZapCtxVar = reqZapCtx{}

func withLogging(log *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log = log.With(stackdriver.LogHTTPRequest(&stackdriver.HTTPRequest{
				Method:    r.Method,
				URL:       r.URL.Path,
				UserAgent: r.Header.Get("user-agent"),
				Referrer:  r.Header.Get("referer"),
				RemoteIP:  r.Header.Get("x-forwarded-for"),
			}))
			ctx := context.WithValue(r.Context(), reqZapCtxVar, log)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func loggerFrom(ctx context.Context) *zap.Logger {
	v := ctx.Value(reqZapCtxVar)
	if v == nil {
		panic("request did not have a logger")
	}
	vv, ok := v.(*zap.Logger)
	if !ok {
		panic(fmt.Sprintf("req ctx had wrong type of logger (%T)", vv))
	}
	return vv
}
