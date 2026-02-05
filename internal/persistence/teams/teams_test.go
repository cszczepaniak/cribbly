package teams

import (
	"cmp"
	"database/sql"
	"testing"

	"github.com/cszczepaniak/cribbly/internal/assert"

	"github.com/cszczepaniak/cribbly/internal/persistence/sqlite"
)

func TestTeamsRepo(t *testing.T) {
	db := sqlite.NewInMemoryForTest(t)
	s := NewRepository(db)
	assert.NoError(t, s.Init(t.Context()))

	team1, err := s.Create(t.Context(), "team1")
	assert.NoError(t, err)
	team2, err := s.Create(t.Context(), "team2")
	assert.NoError(t, err)
	team3, err := s.Create(t.Context(), "team3")
	assert.NoError(t, err)

	teams, err := s.GetAll(t.Context())
	assert.ElementsMatch(
		t,
		[]Team{team1, team2, team3},
		teams,
		func(x, y Team) int { return cmp.Compare(x.ID, y.ID) },
	)

	team, err := s.Get(t.Context(), team1.ID)
	assert.NoError(t, err)
	assert.Equal(t, team1, team)

	team, err = s.Get(t.Context(), team2.ID)
	assert.NoError(t, err)
	assert.Equal(t, team2, team)

	team, err = s.Get(t.Context(), team3.ID)
	assert.NoError(t, err)
	assert.Equal(t, team3, team)

	assert.NoError(t, s.Delete(t.Context(), team2.ID))

	teams, err = s.GetAll(t.Context())
	assert.ElementsMatch(
		t,
		[]Team{team1, team3},
		teams,
		func(x, y Team) int { return cmp.Compare(x.ID, y.ID) },
	)

	team, err = s.Get(t.Context(), team2.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestTeamsRepo_Rename(t *testing.T) {
	db := sqlite.NewInMemoryForTest(t)
	s := NewRepository(db)
	assert.NoError(t, s.Init(t.Context()))

	team, err := s.Create(t.Context(), "team")
	assert.NoError(t, err)
	assert.Equal(t, "team", team.Name)

	assert.NoError(t, s.Rename(t.Context(), team.ID, "New Name"))

	gotTeam, err := s.Get(t.Context(), team.ID)
	assert.Equal(t, team.ID, gotTeam.ID)
	assert.Equal(t, "New Name", gotTeam.Name)
}
