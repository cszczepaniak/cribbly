package users

import (
	"testing"
	"time"

	"github.com/cszczepaniak/cribbly/internal/assert"

	"github.com/cszczepaniak/cribbly/internal/persistence/sqlite"
)

func TestUsers(t *testing.T) {
	db := sqlite.NewInMemoryForTest(t)
	s := NewRepository(db)
	assert.NoError(t, s.Init(t.Context()))

	err := s.CreateUser(t.Context(), "mario@mario.com", "secret")
	assert.NoError(t, err)

	err = s.CreateUser(t.Context(), "luigi@mario.com", "secret")
	assert.NoError(t, err)

	pw, err := s.GetPassword(t.Context(), "mario@mario.com")
	assert.NoError(t, err)
	assert.Equal(t, "secret", pw)

	all, err := s.GetAll(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, []User{{
		Name: "luigi@mario.com",
	}, {
		Name: "mario@mario.com",
	}}, all)

	assert.NoError(t, s.DeleteUser(t.Context(), "mario@mario.com"))
	_, err = s.GetPassword(t.Context(), "mario@mario.com")
	assert.ErrorIs(t, err, ErrUnknownUser)

	all, err = s.GetAll(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, []User{{
		Name: "luigi@mario.com",
	}}, all)
}

func TestSessions(t *testing.T) {
	db := sqlite.NewInMemoryForTest(t)
	s := NewRepository(db)
	assert.NoError(t, s.Init(t.Context()))

	// User must exist
	_, err := s.CreateSession(t.Context(), "who?", time.Hour)
	assert.ErrorIs(t, err, ErrUnknownUser)

	err = s.CreateUser(t.Context(), "mario@mario.com", "secret")
	assert.NoError(t, err)

	sessionID, err := s.CreateSession(t.Context(), "mario@mario.com", time.Hour)
	assert.NoError(t, err)

	sesh, err := s.GetSession(t.Context(), sessionID)
	assert.NoError(t, err)
	assert.False(t, sesh.Expired())

	// Create an already-expired session now
	expiredSessionID, err := s.CreateSession(t.Context(), "mario@mario.com", -time.Hour)
	assert.NoError(t, err)

	_, err = s.GetSession(t.Context(), expiredSessionID)
	assert.ErrorIs(t, err, ErrSessionExpired)

	// We opportunistically delete the session if we notice it's expired
	_, err = s.GetSession(t.Context(), expiredSessionID)
	assert.ErrorIs(t, err, ErrSessionExpired)
}
