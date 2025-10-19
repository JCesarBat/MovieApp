package migrations

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"movieexample.com/pkg/config"
)

func TestMigration(t *testing.T) {
	conf := config.GetConfig()
	conn, err := sql.Open("postgres", conf.GetDBConnectionString())
	require.NoError(t, err)

	ctx := context.Background()
	err = Migrate(ctx, conn)
	require.NoError(t, err)
}
