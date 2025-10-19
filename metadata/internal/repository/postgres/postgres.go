package postgres

import (
	"context"
	"database/sql"

	"movieexample.com/metadata/internal/repository"
	"movieexample.com/metadata/pkg/model"
)

// Repository defines a postgres movie metadata repository.
type Repository struct {
	db *sql.DB
}

// New return a new postgres repository.
func New(datasource string) (*Repository, error) {
	db, err := sql.Open("postgres", datasource)
	if err != nil || db.Ping() != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

// Get retrieves movie metadata for by movie id.s
func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var title, description, director string
	row := r.db.QueryRowContext(ctx, "SELECT * FROM movies where id = ?;", id)

	if err := row.Scan(&title, &description, &director); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &model.Metadata{
		ID:          id,
		Title:       title,
		Description: description,
		Director:    director,
	}, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO movies (id, title, description, director) VALUES (?, ?, ?, ?)", id, metadata.Title, metadata.Description, metadata.Director)
	if err != nil {
		return err
	}
	return nil
}
