package apiserver

import (
	"asyncapi/config"
	"context"
	"log/slog"
	"net/http"
)

type ApiServer struct {
	Config *config.Config
	Slog   *slog.Logger
}

func NewApiServer(config *config.Config, slog *slog.Logger) *ApiServer {
	return &ApiServer{Config: config, Slog: slog}
}

func (s *ApiServer) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (s *ApiServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", s.Ping)

	server := &http.Server{
		Addr:    ":4000",
		Handler: mux,
	}
	go func() {

		s.Slog.Info("listen...")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Slog.Error(err.Error())
		}

	}()
	return server.ListenAndServe()
}
