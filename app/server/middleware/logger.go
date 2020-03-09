package middleware

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Logger logs each requests
// ex) {"latency":"8.8124ms,"level":"info","method":"GET","msg":"request done","name":"GetHealth","time":"2019-06-30T20:30:58+09:00","type":"request","uri":"/api/health"}
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.WithFields(log.Fields{
			"type":    "request",
			"method":  r.Method,
			"uri":     r.RequestURI,
			"name":    name,
			"latency": fmt.Sprint(time.Since(start)),
		}).Info("request done")
	})
}
