package migrations

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/lib/pq"
)

//go:embed movie.sql
var migrations embed.FS

// Migrate migrate the content in the embed file in to
// postgres database
func Migrate(ctx context.Context, db *sql.DB) error {

	data, err := migrations.ReadFile("movie.sql")
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, string(data))
	if err != nil {
		if isTableExistsError(err) {
			return nil
		}
		return err
	}
	return nil
}

// is TableExistsError return true if the migrations do before
// false if the migrations never happend
func isTableExistsError(err error) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {

		return pgErr.Code == "42P07"
	}
	return false
}
