package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type logger struct {
	Timestamp time.Time
	Latency   int64
	Path      string
	OS        string
}

func Logger(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h.ServeHTTP(w, r)

		end := time.Now()
		sub := end.Sub(start)
		latency := int64(sub / time.Millisecond)
		os, ok := r.Context().Value(osKey).(string)
		if !ok {
			fmt.Println("os not found in context")
			return
		}

		logger := &logger{
			Timestamp: start,
			Latency:   latency,
			Path:      r.URL.Path,
			OS:        os,
		}

		json, err := json.Marshal(logger)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(json))

	}
	return http.HandlerFunc(fn)
}
