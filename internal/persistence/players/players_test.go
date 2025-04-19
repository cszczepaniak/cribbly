package players

import (
	"testing"

	"github.com/cszczepaniak/cribbly/internal/persistence/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlayerService(t *testing.T) {
	db := sqlite.NewInMemoryForTest(t)
	s := NewService(db)

	require.NoError(t, s.Init(t.Context()))

	id1, err := s.Create(t.Context(), "Mario Mario")
	require.NoError(t, err)
	id2, err := s.Create(t.Context(), "Luigi")
	require.NoError(t, err)
	id3, err := s.Create(t.Context(), "Waluigi")
	require.NoError(t, err)

	players, err := s.Get(t.Context(), id1, id2)
	require.NoError(t, err)

	assert.ElementsMatch(
		t,
		[]Player{{
			ID:   id1,
			Name: "Mario Mario",
		}, {
			ID:   id2,
			Name: "Luigi",
		}},
		players,
	)

	players, err = s.GetAll(t.Context())
	require.NoError(t, err)

	assert.ElementsMatch(
		t,
		[]Player{{
			ID:   id1,
			Name: "Mario Mario",
		}, {
			ID:   id2,
			Name: "Luigi",
		}, {
			ID:   id3,
			Name: "Waluigi",
		}},
		players,
	)
}

func TestAssigningPlayers(t *testing.T) {
	db := sqlite.NewInMemoryForTest(t)
	s := NewService(db)

	require.NoError(t, s.Init(t.Context()))

	id1, err := s.Create(t.Context(), "Mario Mario")
	require.NoError(t, err)
	id2, err := s.Create(t.Context(), "Luigi")
	require.NoError(t, err)

	players, err := s.GetFreeAgents(t.Context())
	require.NoError(t, err)

	assert.ElementsMatch(
		t,
		[]Player{{
			ID:   id1,
			Name: "Mario Mario",
		}, {
			ID:   id2,
			Name: "Luigi",
		}},
		players,
	)

	mario := players[0]
	teamID := uuid.NewString()

	err = s.AssignToTeam(t.Context(), mario, teamID)
	require.NoError(t, err)

	// Now mario is on a team, so we should only see luigi as a free agent.
	players, err = s.GetFreeAgents(t.Context())
	require.NoError(t, err)

	assert.ElementsMatch(
		t,
		[]Player{{
			ID:   id2,
			Name: "Luigi",
		}},
		players,
	)

	// Assigning mario again is an error!
	err = s.AssignToTeam(t.Context(), mario, teamID)
	assert.ErrorIs(t, err, ErrPlayerAlreadyOnATeam)
}
