package apiserver

import (
	"asyncapi/config"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type JwtManager struct {
	config *config.Config
}

func NewJwtManager(config *config.Config) *JwtManager {
	return &JwtManager{
		config: config,
	}
}

var SiningMethod = jwt.SigningMethodHS256

type JwtPair struct {
	AccessToken  *jwt.Token
	RefreshToken *jwt.Token
}
type CustomClaim struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func (j *JwtManager) IsAccesstoken(token *jwt.Token) bool {

	JwtClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	tokenType, ok := JwtClaims["token_type"].(string)
	if !ok {
		return false
	}
	return tokenType == "Access"
}

func (j *JwtManager) Parse(token string) (*jwt.Token, error) {
	parser := jwt.NewParser()
	JwtToken, err := parser.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method != SiningMethod {
			return nil, fmt.Errorf("invalid signing method: %v", t.Header["alg"])
		}
		return []byte(j.config.JwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	return JwtToken, nil
}

func (j *JwtManager) GenerateTokenPair(userId uuid.UUID) (*JwtPair, error) {
	now := time.Now()
	issuer := "http://" + j.config.ApiServerHost + ":" + j.config.ApiServerPort
	jwtAccessToken := jwt.NewWithClaims(SiningMethod, CustomClaim{
		TokenType: "Access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId.String(),
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	},
	)
	Key := []byte(j.config.JwtSecret)

	SignedAccessToken, err := jwtAccessToken.SignedString(Key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}
	AccessToken, err := j.Parse(SignedAccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}

	jwtRefreshToken := jwt.NewWithClaims(SiningMethod, CustomClaim{
		TokenType: "Refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId.String(),
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * 30)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})

	SignedRefreshToken, err := jwtRefreshToken.SignedString(Key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}
	RefreshToken, err := j.Parse(SignedRefreshToken)

	return &JwtPair{AccessToken: AccessToken,
		RefreshToken: RefreshToken}, nil
}
