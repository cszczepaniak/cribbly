package players

import (
	"cmp"
	"database/sql"
	"testing"

	"github.com/google/uuid"

	"github.com/cszczepaniak/gotest/assert"

	"github.com/cszczepaniak/cribbly/internal/moreiter"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
)

func TestPlayerRepo(t *testing.T) {
	db := database.NewInMemory(t)
	s := NewRepository(db)

	assert.NoError(t, s.Init(t.Context()))

	id1, err := s.Create(t.Context(), "Mario", "Mario")
	assert.NoError(t, err)
	id2, err := s.Create(t.Context(), "Luigi", "Mario")
	assert.NoError(t, err)
	id3, err := s.Create(t.Context(), "Waluigi", "Wario")
	assert.NoError(t, err)

	players, err := s.GetAll(t.Context())
	assert.NoError(t, err)

	assert.SliceElemsMatchFunc(
		t,
		[]Player{{
			ID:        id1,
			FirstName: "Mario",
			LastName:  "Mario",
		}, {
			ID:        id2,
			FirstName: "Luigi",
			LastName:  "Mario",
		}, {
			ID:        id3,
			FirstName: "Waluigi",
			LastName:  "Wario",
		}},
		players,
		func(x, y Player) int { return cmp.Compare(x.ID, y.ID) },
	)
}

func TestUpdateName(t *testing.T) {
	db := database.NewInMemory(t)
	s := NewRepository(db)
	assert.NoError(t, s.Init(t.Context()))

	id, err := s.Create(t.Context(), "A", "B")
	assert.NoError(t, err)
	assert.NoError(t, s.UpdateName(t.Context(), id, "C", "D"))

	p, err := s.Get(t.Context(), id)
	assert.NoError(t, err)
	assert.Equal(t, Player{ID: id, FirstName: "C", LastName: "D"}, p)

	assert.ErrorIs(t, s.UpdateName(t.Context(), uuid.NewString(), "X", "Y"), sql.ErrNoRows)
	assert.Error(t, s.UpdateName(t.Context(), id, "", "Y"))
	assert.Error(t, s.UpdateName(t.Context(), id, "X", ""))
}

func TestAssigningPlayers(t *testing.T) {
	db := database.NewInMemory(t)
	s := NewRepository(db)

	assert.NoError(t, s.Init(t.Context()))

	id1, err := s.Create(t.Context(), "Mario", "Mario")
	assert.NoError(t, err)
	id2, err := s.Create(t.Context(), "Luigi", "Mario")
	assert.NoError(t, err)

	players, err := s.GetFreeAgents(t.Context())
	assert.NoError(t, err)

	assert.SliceElemsMatchFunc(
		t,
		[]Player{{
			ID:        id1,
			FirstName: "Mario",
			LastName:  "Mario",
		}, {
			ID:        id2,
			FirstName: "Luigi",
			LastName:  "Mario",
		}},
		players,
		func(x, y Player) int { return cmp.Compare(x.ID, y.ID) },
	)

	mario := players[0]
	teamID := uuid.NewString()

	err = s.AssignToTeam(t.Context(), mario.ID, teamID)
	assert.NoError(t, err)

	// Now mario is on a team, so we should only see luigi as a free agent.
	players, err = s.GetFreeAgents(t.Context())
	assert.NoError(t, err)
	assert.SliceLen(t, players, 1)
	assert.Equal(t,
		Player{
			ID:        id2,
			FirstName: "Luigi",
			LastName:  "Mario",
		},
		players[0],
	)

	// We should also see mario on the team
	players, err = s.GetForTeam(t.Context(), teamID)
	assert.NoError(t, err)
	assert.SliceLen(t, players, 1)
	assert.Equal(t,
		Player{
			ID:        id1,
			FirstName: "Mario",
			LastName:  "Mario",
			TeamID:    teamID,
		},
		players[0],
	)

	// Assigning mario again is an error!
	err = s.AssignToTeam(t.Context(), mario.ID, teamID)
	assert.ErrorIs(t, err, ErrPlayerAlreadyOnATeam)

	// Unassigning mario from a different team is an error.
	assert.Error(t, s.UnassignFromTeam(t.Context(), "not the team", moreiter.Of(mario.ID)))

	// Unassigning mario should make him a free agent again.
	assert.NoError(t, s.UnassignFromTeam(t.Context(), teamID, moreiter.Of(mario.ID)))

	// Unassigning mario again is an error.
	assert.Error(t, s.UnassignFromTeam(t.Context(), teamID, moreiter.Of(mario.ID)))

	players, err = s.GetFreeAgents(t.Context())
	assert.NoError(t, err)

	assert.SliceElemsMatchFunc(
		t,
		[]Player{{
			ID:        id1,
			FirstName: "Mario",
			LastName:  "Mario",
		}, {
			ID:        id2,
			FirstName: "Luigi",
			LastName:  "Mario",
		}},
		players,
		func(x, y Player) int { return cmp.Compare(x.ID, y.ID) },
	)
}
