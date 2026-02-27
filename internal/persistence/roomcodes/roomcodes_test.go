package roomcodes

import (
	"testing"
	"time"

	"github.com/cszczepaniak/cribbly/internal/assert"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
)

func newTestRepo(t *testing.T) Repository {
	t.Helper()

	db := database.NewInMemory(t)
	repo := NewRepository(db)
	assert.NoError(t, repo.Init(t.Context()))

	return repo
}

func TestCreateAndGetRoomCode(t *testing.T) {
	repo := newTestRepo(t)

	expires := time.Now().Add(time.Hour)
	assert.NoError(t, repo.Create(t.Context(), "ABC123", expires))

	rc, err := repo.Get(t.Context(), "ABC123")
	assert.NoError(t, err)
	assert.Equal(t, "ABC123", rc.Code)
}

func TestGetNonexistentRoomCodeReturnsNotFound(t *testing.T) {
	repo := newTestRepo(t)

	_, err := repo.Get(t.Context(), "MISSING")
	assert.ErrorIs(t, err, ErrCodeNotFound)
}

func TestExpiredRoomCodeReturnsExpiredErrorAndDeletes(t *testing.T) {
	repo := newTestRepo(t)

	expired := time.Now().Add(-time.Hour)
	assert.NoError(t, repo.Create(t.Context(), "OLD123", expired))

	_, err := repo.Get(t.Context(), "OLD123")
	assert.ErrorIs(t, err, ErrCodeExpired)

	ok, err := repo.Validate(t.Context(), "OLD123")
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestValidateRoomCode(t *testing.T) {
	repo := newTestRepo(t)

	expires := time.Now().Add(time.Hour)
	assert.NoError(t, repo.Create(t.Context(), "GOOD1", expires))

	ok, err := repo.Validate(t.Context(), "GOOD1")
	assert.NoError(t, err)
	if !ok {
		t.Fatalf("expected room code to be valid")
	}

	ok, err = repo.Validate(t.Context(), "BAD1")
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestLatestReturnsMostRecentNonExpired(t *testing.T) {
	repo := newTestRepo(t)

	now := time.Now()

	assert.NoError(t, repo.Create(t.Context(), "OLD", now.Add(-time.Hour)))
	assert.NoError(t, repo.Create(t.Context(), "MID", now.Add(time.Hour)))
	assert.NoError(t, repo.Create(t.Context(), "NEW", now.Add(2*time.Hour)))

	rc, err := repo.Latest(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, "NEW", rc.Code)
}
