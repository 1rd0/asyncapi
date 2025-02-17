package apiserver

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ApiResponse[T any] struct {
	Data    *T     `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func (r *SignUpRequest) Validate() error {
	if r.Password == "" {
		return errors.New("password is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

func (s *ApiServer) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingUser, err := s.Store.Users.ByEmail(r.Context(), req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if existingUser != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = s.Store.Users.CreateUser(r.Context(), req.Email, req.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ApiResponse[struct{}]{
		Message: "User created successfully",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r *SignInRequest) Validate() error {
	if r.Password == "" {
		return errors.New("password is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

func (s *ApiServer) SignInHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() // Ensure deferred closure handles edge cases

	var req SignInRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		errors.New("password is required")

	}
	if req.Email == "" {
		errors.New("email is required")
	}
	user, err := s.Store.Users.ByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := user.ComparePassword(req.Password); err != nil {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}

	PairToken, err := s.jwtManager.GenerateTokenPair(user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if PairToken == nil {
		http.Error(w, "error generating token pair", http.StatusInternalServerError)
	}
	_, err = s.Store.RefreshTokenStore.Create(r.Context(), user.Id, PairToken.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(ApiResponse[SignInResponse]{
		Data:    &SignInResponse{RefreshToken: PairToken.RefreshToken.Raw, AccessToken: PairToken.AccessToken.Raw},
		Message: "User Sign In successfully",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r TokenRefreshRequest) Validate() error {
	if r.RefreshToken == "" {
		return errors.New("refresh token is required")
	}
	return nil
}

func (s *ApiServer) TokenRefreshHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TokenRefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		CurrentRefreshToken, err := s.jwtManager.Parse(req.RefreshToken)
		if err != nil {
			slog.Error("failed to parse token: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}
		UserIdstr, err := CurrentRefreshToken.Claims.GetSubject()
		if err != nil {
			slog.Error("failed to get user id from token: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		userId, err := uuid.Parse(UserIdstr)
		if err != nil {
			slog.Error("failed to parse user id: ", err)
			return
		}
		CurrentRefreshTokenRecord, err := s.Store.RefreshTokenStore.ByToken(r.Context(), CurrentRefreshToken, userId)
		if err != nil {
			slog.Error("failed to search refresh token record: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if CurrentRefreshTokenRecord.ExpiresAt.Before(time.Now()) {
			http.Error(w, "refresh token expired", http.StatusBadRequest)
			return
		}

		PairToken, err := s.jwtManager.GenerateTokenPair(userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := s.Store.RefreshTokenStore.Delete(r.Context(), userId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := s.Store.RefreshTokenStore.Create(r.Context(), userId, PairToken.RefreshToken); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(ApiResponse[TokenRefreshResponse]{
			Data:    &TokenRefreshResponse{RefreshToken: PairToken.RefreshToken.Raw, AccessToken: PairToken.AccessToken.Raw},
			Message: "Update token successfully",
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}

}
