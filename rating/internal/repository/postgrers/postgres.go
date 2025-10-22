package postgrers

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"movieexample.com/rating/internal/repository"
	model "movieexample.com/rating/pkg"
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

func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]*model.Rating, error) {
	row, err := r.db.QueryContext(ctx, "SELECT user_id, value FROM ratings WHERE record_id = ? AND record_type = ? ", recordID, recordType)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var res []*model.Rating
	for row.Next() {
		var userID string
		var value int32
		if err := row.Scan(&userID, &value); err != nil {
			return nil, err
		}
		res = append(res, &model.Rating{
			UserID: model.UserID(userID),
			Value:  model.RatingValue(value),
		})
	}
	if len(res) == 0 {
		return nil, repository.ErrNotFound
	}
	return res, nil
}

// Put adds a rating for a given record.\
func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO ratings d(record_id, record_type, user_id, value) VALUES (?, ?, ?, ?)",
		recordID, recordType, rating.UserID, rating.Value)
	return err
}
