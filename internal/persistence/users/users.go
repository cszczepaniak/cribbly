package users

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"slices"
	"time"

	"github.com/cszczepaniak/cribbly/internal/persistence/internal/repo"
	"github.com/google/uuid"
)

var (
	ErrUnknownUser    = errors.New("unknown user")
	ErrSessionExpired = errors.New("session expired")
)

type Repository struct {
	repo.Base
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Base: repo.NewBase(db),
	}
}

func (s Repository) WithTx(tx *sql.Tx) Repository {
	s.Base = s.Base.WithTx(tx)
	return s
}

func (s Repository) Init(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS Users (
			Username TEXT,
			PasswordHash BLOB,
			
			PRIMARY KEY (Username)
		)`)
	if err != nil {
		return err
	}

	_, err = s.DB.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS Sessions (
			ID TEXT,
			Username TEXT,
			Expires DATETIME,
			
			PRIMARY KEY (ID)
		)`)
	return err
}

type User struct {
	Name string
}

func (s Repository) GetAll(ctx context.Context) ([]User, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT Username FROM Users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		users = append(users, User{
			Name: name,
		})
	}

	slices.SortFunc(users, func(a, b User) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return users, nil
}

// CreateUser creates the given user.
func (s Repository) CreateUser(ctx context.Context, username, passwordHash string) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO Users (Username, PasswordHash) VALUES (?, ?)`, username, passwordHash)
	if err != nil {
		return err
	}

	return nil
}

// GetPassword returns the persisted hash of the password for the given user.
func (s Repository) GetPassword(ctx context.Context, username string) (string, error) {
	var pw string
	err := s.DB.QueryRowContext(
		ctx,
		`SELECT PasswordHash FROM Users WHERE Username = ?`,
		username,
	).Scan(&pw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUnknownUser
		}

		return "", err
	}

	return pw, nil
}

func (s Repository) ChangePassword(ctx context.Context, username, newPassHash string) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE Users SET PasswordHash = ? WHERE Username = ?`, newPassHash, username)
	return err
}

func (s Repository) DeleteUser(ctx context.Context, username string) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM Users WHERE Username = ?`, username)
	return err
}

func (s Repository) CreateSession(ctx context.Context, username string, expiresIn time.Duration) (string, error) {
	var exists bool
	err := s.DB.QueryRowContext(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM Users WHERE Username = ?)`,
		username,
	).Scan(&exists)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", ErrUnknownUser
	}

	id := uuid.NewString()
	deadline := time.Now().Add(expiresIn)
	_, err = s.DB.ExecContext(
		ctx,
		`INSERT INTO Sessions (ID, Username, Expires) VALUES (?, ?, ?)`,
		id, username, deadline,
	)
	if err != nil {
		return "", err
	}

	return id, nil
}

type Session struct {
	ID       string
	Username string
	expires  time.Time
}

func (s Session) Expired() bool {
	return time.Now().After(s.expires)
}

func (s Repository) GetSession(ctx context.Context, sessionID string) (Session, error) {
	sesh := Session{
		ID: sessionID,
	}

	err := s.DB.QueryRowContext(
		ctx,
		`SELECT Username, Expires FROM Sessions WHERE ID = ?`,
		sessionID,
	).Scan(&sesh.Username, &sesh.expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Session{}, ErrSessionExpired
		}

		return Session{}, err
	}

	if sesh.Expired() {
		_, deleteErr := s.DB.ExecContext(ctx, `DELETE FROM Sessions WHERE ID = ?`, sessionID)
		return Session{}, errors.Join(ErrSessionExpired, deleteErr)
	}

	return sesh, nil
}
