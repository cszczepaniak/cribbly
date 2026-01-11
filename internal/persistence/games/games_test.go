package games

import (
	"testing"

	"github.com/cszczepaniak/cribbly/internal/persistence/sqlite"
	"github.com/stretchr/testify/require"
)

func TestGames(t *testing.T) {
	db := sqlite.NewInMemoryForTest(t)
	s := NewService(db)
	require.NoError(t, s.Init(t.Context()))

	t1 := "a"
	t2 := "b"
	t3 := "c"

	g1, err := s.Create(t.Context(), t1, t2)
	require.NoError(t, err)

	g2, err := s.Create(t.Context(), t1, t3)
	require.NoError(t, err)

	g3, err := s.Create(t.Context(), t2, t3)
	require.NoError(t, err)

	gs, err := s.GetForTeam(t.Context(), t1)
	require.NoError(t, err)

	require.Len(t, gs, 2)
	require.Contains(t, gs, g1)
	require.Equal(t, [2]Score{
		{GameID: g1, TeamID: "a", Score: 0},
		{GameID: g1, TeamID: "b", Score: 0},
	}, gs[g1])
	require.Contains(t, gs, g2)
	require.Equal(t, [2]Score{
		{GameID: g2, TeamID: "a", Score: 0},
		{GameID: g2, TeamID: "c", Score: 0},
	}, gs[g2])

	gs, err = s.GetForTeam(t.Context(), t2)
	require.NoError(t, err)

	require.Len(t, gs, 2)
	require.Contains(t, gs, g1)
	require.Equal(t, [2]Score{
		{GameID: g1, TeamID: "a", Score: 0},
		{GameID: g1, TeamID: "b", Score: 0},
	}, gs[g1])
	require.Contains(t, gs, g3)
	require.Equal(t, [2]Score{
		{GameID: g3, TeamID: "b", Score: 0},
		{GameID: g3, TeamID: "c", Score: 0},
	}, gs[g3])

	gs, err = s.GetForTeam(t.Context(), t3)
	require.NoError(t, err)

	require.Len(t, gs, 2)
	require.Contains(t, gs, g2)
	require.Equal(t, [2]Score{
		{GameID: g2, TeamID: "a", Score: 0},
		{GameID: g2, TeamID: "c", Score: 0},
	}, gs[g2])
	require.Contains(t, gs, g3)
	require.Equal(t, [2]Score{
		{GameID: g3, TeamID: "b", Score: 0},
		{GameID: g3, TeamID: "c", Score: 0},
	}, gs[g3])
}
