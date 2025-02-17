package apiserver

import (
	"asyncapi/config"
	"asyncapi/store"
	"context"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

type ApiServer struct {
	Config     *config.Config
	Slog       *slog.Logger
	Store      *store.Store
	jwtManager *JwtManager
}

func NewApiServer(config *config.Config, slog *slog.Logger, store *store.Store, manager *JwtManager) *ApiServer {
	return &ApiServer{Config: config, Slog: slog, Store: store, jwtManager: manager}
}

func (s *ApiServer) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (s *ApiServer) Start(ctx context.Context) error {
	mux := http.NewServeMux() // cоздаëм роутер дл ярегистрации обработчиков
	mux.HandleFunc("GET /ping", s.Ping)
	mux.HandleFunc("POST /auth/signup", s.SignUpHandler)
	mux.HandleFunc("POST /auth/signin", s.SignInHandler)
	mux.HandleFunc("POST /auth/refreshTokens", s.TokenRefreshHandler())
	// создаем обработкич
	Middlewares := NewLoggerMidleware(s.Slog) // промежуточной слой для логирвоания
	Middlewares = NewAuthMiddleware(s.jwtManager, s.Store.Users)
	server := &http.Server{
		Addr:    net.JoinHostPort(s.Config.ApiServerHost, s.Config.ApiServerPort),
		Handler: Middlewares(mux),
	}
	go func() {

		s.Slog.Info("listen on...", s.Config.ApiServerHost, s.Config.ApiServerPort)
		s.Slog.Info(s.Config.DatabaseURL())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Slog.Error(err.Error())
		}

	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		s.Slog.Info("shutting down...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.Slog.Error(err.Error())

		}
	}()
	wg.Wait()
	return nil
}
