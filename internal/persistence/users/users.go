package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	errUnknownUser    = errors.New("unknown user")
	errSessionExpired = errors.New("session expired")
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) Service {
	return Service{
		db: db,
	}
}

func (s Service) Init(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS Users (
			Username TEXT,
			PasswordHash BLOB,
			
			PRIMARY KEY (Username)
		)`)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS Sessions (
			ID TEXT,
			Username TEXT,
			Expires DATETIME,
			
			PRIMARY KEY (ID)
		)`)
	return err
}

// CreateUser creates the given user. The user must have already been reserved using [ReserveUser],
// otherwise an error is returned.
func (s Service) CreateUser(ctx context.Context, username string, passwordHash []byte) error {
	res, err := s.db.ExecContext(ctx, `UPDATE Users SET PasswordHash = ? WHERE Username = ?`, passwordHash, username)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return errUnknownUser
	}

	return nil
}

// ReserveUser reserves the given username that can be used to register a user later.
func (s Service) ReserveUser(ctx context.Context, username string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO Users (Username) VALUES (?)`, username)
	return err
}

// GetPassword returns the persisted hash of the password for the given user.
func (s Service) GetPassword(ctx context.Context, username string) ([]byte, error) {
	var pw []byte
	err := s.db.QueryRowContext(
		ctx,
		`SELECT PasswordHash FROM Users WHERE Username = ?`,
		username,
	).Scan(&pw)
	if err != nil {
		return nil, err
	}

	return pw, nil
}

func (s Service) DeleteUser(ctx context.Context, username string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM Users WHERE Username = ?`, username)
	return err
}

func (s Service) CreateSession(ctx context.Context, username string, expiresIn time.Duration) (string, error) {
	var exists bool
	err := s.db.QueryRowContext(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM Users WHERE Username = ?)`,
		username,
	).Scan(&exists)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", errUnknownUser
	}

	id := uuid.NewString()
	deadline := time.Now().Add(expiresIn)
	_, err = s.db.ExecContext(
		ctx,
		`INSERT INTO Sessions (ID, Username, Expires) VALUES (?, ?, ?)`,
		id, username, deadline,
	)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s Service) GetSession(ctx context.Context, sessionID string) (string, time.Time, error) {
	var username string
	var expires time.Time

	err := s.db.QueryRowContext(
		ctx,
		`SELECT Username, Expires FROM Sessions WHERE ID = ?`,
		sessionID,
	).Scan(&username, &expires)
	if err != nil {
		return "", time.Time{}, err
	}

	if time.Now().After(expires) {
		_, deleteErr := s.db.ExecContext(ctx, `DELETE FROM Sessions WHERE ID = ?`, sessionID)
		return "", time.Time{}, errors.Join(errSessionExpired, deleteErr)
	}

	return username, expires, nil
}
