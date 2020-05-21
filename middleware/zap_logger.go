package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// UserAgent is the key
const UserAgent = "User-Agent"

// BuildLogger sets up zap logger config
func BuildLogger() (logger *zap.Logger, err error) {
	config := zap.NewProductionConfig()
	config.Encoding = "json"
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.CallerKey = "logger"
	return config.Build()
}

// ZapLogger is a middleware that logs the start and end of each request by zap
func ZapLogger(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				if !strings.HasPrefix(r.Header.Get(UserAgent), "kube-probe") {
					l.Info("",
						zap.Any("request", map[string]interface{}{
							"return-status": ww.Status(),
							"http-method":   r.Method,
							"headers": map[string]interface{}{
								"Content-Type":   r.Header.Get("Content-Type"),
								"Content-Length": r.Header.Get("Content-Length"),
								"User-Agent":     r.Header.Get(UserAgent),
								"Server":         r.Header.Get("Server"),
								"Via":            r.Header.Get("Via"),
								"Accept":         r.Header.Get("Accept"),
								"Authorization":  r.Header.Get("Authorization"),
							},
						}),

						// Other data
						zap.String("X-FORWARDED-FOR", r.Header.Get("X-FORWARDED-FOR")),
						zap.String("Remote Addr", r.RemoteAddr),
						zap.String("Proto", r.Proto),
						zap.String("Path", r.URL.Path),
						zap.Duration("lat", time.Since(t1)),
						zap.Int("size", ww.BytesWritten()),
						zap.String("reqId", middleware.GetReqID(r.Context())),
					)
				}
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
