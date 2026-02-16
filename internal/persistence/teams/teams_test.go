package teams

import (
	"cmp"
	"database/sql"
	"testing"

	"github.com/cszczepaniak/cribbly/internal/assert"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/google/uuid"
)

func TestTeamsRepo(t *testing.T) {
	db := database.NewInMemory(t)
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
	db := database.NewInMemory(t)
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

func TestTeamsRepo_AssignAndUnassign(t *testing.T) {
	db := database.NewInMemory(t)
	s := NewRepository(db)
	assert.NoError(t, s.Init(t.Context()))

	team1, err := s.Create(t.Context(), "team1")
	assert.NoError(t, err)
	assert.Equal(t, team1.Name, "team1")

	team2, err := s.Create(t.Context(), "team2")
	assert.NoError(t, err)
	assert.Equal(t, team2.Name, "team2")

	divID := uuid.NewString()
	assert.NoError(t, s.AssignToDivision(t.Context(), team1.ID, divID))

	team1WithDiv := team1
	team1WithDiv.DivisionID = divID
	ts, err := s.GetForDivision(t.Context(), divID)
	assert.NoError(t, err)
	assert.Equal(t, ts, []Team{team1WithDiv})

	ts, err = s.GetWithoutDivision(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, ts, []Team{team2})

	assert.NoError(t, s.AssignToDivision(t.Context(), team2.ID, divID))

	team2WithDiv := team2
	team2WithDiv.DivisionID = divID

	ts, err = s.GetForDivision(t.Context(), divID)
	assert.NoError(t, err)
	assert.Equal(t, ts, []Team{team1WithDiv, team2WithDiv})

	ts, err = s.GetWithoutDivision(t.Context())
	assert.NoError(t, err)
	assert.SliceLen(t, ts, 0)

	assert.NoError(t, s.UnassignFromDivision(t.Context(), team1))

	ts, err = s.GetForDivision(t.Context(), divID)
	assert.NoError(t, err)
	assert.Equal(t, ts, []Team{team2WithDiv})

	ts, err = s.GetWithoutDivision(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, ts, []Team{team1})
}
