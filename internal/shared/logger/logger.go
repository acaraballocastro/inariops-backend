package logger

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func Info(msg string, args ...any) {
	log.Printf("[INFO] "+msg, args...)
}

func Error(msg string, args ...any) {
	log.Printf("[ERROR] "+msg, args...)
}

func Debug(msg string, args ...any) {
	log.Printf("[DEBUG] "+msg, args...)
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			status:         200, // default si no se escribe explícitamente
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		log.Printf("%d %s %s %s",
			rw.status,
			r.Method,
			r.RequestURI,
			duration,
		)
	})
}
