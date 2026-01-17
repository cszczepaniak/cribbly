package teams

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/cribbly/internal/persistence/sqlite"
)

func TestTeamsRepo(t *testing.T) {
	db := sqlite.NewInMemoryForTest(t)
	s := NewRepository(db)
	require.NoError(t, s.Init(t.Context()))

	team1, err := s.Create(t.Context())
	require.NoError(t, err)
	team2, err := s.Create(t.Context())
	require.NoError(t, err)
	team3, err := s.Create(t.Context())
	require.NoError(t, err)

	teams, err := s.GetAll(t.Context())
	assert.ElementsMatch(t, []Team{team1, team2, team3}, teams)

	team, err := s.Get(t.Context(), team1.ID)
	require.NoError(t, err)
	assert.Equal(t, team1, team)

	team, err = s.Get(t.Context(), team2.ID)
	require.NoError(t, err)
	assert.Equal(t, team2, team)

	team, err = s.Get(t.Context(), team3.ID)
	require.NoError(t, err)
	assert.Equal(t, team3, team)

	require.NoError(t, s.Delete(t.Context(), team2.ID))

	teams, err = s.GetAll(t.Context())
	assert.ElementsMatch(t, []Team{team1, team3}, teams)

	team, err = s.Get(t.Context(), team2.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestTeamsRepo_Rename(t *testing.T) {
	db := sqlite.NewInMemoryForTest(t)
	s := NewRepository(db)
	require.NoError(t, s.Init(t.Context()))

	team, err := s.Create(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "Unnamed Team", team.Name)

	require.NoError(t, s.Rename(t.Context(), team.ID, "New Name"))

	gotTeam, err := s.Get(t.Context(), team.ID)
	assert.Equal(t, team.ID, gotTeam.ID)
	assert.Equal(t, "New Name", gotTeam.Name)
}
