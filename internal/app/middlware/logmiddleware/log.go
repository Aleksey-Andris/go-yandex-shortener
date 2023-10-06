package logmiddleware

import (
	"net/http"
	"time"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/logger"
)

type responseData struct {
	status int
	size   int
}

type logginResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (res *logginResponseWriter) Write(b []byte) (int, error) {
	size, err := res.ResponseWriter.Write(b)
	res.responseData.size += size
	return size, err
}

func (res *logginResponseWriter) WriteHeader(statusCode int) {
	res.ResponseWriter.WriteHeader(statusCode)
	res.responseData.status = statusCode
}

func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lRes := logginResponseWriter{
			ResponseWriter: res,
			responseData:   responseData,
		}
		next.ServeHTTP(&lRes, req)
		duration := time.Since(start)

		logger.Log().Sugar().Infoln(
			"uri", req.RequestURI,
			"method", req.Method,
			"status", responseData.status,
			"duartion", duration,
			"size", responseData.size,
		)
	})
}
