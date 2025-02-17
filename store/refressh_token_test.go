package store_test

import (
	"asyncapi/apiserver"
	"asyncapi/config"
	"asyncapi/store"
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestRefreshTokenStore(t *testing.T) {
	os.Setenv("ENV", "test")
	conf, err := config.New()
	require.NoError(t, err)
	db, err := store.NewPostgresDb(conf)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE %s;", strings.Join([]string{"users", "refresh_tokens", "reports"}, ",")))
		require.NoError(t, err)

		err = db.Close()
		require.NoError(t, err)
	})
	m, err := migrate.New(
		fmt.Sprintf("file://%s", conf.ProjectRoot),
		conf.DatabaseURL())
	require.NoError(t, err)

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}

	userStore := store.NewUserStore(db)
	ctx := context.Background()
	TokenManager := apiserver.NewJwtManager(conf)

	user, err := userStore.CreateUser(ctx, "Test@mail.ru", "TestPassword")
	require.NoError(t, err)

	TokenPair, err := TokenManager.GenerateTokenPair(user.Id)
	require.NoError(t, err)

	refreshStore := store.NewRefreshTokenStore(db)
	require.NoError(t, err)

	RefreshTokenRecord, err := refreshStore.Create(ctx, user.Id, TokenPair.RefreshToken)
	require.NoError(t, err)

	require.Equal(t, user.Id, RefreshTokenRecord.UserId)
	userByKey, err := refreshStore.ByToken(ctx, TokenPair.RefreshToken, user.Id)
	require.NoError(t, err)
	require.Equal(t, user.Id, userByKey.UserId)
	fmt.Println(userByKey.UserId)

	result, err := refreshStore.Delete(ctx, user.Id)
	require.NoError(t, err)
	RowsEffect, err := result.RowsAffected()
	require.Equal(t, int64(1), RowsEffect)
}
