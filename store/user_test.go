package store_test

import (
	"asyncapi/config"
	"asyncapi/store"
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
	// подключаюсь к базе данных
	db, err := store.NewPostgresDb(conf)
	require.NoError(t, err)
	defer db.Close()
	// создаю обьект своих миграций

	m, err := migrate.New(
		fmt.Sprintf("file://%s", conf.ProjectRoot),
		conf.DatabaseURL())
	require.NoError(t, err)

	//выполняю миграю в перед
	if err := m.Up(); err != nil {
		require.NoError(t, err)
	}

	//создаю обьект который будет вуполнять функционал по USER
	userStore := store.NewUserStore(db)
	ctx := context.Background()
	user, err := userStore.CreateUser(ctx, "third@mail.ru", "testinpassword")
	require.NoError(t, err)

	require.Equal(t, "third@mail.ru", user.Email)
	require.NoError(t, user.ComparePassword("testinpassword"))

}
