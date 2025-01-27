package store

import (
	"asyncapi/config"
	"strings"

	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestUserStore(t *testing.T) {
	os.Setenv("ENV", "test")
	conf, err := config.New()
	require.NoError(t, err)

	db, err := NewPostgresDb(conf)
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

	userStore := NewUserStore(db)
	ctx := context.Background()
	user, err := userStore.CreateUser(ctx, "Test2@mail.ru", "testinpassword")
	require.NoError(t, err)

	require.Equal(t, "Test2@mail.ru", user.Email)
	require.NoError(t, user.ComparePassword("testinpassword"))

	user2, err := userStore.ById(ctx, user.Id)
	require.Equal(t, user, user2)
	user2, err = userStore.ByEmail(ctx, user.Email)
	require.Equal(t, user, user2)

}
