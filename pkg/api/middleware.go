package api

import (
	"log"
	"net/http"
	"time"
)

func getRemoteAddr(r *http.Request) string {
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		return forwardedFor
	}
	return r.RemoteAddr
}

/* Support for response status codes logging */
type loggingResponseWriter struct {
	responseWriter http.ResponseWriter
	statusCode     int
	responseLength int
}

func newLoggingResponseWriter(w http.ResponseWriter) loggingResponseWriter {
	return loggingResponseWriter{
		responseWriter: w,
		statusCode:     http.StatusOK,
		responseLength: 0,
	}
}

func (lrw *loggingResponseWriter) Header() http.Header {
	return lrw.responseWriter.Header()
}

func (lrw *loggingResponseWriter) WriteHeader(status int) {
	lrw.statusCode = status
	lrw.responseWriter.WriteHeader(status)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode == 0 {
		lrw.statusCode = http.StatusOK
	}
	n, err := lrw.responseWriter.Write(b)
	lrw.responseLength += n
	return n, err
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var startTime = time.Now()

		var lrw = newLoggingResponseWriter(w)
		next.ServeHTTP(&lrw, r)

		var elapsed = time.Since(startTime)

		log.Printf("[%s] %s %s: %d %d %f",
			getRemoteAddr(r),
			r.Method,
			r.RequestURI,
			lrw.statusCode,
			lrw.responseLength,
			elapsed.Seconds(),
		)
	})
}
