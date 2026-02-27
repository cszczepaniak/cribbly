package roomcodes

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/cszczepaniak/cribbly/internal/persistence/database"
)

var (
	ErrCodeNotFound = errors.New("room code not found")
	ErrCodeExpired  = errors.New("room code expired")
)

type Repository struct {
	db database.Database
}

func NewRepository(db database.Database) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) Init(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS RoomCodes (
			Code TEXT,
			Expires DATETIME,

			PRIMARY KEY (Code)
		)`)
	return err
}

type RoomCode struct {
	Code    string
	Expires time.Time
}

func (rc RoomCode) Expired() bool {
	return time.Now().After(rc.Expires)
}

// Create inserts a new room code that expires at the given time.
func (r Repository) Create(ctx context.Context, code string, expiresAt time.Time) error {
	return r.db.ExecVoid(ctx, `
		INSERT INTO RoomCodes (Code, Expires) VALUES (?, ?)
	`, code, expiresAt)
}

// Get returns the room code with the given code value.
func (r Repository) Get(ctx context.Context, code string) (RoomCode, error) {
	var rc RoomCode
	err := r.db.QueryRowContext(ctx, `
		SELECT Code, Expires FROM RoomCodes WHERE Code = ?
	`, code).Scan(&rc.Code, &rc.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return RoomCode{}, ErrCodeNotFound
		}
		return RoomCode{}, err
	}

	if rc.Expired() {
		// Best-effort cleanup of expired codes.
		_, delErr := r.db.ExecContext(ctx, `DELETE FROM RoomCodes WHERE Code = ?`, code)
		return RoomCode{}, errors.Join(ErrCodeExpired, delErr)
	}

	return rc, nil
}

// Validate returns true if the given code exists and has not expired.
func (r Repository) Validate(ctx context.Context, code string) (bool, error) {
	_, err := r.Get(ctx, code)
	if err != nil {
		if errors.Is(err, ErrCodeNotFound) || errors.Is(err, ErrCodeExpired) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Latest returns the most recently expiring, non-expired room code, if any.
func (r Repository) Latest(ctx context.Context) (RoomCode, error) {
	var rc RoomCode
	err := r.db.QueryRowContext(ctx, `
		SELECT Code, Expires FROM RoomCodes
		WHERE Expires > CURRENT_TIMESTAMP
		ORDER BY Expires DESC
		LIMIT 1
	`).Scan(&rc.Code, &rc.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return RoomCode{}, ErrCodeNotFound
		}
		return RoomCode{}, err
	}
	return rc, nil
}

