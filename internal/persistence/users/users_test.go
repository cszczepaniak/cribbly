package users

import (
	"database/sql"
	"testing"
	"time"

	"github.com/cszczepaniak/cribbly/internal/persistence/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsers(t *testing.T) {
	db := sqlite.NewInMemoryForTest(t)
	s := NewService(db)
	require.NoError(t, s.Init(t.Context()))

	err := s.CreateUser(t.Context(), "mario", "secret")
	assert.ErrorIs(t, err, ErrUnknownUser)

	require.NoError(t, s.ReserveUser(t.Context(), "mario"))
	err = s.CreateUser(t.Context(), "mario", "secret")
	require.NoError(t, err)

	pw, err := s.GetPassword(t.Context(), "mario")
	require.NoError(t, err)
	assert.Equal(t, "secret", pw)

	require.NoError(t, s.DeleteUser(t.Context(), "mario"))
	_, err = s.GetPassword(t.Context(), "mario")
	assert.ErrorIs(t, err, ErrUnknownUser)
}

func TestSessions(t *testing.T) {
	db := sqlite.NewInMemoryForTest(t)
	s := NewService(db)
	require.NoError(t, s.Init(t.Context()))

	// User must exist
	_, err := s.CreateSession(t.Context(), "who?", time.Hour)
	assert.ErrorIs(t, err, ErrUnknownUser)

	require.NoError(t, s.ReserveUser(t.Context(), "mario"))
	err = s.CreateUser(t.Context(), "mario", "secret")
	require.NoError(t, err)

	sessionID, err := s.CreateSession(t.Context(), "mario", time.Hour)
	require.NoError(t, err)

	expires, err := s.GetSession(t.Context(), sessionID)
	require.NoError(t, err)
	assert.False(t, expires.Before(time.Now()))

	// Create an already-expired session now
	expiredSessionID, err := s.CreateSession(t.Context(), "mario", -time.Hour)
	require.NoError(t, err)

	_, err = s.GetSession(t.Context(), expiredSessionID)
	assert.ErrorIs(t, err, errSessionExpired)

	// We opportunistically delete the session if we notice it's expired
	_, err = s.GetSession(t.Context(), expiredSessionID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}
