package middleware

import (
	"context"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/session"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/utils"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Key int

const CurUserKey Key = 1

func AuthMiddleware(logger *zap.SugaredLogger, sm *session.SessionsManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infow("auth middleware")
		authParam := r.Header.Get("authorization")
		if !strings.HasPrefix(authParam, "Bearer ") || authParam == "" {
			logger.Infow("Wrong token")
			utils.MakeResponse(logger, w, []byte(`{"message": "no access token or token has bad format"}`), 401)
			return
		}
		curSession, err := sm.GetSession(strings.TrimPrefix(authParam, "Bearer "))
		if curSession == nil || err != nil {
			logger.Infow("error in receiving the session")
			utils.MakeResponse(logger, w, []byte(`{"message": "error in receiving session"}`), 401)
			return
		}
		curUser := curSession.User
		ctx := r.Context()
		ctx = context.WithValue(ctx, CurUserKey, curUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
