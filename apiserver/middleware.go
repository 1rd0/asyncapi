package apiserver

import (
	"asyncapi/store"
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strings"
)

func NewLoggerMidleware(Logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Logger.Info("New Request path ", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

type UserCtxKey struct{}

func ContextWithUser(ctx context.Context, user *store.User) context.Context {
	return context.WithValue(ctx, UserCtxKey{}, user)
}

func NewAuthMiddleware(JwtManager *JwtManager, UserStore *store.UserStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/auth") {
				next.ServeHTTP(w, r)
				return
			}
			r.Header.Get("Authorization")
			var token string
			if parts := strings.Split(r.Header.Get("Authorization"), " "); len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
			if token == "" {
				slog.Error("no token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ParsedToken, err := JwtManager.Parse(token)
			if err != nil {
				slog.Error("failed to parse token: ", err)
				w.WriteHeader(http.StatusUnauthorized)
			}

			if !JwtManager.IsAccesstoken(ParsedToken) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("invalid token"))
				return
			}

			UserIdstr, err := ParsedToken.Claims.GetSubject()
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("invalid token claims"))
				return
			}

			UserId, err := uuid.Parse(UserIdstr)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("invalid token parse error"))
			}

			user, err := UserStore.ById(r.Context(), UserId)
			if err != nil {
				slog.Error("failed to get user by id: ", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(ContextWithUser(r.Context(), user)))
		})

	}
}
