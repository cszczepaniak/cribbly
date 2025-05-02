package divisions

import (
	"database/sql"
	"testing"

	"github.com/cszczepaniak/cribbly/internal/persistence/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDivisionsService(t *testing.T) {
	db := sqlite.NewInMemoryForTest(t)
	s := NewService(db)
	require.NoError(t, s.Init(t.Context()))

	division1, err := s.Create(t.Context())
	require.NoError(t, err)
	division2, err := s.Create(t.Context())
	require.NoError(t, err)
	division3, err := s.Create(t.Context())
	require.NoError(t, err)

	divisions, err := s.GetAll(t.Context())
	assert.ElementsMatch(t, []Division{division1, division2, division3}, divisions)

	division, err := s.Get(t.Context(), division1.ID)
	require.NoError(t, err)
	assert.Equal(t, division1, division)

	division, err = s.Get(t.Context(), division2.ID)
	require.NoError(t, err)
	assert.Equal(t, division2, division)

	division, err = s.Get(t.Context(), division3.ID)
	require.NoError(t, err)
	assert.Equal(t, division3, division)

	require.NoError(t, s.Delete(t.Context(), division2.ID))

	divisions, err = s.GetAll(t.Context())
	assert.ElementsMatch(t, []Division{division1, division3}, divisions)

	division, err = s.Get(t.Context(), division2.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestDivisionsService_Rename(t *testing.T) {
	db := sqlite.NewInMemoryForTest(t)
	s := NewService(db)
	require.NoError(t, s.Init(t.Context()))

	division, err := s.Create(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "Unnamed Division", division.Name)

	require.NoError(t, s.Rename(t.Context(), division.ID, "New Name"))

	gotDivision, err := s.Get(t.Context(), division.ID)
	assert.Equal(t, division.ID, gotDivision.ID)
	assert.Equal(t, "New Name", gotDivision.Name)
}
