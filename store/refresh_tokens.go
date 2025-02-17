package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type RefreshTokenStore struct {
	db *sqlx.DB
}

func NewRefreshTokenStore(db *sql.DB) *RefreshTokenStore {
	return &RefreshTokenStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

type RefreshToken struct {
	HashedToken string    `db:"hashed_token"`
	UserId      uuid.UUID `db:"user_id"`
	CreatedAt   time.Time `db:"created_at"`
	ExpiresAt   time.Time `db:"expires_at"`
}

func (f *RefreshTokenStore) Base64Token(token *jwt.Token) (string, error) {
	h := sha256.New()
	h.Write([]byte(token.Raw))
	hashedBytes := h.Sum(nil)
	Base64Tokenhash := base64.StdEncoding.EncodeToString(hashedBytes)
	return Base64Tokenhash, nil
}
func (f *RefreshTokenStore) Create(ctx context.Context, user_id uuid.UUID, token *jwt.Token) (*RefreshToken, error) {
	const insert = `INSERT INTO refresh_tokens (user_id, hashed_token, expires_at) VALUES ($1, $2, $3) RETURNING *;`

	Base64Tokens, err := f.Base64Token(token)
	if err != nil {
		return nil, err
	}
	expires_at, err := token.Claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	var refreshToken RefreshToken
	if err := f.db.GetContext(ctx, &refreshToken, insert, user_id, Base64Tokens, expires_at.Time); err != nil {
		return nil, err
	}
	return &refreshToken, nil

}
func (f *RefreshTokenStore) ByToken(ctx context.Context, token *jwt.Token, user_id uuid.UUID) (*RefreshToken, error) {
	const query = `SELECT * FROM refresh_tokens WHERE user_id = $1 AND hashed_token = $2;`
	var refreshToken RefreshToken
	Base64Tokens, err := f.Base64Token(token)
	if err != nil {
		return nil, err
	}
	if err := f.db.GetContext(ctx, &refreshToken, query, user_id, Base64Tokens); err != nil {
		return nil, err
	}
	if refreshToken.UserId == uuid.Nil {
		return nil, fmt.Errorf("refresh token not found")
	}
	return &refreshToken, nil
}

func (f *RefreshTokenStore) Delete(ctx context.Context, user_id uuid.UUID) (sql.Result, error) {
	const deleteState = `DELETE FROM refresh_tokens WHERE user_id = $1;`

	result, err := f.db.ExecContext(ctx, deleteState, user_id)
	if err != nil {
		return result, fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return result, nil
}
