package games

import (
	"fmt"
	"testing"

	"github.com/cszczepaniak/cribbly/internal/notifier"
	"github.com/cszczepaniak/cribbly/internal/persistence/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGames(t *testing.T) {
	n := &notifier.Notifier{}
	db := sqlite.NewInMemoryForTest(t)
	s := NewRepository(db, n)
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

func TestGames_Notifications(t *testing.T) {
	n := &notifier.Notifier{}
	db := sqlite.NewInMemoryForTest(t)
	s := NewRepository(db, n)
	require.NoError(t, s.Init(t.Context()))

	t1 := "a"
	t2 := "b"
	t3 := "c"

	g1, err := s.Create(t.Context(), t1, t2)
	require.NoError(t, err)

	g2, err := s.Create(t.Context(), t1, t3)
	require.NoError(t, err)

	sub, cancel := n.Subscribe()
	t.Cleanup(cancel)

	require.NoError(t, s.UpdateScore(t.Context(), g1, t1, 100))
	// should get notified
	<-sub

	require.NoError(t, s.UpdateScores(t.Context(), g2, t1, 121, t3, 99))
	// should get notified
	<-sub
}

func TestGames_UpdateScores_TeamsMustExistForGame(t *testing.T) {
	n := &notifier.Notifier{}
	db := sqlite.NewInMemoryForTest(t)
	s := NewRepository(db, n)
	require.NoError(t, s.Init(t.Context()))

	t1 := "a"
	t2 := "b"

	g, err := s.Create(t.Context(), t1, t2)
	require.NoError(t, err)

	err = s.UpdateScores(t.Context(), g, t1, 121, "not a team", 100)
	assert.Error(t, err)

	var exists bool
	err = db.QueryRowContext(
		t.Context(),
		"SELECT EXISTS(SELECT 1 FROM Scores WHERE TeamID = ?)",
		"not a team",
	).Scan(&exists)
	assert.False(t, exists)
}

func TestTournamentGames(t *testing.T) {
	n := &notifier.Notifier{}
	db := sqlite.NewInMemoryForTest(t)
	s := NewRepository(db, n)
	require.NoError(t, s.Init(t.Context()))

	err := s.InitializeTournament(t.Context(), 17)
	assert.Error(t, err)

	require.NoError(t, s.InitializeTournament(t.Context(), 32))

	tourney, err := s.LoadTournament(t.Context())
	require.NoError(t, err)

	require.Len(t, tourney.Rounds, 5)
	assert.Len(t, tourney.Rounds[0].Games, 16)
	assert.Len(t, tourney.Rounds[1].Games, 8)
	assert.Len(t, tourney.Rounds[2].Games, 4)
	assert.Len(t, tourney.Rounds[3].Games, 2)
	assert.Len(t, tourney.Rounds[4].Games, 1)

	for i := range 32 {
		if i%2 == 0 {
			require.NoError(t, s.PutTeam1IntoTournamentGame(t.Context(), 0, i/2, fmt.Sprintf("team%d", i)))
		} else {
			require.NoError(t, s.PutTeam2IntoTournamentGame(t.Context(), 0, i/2, fmt.Sprintf("team%d", i)))
		}
	}

	tourney, err = s.LoadTournament(t.Context())
	require.NoError(t, err)
	require.Len(t, tourney.Rounds, 5)
	for i, g := range tourney.Rounds[0].Games {
		t1 := fmt.Sprintf("team%d", 2*i)
		t2 := fmt.Sprintf("team%d", 2*i+1)
		assert.Equal(t, t1, g.TeamIDs[0])
		assert.Equal(t, t2, g.TeamIDs[1])

		// Advance team1
		require.NoError(t, s.SetTournamentGameWinner(t.Context(), g.Round, i, t1))

		if i%2 == 0 {
			require.NoError(t, s.PutTeam1IntoTournamentGame(t.Context(), g.Round+1, i/2, t1))
		} else {
			require.NoError(t, s.PutTeam2IntoTournamentGame(t.Context(), g.Round+1, i/2, t1))
		}
	}

	tourney, err = s.LoadTournament(t.Context())
	require.NoError(t, err)
	require.Len(t, tourney.Rounds, 5)
	for i, g := range tourney.Rounds[0].Games {
		t1 := fmt.Sprintf("team%d", 2*i)
		assert.Equal(t, t1, g.Winner)
	}
	for i, g := range tourney.Rounds[1].Games {
		t1 := fmt.Sprintf("team%d", 4*i)
		t2 := fmt.Sprintf("team%d", 4*i+2)
		assert.Equal(t, t1, g.TeamIDs[0])
		assert.Equal(t, t2, g.TeamIDs[1])
	}
}
