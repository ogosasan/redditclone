package middleware

import (
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

func PanicMiddleware(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infow("panicMiddleware", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				logger.Infow("recovered", err)
				utils.MakeResponse(logger, w, []byte(`{"message": "panic occurred"}`), 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
