package apiserver

import (
	"asyncapi/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJwtManager(t *testing.T) {

	conf, err := config.New()
	require.NoError(t, err)

	jwtManager := NewJwtManager(conf)
	require.NoError(t, err)
	user_id := uuid.New()
	TokenPair, err := jwtManager.GenerateTokenPair(user_id)
	require.NoError(t, err)

	AccessSubj, err := TokenPair.AccessToken.Claims.GetSubject()
	require.NoError(t, err)
	require.Equal(t, user_id.String(), AccessSubj)

	AccessIssuers, err := TokenPair.AccessToken.Claims.GetIssuer()
	require.NoError(t, err)
	require.Equal(t, "http//"+conf.ApiServerHost+":"+conf.ApiServerPort, AccessIssuers)

	RefreshSubj, err := TokenPair.RefreshToken.Claims.GetSubject()
	require.NoError(t, err)
	require.Equal(t, user_id.String(), RefreshSubj)
	RefreshIssuers, err := TokenPair.AccessToken.Claims.GetIssuer()
	require.NoError(t, err)

	require.Equal(t, "http//"+conf.ApiServerHost+":"+conf.ApiServerPort, RefreshIssuers)

	ParsedAccesToke, err := jwtManager.Parse(TokenPair.AccessToken.Raw)
	require.NoError(t, err)
	require.Equal(t, TokenPair.AccessToken, ParsedAccesToke)

	ParsedRefreshToke, err := jwtManager.Parse(TokenPair.RefreshToken.Raw)
	require.NoError(t, err)
	require.Equal(t, TokenPair.RefreshToken, ParsedRefreshToke)
}
