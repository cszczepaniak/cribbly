package divisions

import (
	"cmp"
	"database/sql"
	"testing"

	"github.com/cszczepaniak/cribbly/internal/assert"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
)

func TestDivisionsRepo(t *testing.T) {
	db := database.NewInMemory(t)
	s := NewRepository(db)
	assert.NoError(t, s.Init(t.Context()))

	division1, err := s.Create(t.Context())
	assert.NoError(t, err)
	division2, err := s.Create(t.Context())
	assert.NoError(t, err)
	division3, err := s.Create(t.Context())
	assert.NoError(t, err)

	divisions, err := s.GetAll(t.Context())
	assert.ElementsMatch(
		t,
		[]Division{division1, division2, division3},
		divisions,
		func(x, y Division) int { return cmp.Compare(x.ID, y.ID) },
	)

	division, err := s.Get(t.Context(), division1.ID)
	assert.NoError(t, err)
	assert.Equal(t, division1, division)

	division, err = s.Get(t.Context(), division2.ID)
	assert.NoError(t, err)
	assert.Equal(t, division2, division)

	division, err = s.Get(t.Context(), division3.ID)
	assert.NoError(t, err)
	assert.Equal(t, division3, division)

	assert.NoError(t, s.Delete(t.Context(), division2.ID))

	divisions, err = s.GetAll(t.Context())
	assert.ElementsMatch(
		t,
		[]Division{division1, division3},
		divisions,
		func(x, y Division) int { return cmp.Compare(x.ID, y.ID) },
	)

	division, err = s.Get(t.Context(), division2.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestDivisionsRepo_Rename(t *testing.T) {
	db := database.NewInMemory(t)
	s := NewRepository(db)
	assert.NoError(t, s.Init(t.Context()))

	division, err := s.Create(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, "Unnamed Division", division.Name)

	assert.NoError(t, s.Rename(t.Context(), division.ID, "New Name"))

	gotDivision, err := s.Get(t.Context(), division.ID)
	assert.Equal(t, division.ID, gotDivision.ID)
	assert.Equal(t, "New Name", gotDivision.Name)
}
