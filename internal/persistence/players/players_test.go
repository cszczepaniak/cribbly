package players

import (
	"testing"

	"github.com/cszczepaniak/cribbly/internal/persistence/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlayerService(t *testing.T) {
	db := sqlite.NewInMemory(t)
	s := NewService(db)

	require.NoError(t, s.Init(t.Context()))

	require.NoError(t, s.Create(t.Context(), "id1", "Mario Mario"))
	require.NoError(t, s.Create(t.Context(), "id2", "Luigi"))
	require.NoError(t, s.Create(t.Context(), "id3", "Waluigi"))

	players, err := s.Get(t.Context(), "id1", "id2", "id3")
	require.NoError(t, err)

	assert.ElementsMatch(
		t,
		[]Player{{
			ID:   "id1",
			Name: "Mario Mario",
		}, {
			ID:   "id2",
			Name: "Luigi",
		}, {
			ID:   "id3",
			Name: "Waluigi",
		}},
		players,
	)
}
