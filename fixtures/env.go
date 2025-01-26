package fixtures

import (
	"asyncapi/config"
	"asyncapi/store"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

type TestEnv struct {
	DB     *sql.DB
	Config *config.Config
}

func NewTestEnv(t *testing.T) *TestEnv {
	os.Setenv("ENV", "test")

	conf, err := config.New()
	require.NoError(t, err)

	db, err := store.NewPostgresDb(conf)
	require.NoError(t, err)

	return &TestEnv{
		Config: conf,
		DB:     db,
	}
}

func (e *TestEnv) SetUpDatabase(t *testing.T) func(t *testing.T) {
	t.Logf("Строка подключения к базе данных: %s", e.Config.DatabaseURL())

	m, err := migrate.New(
		fmt.Sprintf("file://%s", e.Config.ProjectRoot),
		e.Config.DatabaseURL())
	require.NoError(t, err)
	//выполняю миграю в перед
	if err := m.Up(); err != nil {

		require.NoError(t, err, "Ошибка выполнения миграций")
		 
	}

	return e.TearDownDatabase
}

func (e *TestEnv) TearDownDatabase(t *testing.T) {
	_, err := e.DB.Exec(fmt.Sprintf("DROP TABLE %s;", strings.Join([]string{"users", "refresh_tokens", "reports"}, ",")))
	require.NoError(t, err)

}
